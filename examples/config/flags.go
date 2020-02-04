// Code generated by github.com/scottshotgg/configsquared; DO NOT EDIT.

// It is in your best interest NOT to edit this file as it will be removed if re-generated.

package config

import "flag"

type flags struct {
	port  stringFlag
	ports stringArrayFlag
}

func newFlags() *flags {
	var f flags

	flag.Var(&f.port, "port", "")
	flag.Var(&f.ports, "ports", "")

	// nested struct fields

	return &f
}

func (f *flags) required() {

}

func (f *flags) defaults() {

}

func (f *flags) toConfig() Config {
	return Config{
		port:  f.port.value,
		ports: f.ports.value,
	}
}
