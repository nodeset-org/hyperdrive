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
	cfg, _, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading Hyperdrive config: %w", err)
	}
	if !cfg.StakeWise.Enabled.Value {
		fmt.Println("The StakeWise module is not enabled in your Hyperdrive configuration.")
		return nil
	}

	// Upload to the server
	_, err = swcmdutils.UploadDepositData(c, hd, sw)
	return err
}
