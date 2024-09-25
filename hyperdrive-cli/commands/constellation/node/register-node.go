package node

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
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
	csRegResponse, err := cs.Api.Node.GetRegistrationStatus()
	if err != nil {
		return err
	}
	if csRegResponse.Data.Registered {
		fmt.Println("Your node is already registered with Constellation.")
		return nil
	}

	// Check if the node's registered with NodeSet
	nsRegResponse, err := hd.Api.NodeSet.GetRegistrationStatus()
	if err != nil {
		return err
	}
	if nsRegResponse.Data.Status != api.NodeSetRegistrationStatus_Registered {
		fmt.Println("Your node is not registered with NodeSet. Please register it with `hyperdrive nodeset register-node` first.")
		return nil
	}

	// Print the notice
	fmt.Printf("%sNOTE:\n", terminal.ColorYellow)
	fmt.Println("Your NodeSet account can only have one node registered with Constellation at a time.")
	fmt.Println("Registration requires a special off-chain signature from the Constellation administrator.")
	fmt.Printf("If you proceed, Hyperdrive will retrieve this signature for you automatically which will lock this node as your account's Constellation node (even if you aren't ready to submit the registration transaction yet).\n\n%s", terminal.ColorReset)

	// Prompt for confirmation
	if !(c.Bool("yes") || utils.ConfirmWithIAgree("Are you ready to assign this node as your account's Constellation node?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Get the registration TX
	response, err := cs.Api.Node.Register()
	if err != nil {
		return err
	}
	fmt.Println("Signature retrieved. This node is now locked in as your account's Constellation node.")

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
		"Are you ready to register this node with Constellation?",
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
