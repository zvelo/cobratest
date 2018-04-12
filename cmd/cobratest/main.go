package main

import "zvelo.io/cobratest/cmd/cobratest/cmd"

var (
	// these values should be set by the linker as args to `go build`
	version   = "unknown"
	gitCommit = "unknown"
	buildDate = "unknown"
)

// main is a simple one liner
func main() {
	cmd.Execute(version, gitCommit, buildDate)
}
