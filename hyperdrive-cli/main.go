package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-cli/commands/rocketpool"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-cli/commands/service"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-cli/commands/stakewise"
	"github.com/spf13/cobra"
)

func main() {
	// Master command for the binary
	rootCmd := &cobra.Command{
		Short: "Hyperdrive initialization and Rocketpool service status check",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hyperdrive")
		},
		SilenceUsage: true,
	}

	// Set up flags
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleError(fmt.Errorf("error getting user's home directory for default installation path: %w", err))
		os.Exit(1)
	}
	installPath := rootCmd.Flags().StringP("install-path", "p", filepath.Join(homeDir, ".hyperdrive"), "Location of the Hyperdrive install folder")

	// Register the subcommands
	rocketpool.RegisterCommands(rootCmd, *installPath)
	service.RegisterCommands(rootCmd, *installPath)
	stakewise.RegisterCommands(rootCmd, *installPath)

	// Run - this automatically prints errors with the help text
	rootCmd.Execute()
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
