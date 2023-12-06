package main

import (
	"fmt"
	"os"

	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Hyperdrive daemon init",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hyperdrive daemon")
	},
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "todo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("daemon init\n")
		apiManager := api.NewApiManager()
		err := apiManager.Start()
		if err != nil {
			fmt.Printf("error starting API server: %w", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(InitCmd)
}

// This function is executed prior to any Cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
