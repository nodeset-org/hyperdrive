package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/ethclient"
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

	var ethClientCmd = &cobra.Command{
		Use:   "ethclient",
		Short: "Connect to Ethereum node and print block number",
		Run:   connectEthClient,
	}

	rootCmd.AddCommand(helloCmd, startCmd, ethClientCmd) // Add the netstatusCmd here

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func connectEthClient(cmd *cobra.Command, args []string) {
	client, err := ethclient.Dial("http://127.0.0.1:8551")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to retrieve the latest Ethereum block header: %v", err)
	}

	fmt.Println("Current block number:", header.Number.String())
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
