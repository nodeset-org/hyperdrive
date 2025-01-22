package utils

import (
	"github.com/urfave/cli/v2"
)

const (
	NoRestartFlag string = "no-restart"
	MnemonicFlag  string = "mnemonic"
)

var (
	YesFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "yes",
		Aliases: []string{"y"},
		Usage:   "Automatically confirm all interactive questions",
	}
	ComposeFileFlag *cli.StringSliceFlag = &cli.StringSliceFlag{
		Name:    "compose-file",
		Aliases: []string{"f"},
		Usage:   "Supplemental Docker compose files for custom containers to include when performing service commands such as 'start' and 'stop'; this flag may be defined multiple times",
	}
	UserDirPathFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "config-path",
		Aliases: []string{"c"},
		Usage:   "Directory to install and save all of Hyperdrive's configuration and data to",
	}
)
