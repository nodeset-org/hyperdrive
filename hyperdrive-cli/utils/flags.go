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
		Name:    "config-dir",
		Aliases: []string{"c"},
		Usage:   "Directory to install and save all of Hyperdrive's configuration and data to",
	}
	SystemDirPathFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "system-dir",
		Aliases: []string{"s"},
		Usage:   "Directory where Hyperdrive's system files are installed on the system. If blank, this will use the same default as Hyperdrive's installation package for your Operating System.",
	}
	AllowRootFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "allow-root",
		Aliases: []string{"r"},
		Usage:   "Allow hyperdrive to be run as the root user",
	}
)
