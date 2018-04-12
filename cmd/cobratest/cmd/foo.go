package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"zvelo.io/cobratest/sub"
)

// fooCmd contains the configuration for the foo command and does not pollute
// the global namespace
type fooCmd struct {
	Addr string
	Sub  *sub.Config
}

func init() {
	// initialize a new fooCmd with defaults
	f := fooCmd{
		Addr: ":http",
		Sub:  sub.DefaultConfig(),
	}

	// create the cobra command
	cmd := cobra.Command{
		Use:   "foo",
		Short: "bar",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			// this is necessary to set the config values from env and flags
			// only publicly exported struct fields can be configured
			return viper.Unmarshal(&f)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return f.run()
		},
	}

	// add command flags, the name of the flag must match the name of the struct
	// field for viper.Unmarshal to work. alternatively, the mapstructure struct
	// "tag" must match the flag name in which case the struct field name can be
	// anything you want. github.com/mitchellh/mapstructure
	// also note that the Var varients (e.g. StringVar) aren't particularly
	// helpful as, even if they set the destination properly, the value isn't
	// processed by viper. just stick with the viper.Unmarshal instead.
	cmd.Flags().StringP("addr", "a", f.Addr, "listening address")

	// add flags exported by the subpackage. the prefix must match the struct
	// field name, or the mapstructure struct tag
	// cmd.Flags().AddGoFlagSet(f.Sub.FlagSet("sub"))
	cmd.Flags().AddFlagSet(f.Sub.PFlagSet("sub"))

	// cmd.MarkFlagRequired("sub.url")

	// this is how we tell viper what the configuration keys are
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	// register this command with the root command
	rootCmd.AddCommand(&cmd)
}

func (f *fooCmd) run() error {
	fmt.Println("addr:", f.Addr)

	s := sub.New(sub.WithConfig(f.Sub))
	s.Do()

	return nil
}
