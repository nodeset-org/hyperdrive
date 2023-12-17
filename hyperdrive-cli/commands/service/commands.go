package service

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/spf13/cobra"
)

func RegisterCommands(rootCmd *cobra.Command, installPath string) {
	// Master group for service commands
	serviceCmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"s"},
		Short:   "Manage the Hyperdrive configuration and Docker containers",
	}
	composeFiles := serviceCmd.Flags().StringSliceP("compose-file", "f", []string{}, "Supplemental docker-compose file to add to the service commands. You can use this once per file, or once with all files using commas to separate them.")

	// service config
	configCmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c"},
		Short:   "Creates the Hyperdrive configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return configHyperdrive(installPath)
		},
	}
	serviceCmd.AddCommand(configCmd)

	// service start
	startCmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"s"},
		Short:   "Starts the Hyperdrive Docker containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.NewClient(installPath)
			if err != nil {
				return fmt.Errorf("error running start: %w", err)
			}
			return client.StartService(*composeFiles)
		},
	}
	serviceCmd.AddCommand(startCmd)

	// service stop
	stopCmd := &cobra.Command{
		Use:     "stop",
		Aliases: []string{"p"},
		Short:   "Stops the Hyperdrive Docker containers without destroying them",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.NewClient(installPath)
			if err != nil {
				return fmt.Errorf("error running stop: %w", err)
			}
			return client.StopService(*composeFiles)
		},
	}
	serviceCmd.AddCommand(stopCmd)

	// Add service to the root command
	rootCmd.AddCommand(serviceCmd)
}
