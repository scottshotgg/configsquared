package config

import "time"

// this would be package config

type Config struct {
	port  string
	mongo *Mongo
}

type Mongo struct {
	port    string
	addr    string
	timeout time.Duration
}

// embedded structs are the same to parse and generate config for
// same for go struct generation
// only the flags are different
