# Go-Code-Scanner
Go tool to scan Github repositories for Access Key IDs and Secret Tokens

## How it works


* The program first accepts the repository link as a command line argument

* Then, the repository is cloned locally under the root dir by running `git clone <repo-link>` within the prog.

* The prog then cd's into the repo directory and runs the `git log -p` flag to access all previous commits.

* Regex is used to match against AWS access IDs and secret tokens by setting up the `access_regex` and `secret_regex` variables. These scan `git log -p` 's output.

* Hence we get 2 slices: `access_key_matches` and `secret_key_matches` which contain the matched items

* All possible combinations of access key id and secret token are passed to the `checkKeys()` function. This function sets up an aws session using the credentials passed to it and then attempts to make requests to `sts.GetCallerIdentity()` method which returns an output if the credentials work else an error which isn't displayed in this case.

* Then, all access key IDs are checked individually by passing them to the `checkAccessKeys()` func which uses the `sts.GetAccessKeyInfo()` method internally. Each valid access key ID is then printed out.


## How to run the code


* Run the code from the proj's root dir by running: `go run main.go <repo-link>`
