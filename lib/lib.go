package lib

import "github.com/spf13/pflag"

const (
	DefaultTest = "defaultTest"
)

type Lib struct {
	Test string
}

func New() *Lib {
	return &Lib{
		Test: DefaultTest,
	}
}

func (c *Lib) Flags() (f *pflag.FlagSet) {
	return c.PrefixedFlags("lib")
}

func (c *Lib) PrefixedFlags(x string) (f *pflag.FlagSet) {
	prefix := func(f string) string {
		if x != "" {
			return x + "." + f
		}
		return f
	}
	f = pflag.NewFlagSet("lib", pflag.ContinueOnError)
	f.StringVar(&c.Test,
		prefix("test"), c.Test, "test lib flag")
	return
}
