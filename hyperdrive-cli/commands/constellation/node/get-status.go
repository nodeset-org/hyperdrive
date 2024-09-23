package node

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

func getStatus(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Get the node status
	response, err := cs.Api.Node.GetRegistrationStatus()
	if err != nil {
		return err
	}

	// Print the status
	if response.Data.Registered {
		fmt.Println("Node is registered with Constellation.")
	} else {
		fmt.Println("Node is not registered with Constellation yet.")
	}
	return nil
}
