package node

import (
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-cli/client"
	"github.com/spf13/cobra"
)

func RegisterCommands(parentCmd *cobra.Command, installPath string) {
	// Master group for service commands
	nodeCmd := &cobra.Command{
		Use:     "node",
		Aliases: []string{"n"},
		Short:   "Run commands related to Stakewise node operation",
	}

	// example call-daemon
	uploadDepositDataCmd := &cobra.Command{
		Use:     "upload-deposit-data",
		Aliases: []string{"u"},
		Short:   "Upload the Stakewise deposit data to the NodeSet vault service",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.NewClient(installPath)
			if err != nil {
				return err
			}
			return uploadDepositData(client, args[0])
		},
	}
	nodeCmd.AddCommand(uploadDepositDataCmd)

	// Add service to the root command
	parentCmd.AddCommand(nodeCmd)
}
