package adapter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/shared/utils/command"
)

const (
	RunCommandString string = HyperdriveModuleCommand + " run"
)

// Request format for `run`
type RunRequest struct {
	KeyedRequest

	// The command to run
	Command string `json:"command"`
}

// Run a command on the adapter
func (c *AdapterClient) RunNoninteractive(ctx context.Context, logger *slog.Logger, command string) (string, string, error) {
	request := &RunRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Command: command,
	}

	// Start an exec command
	cmdArray := strings.Split(command, " ")
	cmdArray = append(c.entrypoint, cmdArray...)
	idResponse, err := c.dockerClient.ContainerExecCreate(ctx, c.containerName, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          cmdArray,
	})
	if err != nil {
		return "", "", fmt.Errorf("error creating exec command [%s]: %w", command, err)
	}
	if logger != nil {
		logger.Debug(
			"Exec command created",
			"cmd", strings.Join(cmdArray, " "),
			"id", idResponse.ID,
		)
	}

	// Attach reader/writer to the exec command
	execResponse, err := c.dockerClient.ContainerExecAttach(ctx, idResponse.ID, container.ExecAttachOptions{
		Tty: false,
	})
	if err != nil {
		return "", "", fmt.Errorf("error attaching to exec command [%s]: %w", command, err)
	}
	defer execResponse.Close()
	if logger != nil {
		logger.Debug("Attached to exec command")
	}

	// Send the request down via stdin
	err = json.NewEncoder(execResponse.Conn).Encode(request)
	if err != nil {
		return "", "", fmt.Errorf("error sending request for command [%s]: %w", command, err)
	}

	// Start the output reader
	var outBuf, errBuf bytes.Buffer
	exited := make(chan error)
	go func() {
		// Docker output demuxer to separate stdout and stderr, blocks until EOF
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, execResponse.Reader)
		exited <- err
	}()

	// Start the input reader
	inStopped := make(chan error)
	go func() {
		// Copy the input to the exec command
		_, err = io.Copy(execResponse.Conn, os.Stdin)
		inStopped <- err
	}()

	// Run until something stops
	select {
	case <-ctx.Done():
		// The context was cancelled first
		return "", "", ctx.Err()

	case exitErr := <-exited:
		// The command exited first
		if exitErr != nil {
			return "", "", fmt.Errorf("error reading command [%s] response: %w", command, exitErr)
		}
		if logger != nil {
			logger.Debug("Exec command exited")
		}
		break
		/*
			case inErr := <-inStopped:
				// The input stopped first
				if inErr != nil {
					return fmt.Errorf("error sending input for command [%s]: %w", command, inErr)
				}
				logger.Debug("Input stopped")
				break
		*/
	}

	// Read the stdout and stderr
	stdout, err := io.ReadAll(&outBuf)
	if err != nil {
		return "", "", fmt.Errorf("error reading stdout for command [%s]: %w", command, err)
	}
	stderr, err := io.ReadAll(&errBuf)
	if err != nil {
		return "", "", fmt.Errorf("error reading stderr for command [%s]: %w", command, err)
	}

	// Get the exit code
	inspectResponse, err := c.dockerClient.ContainerExecInspect(ctx, idResponse.ID)
	if err != nil {
		return "", "", fmt.Errorf("error inspecting exec command result for [%s]: %w", command, err)
	}

	// Print the output
	if len(stdout) > 0 && logger != nil {
		logger.Info("Command output", "stdout", string(stdout))
	}
	if len(stderr) > 0 && logger != nil {
		logger.Error("Command error", "stderr", string(stderr))
	}
	if inspectResponse.ExitCode != 0 {
		return "", "", fmt.Errorf("command [%s] errored with code %d", command, inspectResponse.ExitCode)
	}
	return string(stdout), string(stderr), nil
}

// Run a command on the adapter
func (c *AdapterClient) Run2(ctx context.Context, cmd string, interactive bool) error {
	request := &RunRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Command: cmd,
	}

	// Start an exec command via the shell
	entrypoint := strings.Join(c.entrypoint, " ")
	cmdString := "docker exec -it " + c.containerName + " " + entrypoint + " hd-module run"
	fullCmd := command.NewCommand(cmdString)

	// Send the request down via stdin
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(request)
	if err != nil {
		return fmt.Errorf("error sending request for command [%s]: %w", cmdString, err)
	}
	buf.Write([]byte("\n"))
	fullCmd.SetStdin(buf)

	outbuf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	fullCmd.SetStdout(outbuf)
	fullCmd.SetStderr(errbuf)

	// Start the command
	if err := fullCmd.Start(); err != nil {
		return err
	}

	// Wait for the command to exit
	err = fullCmd.Wait()

	// Print the output
	if outbuf.Len() > 0 {
		fmt.Println(outbuf.String())
	}
	if errbuf.Len() > 0 {
		fmt.Println(errbuf.String())
	}
	if err != nil {
		return fmt.Errorf("error running command [%s]: %w", cmdString, err)
	}
	return nil
}
