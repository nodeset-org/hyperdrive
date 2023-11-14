package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/nodeset-org/hyperdrive/cmd/constellation"
)

var rootCmd = &cobra.Command{
	Use:   "hyperdrive",
	Short: "Hyperdrive initialization and Rocketpool service status check",
  	Run: func(cmd *cobra.Command, args []string) {
  		fmt.Println("hyperdrive")
  	},
}

// This function is used to manage/add sub-commands.
// Parent init functions should be adding sub-commands.
func init() {
  	rootCmd.AddCommand(constellation.ConstellationCmd)
}

// This function is executed prior to any Cobra command
func Execute() {
  	if err := rootCmd.Execute(); err != nil {
    	fmt.Fprintln(os.Stderr, err)
    	os.Exit(1)
  	}
}
