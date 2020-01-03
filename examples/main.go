package main

import (
	"fmt"

	config "github.com/scottshotgg/config/examples/config"
)

// Make a new type so that we can implement config.Validator
type validator struct{}

// Implement functions as required; only function is ValidateTimeout for this example
func (v *validator) ValidateRetries(retries int) error {
	if retries < 0 || retries > 10 {
		return fmt.Errorf("invalid retries: amount must be between 0 and 10; %d", retries)
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

	fmt.Println("Amount of retries available:", c.Retries())

	fmt.Println("HTTP timeout is", c.Timeout())

	fmt.Println("At:", c.At())

	fmt.Println("When:", c.When())

	fmt.Printf("Remote: %+v\n", c.Remote())

	fmt.Println("Idc about this value:", c.Idc())
}
