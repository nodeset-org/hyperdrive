package utils

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/utils/input"
	"github.com/urfave/cli/v2"
)

// Print a TX's details to the console.
func PrintTransactionHash(hd *client.HyperdriveClient, hash common.Hash) {
	finalMessage := "Waiting for the transaction to be included in a block... you may wait here for it, or press CTRL+C to exit and return to the terminal.\n\n"
	printTransactionHashImpl(hd, hash, finalMessage)
}

// Print a TX's details to the console, but inform the user NOT to cancel it.
func PrintTransactionHashNoCancel(hd *client.HyperdriveClient, hash common.Hash) {
	finalMessage := "Waiting for the transaction to be included in a block... **DO NOT EXIT!** This transaction is one of several that must be completed.\n\n"
	printTransactionHashImpl(hd, hash, finalMessage)
}

// Print a batch of transaction hashes to the console.
func PrintTransactionBatchHashes(hd *client.HyperdriveClient, hashes []common.Hash) {
	finalMessage := "Waiting for the transactions to be included in one or more blocks... you may wait here for them, or press CTRL+C to exit and return to the terminal.\n\n"

	// Print the hashes
	fmt.Println("Transactions have been submitted with the following hashes:")
	hashStrings := make([]string, len(hashes))
	for i, hash := range hashes {
		hashString := hash.String()
		hashStrings[i] = hashString
		fmt.Println(hashString)
	}
	fmt.Println()

	txWatchUrl := getTxWatchUrl(hd)
	if txWatchUrl != "" {
		fmt.Println("You may follow their progress by visiting the following URLs in sequence:")
		for _, hash := range hashStrings {
			fmt.Printf("%s/%s\n", txWatchUrl, hash)
		}
	}
	fmt.Println()

	fmt.Print(finalMessage)
}

// Print a warning to the console if the user set a custom nonce, but this operation involves multiple transactions
func PrintMultiTransactionNonceWarning() {
	fmt.Printf("%sNOTE: You have specified the `nonce` flag to indicate a custom nonce for this transaction.\n"+
		"However, this operation requires multiple transactions.\n"+
		"Hyperdrive will use your custom value as a basis, and increment it for each additional transaction.\n"+
		"If you have multiple pending transactions, this MAY OVERRIDE more than the one that you specified.%s\n\n", terminal.ColorYellow, terminal.ColorReset)
}

// Implementation of PrintTransactionHash and PrintTransactionHashNoCancel
func printTransactionHashImpl(hd *client.HyperdriveClient, hash common.Hash, finalMessage string) {
	txWatchUrl := getTxWatchUrl(hd)
	hashString := hash.String()
	fmt.Printf("Transaction has been submitted with hash %s.\n", hashString)
	if txWatchUrl != "" {
		fmt.Printf("You may follow its progress by visiting:\n")
		fmt.Printf("%s/%s\n\n", txWatchUrl, hashString)
	}
	fmt.Print(finalMessage)
}

// Get the URL for watching the transaction in a block explorer
func getTxWatchUrl(hd *client.HyperdriveClient) string {
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: couldn't read config file so the transaction URL will be unavailable (%s).\n", err)
		return ""
	}

	if isNew {
		fmt.Print("Settings file not found. Please run `hyperdrive service config` to set up Hyperdrive.")
		return ""
	}
	resources := cfg.Hyperdrive.GetNetworkResources()
	return resources.TxWatchUrl
}

// Convert a Unix datetime to a string, or `---` if it's zero
func GetDateTimeString(dateTime uint64) string {
	return GetDateTimeStringOfTime(time.Unix(int64(dateTime), 0))
}

// Convert a Unix datetime to a string, or `---` if it's zero
func GetDateTimeStringOfTime(dateTime time.Time) string {
	timeString := dateTime.Format(time.RFC822)
	if dateTime == time.Unix(0, 0) {
		timeString = "---"
	}
	return timeString
}

// Gets the hex string of an address, or "none" if it was the 0x0 address
func GetPrettyAddress(address common.Address) string {
	addressString := address.Hex()
	if addressString == "0x0000000000000000000000000000000000000000" {
		return "<none>"
	}
	return addressString
}

// Prints an error message when the Beacon client is not using the deposit contract address that Hyperdrive expects
func PrintDepositMismatchError(rpNetwork, beaconNetwork uint64, rpDepositAddress, beaconDepositAddress common.Address) {
	fmt.Printf("%s***ALERT***\n", terminal.ColorRed)
	fmt.Println("YOUR ETH2 CLIENT IS NOT CONNECTED TO THE SAME NETWORK THAT HYPERDRIVE IS USING!")
	fmt.Println("This is likely because your ETH2 client is using the wrong configuration.")
	fmt.Println("For the safety of your funds, Hyperdrive will not let you deposit your ETH until this is resolved.")
	fmt.Println()
	fmt.Println("To fix it if you are in Docker mode:")
	fmt.Println("\t1. Run 'hyperdrive service install -d' to get the latest configuration")
	fmt.Println("\t2. Run 'hyperdrive service stop' and 'hyperdrive service start' to apply the configuration.")
	fmt.Println("If you are using Hybrid or Native mode, please correct the network flags in your ETH2 launch script.")
	fmt.Println()
	fmt.Println("Details:")
	fmt.Printf("\tHyperdrive expects deposit contract %s on chain %d.\n", rpDepositAddress.Hex(), rpNetwork)
	fmt.Printf("\tYour Beacon client is using deposit contract %s on chain %d.%s\n", beaconDepositAddress.Hex(), beaconNetwork, terminal.ColorReset)
}

// Prints what network you're currently on
func PrintNetwork(currentNetwork config.Network, isNew bool) error {
	if isNew {
		return fmt.Errorf("Settings file not found. Please run `hyperdrive service config` to set up Hyperdrive.")
	}

	switch currentNetwork {
	case config.Network_Mainnet:
		fmt.Printf("Hyperdrive is currently using the %sEthereum Mainnet.%s\n\n", terminal.ColorGreen, terminal.ColorReset)
	case hdconfig.Network_HoleskyDev:
		fmt.Printf("Hyperdrive is currently using the %sHolesky Development Network.%s\n\n", terminal.ColorYellow, terminal.ColorReset)
	case config.Network_Holesky:
		fmt.Printf("Hyperdrive is currently using the %sHolesky Test Network.%s\n\n", terminal.ColorYellow, terminal.ColorReset)
	default:
		fmt.Printf("%sYou are on an unexpected network [%v].%s\n\n", terminal.ColorYellow, currentNetwork, terminal.ColorReset)
	}

	return nil
}

// Parses a string representing either a floating point value or a raw wei amount into a *big.Int
func ParseFloat(c *cli.Context, name string, value string, isFraction bool) (*big.Int, error) {
	var floatValue float64
	if c.Bool(RawFlag.Name) {
		val, err := input.ValidateBigInt(name, value)
		if err != nil {
			return nil, err
		}
		return val, nil
	} else if isFraction {
		val, err := input.ValidateFraction(name, value)
		if err != nil {
			return nil, err
		}
		floatValue = val
	} else {
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		floatValue = val
	}

	trueVal := eth.EthToWei(floatValue)
	fmt.Printf("Your value will be multiplied by 10^18 to be used in the contracts, which results in:\n\n\t[%s]\n\n", trueVal.String())
	if !(c.Bool("yes") || Confirm("Please make sure this is what you want and does not have any floating-point errors.\n\nIs this result correct?")) {
		fmt.Printf("Cancelled. Please try again with the '--%s' flag and provide an explicit value instead.\n", RawFlag.Name)
		return nil, nil
	}
	return trueVal, nil
}

// Clear terminal output
func ClearTerminal() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
