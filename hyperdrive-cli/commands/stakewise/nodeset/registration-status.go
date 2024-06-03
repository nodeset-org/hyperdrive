package nodeset

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	cliutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

func registrationStatus(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	return CheckRegistrationStatus(c, hd, sw)
}

func CheckRegistrationStatus(c *cli.Context, hd *client.HyperdriveClient, sw *client.StakewiseClient) error {
	// Get wallet response
	response, err := hd.Api.Wallet.Status()
	if err != nil {
		return err
	}

	// Make sure we have a wallet loaded
	if !response.Data.WalletStatus.Wallet.IsLoaded {
		fmt.Println("The node wallet has not been initialized yet. Please run `hyperdrive wallet status` to learn more.")
		return nil
	}

	// Get the registration status
	resp, err := sw.Api.Nodeset.RegistrationStatus()
	if err != nil {
		return err
	}
	if resp.Data.Registered {
		fmt.Println("Your node is registered.")
		return nil
	}

	fmt.Println("Your node is not currently registered.")
	if cliutils.Confirm("Would you like to register now so you can upload your validator keys to NodeSet?") {
		return registerNode(c)
	}

	return nil
}
