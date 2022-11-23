# avocet
Avocet is a cli tool for packing directory and appending manifest to file header.

Install avocet command to the $GOPATH/bin directory:
```bash
go install
```
Run command:
```bash
# package current directory
avocet create .
# this will produce a file named {cwd}.zip
```