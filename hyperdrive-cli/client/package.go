package client

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// Starts Hyperdrive as the provided user
func StartServiceAsUser(uid uint32, hyperdriveBin string, configPath string) bool {
	// Get the user from the uid
	user, err := user.LookupId(fmt.Sprintf("%d", uid))
	if err != nil {
		fmt.Printf("ERROR: Could not find user with ID %d: %s\n", uid, err.Error())
		return false
	}
	name := user.Username

	// Prep the hyperdrive command
	hyperdriveCmd := []string{hyperdriveBin}
	if uid == 0 {
		hyperdriveCmd = append(hyperdriveCmd, "--allow-root")
	}
	hyperdriveCmd = append(hyperdriveCmd, "--config-path", configPath, "service", "start", "-y")

	// Check for sudo first
	sudo := "sudo"
	exists, _ := checkIfCommandExists(sudo)
	if exists {
		err := startServiceAsUserViaSudo(name, hyperdriveCmd)
		if err == nil {
			return true
		}
		fmt.Printf("WARN: Error starting Hyperdrive as user %d via sudo: %s\n", uid, err.Error())
	}

	// Check for su next
	su := "su"
	exists, _ = checkIfCommandExists(su)
	if exists {
		err := startServiceAsUserViaSu(name, hyperdriveCmd)
		if err == nil {
			return true
		}
		fmt.Printf("WARN: Error starting Hyperdrive as user %d via su: %s\n", uid, err.Error())
	}

	/*
		// Check for doas next
		doas := "doas"
		exists, _ = checkIfCommandExists(doas)
		if exists {
			err := startServiceAsUserViaDoas(name, hyperdriveCmd)
			if err == nil {
				return true
			}
			fmt.Printf("WARN: Error starting Hyperdrive as user %d via doas: %s\n", uid, err.Error())
		}
	*/

	fmt.Printf("ERROR: All privilege escalation commands to start Hyperdrive as user %d failed or were not found.\n", uid)
	return false
}

// Starts Hyperdrive as the provided user via `sudo`
func startServiceAsUserViaSudo(name string, hyperdriveCmd []string) error {
	// Set up the args for sudo
	args := []string{
		"-i",
		"-u",
		name,
	}
	args = append(args, hyperdriveCmd...)

	// Create the command and set the output to the terminal
	cmd := &command{
		cmd: exec.Command("sudo", args...),
	}
	return runStartServiceCommand(cmd)
}

// Starts Hyperdrive as the provided user via `su`
func startServiceAsUserViaSu(name string, hyperdriveCmd []string) error {
	// Set up the args for su
	args := []string{
		"-l",
		name,
		"-c",
	}
	hyperdriveCmdStr := fmt.Sprintf("'%s'", strings.Join(hyperdriveCmd, " "))
	args = append(args, hyperdriveCmdStr)

	// Create the command and set the output to the terminal
	cmd := &command{
		cmd: exec.Command("su", args...),
	}
	return runStartServiceCommand(cmd)
}

// Runs the given service start command
func runStartServiceCommand(cmd *command) error {
	cmd.SetStdin(os.Stdin)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	// Run the command
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
