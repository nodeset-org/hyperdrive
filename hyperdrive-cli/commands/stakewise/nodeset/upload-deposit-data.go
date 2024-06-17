package nodeset

import (
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
	_, err = swcmdutils.UploadDepositData(c, sw)
	return err
}
