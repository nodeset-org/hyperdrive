package nodeset

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/urfave/cli/v2"
)

var (
	RegisterEmailFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "email",
		Aliases: []string{"e"},
		Usage:   "Email address to register with NodeSet.",
	}
)

func registerNode(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}

	// Check if doesn't have a wallet it's already registered
	hasWallet, shouldRegister, err := checkRegistrationStatusImpl(hd)
	if err != nil {
		return err
	}
	if !hasWallet || !shouldRegister {
		return nil
	}

	return registerNodeImpl(c, hd)
}
