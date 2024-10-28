package minipool

import (
	"fmt"

	csapi "github.com/nodeset-org/hyperdrive-constellation/shared/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	ncli "github.com/rocket-pool/node-manager-core/cli/utils"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/urfave/cli/v2"
)

var (
	stakeMinipoolsFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "minipools",
		Aliases: []string{"m"},
		Usage:   "A comma-separated list of addresses for minipools to stake (or 'all' to stake all available minipools)",
	}
)

func stakeMinipools(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading Hyperdrive config: %w", err)
	}
	if !cfg.Constellation.Enabled.Value {
		fmt.Println("The Constellation module is not enabled in your Hyperdrive configuration.")
		return nil
	}

	// Build the TX
	response, err := cs.Api.Minipool.Stake()
	if err != nil {
		return err
	}

	// Get stakeable minipools
	stakeableMinipools := []csapi.MinipoolStakeDetails{}
	scrubPeriodMinipools := []csapi.MinipoolStakeDetails{}
	for _, minipool := range response.Data.Details {
		if minipool.CanStake {
			stakeableMinipools = append(stakeableMinipools, minipool)
		} else if minipool.StillInScrubPeriod {
			scrubPeriodMinipools = append(scrubPeriodMinipools, minipool)
		}
	}

	if len(scrubPeriodMinipools) > 0 {
		fmt.Println("The following minipools are still in the scrub period and cannot be staked yet:")
		for _, mp := range scrubPeriodMinipools {
			fmt.Printf("%s (%s left)\n", mp.Address.Hex(), mp.RemainingTime)
		}
		fmt.Println()
	}

	// Check for stakeable minipools
	if len(stakeableMinipools) == 0 {
		fmt.Println("No minipools can be staked.")
		return nil
	}

	// Get selected minipools
	options := make([]ncli.SelectionOption[csapi.MinipoolStakeDetails], len(stakeableMinipools))
	for i, mp := range stakeableMinipools {
		option := &options[i]
		option.Element = &stakeableMinipools[i]
		option.ID = fmt.Sprint(mp.Address)
		option.Display = fmt.Sprintf("%s (%s until dissolved)", mp.Address.Hex(), mp.TimeUntilDissolve)
	}
	selectedMinipools, err := utils.GetMultiselectIndices(c, stakeMinipoolsFlag.Name, options, "Please select a minipool to stake:")
	if err != nil {
		return fmt.Errorf("error determining minipool selection: %w", err)
	}

	// Validation
	txInfos := make([]*eth.TransactionInfo, len(selectedMinipools))
	for i, mp := range selectedMinipools {
		txInfos[i] = mp.TxInfo
	}

	fmt.Println()
	fmt.Println("NOTE: Your Constellation Validator Client must be restarted after this process so it loads the new validator keys.")
	fmt.Println("Since you are manually staking the minipools, this must be done manually.")
	fmt.Println("When you have finished staking all your minipools, please restart your Constellation Validator Client.")
	fmt.Println()

	// Run the TXs
	validated, err := tx.HandleTxBatch(c, hd, txInfos,
		fmt.Sprintf("Are you sure you want to stake %d minipools?", len(selectedMinipools)),
		func(i int) string {
			return fmt.Sprintf("stake of minipool %s", selectedMinipools[i].Address.Hex())
		},
		"Staking minipools...",
	)
	if err != nil {
		return err
	}
	if !validated {
		return nil
	}

	// Log & return
	fmt.Println("Successfully staked all selected minipools.")
	return nil
}
