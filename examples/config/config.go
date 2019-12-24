// Code generated by github.com/scottshotgg/configsquared; DO NOT EDIT.

// It is in your best interest NOT to edit this file as it will be removed if re-generated.

package config

import (
	"flag"
)

type Config struct {
	port    string
	mock    bool
	timeout int
}

var (
	// Internal singleton config to avoid re-parse if called multiple times
	c Config
)

func (c *Config) Port() string {
	return c.port
}
func (c *Config) Mock() bool {
	return c.mock
}
func (c *Config) Timeout() int {
	return c.timeout
}

func Parse(v Validator) (*Config, error) {
	if !flag.Parsed() {
		var f = newFlags()

		flag.Parse()

		f.required()
		f.defaults()

		c = f.toConfig()
		return &c, c.validate(v)
	}

	return &c, nil
}

func (c *Config) validate(v Validator) error {
	var err error

	err = v.ValidateTimeout(c.timeout)
	if err != nil {
		return err
	}

	return nil
}
