package example

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/spf13/cobra"
)

func RegisterCommands(rootCmd *cobra.Command, installPath string) {
	// Master group for service commands
	exampleCmd := &cobra.Command{
		Use:     "example",
		Aliases: []string{"e"},
		Short:   "Run the example commands for development exploration",
	}

	// example call-daemon
	callDaemonCmd := &cobra.Command{
		Use:     "call-daemon",
		Aliases: []string{"c"},
		Short:   "Call the Rocket Pool daemon to invoke a subset of network commands",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.NewClient(installPath)
			if err != nil {
				return err
			}
			return callDaemon(client, args[0])
		},
	}
	exampleCmd.AddCommand(callDaemonCmd)

	// Add service to the root command
	rootCmd.AddCommand(exampleCmd)
}
