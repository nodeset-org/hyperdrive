package minipool

import (
	"crypto/rand"
	"fmt"
	"math/big"

	csconfig "github.com/nodeset-org/hyperdrive-constellation/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/utils/math"
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
			fmt.Println("- Your node is not registered with NodeSet. Please whitelist your node with your nodeset.io account, register with `hyperdrive ns r`, then try again.")
		}
		if response.Data.NotWhitelistedWithConstellation {
			fmt.Println("- Your node is not registered with Constellation. Please register it with `hyperdrive cs n r`, then try again.")
		}
		if response.Data.InsufficientBalance {
			additionalEthRequired := new(big.Int).Sub(response.Data.LockupAmount, response.Data.NodeBalance)
			fmt.Printf("- You don't have enough ETH in your node wallet to make a new minipool. Your node requires at least %.6f more ETH (plus enough for gas).\n", eth.WeiToEth(additionalEthRequired))
		}
		if response.Data.MaxMinipoolsReached {
			fmt.Println("- You have reached the maximum number of minipools you can create.")
		}
		if response.Data.InsufficientLiquidity {
			fmt.Println("- Constellation doesn't have enough ETH or RPL liquidity in its vaults to fund a new minipool. Please wait for more deposits to its vaults.")
		}
		if response.Data.MissingExitMessage {
			fmt.Println("- nodeset.io is missing a signed exit message for at least one of your previous validators. If you recently created a new minipool, you'll have to wait until it's been given an index on the Beacon Chain; Hyperdrive will upload a signed exit message automatically once an index is available.")
		}
		if response.Data.NodeSetDepositingDisabled {
			fmt.Println("- NodeSet has currently disabled new minipool creation.")
		}
		if response.Data.RocketPoolDepositingDisabled {
			fmt.Println("- Rocket Pool has currently disabled new minipool creation.")
		}
		if response.Data.IncorrectNodeAddress {
			fmt.Println("- You have a different node registered for Constellation. You can only create minipools from that node.")
		}
		if response.Data.InvalidPermissions {
			fmt.Println("- Your user account does not have the required permissions to use this Constellation deployment. Please reach out to the NodeSet administrators for help.")
		}

		return nil
	}

	// Print a note about requirements
	fmt.Printf("%sNOTE: Creating a new minipool will require a temporary deposit of %.2f ETH. It will be returned to you when your minipool passes the scrub check and your node issues its second deposit (or you call `stake` manually with the `hyperdrive cs m k` command).\n%s", terminal.ColorYellow, eth.WeiToEth(response.Data.LockupAmount), terminal.ColorReset)

	// Prompt for confirmation
	if !(c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to continue?")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Print salt and minipool address info
	if c.String(saltFlag.Name) != "" {
		fmt.Printf("Using custom salt %s, your minipool address will be %s.\n\n", c.String(saltFlag.Name), response.Data.MinipoolAddress.Hex())
	}

	// Run the TX
	validated, err := tx.HandleTx(c, hd, response.Data.TxInfo,
		"Exiting this minipool capital cannot be done until your minipool has been *active* on the Beacon Chain for 256 epochs (approx. 27 hours).",
		"creating minipool",
		"Creating minipool...",
	)
	if err != nil {
		return err
	}
	if !validated {
		// Prompt for saving the key
		if c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to save the private key for this validator to disk? You'll need to do this if you plan to submit that transaction later and want to be able to attest with this validator.") {
			_, err = cs.Api.Wallet.CreateValidatorKey(response.Data.ValidatorPubkey, response.Data.Index, 1)
			if err != nil {
				return fmt.Errorf("error saving validator key to disk: %w", err)
			}
		}
		return nil
	}

	// Save the validator key to disk
	_, err = cs.Api.Wallet.CreateValidatorKey(response.Data.ValidatorPubkey, response.Data.Index, 1)
	if err != nil {
		fmt.Printf("%sError saving validator key to disk: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
		fmt.Println("Please rebuild your Constellation wallet manually to regenerate and save the validator key. Otherwise you will not be able to validate when it becomes active on the Beacon Chain!")
		return nil
	}

	// Log & return
	fmt.Printf("Minipool created successfully! You have temporarily locked up %.2f ETH.\n", math.RoundDown(eth.WeiToEth(response.Data.LockupAmount), 6))
	fmt.Printf("Your new minipool's address is: %s\n", response.Data.MinipoolAddress)
	fmt.Printf("The validator pubkey is: %s\n\n", response.Data.ValidatorPubkey.HexWithPrefix())

	fmt.Println("Your minipool is now in Initialized status.")
	fmt.Println("Once the remaining ETH has been assigned to your minipool from Rocket Pool's staking pool, it will move to Prelaunch status.")
	fmt.Printf("After that, it will move to Staking status once %s have passed.\n", response.Data.ScrubPeriod)
	fmt.Println("You can watch its progress using `hyperdrive s dl cs-tasks`.")

	fmt.Println()

	fmt.Println("Your Constellation Validator Client must be restarted in order to load the new validator key so it can begin attesting once it has been activated on the Beacon Chain.")
	if c.Bool(utils.YesFlag.Name) || utils.Confirm("Would you like to restart the Constellation Validator Client now?") {
		_, err := hd.Api.Service.RestartContainer(string(csconfig.ContainerID_ConstellationValidator))
		if err != nil {
			fmt.Printf("%sWARNING: Error restarting Constellation Validator Client: %s%s\n", terminal.ColorRed, err.Error(), terminal.ColorReset)
			fmt.Println("Please restart the Constellation Validator Client manually before your validator becomes active in order to load the new validator key.")
			fmt.Printf("%sIf you don't restart it, you will miss attestations and lose ETH!%s\n", terminal.ColorYellow, terminal.ColorReset)
		} else {
			fmt.Println("Successfully restarted the Constellation Validator Client. Your new validator key is now loaded.")
			return nil
		}
	} else {
		fmt.Println("Please restart the Constellation Validator Client manually before your validator becomes active in order to load the new validator key.")
		fmt.Printf("%sIf you don't restart it, you will miss attestations and lose ETH, and may be ejected from NodeSet!%s\n", terminal.ColorYellow, terminal.ColorReset)
	}

	return nil
}
