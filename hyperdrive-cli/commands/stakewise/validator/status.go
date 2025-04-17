package validator

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	swapi "github.com/nodeset-org/hyperdrive-stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/utils/input"
	"github.com/urfave/cli/v2"
)

var (
	vaultFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "vault",
		Aliases: []string{"v"},
		Usage:   "Provide the address of a vault if you only want information for validators in that specific vault, instead of all vaults",
	}
)

func getStatus(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading Hyperdrive config: %w", err)
	}
	if !cfg.StakeWise.Enabled.Value {
		fmt.Println("The StakeWise module is not enabled in your Hyperdrive configuration.")
		return nil
	}

	// Check wallet status
	_, ready, err := utils.CheckIfWalletReady(hd)
	if err != nil {
		return err
	}
	if !ready {
		return nil
	}

	// Get validator status
	var vault *common.Address
	if c.IsSet(vaultFlag.Name) {
		vaultString := c.String(vaultFlag.Name)
		vaultAddress, err := input.ValidateAddress("vault", vaultString)
		if err != nil {
			return fmt.Errorf("invalid vault address [%s]: %w", vaultString, err)
		}
		vault = &vaultAddress
	}
	statusResponse, err := sw.Api.Validator.Status(vault)
	if err != nil {
		return fmt.Errorf("error while getting validator status: %w", err)
	}
	if statusResponse.Data.NotRegisteredWithNodeSet {
		fmt.Println("You are not registered with NodeSet yet.")
		fmt.Println("Please register with NodeSet using the `hyperdrive nodeset register-node` command.")
		return nil
	}
	if statusResponse.Data.InvalidPermissions {
		fmt.Println("Your node doesn't have permission to use the StakeWise module yet.")
		return nil
	}
	if len(statusResponse.Data.Vaults) == 0 {
		fmt.Println("There are no StakeWise vaults for this deployment yet.")
		return nil
	}

	// Print status info
	totalInfo := make([]*swapi.ValidatorInfo, 0)
	vaultMap := map[beacon.ValidatorPubkey]*swapi.VaultInfo{}
	for _, vault := range statusResponse.Data.Vaults {
		for _, validator := range vault.Validators {
			vaultMap[validator.Pubkey] = vault
			totalInfo = append(totalInfo, validator)
		}
	}
	if len(totalInfo) == 0 {
		fmt.Println("You don't have any validators registered with StakeWise yet.")
		return nil
	}

	for _, validator := range totalInfo {
		vault := vaultMap[validator.Pubkey]
		var seenString string
		if validator.HasBeaconIndex {
			seenString = fmt.Sprintf("Yes (Index %s)", validator.Index)
		} else {
			seenString = "No"
		}
		fmt.Printf("Validator %s%s%s\n", terminal.ColorBlue, validator.Pubkey, terminal.ColorReset)
		fmt.Printf("  Vault: %s (%s%s%s)\n", vault.Name, terminal.ColorGreen, vault.Address.Hex(), terminal.ColorReset)
		fmt.Printf("  Seen on Beacon Yet: %s\n", seenString)
		if validator.HasBeaconIndex {
			fmt.Printf("  Status: %s\n", getBeaconStatusLabel(validator.State))
			fmt.Printf("  Balance: %.6f\n", eth.GweiToEth(float64(validator.Balance)))
		}
		fmt.Println()
	}

	return nil
}

func getBeaconStatusLabel(state beacon.ValidatorState) string {
	switch state {
	case beacon.ValidatorState_ActiveExiting:
		return "Active (Exiting in Progress)"
	case beacon.ValidatorState_ActiveOngoing:
		return "Active"
	case beacon.ValidatorState_ActiveSlashed:
		return "Slashed (Exit in Progress)"
	case beacon.ValidatorState_ExitedSlashed:
		return "Slashed (Exited)"
	case beacon.ValidatorState_ExitedUnslashed:
		return "Exited (Withdrawal Pending)"
	case beacon.ValidatorState_PendingInitialized:
		return "Seen on Beacon, Waiting for More Deposits"
	case beacon.ValidatorState_PendingQueued:
		return "In Beacon Activation Queue"
	case beacon.ValidatorState_WithdrawalDone:
		return "Exited and Withdrawn"
	case beacon.ValidatorState_WithdrawalPossible:
		return "Exited (Waiting for Wihdrawal)"
	default:
		return fmt.Sprintf("<Unknown Beacon Status: %s>", state)
	}
}
