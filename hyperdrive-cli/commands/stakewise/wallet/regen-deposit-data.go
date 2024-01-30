package wallet

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	swcmdutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

var (
	regenDepositDataNoRestartFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:  "no-restart",
		Usage: fmt.Sprintf("Don't automatically restart the Stakewise Operator containers after regenerating the deposit data. %sOnly use this if you know what you're doing and can restart it manually.%s", terminal.ColorYellow, terminal.ColorReset),
	}
)

func regenerateDepositData(c *cli.Context) error {
	hd := client.NewClientFromCtx(c)
	sw := client.NewStakewiseClientFromCtx(c)
	noRestart := c.Bool(regenDepositDataNoRestartFlag.Name)

	// Regenerate the deposit data
	err := swcmdutils.RegenDepositData(hd, sw, noRestart)
	if err != nil {
		fmt.Println("error")
		return err
	}

	// Upload to the server
	err = swcmdutils.UploadDepositData(sw)
	if err != nil {
		return err
	}

	if !noRestart {
		fmt.Println()
		fmt.Println("Your node's keys are now ready for use. When NodeSet selects one of them for a new deposit, your system will deposit it and begin attesting automatically.")
	} else {
		fmt.Println("Your node's keys are uploaded, but you *must* restart your Validator Client and Stakewise Operator service at your earliest convenience to begin attesting.")
	}

	return nil
}
