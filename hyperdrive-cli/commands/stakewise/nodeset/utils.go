package nodeset

import (
	"fmt"
	"net/mail"

	swapi "github.com/nodeset-org/hyperdrive-stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// CheckRegistrationStatus checks the registration status of the node with NodeSet and prompts the user to register if not already done
func CheckRegistrationStatus(c *cli.Context, hd *client.HyperdriveClient, sw *client.StakewiseClient) error {
	// Check if the node is already registered
	shouldRegister, err := checkRegistrationStatusImpl(hd, sw)
	if err != nil {
		return err
	}
	if !shouldRegister {
		return nil
	}

	// Prompt for registration
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to register your node now?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	return RegisterNodeImpl(c, sw)
}

// Returns true if the node should register because it hasn't yet and is able to
func checkRegistrationStatusImpl(hd *client.HyperdriveClient, sw *client.StakewiseClient) (bool, error) {
	// Get wallet response
	response, err := hd.Api.Wallet.Status()
	if err != nil {
		return false, err
	}

	// Make sure we have a wallet loaded
	if !response.Data.WalletStatus.Wallet.IsLoaded {
		fmt.Println("The node wallet has not been initialized yet. Please run `hyperdrive wallet status` to learn more.")
		return false, nil
	}

	// Get the registration status
	resp, err := sw.Api.Nodeset.RegistrationStatus()
	if err != nil {
		return false, err
	}
	switch resp.Data.Status {
	case swapi.NodesetRegistrationStatus_Unknown:
		fmt.Println("Hyperdrive couldn't check your node's registration status:")
		fmt.Println(resp.Data.ErrorMessage)
		fmt.Println("Please try again later.")
	case swapi.NodesetRegistrationStatus_NoWallet:
		fmt.Println("Your node can't be registered until you have a node wallet initialized. Please run `hyperdrive wallet init` or `hyperdrive wallet recover` first.")
	case swapi.NodesetRegistrationStatus_Unregistered:
		fmt.Println("Your node is not currently registered with NodeSet.")
		return true, nil
	case swapi.NodesetRegistrationStatus_Registered:
		fmt.Println("Your node is registered with NodeSet.")
	}
	return false, nil
}

// Registers the node with NodeSet
func RegisterNodeImpl(c *cli.Context, sw *client.StakewiseClient) error {
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
	response, err := sw.Api.Nodeset.RegisterNode(email)
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
