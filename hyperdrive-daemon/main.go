package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/services"
	"github.com/spf13/cobra"
)

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Short: "Hyperdrive daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hyperdrive daemon")
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
	rootCmd.AddCommand(initCmd)

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
	socketPath := cfg.DaemonSocketPath.Value
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
