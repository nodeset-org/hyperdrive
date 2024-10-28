package utils

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
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
	PrintTxDataFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "print-tx-data",
		Aliases: []string{"pt"},
		Usage:   "Print any TX data for transactions without signing or submitting them. Useful for masquerade mode or offline wallet operations.",
	}
	SignTxOnlyFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "sign-tx-only",
		Aliases: []string{"st"},
		Usage:   "Sign any TXs and print the results, but don't submit it to the network. Useful if you want to save a TX for later or bundle it up with a service like Flashbots.",
	}
	IgnoreTxSimFailureFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "ignore-tx-sim-failure",
		Aliases: []string{"itsf"},
		Usage:   fmt.Sprintf("Ignore any transaction simulation failures and sign / submit transactions even if they will revert. %sThis can result in failed transactions and loss of funds. Only use this if you absolutely know what you're doing.%s", terminal.ColorRed, terminal.ColorReset),
	}
	ForceGasLimitFlag *cli.Uint64Flag = &cli.Uint64Flag{
		Name:    "force-gas-limit",
		Aliases: []string{"fgl"},
		Usage:   fmt.Sprintf("Force a specific gas limit for all transactions produced by the command being run, overriding the gas limit estimated by transaction simulation. %sOnly use this if you know what you're doing.%s", terminal.ColorRed, terminal.ColorReset),
	}
	RawFlag *cli.BoolFlag = &cli.BoolFlag{
		Name: "raw",
	}
	ComposeFileFlag *cli.StringSliceFlag = &cli.StringSliceFlag{
		Name:    "compose-file",
		Aliases: []string{"f"},
		Usage:   "Supplemental Docker compose files for custom containers to include when performing service commands such as 'start' and 'stop'; this flag may be defined multiple times",
	}
)

func InstantiateFlag[FlagType cli.Flag](prototype FlagType, description string) cli.Flag {
	switch typedProto := any(prototype).(type) {
	case *cli.BoolFlag:
		return &cli.BoolFlag{
			Name:    typedProto.Name,
			Aliases: typedProto.Aliases,
			Usage:   description,
		}
	case *cli.Uint64Flag:
		return &cli.Uint64Flag{
			Name:    typedProto.Name,
			Aliases: typedProto.Aliases,
			Usage:   description,
		}
	case *cli.StringFlag:
		return &cli.StringFlag{
			Name:    typedProto.Name,
			Aliases: typedProto.Aliases,
			Usage:   description,
		}
	case *cli.Float64Flag:
		return &cli.Float64Flag{
			Name:    typedProto.Name,
			Aliases: typedProto.Aliases,
			Usage:   description,
		}
	default:
		panic("unsupported flag type")
	}
}
