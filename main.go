package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func appendIfNotExists(slice []string, item string) []string {
	for _, el := range slice {
		if el == item {
			return slice
		}
	}
	return append(slice, item)
}

func main() {
	// Get repository link from command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <repository link>")
		os.Exit(1)
	}
	repoLink := os.Args[1]

	// Clone repository
	dir, err := filepath.Abs(filepath.Base(repoLink))
	if err != nil {
		fmt.Println("Error: Failed to get absolute path for directory")
		os.Exit(1)
	}
	cmd := exec.Command("git", "clone", repoLink, dir)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Failed to clone repository %s\n", repoLink)
		os.Exit(1)
	}

	// Navigate into directory
	if err := os.Chdir(dir); err != nil {
		fmt.Println("Error: Failed to navigate to directory")
		os.Exit(1)
	}

	// Use git log -p output
	cmd = exec.Command("git", "log", "-p")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error: Failed to run git log -p")
		os.Exit(1)
	}

	access_regex := regexp.MustCompile("AKIA[0-9A-Z]{16}")
	access_key_matches := make([]string, 0)

	secret_regex := regexp.MustCompile("[0-9a-zA-Z/+]{40}")
	secret_key_matches := make([]string, 0)

	reader := bytes.NewReader(out)
	scanner := bufio.NewScanner(reader)

	// Scan for access keys and secret keys
	for scanner.Scan() {
		line := scanner.Text()
		matches := access_regex.FindAllString(line, -1)
		for _, match := range matches {
			access_key_matches = appendIfNotExists(access_key_matches, match)
		}

		matches = secret_regex.FindAllString(line, -1)
		for _, match := range matches {
			secret_key_matches = appendIfNotExists(secret_key_matches, match)
		}
	}

	fmt.Print("Printing access keys and secret keys, if any...\n\n")
	// Check all combinations of access keys and secret keys and print the valid ones
	if len(access_key_matches) > 0 && len(secret_key_matches) > 0 {
		for _, match := range access_key_matches {
			for _, match2 := range secret_key_matches {
				checkKeys(match, match2)
			}
		}
	}

	fmt.Print("\nPrinting access keys, if any...\n\n")

	// Print the valid access keys
	if len(access_key_matches) > 0 {
		for _, match := range access_key_matches {
			checkAccessKeys(match)
		}
	}

}

// Func to check if access key is valid
func checkAccessKeys(accessKeyID string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session: ", err)
		return
	}

	svc := sts.New(sess)

	// Call GetAccessKeyInfo API
	input := &sts.GetAccessKeyInfoInput{
		AccessKeyId: &accessKeyID,
	}

	_, err = svc.GetAccessKeyInfo(input)

	if err != nil {
		return
	}

	fmt.Println(accessKeyID, "is an access key present in the repository")
}

// Func to check if key pair is valid
func checkKeys(accessKeyID string, secretToken string) {

	creds := credentials.NewStaticCredentials(accessKeyID, secretToken, "")

	// Create a new session with AWS credentials
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: creds,
	})
	if err != nil {
		fmt.Println("Error creating AWS session: ", err)
		return
	}

	// Create a new STS client
	svc := sts.New(sess)

	input := &sts.GetCallerIdentityInput{}
	_, err = svc.GetCallerIdentity(input)

	if err != nil {
		return
	}

	fmt.Println(accessKeyID, "and", secretToken, "is a key pair present in the repository")
}
