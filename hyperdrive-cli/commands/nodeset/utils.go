package nodeset

import (
	"fmt"
	"net/mail"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/urfave/cli/v2"
)

// CheckRegistrationStatus checks the registration status of the node with NodeSet and prompts the user to register if not already done
// Returns whether or not the caller should continue with its operation after this check completes, or if it should exit
func CheckRegistrationStatus(c *cli.Context, hd *client.HyperdriveClient) (bool, error) {
	// Check if the node is already registered
	hasWallet, shouldRegister, err := checkRegistrationStatusImpl(hd)
	if err != nil {
		return false, err
	}
	if !shouldRegister {
		return hasWallet, nil
	}

	// Prompt for registration
	if c.Bool(utils.YesFlag.Name) {
		return false, nil
	}
	if !utils.Confirm("Would you like to register your node now?") {
		fmt.Println("Cancelled.")
		return false, nil
	}

	return hasWallet, registerNodeImpl(c, hd)
}

// Returns true if the node should register because it hasn't yet and is able to
func checkRegistrationStatusImpl(hd *client.HyperdriveClient) (bool, bool, error) {
	// Check wallet status
	_, ready, err := utils.CheckIfWalletReady(hd)
	if err != nil {
		return false, false, err
	}
	if !ready {
		return false, false, nil
	}

	// Get the registration status
	resp, err := hd.Api.NodeSet.GetRegistrationStatus()
	if err != nil {
		return false, false, err
	}
	switch resp.Data.Status {
	case api.NodeSetRegistrationStatus_Unknown:
		fmt.Println("Hyperdrive couldn't check your node's registration status:")
		fmt.Println(resp.Data.ErrorMessage)
		fmt.Println("Please try again later.")
	case api.NodeSetRegistrationStatus_NoWallet:
		fmt.Println("Your node can't be registered until you have a node wallet initialized. Please run `hyperdrive wallet init` or `hyperdrive wallet recover` first.")
	case api.NodeSetRegistrationStatus_Unregistered:
		fmt.Println("Your node is not currently registered with NodeSet.")
		return true, true, nil
	case api.NodeSetRegistrationStatus_Registered:
		fmt.Println("Your node is registered with NodeSet.")
	}
	return true, false, nil
}

// Registers the node with NodeSet
func registerNodeImpl(c *cli.Context, hd *client.HyperdriveClient) error {
	// Get the email
	email := c.String(RegisterEmailFlag.Name)
	if email == "" {
		for {
			email = utils.Prompt("Enter the email address you'd like to register with NodeSet:", "^.*$", "Invalid email address, try again")
			_, err := mail.ParseAddress(email)
			if err == nil {
				break
			}
			fmt.Println("Invalid email address, try again")
		}
	}

	// Register the node
	response, err := hd.Api.NodeSet.RegisterNode(email)
	if err != nil {
		return fmt.Errorf("error registering node: %w", err)
	}

	// Validation
	if response.Data.NotWhitelisted {
		fmt.Printf("Your node has not been whitelisted in the NodeSet account for email address [%s]. Please go to the NodeSet website and add your node to your account's whitelist.\n", email)
		return nil
	}

	fmt.Println("Node successfully registered.")
	return err
}
