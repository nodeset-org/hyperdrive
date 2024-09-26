package service

import (
	"fmt"
)

// View the Hyperdrive service stats
func serviceStats() error {
	fmt.Println("No longer supported - please run `docker stats -a` instead.")
	return nil
}
