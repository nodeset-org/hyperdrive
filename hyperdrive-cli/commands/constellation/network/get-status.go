package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/urfave/cli/v2"
)

func getStats(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	cs, err := client.NewConstellationClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Get the network stats
	response, err := cs.Api.Network.Stats()
	if err != nil {
		return err
	}

	oneEth := big.NewInt(1e18)
	rplStakeInEth := big.NewInt(0).Set(response.Data.SuperNodeRplStake)
	rplStakeInEth.Mul(rplStakeInEth, response.Data.RplPrice)
	rplStakeInEth.Div(rplStakeInEth, oneEth)

	ethBond := big.NewInt(int64(response.Data.ActiveMinipoolCount * 8)) // Hard-code LEB8s for now
	bondedRatio := big.NewInt(0)
	borrowedRatio := big.NewInt(0)

	if ethBond.Cmp(common.Big0) > 0 {
		bondedRatio.Div(rplStakeInEth, ethBond)
		borrowedRatio.Div(bondedRatio, big.NewInt(24/8))
	}

	// Print the stats
	fmt.Printf("%s======== Network Settings =========%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Printf("Active Minipool Limit:    %d per node\n", response.Data.ValidatorLimit)
	fmt.Println()

	fmt.Printf("%s========== Deposit Pools ==========%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Printf("Constellation ETH Pool:   %.6f ETH\n", eth.WeiToEth(response.Data.ConstellationEthBalance))
	fmt.Printf("Constellation RPL Pool:   %.6f RPL\n", eth.WeiToEth(response.Data.ConstellationRplBalance))
	fmt.Printf("Rocket Pool Deposit Pool: %.6f ETH\n", eth.WeiToEth(response.Data.RocketPoolEthBalance))
	fmt.Printf("Rocket Pool Utilization:  %.2f%%\n", eth.WeiToEth(response.Data.RocketPoolEthUtilizationRate)*100)
	fmt.Printf("Minipool Queue Length:    %d\n", response.Data.MinipoolQueueLength)
	fmt.Printf("Minipool Queue Capacity:  %.6f ETH\n", eth.WeiToEth(response.Data.MinipoolQueueCapacity))
	fmt.Printf("RPL Price:                %.6f ETH\n", eth.WeiToEth(response.Data.RplPrice))
	fmt.Println()

	fmt.Printf("%s=========== Super Node ===========%s\n", terminal.ColorGreen, terminal.ColorReset)
	fmt.Printf("Address:              %s%s%s\n", terminal.ColorBlue, response.Data.SuperNodeAddress, terminal.ColorReset)
	fmt.Printf("RPL Staked:           %.6f RPL\n", eth.WeiToEth(response.Data.SuperNodeRplStake))
	fmt.Printf("Bonded Coll. Ratio:   %.2f%%\n", eth.WeiToEth(bondedRatio)*100)
	fmt.Printf("Borrowed Coll. Ratio: %.2f%%\n", eth.WeiToEth(borrowedRatio)*100)
	fmt.Printf("Subnodes:             %d\n", response.Data.SubnodeCount)
	fmt.Printf("Active Minipools:     %d\n", response.Data.ActiveMinipoolCount)
	fmt.Printf("    Initialized:      %d\n", response.Data.InitializedMinipoolCount)
	fmt.Printf("    Prelaunch:        %d\n", response.Data.PrelaunchMinipoolCount)
	fmt.Printf("    Staking:          %d\n", response.Data.StakingMinipoolCount)
	fmt.Printf("    Dissolved:        %d\n", response.Data.DissolvedMinipoolCount)
	fmt.Printf("Finalized Minipools:  %d\n", response.Data.FinalizedMinipoolCount)

	return nil
}
