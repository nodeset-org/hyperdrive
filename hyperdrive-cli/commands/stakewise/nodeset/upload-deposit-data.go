package nodeset

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	swcmdutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/utils"
	"github.com/urfave/cli/v2"
)

func uploadDepositData(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Upload to the server
	_, err = swcmdutils.UploadDepositData(sw)
	return err
}

func blah(c *cli.Context) error {
	// Get the client
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}
	fmt.Printf("Calling Blah\n")
	resp, err := sw.Api.Nodeset.Blah("nodeAddress", "email", "signature")
	if err != nil {
		return err
	}
	fmt.Printf("Response: %v\n", resp)
	return nil
}
