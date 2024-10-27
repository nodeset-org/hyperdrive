package service

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"

	dtc "github.com/docker/docker/api/types/container"
	dmount "github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
)

const (
	daemonBinLocation   string = "/usr/bin/hyperdrive-daemon"
	daemonImageRegexStr string = "nodeset/hyperdrive:v.*"
	daemonUserDirFlag   string = "--user-dir"
)

var (
	daemonImageRegex *regexp.Regexp = regexp.MustCompile(daemonImageRegexStr)
)

// Called by package managers to restart the service after installation if it was already running
// Note this is going to run as a superuser, not necessarily the user running Hyperdrive, so we can't rely on the default context flags
func safeStartAfterInstall(systemDir string) {
	hyperdriveBinPath := os.Args[0]

	// Get a Docker client
	d, err := docker.NewClientWithOpts(docker.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Error creating Docker client: %s\n", err.Error())
		return
	}

	// Get a list of all running images
	cl, err := d.ContainerList(context.Background(), dtc.ListOptions{All: false})
	if err != nil {
		fmt.Printf("Error getting running Docker container list: %s\n", err.Error())
		return
	}

	// Find the daemon containers
	daemonIds := []string{}
	for _, container := range cl {
		isDaemon := daemonImageRegex.MatchString(container.Image)
		if !isDaemon {
			continue
		}
		name := strings.TrimPrefix(container.Names[0], "/")
		fmt.Printf("Found running daemon container [%s] with image [%s]\n", name, container.Image)
		daemonIds = append(daemonIds, container.ID)
	}

	// Run a service start for each daemon container
	for _, id := range daemonIds {
		containerInfo, err := d.ContainerInspect(context.Background(), id)
		if err != nil {
			fmt.Printf("Error inspecting container %s: %s\n", id, err.Error())
			continue
		}
		containerName := strings.TrimPrefix(containerInfo.Name, "/")

		// Make sure it has the system dir mounted
		usesSystemDir := false
		for _, mount := range containerInfo.Mounts {
			if mount.Type != dmount.TypeBind {
				continue
			}
			if !strings.HasPrefix(mount.Source, systemDir) {
				continue
			}
			usesSystemDir = true
			break
		}
		if !usesSystemDir {
			fmt.Printf("Daemon container [%s] does not have the system dir mounted, skipping...\n", containerName)
			continue
		}

		// Get the user dir from the args
		userDir := ""
		for i := 0; i < len(containerInfo.Args); i++ {
			arg := containerInfo.Args[i]
			if arg != daemonUserDirFlag {
				continue
			}
			if i+1 >= len(containerInfo.Args) {
				break
			}
			userDir = containerInfo.Args[i+1]
			break
		}
		if userDir == "" {
			// Couldn't get the user dir from the args, so skip this container
			fmt.Printf("Daemon container [%s] does not have a user dir arg, skipping...\n", containerName)
			continue
		}

		// Get the owner of the user dir
		userDirStat := syscall.Stat_t{}
		err = syscall.Stat(userDir, &userDirStat)
		if err != nil {
			fmt.Printf("Error getting user dir [%s] stat for [%s]: %s\n", userDir, containerName, err.Error())
			continue
		}
		owner := userDirStat.Uid

		// Start the service
		success := client.StartServiceAsUser(owner, hyperdriveBinPath, userDir)
		if !success {
			fmt.Printf("WARN: starting service for container [%s] failed.\n", containerName)
			fmt.Println("Please restart the service manually with `hyperdrive service start` to update the daemon services.")
			continue
		}
	}
}
