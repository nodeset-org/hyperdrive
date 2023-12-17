package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/shared/config"
)

func configHyperdrive(installPath string) error {
	mgr := config.NewConfigManager(installPath)

	_, alreadyExists, err := mgr.LoadOrCreateConfig()
	if err != nil {
		return fmt.Errorf("error configuring Hyperdrive: %w", err)
	}

	if alreadyExists {
		fmt.Println("Configuration already exists.")
		return nil
	}
	fmt.Println("Config file saved.")
	return nil
}
