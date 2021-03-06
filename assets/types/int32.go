// Code generated by github.com/scottshotgg/configsquared; DO NOT EDIT.

// It is in your best interest NOT to edit this file as it will be removed if re-generated.

package config

import (
	"strconv"
)

type int32Flag struct {
	set    bool
	value  int32
	sValue string
}

// If flag is not provided it will not get to this function
func (i *int32Flag) Set(x string) error {
	// Set the string value
	i.sValue = x

	// Parse the value from the provided string
	var value, err = strconv.ParseInt(i.sValue, 10, 32)
	if err != nil {
		// TODO: test this out
		return err
	}

	// Set the actual value
	i.value = int32(value)

	// Mark the flag as set
	i.set = true

	return nil
}

func (i *int32Flag) String() string {
	return i.sValue
}
