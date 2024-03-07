package wallet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	nmc_utils "github.com/rocket-pool/node-manager-core/utils"
	"github.com/rocket-pool/node-manager-core/wallet"
	"github.com/urfave/cli/v2"
)

const (
	signatureVersion int = 1
)

type PersonalSignature struct {
	Address   common.Address `json:"address"`
	Message   string         `json:"msg"`
	Signature string         `json:"sig"`
	Version   string         `json:"version"` // beaconcha.in expects a string
}

var (
	signMessageFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "message",
		Aliases: []string{"m"},
		Usage:   "The 'quoted message' to be signed",
	}
)

func signMessage(c *cli.Context) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Get & check wallet status
	status, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}
	if !wallet.IsWalletReady(status.Data.WalletStatus) {
		fmt.Println("The node wallet is not loaded or your node is in read-only mode. Please run `hyperdrive wallet status` for more details.")
		return nil
	}

	// Get the message
	message := c.String(signMessageFlag.Name)
	for message == "" {
		message = utils.Prompt("Please enter the message you want to sign: (EIP-191 personal_sign)", "^.+$", "Please enter the message you want to sign: (EIP-191 personal_sign)")
	}

	// Build the TX
	response, err := hd.Api.Wallet.SignMessage([]byte(message))
	if err != nil {
		return err
	}

	// Print the signature
	formattedSignature := PersonalSignature{
		Address:   status.Data.WalletStatus.Wallet.WalletAddress,
		Message:   message,
		Signature: nmc_utils.EncodeHexWithPrefix(response.Data.SignedMessage),
		Version:   fmt.Sprint(signatureVersion),
	}
	bytes, err := json.MarshalIndent(formattedSignature, "", "    ")
	if err != nil {
		return err
	}

	fmt.Printf("Signed Message:\n\n%s\n", string(bytes))
	return nil
}
