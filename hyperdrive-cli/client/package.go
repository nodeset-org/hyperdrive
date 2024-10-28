package client

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// Starts Hyperdrive as the provided user via `su`
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
	hyperdriveCmdStr := fmt.Sprintf("%s", strings.Join(hyperdriveCmd, " "))

	// Run `su` to start Hyperdrive as the user
	cmd := &command{
		cmd: exec.Command("su", "-l", name, "-c", hyperdriveCmdStr),
	}
	err = runStartServiceCommand(cmd)
	if err == nil {
		return true
	}

	fmt.Printf("WARN: Error starting Hyperdrive as user '%s' (%d): %s\n", name, uid, err.Error())
	return false
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
