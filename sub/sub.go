package sub

import (
	"flag"
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

const (
	// DefaultURL is the default URL
	DefaultURL = "https://zvelo.com/"

	// DefaultTTL is the default TTL
	DefaultTTL = 1 * time.Minute
)

// Config configures Sub
type Config struct {
	URL                 string
	TTL                 time.Duration
	Bool0, Bool1, Bool2 bool
}

// DefaultConfig returns a Config set to default values
func DefaultConfig() *Config {
	return &Config{
		URL: DefaultURL,
		TTL: DefaultTTL,
	}
}

// FlagSet returns the flags that can be used to configure Sub with defaults
// taken from c
func (c Config) FlagSet(prefix string) *flag.FlagSet {
	fs := flag.NewFlagSet("sub", flag.ContinueOnError)

	fs.String(prefix+".url", c.URL, "the sub url")
	fs.Duration(prefix+".ttl", c.TTL, "the sub ttl")
	fs.Bool(prefix+".bool0", c.Bool0, "sub bool 0")
	fs.Bool(prefix+".bool1", c.Bool1, "sub bool 1")
	fs.Bool(prefix+".bool2", c.Bool2, "sub bool 2")

	return fs
}

// PFlagSet returns the flags that can be used to configure Sub with defaults
// taken from c
// pflag is much more featureful than the standard go "flag" package
func (c Config) PFlagSet(prefix string) *pflag.FlagSet {
	fs := pflag.NewFlagSet("sub", pflag.ContinueOnError)

	fs.StringP(prefix+".url", "u", c.URL, "the sub url")

	fs.DurationP(prefix+".ttl", "t", c.TTL, "the sub ttl")
	fs.MarkShorthandDeprecated(prefix+".ttl", "please use --"+prefix+".ttl only") // #nosec

	fs.BoolP(prefix+".bool0", "b", c.Bool0, "sub bool 0")
	fs.MarkDeprecated(prefix+".bool0", "use "+prefix+".bool2 instead") // #nosec

	fs.BoolP(prefix+".bool1", "o", c.Bool1, "sub bool 1")
	fs.MarkHidden(prefix + ".bool1") // #nosec

	fs.BoolP(prefix+".bool2", "l", c.Bool2, "sub bool 2")

	return fs
}

// Sub is an example subpackage
type Sub struct {
	Config *Config
}

// Option is used to configure Sub
type Option func(*Config)

// WithConfig returns an option that copies values from val to the Sub Config
func WithConfig(val *Config) Option {
	return func(c *Config) {
		*c = *val
	}
}

// New constructs a new Sub
func New(opts ...Option) *Sub {
	c := DefaultConfig()

	for _, o := range opts {
		o(c)
	}

	return &Sub{Config: c}
}

// Do does something
func (s *Sub) Do() {
	fmt.Println("sub.url:", s.Config.URL)
	fmt.Println("sub.ttl:", s.Config.TTL)
	fmt.Println("sub.bool0:", s.Config.Bool0)
	fmt.Println("sub.bool1:", s.Config.Bool1)
	fmt.Println("sub.bool2:", s.Config.Bool2)
}
