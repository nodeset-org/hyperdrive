package constellation

import (
	"fmt"
	"math/big"

	"github.com/nodeset-org/hyperdrive/services"
	"github.com/nodeset-org/hyperdrive/validator"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "todo",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Fetch URL from user-settings.yml
		// TODO: Setup Proteus to get Beacon Client node
		beaconClientManager, error := services.NewBeaconClientManager("https://eth.llamarpc.com")
		if error != nil {
			fmt.Printf("Error: %s\n", error)
		}

		// EC Manager
		// ecManager, err := NewExecutionClientManager(cfg)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating executon client manager: %w", err)
		// }

		// Rocket Pool
		// rp, err := rocketpool.NewRocketPool(
		// 	ecManager,
		// 	common.HexToAddress(cfg.Smartnode.GetStorageAddress()),
		// 	common.HexToAddress(cfg.Smartnode.GetMulticallAddress()),
		// 	common.HexToAddress(cfg.Smartnode.GetBalanceBatcherAddress()),
		// )

		depositContext := &validator.DepositContext{}
		depositContext.Bc = beaconClientManager.PrimaryBc
		depositContext.Salt = big.NewInt(0)

		fmt.Printf("depositContext: %v\n", depositContext)
		depositContext.SubmitInitDeposit()
	},
}
