package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/api"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/common/services"
	"github.com/spf13/cobra"
)

const (
	DaemonSocketPath string = "/.hyperdrive/data/sockets/daemon.sock"
)

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Short: "Hyperdrive daemon for Stakewise",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hyperdrive daemon for Stakewise")
		},
	}

	initCmd := &cobra.Command{
		Use:   "run hyperdrive-config-path",
		Short: "Run the Daemon",
		Args: func(cmd *cobra.Command, args []string) error {
			err := cobra.MinimumNArgs(1)(cmd, args)
			if err != nil {
				return err
			}
			snSettingsPath := args[0]
			if !strings.HasPrefix(snSettingsPath, "/") {
				return fmt.Errorf("hyperdrive settings path must be an absolute path")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			handleError("init", runDaemonApi(args[0]))
		},
	}

	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "todo",
		Args: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("register")

			// TODO: Do as API call but as an example of what should be called:

			// 1. Create RP node contract
			// node.RegisterNode()

			// 2. Set withdrawal address (use constellation smart contract - parameter from .env - by deployment vs network???)
			// node.SetWithdrawalAddress()
		},
	}

	prepareValidator := &cobra.Command{
		Use:   "prepare-validator",
		Short: "todo",
		Args: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		// TODO: Run this as a post-hook to the register command
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("prepare-validator")
			// 1. Create new validator key (derived from node wallet key)
			// 2. Create initial DepositData for pubkey (1 eth deposit)
			// 3. Submit deposit data + exit message to Pyxis
			// 4. Send Pyxis response (NewValidatorAuthorization) to constellation smart contract
			// 5. Loop over and create as many minipools as possible (1 eth node deposit)
		},
	}

	nodeFinalize := &cobra.Command{
		Use:   "node-final-deposit",
		Short: "todo",
		Args: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		// TODO: Run this as a post-hook to the nodeInitDeposit command
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("node-final-deposit")
			// 1. Submit 2nd deposit after 12 hr scrub-check
			// 2. Create ExitMessage for new key (This happens later - discuss more with Joe/Mike)

		},
	}
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(prepareValidator)
	rootCmd.AddCommand(nodeFinalize)

	// Run the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func handleError(command string, err error) {
	if err != nil {
		fmt.Printf("Error running %s: %s\n", command, err.Error())
		os.Exit(1)
	}
}

func runDaemonApi(cfgPath string) error {
	// Wait group to handle the daemon
	wg := new(sync.WaitGroup)
	wg.Add(1)

	// Create the service provider
	sp, err := services.NewServiceProvider(cfgPath)
	if err != nil {
		return fmt.Errorf("error creating service provider: %w", err)
	}

	// Start the watchtower relayer
	cfg := sp.GetConfig()
	socketPath := DaemonSocketPath
	debug := cfg.DebugMode.Value
	manager := api.NewApiManager(sp, socketPath, debug)
	err = manager.Start(wg)
	if err != nil {
		return fmt.Errorf("error starting daemon: %w", err)
	}

	// Handle process closures
	termListener := make(chan os.Signal, 1)
	signal.Notify(termListener, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-termListener
		fmt.Println("Shutting down daemon...")
		err := manager.Stop()
		if err != nil {
			fmt.Printf("WARNING: daemon didn't shutdown cleanly: %s\n", err.Error())
			wg.Done()
		}
	}()

	// Run the daemon until closed
	fmt.Printf("Started daemon on %s.\n", socketPath)
	wg.Wait()
	fmt.Println("Daemon stopped.")
	return nil
}
