package nodeset

import (
	"fmt"
	"net/mail"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	cliutils "github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
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
	sw, err := client.NewStakewiseClientFromCtx(c, hd)
	if err != nil {
		return err
	}

	// Get the email
	email := RegisterEmailFlag.Name
	if email == "" {
		for {
			email = cliutils.Prompt("Enter the email address you'd like to register with NodeSet:", "^.*$", "Invalid email address, try again")
			_, err := mail.ParseAddress(email)
			if err == nil {
				break
			}
			fmt.Println("Invalid email address, try again")
		}
	}

	_, err = sw.Api.Nodeset.RegisterNode(email)

	return err
}
