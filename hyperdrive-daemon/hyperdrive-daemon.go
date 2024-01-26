package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/services"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/tasks"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

// Run
func main() {
	// Add logo and attribution to application help template
	attribution := "ATTRIBUTION:\n   Adapted from the Rocket Pool Smart Node (https://github.com/rocketpool/smartnode) with love."
	cli.AppHelpTemplate = fmt.Sprintf("\n%s\n\n%s\n%s\n", shared.Logo, cli.AppHelpTemplate, attribution)
	cli.CommandHelpTemplate = fmt.Sprintf("%s\n%s\n", cli.CommandHelpTemplate, attribution)
	cli.SubcommandHelpTemplate = fmt.Sprintf("%s\n%s\n", cli.SubcommandHelpTemplate, attribution)

	// Initialise application
	app := cli.NewApp()

	// Set application info
	app.Name = "hyperdrive-daemon"
	app.Usage = "Hyperdrive Daemon for NodeSet Node Operator Management"
	app.Version = shared.HyperdriveVersion
	app.Authors = []*cli.Author{
		{
			Name:  "Joe Clapis",
			Email: "joe@nodeset.io",
		},
	}
	app.Copyright = "(C) 2024 NodeSet LLC"
	app.Flags = []cli.Flag{}
	app.Action = func(ctx *cli.Context) error {
		// Wait group to handle the API server (separate because of error handling)
		apiWg := new(sync.WaitGroup)
		apiWg.Add(1)

		// Wait group to handle the task loop server
		taskWg := new(sync.WaitGroup)
		taskWg.Add(1)

		// Create the service provider
		sp, err := services.NewServiceProvider(config.DaemonConfigPath)
		if err != nil {
			return fmt.Errorf("error creating service provider: %w", err)
		}

		// Get the owner of the config file
		var cfgFileStat syscall.Stat_t
		err = syscall.Stat(config.DaemonConfigPath, &cfgFileStat)
		if err != nil {
			return fmt.Errorf("error getting config file [%s] info: %w", config.DaemonConfigPath, err)
		}

		// Start the API manager
		apiMgr := api.NewApiManager(sp)
		err = apiMgr.Start(apiWg, cfgFileStat.Uid, cfgFileStat.Gid)
		if err != nil {
			return fmt.Errorf("error starting API manager: %w", err)
		}

		// Start the task loop
		taskLoop := tasks.NewTaskLoop(sp, taskWg)
		err = taskLoop.Run()
		if err != nil {
			return fmt.Errorf("error starting task loop: %w", err)
		}

		// TODO: Metrics manager

		// Handle process closures
		termListener := make(chan os.Signal, 1)
		signal.Notify(termListener, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-termListener
			fmt.Println("Shutting down daemon...")
			err := apiMgr.Stop()
			if err != nil {
				fmt.Printf("WARNING: daemon didn't shutdown cleanly: %s\n", err.Error())
				apiWg.Done()
			}
			taskLoop.Stop()
		}()

		// Run the daemon until closed
		fmt.Printf("Started daemon on %s.\n", config.DaemonSocketPath)
		apiWg.Wait()
		taskWg.Wait()
		fmt.Println("Daemon stopped.")
		return nil
	}

	// Run application
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
