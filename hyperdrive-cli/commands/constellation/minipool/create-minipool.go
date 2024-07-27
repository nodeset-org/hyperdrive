package minipool

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/urfave/cli/v2"
)

var (
	saltFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "salt",
		Aliases: []string{"l"},
		Usage:   "An optional seed to use when generating the new minipool's address. Use this if you want it to have a custom vanity address.",
	}
)

func createMinipool(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Get the minipool salt
	var salt *big.Int
	saltString := c.String(saltFlag.Name)
	if saltString != "" {
		var success bool
		salt, success = big.NewInt(0).SetString(saltString, 0)
		if !success {
			return fmt.Errorf("invalid minipool salt: %s", saltString)
		}
	} else {
		buffer := make([]byte, 32)
		_, err = rand.Read(buffer)
		if err != nil {
			return fmt.Errorf("error generating random salt: %w", err)
		}
		salt = big.NewInt(0).SetBytes(buffer)
	}

	// Build the TX
	response, err := cs.Api.Minipool.Create(salt)
	if err != nil {
		return err
	}

	// Verify
	if !response.Data.CanCreate {
		fmt.Println("Cannot create new minipool:")
		if response.Data.NotRegisteredWithNodeSet {
			fmt.Println("Your node is not registered with NodeSet. Please whitelist your node with your nodeset.io account, register with `hyperdrive ns r`, then try again.")
		}
		if response.Data.NotWhitelistedWithConstellation {
			fmt.Println("Your node is not registered with Constellation. Please register it with `hyperdrive cs n r`, then try again.")
		}
		if response.Data.InsufficientBalance {
			additionalEthRequired := new(big.Int).Sub(response.Data.NodeBalance, response.Data.LockupAmount)
			fmt.Printf("You don't have enough ETH in your node wallet to make a new minipool. Your node requires at least %.6f more ETH.\n", eth.WeiToEth(additionalEthRequired))
		}
		if response.Data.InsufficientLiquidity {
			fmt.Println("Constellation doesn't have enough ETH or RPL liquidity in its vaults to fund a new minipool. Please wait for more deposits to its vaults.")
		}
		if response.Data.InsufficientMinipoolCount {
			fmt.Println("Your node is not allowed to make any more minipools. Ensure you have uploaded signed exit messages for your existing minipools first, then try again.")
		}
		if response.Data.NodeSetDepositingDisabled {
			fmt.Println("NodeSet has currently disabled new minipool creation.")
		}
		if response.Data.RocketPoolDepositingDisabled {
			fmt.Println("Rocket Pool has currently disabled new minipool creation.")
		}
		return nil
	}

	// Print a note about requirements
	fmt.Printf("%sNOTE: Creating a new minipool will require a temporary deposit of %.2f ETH. It will be returned to you when your minipool passes the scrub check and your node issues its second deposit (or you call `stake` manually with the `hyperdrive cs m k` command).\n", terminal.ColorYellow, eth.WeiToEth(response.Data.LockupAmount), terminal.ColorReset)

	// Prompt for confirmation
	if !(c.Bool("yes") || utils.Confirm("Would you like to continue?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	return nil
}
