package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
)

func start(installPath string, composeFiles []string) error {
	client, err := client.NewClient(installPath)
	if err != nil {
		return fmt.Errorf("error running start: %w", err)
	}
	return client.StartService(composeFiles)
}
