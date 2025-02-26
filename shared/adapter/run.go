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
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive/shared/adapter/interactive"
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

// Run a command in "interactive" mode, effectively proxying the current terminal into the adapter
func (c *AdapterClient) Run(ctx context.Context, logger *slog.Logger, command string, isInteractive bool) error {
	request := &RunRequest{
		KeyedRequest: KeyedRequest{
			Key: c.key,
		},
		Command: command,
	}

	// Start an exec command
	cmdArray := append(c.entrypoint, strings.Split(RunCommandString, " ")...)
	idResponse, err := c.dockerClient.ContainerExecCreate(ctx, c.containerName, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          isInteractive,
		Cmd:          cmdArray,
	})
	if err != nil {
		return fmt.Errorf("error creating exec command [%s]: %w", command, err)
	}
	if logger != nil {
		logger.Debug(
			"Exec command created",
			"cmd", strings.Join(cmdArray, " "),
			"id", idResponse.ID,
		)
	}

	// Attach to the exec command
	streams := interactive.NewStandardStreamWrapper(isInteractive)
	var consoleSize *[2]uint
	if isInteractive {
		height, width := streams.Out().GetTtySize()
		consoleSize = &[2]uint{height, width}
	}
	execResponse, err := c.dockerClient.ContainerExecAttach(ctx, idResponse.ID, container.ExecAttachOptions{
		Tty:         isInteractive,
		ConsoleSize: consoleSize,
	})
	if err != nil {
		return fmt.Errorf("error attaching to exec command [%s]: %w", command, err)
	}
	defer execResponse.Close()
	if logger != nil {
		logger.Debug("Attached to exec command")
	}

	// Send the request down via stdin before attaching streams
	buffer := new(bytes.Buffer)
	err = json.NewEncoder(buffer).Encode(request)
	if err != nil {
		return fmt.Errorf("error encoding request for command [%s]: %w", command, err)
	}
	bufferSize := buffer.Len()
	_, err = io.Copy(execResponse.Conn, buffer)
	if err != nil {
		return fmt.Errorf("error sending request for command [%s]: %w", command, err)
	}
	readBuffer := make([]byte, bufferSize)
	totalBytesRead := 0

	// NOTE: for some reason the request we just sent is going to get echoed back to us
	// so we need to read it first to clear the hijacked connection before attaching the streams.
	// That way it doesn't get printed to the user's terminal.
	for totalBytesRead < bufferSize {
		bytesRead, err := execResponse.Conn.Read(readBuffer)
		if err != nil {
			return fmt.Errorf("error reading response for command [%s]: %w", command, err)
		}
		readBuffer = readBuffer[bytesRead:]
		totalBytesRead += bytesRead
	}
	if totalBytesRead != bufferSize {
		return fmt.Errorf("error reading response for command [%s]: expected %d bytes, got %d", command, bufferSize, totalBytesRead)
	}

	// Attach streams to the container
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- func() error {
			streamer := interactive.NewHijackedIOStreamer(isInteractive, logger, streams, execResponse)
			return streamer.Stream(ctx)
		}()
	}()

	// Handle terminal resizing
	if isInteractive && streams.In().IsTerminal() {
		if err := interactive.MonitorTtySize(ctx, logger, c.dockerClient, streams, idResponse.ID, true); err != nil {
			_, _ = fmt.Fprintln(streams.Err(), "Error monitoring TTY size:", err)
		}
	}

	// Wait for the command to exit
	if err := <-errCh; err != nil {
		logger.Debug(fmt.Sprintf("Error hijack: %s", err))
		return err
	}
	return getExecExitStatus(ctx, c.dockerClient, idResponse.ID)
}

func getExecExitStatus(ctx context.Context, docker client.ContainerAPIClient, execID string) error {
	resp, err := docker.ContainerExecInspect(ctx, execID)
	if err != nil {
		// If we can't connect, then the daemon probably died.
		if !client.IsErrConnectionFailed(err) {
			return err
		}
		return fmt.Errorf("error inspecting container process: %w", err)
	}
	status := resp.ExitCode
	if status != 0 {
		return fmt.Errorf("container process exited with status code %d", status)
	}
	return nil
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
