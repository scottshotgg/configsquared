// Code generated by github.com/scottshotgg/configsquared; DO NOT EDIT.

// It is in your best interest NOT to edit this file as it will be removed if re-generated.

package config

import (
	"errors"
)

type complex128Flag struct {
	set    bool
	value  complex64
	sValue string
}

// TODO: flesh this out more

// If flag is not provided it will not get to this function
func (c *complex128Flag) Set(x string) error {
	// TODO: implement this
	return errors.New("type not implemented; complex128")

	// if len(x) > 1 {
	// 	return errors.New("byte flag only accepts one byte")
	// }

	// b.value = byte(x[0])

	// b.set = true

	// return nil
}

func (c *complex128Flag) String() string {
	return c.sValue
}
