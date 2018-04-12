package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd is the _only_ global variable and it must be so that other commands
// can use rootCmd.AddCommand in their init() funcs
var rootCmd = cobra.Command{
	Use:   "cobratest",
	Short: "example app",
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
func Execute(version, commit, built string) {
	rootCmd.Version = fmt.Sprintf("%s (commit %s; built %s; %s)", version, commit, built, runtime.Version())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
