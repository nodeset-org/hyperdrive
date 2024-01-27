package wallet

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/tx"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	"github.com/urfave/cli/v2"
)

func sendMessage(c *cli.Context, toAddressOrEns string, message []byte) error {
	// Get RP client
	hd, err := client.NewClientFromCtx(c).WithReady()
	if err != nil {
		return err
	}

	// Get the address
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
		toAddress, err = input.ValidateAddress("to address", toAddressOrEns)
		if err != nil {
			return err
		}
		toAddressString = toAddress.Hex()
	}

	// Build the TX
	response, err := hd.Api.Wallet.SendMessage(message, toAddress)
	if err != nil {
		return err
	}

	// Run the TX
	err = tx.HandleTx(c, hd, response.Data.TxInfo,
		fmt.Sprintf("Are you sure you want to send a message to %s?", toAddressString),
		"sending message",
		fmt.Sprintf("Sending message to %s...", toAddressString),
	)
	if err != nil {
		return err
	}

	// Log & return
	fmt.Printf("Successfully sent message to %s.\n", toAddressString)
	return nil
}
