package stakewise

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/stakewise/node"
	"github.com/spf13/cobra"
)

func RegisterCommands(parentCmd *cobra.Command, installPath string) {
	// Stakewise group command
	stakewiseCmd := &cobra.Command{
		Use:     "stakewise",
		Aliases: []string{"sw"},
		Short:   "Stakewise commands",
	}

	// Register subcommands
	node.RegisterCommands(stakewiseCmd, installPath)

	// Add this to the parent
	parentCmd.AddCommand(stakewiseCmd)
}
