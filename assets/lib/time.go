// Code generated by github.com/scottshotgg/configsquared; DO NOT EDIT.

// It is in your best interest NOT to edit this file as it will be removed if re-generated.

package config

import (
	"errors"
	"time"
)

// TODO: will need a format specifier

type timeFlag struct {
	set    bool
	value  time.Time
	sValue string
}

// If flag is not provided it will not get to this function
func (t *timeFlag) Set(x string) error {
	// TODO: implement this
	return errors.New("type not implemented; time")

	// // Set the string value
	// d.sValue = x

	// // Parse the value from the provided string

	// // Set the actual value
	// // t.value = int(value)

	// // Mark the flag as set
	// d.set = true

	// return nil
}

func (t *timeFlag) String() string {
	return t.sValue
}
