package gas

import (
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/gas"
	"github.com/urfave/cli/v2"
)

func GetMaxFees(c *cli.Context, hd *client.HyperdriveClient, simResult eth.SimulationResult) (*big.Int, *big.Int, error) {
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting Hyperdrive configuration: %w", err)
	}
	if isNew {
		return nil, nil, fmt.Errorf("Settings file not found. Please run `hyperdrive service config` to set up Hyperdrive.")
	}

	// Get the max fee - prioritize the CLI arguments, default to the config file setting
	maxFeeGwei := hd.Context.MaxFee

	// Get the priority fee - prioritize the CLI arguments, default to the config file setting
	maxPriorityFeeGwei := hd.Context.MaxPriorityFee
	if maxPriorityFeeGwei == 0 {
		maxPriorityFee := eth.GweiToWei(cfg.Hyperdrive.MaxPriorityFee.Value)
		if maxPriorityFee == nil || maxPriorityFee.Uint64() == 0 {
			defaultFee := cfg.Hyperdrive.MaxPriorityFee.Default[cfg.Hyperdrive.Network.Value]
			fmt.Printf("%sNOTE: max priority fee not set or set to 0, defaulting to %.4f gwei%s\n", terminal.ColorYellow, defaultFee, terminal.ColorReset)
			maxPriorityFeeGwei = defaultFee
		} else {
			maxPriorityFeeGwei = eth.WeiToGwei(maxPriorityFee)
		}
	}

	// Use the requested max fee and priority fee if provided
	if maxFeeGwei != 0 {
		if maxPriorityFeeGwei > maxFeeGwei {
			fmt.Printf("NOTE: Priority fee cannot be greater than max fee. Lowering priority fee to %.2f gwei.\n", maxFeeGwei)
			maxPriorityFeeGwei = maxFeeGwei
		}
		fmt.Printf("%sUsing the requested max fee of %.2f gwei (including a max priority fee of %.2f gwei).\n", terminal.ColorYellow, maxFeeGwei, maxPriorityFeeGwei)
		lowLimit := maxFeeGwei / eth.WeiPerGwei * float64(simResult.EstimatedGasLimit)
		highLimit := maxFeeGwei / eth.WeiPerGwei * float64(simResult.SafeGasLimit)
		fmt.Printf("Total cost: %.4f to %.4f ETH%s\n", lowLimit, highLimit, terminal.ColorReset)
	} else {
		if c.Bool(utils.YesFlag.Name) {
			maxFeeWei, err := GetHeadlessMaxFeeWei()
			if err != nil {
				return nil, nil, err
			}
			maxFeeGwei = eth.WeiToGwei(maxFeeWei)
		} else {
			// Try to get the latest gas prices from Etherchain
			etherchainData, err := gas.GetEtherchainGasPrices()
			if err == nil {
				// Print the Etherchain data and ask for an amount
				maxFeeGwei = handleEtherchainGasPrices(etherchainData, simResult, maxPriorityFeeGwei, 0)

			} else {
				// Fallback to Etherscan
				fmt.Printf("%sWarning: couldn't get gas estimates from Etherchain - %s\nFalling back to Etherscan%s\n", terminal.ColorYellow, err.Error(), terminal.ColorReset)
				etherscanData, err := gas.GetEtherscanGasPrices()
				if err == nil {
					// Print the Etherscan data and ask for an amount
					maxFeeGwei = handleEtherscanGasPrices(etherscanData, simResult, maxPriorityFeeGwei, 0)
				} else {
					return nil, nil, fmt.Errorf("Error getting gas price suggestions: %w", err)
				}
			}
		}
		if maxPriorityFeeGwei > maxFeeGwei {
			fmt.Printf("NOTE: Priority fee cannot be greater than max fee. Lowering priority fee to %.2f gwei.\n", maxFeeGwei)
			maxPriorityFeeGwei = maxFeeGwei
		}
		fmt.Printf("%sUsing a max fee of %.2f gwei and a priority fee of %.2f gwei.\n%s", terminal.ColorBlue, maxFeeGwei, maxPriorityFeeGwei, terminal.ColorReset)
	}

	// Verify the node has enough ETH to use this max fee
	maxFee := eth.GweiToWei(maxFeeGwei)
	ethRequired := big.NewInt(0).Mul(maxFee, big.NewInt(int64(simResult.SafeGasLimit)))
	response, err := hd.Api.Wallet.Balance()
	if err != nil {
		fmt.Printf("%sWARNING: couldn't check the ETH balance of the node (%s)\nPlease ensure your node wallet has enough ETH to pay for this transaction.%s\n\n", terminal.ColorYellow, err.Error(), terminal.ColorReset)
	} else if response.Data.Balance.Cmp(ethRequired) < 0 {
		return nil, nil, fmt.Errorf("Your node has %.6f ETH in its wallet, which is not enough to pay for this transaction with a max fee of %.4f gwei; you require at least %.6f more ETH.", eth.WeiToEth(response.Data.Balance), maxFeeGwei, eth.WeiToEth(big.NewInt(0).Sub(ethRequired, response.Data.Balance)))
	}
	maxPriorityFee := eth.GweiToWei(maxPriorityFeeGwei)

	return maxFee, maxPriorityFee, nil
}

// Get the suggested max fee for service operations
func GetHeadlessMaxFeeWei() (*big.Int, error) {
	etherchainData, err := gas.GetEtherchainGasPrices()
	if err == nil {
		return etherchainData.RapidWei, nil
	}

	fmt.Printf("%sWARNING: couldn't get gas estimates from Etherchain - %s\nFalling back to Etherscan%s\n", terminal.ColorYellow, err.Error(), terminal.ColorReset)
	etherscanData, err := gas.GetEtherscanGasPrices()
	if err == nil {
		return eth.GweiToWei(etherscanData.FastGwei), nil
	}

	return nil, fmt.Errorf("error getting gas price suggestions: %w", err)
}

func handleEtherchainGasPrices(gasSuggestion gas.EtherchainGasFeeSuggestion, simResult eth.SimulationResult, priorityFee float64, gasLimit uint64) float64 {
	rapidGwei := math.Ceil(eth.WeiToGwei(gasSuggestion.RapidWei) + priorityFee)
	rapidEth := eth.WeiToEth(gasSuggestion.RapidWei)

	var rapidLowLimit float64
	var rapidHighLimit float64
	if gasLimit == 0 {
		rapidLowLimit = rapidEth * float64(simResult.EstimatedGasLimit)
		rapidHighLimit = rapidEth * float64(simResult.SafeGasLimit)
	} else {
		rapidLowLimit = rapidEth * float64(gasLimit)
		rapidHighLimit = rapidLowLimit
	}

	fastGwei := math.Ceil(eth.WeiToGwei(gasSuggestion.FastWei) + priorityFee)
	fastEth := eth.WeiToEth(gasSuggestion.FastWei)

	var fastLowLimit float64
	var fastHighLimit float64
	if gasLimit == 0 {
		fastLowLimit = fastEth * float64(simResult.EstimatedGasLimit)
		fastHighLimit = fastEth * float64(simResult.SafeGasLimit)
	} else {
		fastLowLimit = fastEth * float64(gasLimit)
		fastHighLimit = fastLowLimit
	}

	standardGwei := math.Ceil(eth.WeiToGwei(gasSuggestion.StandardWei) + priorityFee)
	standardEth := eth.WeiToEth(gasSuggestion.StandardWei)

	var standardLowLimit float64
	var standardHighLimit float64
	if gasLimit == 0 {
		standardLowLimit = standardEth * float64(simResult.EstimatedGasLimit)
		standardHighLimit = standardEth * float64(simResult.SafeGasLimit)
	} else {
		standardLowLimit = standardEth * float64(gasLimit)
		standardHighLimit = standardLowLimit
	}

	slowGwei := math.Ceil(eth.WeiToGwei(gasSuggestion.SlowWei) + priorityFee)
	slowEth := eth.WeiToEth(gasSuggestion.SlowWei)

	var slowLowLimit float64
	var slowHighLimit float64
	if gasLimit == 0 {
		slowLowLimit = slowEth * float64(simResult.EstimatedGasLimit)
		slowHighLimit = slowEth * float64(simResult.SafeGasLimit)
	} else {
		slowLowLimit = slowEth * float64(gasLimit)
		slowHighLimit = slowLowLimit
	}

	fmt.Printf("%s+============== Suggested Gas Prices ==============+\n", terminal.ColorBlue)
	fmt.Println("| Avg Wait Time |  Max Fee  |    Total Gas Cost    |")
	fmt.Printf("| %-13s | %-9s | %.4f to %.4f ETH |\n",
		gasSuggestion.RapidTime, fmt.Sprintf("%d gwei", int(rapidGwei)), rapidLowLimit, rapidHighLimit)
	fmt.Printf("| %-13s | %-9s | %.4f to %.4f ETH |\n",
		gasSuggestion.FastTime, fmt.Sprintf("%d gwei", int(fastGwei)), fastLowLimit, fastHighLimit)
	fmt.Printf("| %-13s | %-9s | %.4f to %.4f ETH |\n",
		gasSuggestion.StandardTime, fmt.Sprintf("%d gwei", int(standardGwei)), standardLowLimit, standardHighLimit)
	fmt.Printf("| %-13s | %-9s | %.4f to %.4f ETH |\n",
		gasSuggestion.SlowTime, fmt.Sprintf("%d gwei", int(slowGwei)), slowLowLimit, slowHighLimit)
	fmt.Printf("+==================================================+\n\n%s", terminal.ColorReset)

	fmt.Printf("These prices include a maximum priority fee of %.2f gwei.\n", priorityFee)

	for {
		desiredPrice := utils.Prompt(
			fmt.Sprintf("Please enter your max fee (including the priority fee) or leave blank for the default of %d gwei:", int(fastGwei)),
			"^(?:[1-9]\\d*|0)?(?:\\.\\d+)?$",
			"Not a valid gas price, try again:")

		if desiredPrice == "" {
			return fastGwei
		}

		desiredPriceFloat, err := strconv.ParseFloat(desiredPrice, 64)
		if err != nil {
			fmt.Printf("Not a valid gas price (%s), try again.", err.Error())
			fmt.Println()
			continue
		}
		if desiredPriceFloat <= 0 {
			fmt.Println("Max fee must be greater than zero.")
			continue
		}

		return desiredPriceFloat
	}
}

func handleEtherscanGasPrices(gasSuggestion gas.EtherscanGasFeeSuggestion, simResult eth.SimulationResult, priorityFee float64, gasLimit uint64) float64 {
	fastGwei := math.Ceil(gasSuggestion.FastGwei + priorityFee)
	fastEth := gasSuggestion.FastGwei / eth.WeiPerGwei

	var fastLowLimit float64
	var fastHighLimit float64
	if gasLimit == 0 {
		fastLowLimit = fastEth * float64(simResult.EstimatedGasLimit)
		fastHighLimit = fastEth * float64(simResult.SafeGasLimit)
	} else {
		fastLowLimit = fastEth * float64(gasLimit)
		fastHighLimit = fastLowLimit
	}

	standardGwei := math.Ceil(gasSuggestion.StandardGwei + priorityFee)
	standardEth := gasSuggestion.StandardGwei / eth.WeiPerGwei

	var standardLowLimit float64
	var standardHighLimit float64
	if gasLimit == 0 {
		standardLowLimit = standardEth * float64(simResult.EstimatedGasLimit)
		standardHighLimit = standardEth * float64(simResult.SafeGasLimit)
	} else {
		standardLowLimit = standardEth * float64(gasLimit)
		standardHighLimit = standardLowLimit
	}

	slowGwei := math.Ceil(gasSuggestion.SlowGwei + priorityFee)
	slowEth := gasSuggestion.SlowGwei / eth.WeiPerGwei

	var slowLowLimit float64
	var slowHighLimit float64
	if gasLimit == 0 {
		slowLowLimit = slowEth * float64(simResult.EstimatedGasLimit)
		slowHighLimit = slowEth * float64(simResult.SafeGasLimit)
	} else {
		slowLowLimit = slowEth * float64(gasLimit)
		slowHighLimit = slowLowLimit
	}

	fmt.Printf("%s+============ Suggested Gas Prices ============+\n", terminal.ColorBlue)
	fmt.Println("|   Speed   |  Max Fee  |    Total Gas Cost    |")
	fmt.Printf("| Fast      | %-9s | %.4f to %.4f ETH |\n",
		fmt.Sprintf("%d gwei", int(fastGwei)), fastLowLimit, fastHighLimit)
	fmt.Printf("| Standard  | %-9s | %.4f to %.4f ETH |\n",
		fmt.Sprintf("%d gwei", int(standardGwei)), standardLowLimit, standardHighLimit)
	fmt.Printf("| Slow      | %-9s | %.4f to %.4f ETH |\n",
		fmt.Sprintf("%d gwei", int(slowGwei)), slowLowLimit, slowHighLimit)
	fmt.Printf("+==============================================+\n\n%s", terminal.ColorReset)

	fmt.Printf("These prices include a maximum priority fee of %.2f gwei.\n", priorityFee)

	for {
		desiredPrice := utils.Prompt(
			fmt.Sprintf("Please enter your max fee (including the priority fee) or leave blank for the default of %d gwei:", int(fastGwei)),
			"^(?:[1-9]\\d*|0)?(?:\\.\\d+)?$",
			"Not a valid gas price, try again:")

		if desiredPrice == "" {
			return fastGwei
		}

		desiredPriceFloat, err := strconv.ParseFloat(desiredPrice, 64)
		if err != nil {
			fmt.Printf("Not a valid gas price (%s), try again.", err.Error())
			fmt.Println()
			continue
		}
		if desiredPriceFloat <= 0 {
			fmt.Println("Max fee must be greater than zero.")
			continue
		}

		return desiredPriceFloat
	}
}
