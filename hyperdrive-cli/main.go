package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hyperdrive",
	Short: "Hyperdrive initialization and Rocketpool service status check",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hyperdrive")
	},
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "todo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("hello world\n")
		// TODO: Fetch URL from user-settings.yml
		// TODO: Setup Proteus to get Beacon Client node
		// beaconClientManager, error := services.NewBeaconClientManager("https://eth.llamarpc.com")
		// if error != nil {
		// 	fmt.Printf("Error: %s\n", error)
		// }

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

		// depositContext := &validator.DepositContext{}
		// depositContext.Bc = beaconClientManager.PrimaryBc
		// depositContext.Salt = big.NewInt(0)

		// fmt.Printf("depositContext: %v\n", depositContext)
		// depositContext.SubmitInitDeposit()
	},
}

func init() {
	rootCmd.AddCommand(InitCmd)
}

// This function is executed prior to any Cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
