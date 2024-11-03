package utils

import (
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/context"
	"github.com/urfave/cli/v2"
)

// Set up the context with the config path since bash completion bypasses the app's Before function
func BootstrapCliForBashCompletion(c *cli.Context) error {
	configPath := c.String(UserDirPathFlag.Name)
	path, err := homedir.Expand(strings.TrimSpace(configPath))
	if err != nil {
		return err
	}
	hdCtx := context.NewHyperdriveContext(path, nil)
	context.SetHyperdriveContext(c, hdCtx)
	return nil
}
