package sub

import (
	"flag"
	"fmt"
	"time"
)

const (
	// DefaultURL is the default URL
	DefaultURL = "https://zvelo.com/"

	// DefaultTTL is the default TTL
	DefaultTTL = 1 * time.Minute
)

// Config configures Sub
type Config struct {
	URL string
	TTL time.Duration
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
}
