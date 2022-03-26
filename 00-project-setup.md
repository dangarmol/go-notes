# Project Setup

## Some .zshrc config

```bash
# Go configuration
# export GOROOT=/usr/local/go  # Not needed unless using non-default install path
export GOPATH=/Users/daniel/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

## Go Commands

- `go get <repo_url>` installs the contents of a repo on the GOPATH folder.
- `go run <.go file>` builds and runs file, but doesn’t save binary.
- `go build <folder>` creates a binary on the current path.
- `go install <folder>` creates a binary on the $GOPATH/bin folder.
- The source code for the Go standard library packages is on `/usr/local/go/src`.
