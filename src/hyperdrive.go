package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hyperdrive",
		Short: "Hyperdrive initialization and Rocketpool service status check",
		Run:   runHyperdrive,
	}

	var helloCmd = &cobra.Command{
		Use:   "hello",
		Short: "Prints 'Hello, World!'",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the Rocketpool service",
		Run:   runRocketpoolStart,
	}

	rootCmd.AddCommand(helloCmd, startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func runHyperdrive(cmd *cobra.Command, args []string) {
	fmt.Println("Hyperdrive initializing...")
	rpServiceStatus := exec.Command("rocketpool", "service", "status")
	out, err := rpServiceStatus.Output()
	if err != nil {
		fmt.Println("Error checking rp service status:", err)
		return
	}
	fmt.Println(string(out))
}

func runRocketpoolStart(cmd *cobra.Command, args []string) {
	fmt.Println("Starting Rocketpool service...")
	rpServiceStart := exec.Command("rocketpool", "service", "start")
	out, err := rpServiceStart.CombinedOutput()
	if err != nil {
		fmt.Println("Error starting rp service:", err)
		return
	}
	fmt.Println(string(out))
}
