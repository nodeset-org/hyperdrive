package wallet

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/urfave/cli/v2"

	"github.com/rocket-pool/node-manager-core/utils/input"
	"github.com/rocket-pool/node-manager-core/utils/math"
)

func nodeSend(c *cli.Context, amount float64, token string, toAddressOrEns string) error {
	// Get RP client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Get amount in wei
	amountWei := eth.EthToWei(amount)

	// Get the recipient
	var toAddress common.Address
	var toAddressString string
	if strings.Contains(toAddressOrEns, ".") {
		response, err := hd.Api.Utils.ResolveEns(common.Address{}, toAddressOrEns)
		if err != nil {
			return err
		}
		toAddress = response.Data.Address
		toAddressString = fmt.Sprintf("%s (%s)", toAddressOrEns, toAddress.Hex())
	} else {
		var err error
		toAddress, err = input.ValidateAddress("to address", toAddressOrEns)
		if err != nil {
			return err
		}
		toAddressString = toAddress.Hex()
	}

	// Build the TX
	response, err := hd.Api.Wallet.Send(amountWei, token, toAddress)
	if err != nil {
		return err
	}
	tokenString := fmt.Sprintf("%s (%s)", response.Data.TokenSymbol, token)

	// Verify
	if !response.Data.CanSend {
		fmt.Println("Cannot send tokens:")
		if response.Data.InsufficientBalance {
			if strings.HasPrefix(token, "0x") {
				fmt.Printf("The node's balance of %s is insufficient.\n", tokenString)
			} else {
				fmt.Printf("The node's %s balance is insufficient.\n", token)
			}
		}
		return nil
	}

	// Print details and create confirm message
	var confirmMsg string
	var submitMsg string
	var successMsg string
	if strings.HasPrefix(token, "0x") {
		fmt.Printf("Token address:   %s\n", token)
		fmt.Printf("Token name:      %s\n", response.Data.TokenName)
		fmt.Printf("Token symbol:    %s\n", response.Data.TokenSymbol)
		fmt.Printf("Node balance:    %.6f %s\n\n", eth.WeiToEth(response.Data.Balance), response.Data.TokenSymbol)
		fmt.Printf("%sWARNING: Please confirm that the above token is the one you intend to send before confirming below!%s\n\n", terminal.ColorYellow, terminal.ColorReset)

		confirmMsg = fmt.Sprintf("Are you sure you want to send %.6f of %s to %s? This action cannot be undone!", math.RoundDown(eth.WeiToEth(amountWei), 6), tokenString, toAddressString)
		submitMsg = fmt.Sprintf("Sending %s to %s...\n", tokenString, toAddressString)
		successMsg = fmt.Sprintf("Successfully sent %.6f of %s to %s.", math.RoundDown(eth.WeiToEth(amountWei), 6), tokenString, toAddressString)
	} else {
		fmt.Printf("Node balance:    %.6f %s\n\n", eth.WeiToEth(response.Data.Balance), token)
		confirmMsg = fmt.Sprintf("Are you sure you want to send %.6f %s to %s? This action cannot be undone!", math.RoundDown(eth.WeiToEth(amountWei), 6), token, toAddressString)
		submitMsg = fmt.Sprintf("Sending %s to %s...\n", token, toAddressString)
		successMsg = fmt.Sprintf("Successfully sent %.6f %s to %s.", math.RoundDown(eth.WeiToEth(amountWei), 6), token, toAddressString)
	}

	// Run the TX
	validated, err := tx.HandleTx(c, hd, response.Data.TxInfo,
		confirmMsg,
		"sending tokens",
		submitMsg,
	)
	if err != nil {
		return err
	}
	if !validated {
		return nil
	}

	// Log & return
	fmt.Println(successMsg)
	return nil
}
