package nodeset

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	swcmdutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/utils"
	"github.com/urfave/cli/v2"
)

func uploadDepositData(c *cli.Context, forceUploadFlag bool) error {
	// Get the client
	sw := client.NewStakewiseClientFromCtx(c)

	// Upload to the server
	_, err := swcmdutils.UploadDepositData(sw, forceUploadFlag)
	return err
}
