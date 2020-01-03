// Code generated by github.com/scottshotgg/configsquared; DO NOT EDIT.

// It is in your best interest NOT to edit this file as it will be removed if re-generated.

package config

import (
	"errors"
)

type byteFlag struct {
	set    bool
	value  byte
	sValue string
}

// TODO: flesh this out more

// If flag is not provided it will not get to this function
func (b *byteFlag) Set(x string) error {
	var bs = []byte(x)

	if len(bs) > 1 {
		return errors.New("byte flag only accepts one byte")
	}

	b.value = bs[0]

	b.set = true

	return nil
}

func (b *byteFlag) String() string {
	return b.sValue
}
