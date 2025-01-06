package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// ===============
// === Command ===
// ===============

// A Command to be executed either locally or remotely
type Command struct {
	cmd *exec.Cmd
}

// Create a command to be run by Hyperdrive
func NewCommand(cmdText string) *Command {
	return &Command{
		cmd: exec.Command("sh", "-c", cmdText),
	}
}

// Run the command
func (c *Command) Run() error {
	return c.cmd.Run()
}

// Start executes the command. Don't forget to call Wait
func (c *Command) Start() error {
	return c.cmd.Start()
}

// Wait for the command to exit
func (c *Command) Wait() error {
	return c.cmd.Wait()
}

func (c *Command) SetStdin(r io.Reader) {
	c.cmd.Stdin = r
}

func (c *Command) SetStdout(w io.Writer) {
	c.cmd.Stdout = w
}

func (c *Command) SetStderr(w io.Writer) {
	c.cmd.Stderr = w
}

// Run the command and return its output
func (c *Command) Output() ([]byte, error) {
	return c.cmd.Output()
}

// Get a pipe to the command's stdout
func (c *Command) StdoutPipe() (io.Reader, error) {
	return c.cmd.StdoutPipe()
}

// Get a pipe to the command's stderr
func (c *Command) StderrPipe() (io.Reader, error) {
	return c.cmd.StderrPipe()
}

// OutputPipes pipes for stdout and stderr
func (c *Command) OutputPipes() (io.Reader, io.Reader, error) {
	cmdOut, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	cmdErr, err := c.StderrPipe()

	return cmdOut, cmdErr, err
}

// =============
// === Utils ===
// =============

// Get the command used to escalate privileges on the system
func GetEscalationCommand() (string, error) {
	// Check for sudo first
	sudo := "sudo"
	exists, err := CheckIfCommandExists(sudo)
	if err != nil {
		return "", fmt.Errorf("error checking if %s exists: %w", sudo, err)
	}
	if exists {
		return sudo, nil
	}

	// Check for doas next
	doas := "doas"
	exists, err = CheckIfCommandExists(doas)
	if err != nil {
		return "", fmt.Errorf("error checking if %s exists: %w", doas, err)
	}
	if exists {
		return doas, nil
	}

	return "", fmt.Errorf("no privilege escalation command found")
}

// Checks if a command exists on the system
func CheckIfCommandExists(command string) (bool, error) {
	// Run `type` to check for existence
	cmd := fmt.Sprintf("type %s", command)
	output, err := ReadOutput(cmd)

	if err != nil {
		exitErr, isExitErr := err.(*exec.ExitError)
		if isExitErr && exitErr.ProcessState.ExitCode() == 127 {
			// Command not found
			return false, nil
		} else {
			return false, fmt.Errorf("error checking if %s exists: %w", command, err)
		}
	} else {
		if strings.Contains(string(output), fmt.Sprintf("%s is", command)) {
			return true, nil
		} else {
			return false, fmt.Errorf("unexpected output when checking for %s: %s", command, string(output))
		}
	}
}

// Run a command and print its output
func PrintOutput(cmdText string) error {
	// Initialize command
	cmd := NewCommand(cmdText)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Wait for the command to exit
	return cmd.Wait()
}

// Run a command and return its output
func ReadOutput(cmdText string) ([]byte, error) {
	// Initialize command
	cmd := NewCommand(cmdText)

	// Run command and return output
	return cmd.Output()
}
