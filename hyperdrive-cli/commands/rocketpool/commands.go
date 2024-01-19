package rocketpool

import (
	"github.com/spf13/cobra"
)

func RegisterCommands(parentCmd *cobra.Command, installPath string) {
	// Rocket Pool group command
	rocketpoolCmd := &cobra.Command{
		Use:     "rocketpool",
		Aliases: []string{"rp"},
		Short:   "Rocket Pool commands",
	}

	// Register subcommands

	// Add this to the parent
	parentCmd.AddCommand(rocketpoolCmd)
}
