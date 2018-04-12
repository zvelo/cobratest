package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const name = "cobratest"

var (
	// these values should be set by the linker as args to `go build`
	version   string
	gitCommit string
	buildDate string
)

// rootCmd is the _only_ global variable and it must be so that other commands
// can use rootCmd.AddCommand in their init() funcs
var rootCmd = cobra.Command{
	Use:     name,
	Short:   "example app",
	Version: fmt.Sprintf("%s (commit %s; built %s; %s)", version, gitCommit, buildDate, runtime.Version()),
}

func init() {
	// this configures viper to automatically check env vars for all keys it
	// knows about
	viper.AutomaticEnv()

	// this tells viper to make these replacements when trying to figure out
	// what the env var name for a key should be
	viper.SetEnvKeyReplacer(strings.NewReplacer(
		"-", "_",
		".", "_",
	))
}

// Execute is the main entrypoint into the app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
