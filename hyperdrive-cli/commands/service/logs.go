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
		case "ec", "el", "execution":
			trueName = "eth1"
		case "cc", "cl", "bc", "bn", "beacon", "consensus":
			trueName = "eth2"
		case "vc":
			trueName = "validator"
		}
		serviceNames = append(serviceNames, trueName)
	}

	// Get RP client
	hd := client.NewClientFromCtx(c)

	// Print service logs
	return hd.PrintServiceLogs(getComposeFiles(c), c.String("tail"), serviceNames...)
}
