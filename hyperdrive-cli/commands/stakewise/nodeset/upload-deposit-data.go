package nodeset

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	swcmdutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/utils"
	"github.com/urfave/cli/v2"
)

func uploadDepositData(c *cli.Context) error {
	// Get the client
	sw := client.NewStakewiseClientFromCtx(c)

	// Upload to the server
	err := swcmdutils.UploadDepositData(sw)
	if err != nil {
		return err
	}

	fmt.Println("done!")
	fmt.Println("Your Stakewise Operator container should reflect the new deposit data shortly.")
	return nil
}
