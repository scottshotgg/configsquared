// Code generated by github.com/scottshotgg/configsquared; DO NOT EDIT.

// It is in your best interest NOT to edit this file as it will be removed if re-generated.

package config

import (
	"errors"
	"net"
)

type ipFlag struct {
	set    bool
	value  net.IP
	sValue string
}

// If flag is not provided it will not get to this function
func (i *ipFlag) Set(x string) error {
	// Set the string value
	i.sValue = x

	// Parse the value from the provided string
	var value = net.ParseIP(i.sValue)
	if value == nil {
		// TODO: test this out
		return errors.New("invalid ip: " + x)
	}

	// Set the actual value
	i.value = value

	// Mark the flag as set
	i.set = true

	return nil
}

func (i *ipFlag) String() string {
	return i.sValue
}
