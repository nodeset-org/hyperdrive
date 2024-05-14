package service

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

// View the Hyperdrive service logs
func serviceLogs(c *cli.Context, aliasedNames ...string) error {
	// Handle name aliasing
	serviceNames := []string{}
	for _, name := range aliasedNames {
		trueName := name
		switch name {
		case "eth1", "el", "execution":
			trueName = "ec"
		case "cc", "cl", "bc", "eth2", "beacon", "consensus":
			trueName = "bn"
		case "vc":
			trueName = "validator"
		}
		serviceNames = append(serviceNames, trueName)
	}

	// Get Hyperdrive client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Print service logs
	return hd.PrintServiceLogs(getComposeFiles(c), c.String("tail"), serviceNames...)
}
