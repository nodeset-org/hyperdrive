package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/nodeset-org/hyperdrive/daemon-utils/services"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/server"
	swtasks "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/tasks"
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
	app.Name = "stakewise-daemon"
	app.Usage = "Hyperdrive Daemon for NodeSet Stakewise Module Management"
	app.Version = shared.HyperdriveVersion
	app.Authors = []*cli.Author{
		{
			Name:  "Nodeset",
			Email: "info@nodeset.io",
		},
	}
	app.Copyright = "(C) 2024 NodeSet LLC"

	moduleDirFlag := &cli.StringFlag{
		Name:     "module-dir",
		Aliases:  []string{"d"},
		Usage:    "The path to the Stakewise module data directory",
		Required: true,
	}

	app.Flags = []cli.Flag{
		moduleDirFlag,
	}
	app.Action = func(c *cli.Context) error {
		// Get the config file
		moduleDir := c.String(moduleDirFlag.Name)
		hyperdriveSocketPath := filepath.Join(moduleDir, config.HyperdriveCliSocketFilename)
		_, err := os.Stat(hyperdriveSocketPath)
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("Hyperdrive socket not found at [%s].", hyperdriveSocketPath)
			os.Exit(1)
		}

		// Wait group to handle the API server (separate because of error handling)
		stopWg := new(sync.WaitGroup)
		stopWg.Add(1)

		// Create the service provider
		sp, err := services.NewServiceProvider(moduleDir, swconfig.ModuleName, swconfig.NewStakewiseConfig, config.ClientTimeout)
		if err != nil {
			return fmt.Errorf("error creating service provider: %w", err)
		}
		stakewiseSp, err := swcommon.NewStakewiseServiceProvider(sp)
		if err != nil {
			return fmt.Errorf("error creating Stakewise service provider: %w", err)
		}

		// Get the owner of the Hyperdrive socket
		var hdSocketStat syscall.Stat_t
		err = syscall.Stat(hyperdriveSocketPath, &hdSocketStat)
		if err != nil {
			return fmt.Errorf("error getting Hyperdrive socket file [%s] info: %w", hyperdriveSocketPath, err)
		}

		// Start the server
		apiServer, err := server.NewStakewiseServer(stakewiseSp)
		if err != nil {
			return fmt.Errorf("error creating Stakewise server: %w", err)
		}
		err = apiServer.Start(stopWg, hdSocketStat.Uid, hdSocketStat.Gid)
		if err != nil {
			return fmt.Errorf("error starting API manager: %w", err)
		}
		fmt.Printf("Started daemon on %s.\n", apiServer.GetSocketPath())

		// Start the task loop
		taskLoop := swtasks.NewTaskLoop(stakewiseSp, stopWg)
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
			err := apiServer.Stop()
			if err != nil {
				fmt.Printf("WARNING: daemon didn't shutdown cleanly: %s\n", err.Error())
				stopWg.Done()
			}
			taskLoop.Stop()
		}()

		// Run the daemon until closed
		fmt.Println("Daemon online.")
		stopWg.Wait()
		fmt.Println("Daemon stopped.")
		return nil
	}

	// Run application
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
