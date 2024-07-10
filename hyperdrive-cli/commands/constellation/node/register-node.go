package node

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	"github.com/urfave/cli/v2"
)

func registerNode(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Check if the node's already registered
	statusResponse, err := cs.Api.Node.GetRegistrationStatus()
	if err != nil {
		return err
	}
	if statusResponse.Data.Registered {
		fmt.Println("Your node is already registered with Constellation.")
		return nil
	}

	// Get the registration TX
	response, err := cs.Api.Node.Register()
	if err != nil {
		return err
	}

	// Check for status issues
	if response.Data.NotRegisteredWithNodeSet {
		fmt.Println("Your node has not been registered with your NodeSet account yet. Please whitelist your node's address in your nodeset.io dashboard, then run `hyperdrive nodeset register-node`.")
		return nil
	}
	if response.Data.NotAuthorized {
		fmt.Println("Your NodeSet account is not permitted to register with Constellation yet.")
		return nil
	}

	// Run the TX
	validated, err := tx.HandleTx(c, hd, response.Data.TxInfo,
		"Are you sure you register this node with Constellation?",
		"registering with Constellation",
		"Registering with Constellation...",
	)
	if err != nil {
		return err
	}
	if !validated {
		return nil
	}

	// Log & return
	fmt.Println("Your node successfully registered with Constellation. You can now create and run minipools.")
	return nil
}
