package adapter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/goccy/go-json"
)

const (
	// The universal prefix for all hyperdrive module commands run in the adapter
	HyperdriveModuleCommand string = "hd-module"
)

var (
	// Error when the container isn't found
	ErrContainerNotFound error = errors.New("container not found")

	// Error when the adapter container is offline
	ErrAdapterContainerOffline error = errors.New("adapter container is offline")
)

type AdapterClient struct {
	// The name of the container
	containerName string

	// The ID of the container, assuming it was found
	containerID string

	// Command prefix to use for running the adapter
	entrypoint []string

	// The Docker client
	dockerClient *client.Client

	// The key for authenticated requests
	key string
}

// Creates a new AdapterClient instance
func NewAdapterClient(containerName string, key string) (*AdapterClient, error) {
	dockerClient, err := client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	// Make sure the container exists
	id, err := getContainerID(dockerClient, containerName)
	if err != nil {
		return nil, err
	}

	// Check if the adapter container is running
	containerJson, err := dockerClient.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("error checking if adapter container is running: %w", err)
	}
	if !containerJson.State.Running {
		return nil, ErrAdapterContainerOffline
	}

	return &AdapterClient{
		containerName: containerName,
		containerID:   id,
		entrypoint:    containerJson.Config.Entrypoint,
		dockerClient:  dockerClient,
		key:           key,
	}, nil
}

// Get the ID of a container by its name
func getContainerID(dockerClient *client.Client, containerName string) (string, error) {
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return "", fmt.Errorf("error listing containers: %w", err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			trimmedName := strings.TrimPrefix(name, "/")
			if trimmedName == containerName {
				return container.ID, nil
			}
		}
	}
	return "", ErrContainerNotFound
}

// Run a docker exec command in the adapter container and get the result
func runCommand[RequestType any, ResponseType any](
	c *AdapterClient,
	ctx context.Context,
	command string,
	request *RequestType,
	response *ResponseType,
) error {
	// Start an exec command
	cmdArray := strings.Split(command, " ")
	cmdArray = append(c.entrypoint, cmdArray...)
	idResponse, err := c.dockerClient.ContainerExecCreate(ctx, c.containerID, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          cmdArray,
	})
	if err != nil {
		return fmt.Errorf("error creating exec command [%s]: %w", command, err)
	}

	// Attach reader/writer to the exec command
	execResponse, err := c.dockerClient.ContainerExecAttach(ctx, idResponse.ID, container.ExecAttachOptions{
		Tty: false,
	})
	if err != nil {
		return fmt.Errorf("error attaching to exec command [%s]: %w", command, err)
	}
	defer execResponse.Close()

	// Send the request down via stdin
	if request != nil {
		err = json.NewEncoder(execResponse.Conn).Encode(request)
		if err != nil {
			return fmt.Errorf("error sending request for command [%s]: %w", command, err)
		}
	}

	// Get the response
	var outBuf, errBuf bytes.Buffer
	exited := make(chan error)
	go func() {
		// Docker output demuxer to separate stdout and stderr, blocks until EOF
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, execResponse.Reader)
		exited <- err
	}()

	// Wait for an exit condition
	select {
	case <-ctx.Done():
		// The context was cancelled first
		return ctx.Err()

	case exitErr := <-exited:
		// The command exited first
		if exitErr != nil {
			return fmt.Errorf("error reading command [%s] response: %w", command, exitErr)
		}
		break
	}

	// Read the stdout and stderr
	stdout, err := io.ReadAll(&outBuf)
	if err != nil {
		return fmt.Errorf("error reading stdout for command [%s]: %w", command, err)
	}
	stderr, err := io.ReadAll(&errBuf)
	if err != nil {
		return fmt.Errorf("error reading stderr for command [%s]: %w", command, err)
	}

	// Get the exit code
	inspectResponse, err := c.dockerClient.ContainerExecInspect(ctx, idResponse.ID)
	if err != nil {
		return fmt.Errorf("error inspecting exec command result for [%s]: %w", command, err)
	}

	// Handle errors
	if inspectResponse.ExitCode != 0 {
		trimmedErr := strings.TrimSpace(string(stderr))
		if len(trimmedErr) > 0 {
			return fmt.Errorf("command [%s] errored with code %d: %s", command, inspectResponse.ExitCode, trimmedErr)
		}
		return fmt.Errorf("command [%s] errored with code %d and no error message", command, inspectResponse.ExitCode)
	}

	// Handle output if no response is expected
	if response == nil {
		trimmedErr := strings.TrimSpace(string(stderr))
		if len(trimmedErr) > 0 {
			fmt.Println(trimmedErr)
		}
		trimmedOut := strings.TrimSpace(string(stdout))
		if len(trimmedOut) > 0 {
			fmt.Println(trimmedOut)
		}
		return nil
	}

	// Handle output if a response is expected
	trimmedResult := strings.TrimSpace(string(stdout))
	if len(trimmedResult) == 0 {
		return fmt.Errorf("command [%s] returned an empty response with exit code %d and no error message", command, inspectResponse.ExitCode)
	}
	err = json.Unmarshal([]byte(trimmedResult), response)
	if err != nil {
		return fmt.Errorf("error unmarshalling response for command [%s]: %w", command, err)
	}
	return nil
}
