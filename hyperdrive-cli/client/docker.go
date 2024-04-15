package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	dt "github.com/docker/docker/api/types"
	dtc "github.com/docker/docker/api/types/container"
	"github.com/nodeset-org/hyperdrive/shared/config"
)

// Get the current Docker image used by the given container
func (c *HyperdriveClient) GetDockerImage(containerName string) (string, error) {
	ci, err := inspectContainer(c, containerName)
	if err != nil {
		return "", err
	}
	return ci.Config.Image, nil
}

// Get the Docker images with the project ID as a prefix that run the VC start script in their command line arguments
func (c *HyperdriveClient) GetValidatorContainers(projectName string) ([]string, error) {
	d, err := c.GetDocker()
	if err != nil {
		return nil, err
	}
	cl, err := d.ContainerList(context.Background(), dtc.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("error getting container list: %w", err)
	}

	// Find all of them that belong to the project
	containers := []string{}
	for _, container := range cl {
		isProjectContainer := false
		for _, name := range container.Names {
			name = strings.TrimPrefix(name, "/") // Docker throws a leading / on names
			if strings.HasPrefix(name, projectName) {
				isProjectContainer = true
			}
		}

		// This container belongs to the project
		if isProjectContainer && strings.Contains(container.Command, config.VcStartScript) {
			name := strings.TrimPrefix(container.Names[0], "/")
			containers = append(containers, name)
		}
	}
	return containers, nil
}

// Get the current Docker image used by the given container
func (c *HyperdriveClient) GetDockerStatus(containerName string) (string, error) {
	ci, err := inspectContainer(c, containerName)
	if err != nil {
		return "", err
	}
	return ci.State.Status, nil
}

// Get the time that the given container shut down
func (c *HyperdriveClient) GetDockerContainerShutdownTime(containerName string) (time.Time, error) {
	ci, err := inspectContainer(c, containerName)
	if err != nil {
		return time.Time{}, err
	}

	// Parse the time
	finishTime, err := time.Parse(time.RFC3339, strings.TrimSpace(ci.State.FinishedAt))
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing container [%s] exit time [%s]: %w", containerName, ci.State.FinishedAt, err)
	}
	return finishTime, nil
}

// Shut down a container
func (c *HyperdriveClient) StopContainer(containerName string) error {
	d, err := c.GetDocker()
	if err != nil {
		return err
	}
	return d.ContainerStop(context.Background(), containerName, dtc.StopOptions{})
}

// Start a container
func (c *HyperdriveClient) StartContainer(containerName string) error {
	d, err := c.GetDocker()
	if err != nil {
		return err
	}
	return d.ContainerStart(context.Background(), containerName, dtc.StartOptions{})
}

// Restart a container
func (c *HyperdriveClient) RestartContainer(containerName string) error {
	d, err := c.GetDocker()
	if err != nil {
		return err
	}
	return d.ContainerRestart(context.Background(), containerName, dtc.StopOptions{})
}

// Deletes a container
func (c *HyperdriveClient) RemoveContainer(containerName string) error {
	d, err := c.GetDocker()
	if err != nil {
		return err
	}
	return d.ContainerRemove(context.Background(), containerName, dtc.RemoveOptions{})
}

// Deletes a volume
func (c *HyperdriveClient) DeleteVolume(volumeName string) error {
	d, err := c.GetDocker()
	if err != nil {
		return err
	}
	return d.VolumeRemove(context.Background(), volumeName, false)
}

// Gets the absolute file path of the client volume
func (c *HyperdriveClient) GetClientVolumeSource(containerName string, volumeTarget string) (string, error) {
	ci, err := inspectContainer(c, containerName)
	if err != nil {
		return "", err
	}

	// Find the mount with the provided destination
	for _, mount := range ci.Mounts {
		if mount.Destination == volumeTarget {
			return mount.Source, nil
		}
	}
	return "", fmt.Errorf("container [%s] doesn't have a volume with [%s] as a destination", containerName, volumeTarget)
}

// Gets the name of the client volume
func (c *HyperdriveClient) GetClientVolumeName(containerName, volumeTarget string) (string, error) {
	ci, err := inspectContainer(c, containerName)
	if err != nil {
		return "", err
	}

	// Find the mount with the provided destination
	for _, mount := range ci.Mounts {
		if mount.Destination == volumeTarget {
			return mount.Name, nil
		}
	}
	return "", fmt.Errorf("container [%s] doesn't have a volume with [%s] as a destination", containerName, volumeTarget)
}

// Gets the disk usage of the given volume
func (c *HyperdriveClient) GetVolumeSize(volumeName string) (int64, error) {
	d, err := c.GetDocker()
	if err != nil {
		return 0, err
	}

	du, err := d.DiskUsage(context.Background(), dt.DiskUsageOptions{})
	if err != nil {
		return 0, fmt.Errorf("error getting disk usage: %w", err)
	}
	for _, volume := range du.Volumes {
		if volume.Name == volumeName {
			return volume.UsageData.Size, nil
		}
	}
	return 0, fmt.Errorf("couldn't find a volume named [%s]", volumeName)
}

// Inspect a Docker container
func inspectContainer(c *HyperdriveClient, container string) (dt.ContainerJSON, error) {
	d, err := c.GetDocker()
	if err != nil {
		return dt.ContainerJSON{}, err
	}
	ci, err := d.ContainerInspect(context.Background(), container)
	if err != nil {
		return dt.ContainerJSON{}, fmt.Errorf("error inspecting container [%s]: %w", container, err)
	}
	return ci, nil
}
