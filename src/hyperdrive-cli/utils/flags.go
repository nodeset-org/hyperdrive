package utils

import (
	"github.com/urfave/cli/v2"
)

const (
	SignTxOnlyFlag string = "sign-tx-only"
	NoRestartFlag  string = "no-restart"
	MnemonicFlag   string = "mnemonic"
)

var (
	YesFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "yes",
		Aliases: []string{"y"},
		Usage:   "Automatically confirm all interactive questions",
	}
	PrintTxDataFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:    "print-tx-data",
		Aliases: []string{"pd"},
		Usage:   "Print the TX data for transactions without signing or submitting them. Useful for masquerade mode or offline wallet operations.",
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
