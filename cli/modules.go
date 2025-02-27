package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"github.com/nodeset-org/hyperdrive/cli/client"
	"github.com/nodeset-org/hyperdrive/shared/utils"
	"github.com/urfave/cli/v2"
)

func HandleCommandNotFound(c *cli.Context, command string) {
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		fmt.Printf("Error creating Hyperdrive client: %v\n", err)
		return
	}

	// Make sure Hyperdrive has been configured already
	cfg, isNew, err := hd.LoadMainSettingsFile()
	if err != nil {
		fmt.Printf("Cannot run command [%s]: error loading main settings file: %v\n", command, err)
		return
	}
	if isNew {
		fmt.Println("Hyperdrive has not been configured yet. Please run 'hyperdrive service configure' first.")
		return
	}

	// Get the list of modules
	//fmt.Println("Loading modules...")
	results, err := hd.LoadModules()
	if err != nil {
		fmt.Printf("WARNING: Modules could not be loaded: %v\n", err)
		fmt.Println("Module commands will not be available until this is resolved.")
		c.App.CommandNotFound = nil
		fmt.Println(cli.ShowCommandHelp(c, command))
		return
	}
	if len(results) == 0 {
		c.App.CommandNotFound = nil
		fmt.Println(cli.ShowCommandHelp(c, command))
		return
	}

	// Organize modules by load status
	failedModules := map[string]*utils.ModuleInfoLoadResult{}
	succeededModules := map[string]*utils.ModuleInfoLoadResult{}
	for _, result := range results {
		if result.LoadError != nil {
			failedModules[string(result.Info.Descriptor.Shortcut)] = result
			continue
		}
		succeededModules[string(result.Info.Descriptor.Shortcut)] = result
	}

	// Check if the command belongs to a failed module
	if mod, exists := failedModules[command]; exists {
		fmt.Printf("Module %s failed to load: %s\n", mod.Info.Descriptor.Name, mod.LoadError)
		fmt.Printf("The [%s] command is not available until this is resolved.\n", command)
		return
	}

	// Ignore commands that don't belong to modules
	mod, exists := succeededModules[command]
	if !exists {
		c.App.CommandNotFound = nil
		fmt.Println(cli.ShowCommandHelp(c, command))
		return
	}

	// Make sure it's configured and enabled
	instance, exists := cfg.Modules[mod.Info.Descriptor.GetFullyQualifiedModuleName()]
	if !exists {
		fmt.Printf("Module %s has not been configured. Please run 'hyperdrive service configure' first.\n", mod.Info.Descriptor.Name)
		return
	}
	if !instance.Enabled {
		fmt.Printf("Module %s is disabled. Please enable it in your service configuration first.\n", mod.Info.Descriptor.Name)
		return
	}

	// Make sure the container exists before running the command
	// TODO: break this into a single utility function to dedup
	projectAdapterName := utils.GetProjectAdapterContainerName(&mod.Info.Descriptor, cfg.ProjectName)
	docker, err := dockerClient.NewClientWithOpts(
		dockerClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}
	containers, err := docker.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	id := ""
	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == "/"+projectAdapterName {
				id = container.ID
				break
			}
		}
		if id != "" {
			break
		}
	}
	if id == "" {
		fmt.Printf("Module %s does not have a project adapter container yet. Please run 'hyperdrive service start' first.\n", mod.Info.Descriptor.Name)
		return
	}

	// Make sure the container is running
	containerInfo, err := docker.ContainerInspect(context.Background(), id)
	if err != nil {
		fmt.Printf("Error inspecting container %s: %v\n", projectAdapterName, err)
		return
	}
	if !containerInfo.State.Running {
		fmt.Printf("The project adapter for module %s is not running yet. Please run 'hyperdrive service start' first.\n", mod.Info.Descriptor.Name)
		return
	}

	// Get the adapter client
	pac, err := hd.GetModuleManager().GetProjectAdapterClient(cfg.ProjectName, mod.Info.Descriptor.GetFullyQualifiedModuleName())
	if err != nil {
		fmt.Printf("Error getting project adapter client: %v\n", err)
		return
	}

	// Run the command
	args := []string{}
	for i, arg := range os.Args {
		if arg == command && len(os.Args) > i+1 {
			args = os.Args[i+1:]
			break
		}
	}
	logger := slog.Default()
	err = pac.Run(context.Background(), logger, strings.Join(args, " "), true)
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		return
	}
}
