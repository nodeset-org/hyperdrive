package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/config"
)

func configHyperdrive(installPath string) error {
	mgr := config.NewConfigManager(installPath)

	// Try to load an existing config
	cfg, alreadyExists, err := mgr.LoadOrCreateConfig(false)
	if err != nil {
		return fmt.Errorf("error configuring Hyperdrive: %w", err)
	}

	// Done if it already exists
	if alreadyExists {
		fmt.Println("Configuration already exists.")
		return nil
	}

	// Save it
	err = mgr.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("error saving Hyperdrive configuration: %w", err)
	}

	fmt.Println("Config file saved.")
	return nil
}
