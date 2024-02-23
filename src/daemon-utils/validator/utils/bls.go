package utils

import (
	"sync"

	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

// Initialize BLS support
var initBLS sync.Once

func InitializeBls() error {
	var err error
	initBLS.Do(func() {
		err = eth2types.InitBLS()
	})
	return err
}
