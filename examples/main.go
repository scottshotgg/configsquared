package main

import (
	"fmt"

	config "github.com/scottshotgg/config/examples/config"
)

// Make a new type so that we can implement config.Validator
type validator struct{}

// Implement functions as required; only function is ValidateTimeout for this example
func (v *validator) ValidateTimeout(timeout int) error {
	if timeout < 0 || timeout > 600 {
		return fmt.Errorf("invalid timeout: timeout must be between 0 and 600; %d", timeout)
	}

	return nil
}

func main() {
	// Since we have required validation, we need to pass our validator struct on Parse.
	// Otherwise, config will Parse without validation and return only a Config struct pointer
	var c, err = config.Parse(&validator{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Using port:", c.Port())

	if c.Mock() {
		fmt.Println("Mocking db")
	}

	fmt.Println("HTTP timeout is", c.Timeout(), "seconds")

	fmt.Println(c.At())

	fmt.Printf("%+v\n", c.Remote())
}
