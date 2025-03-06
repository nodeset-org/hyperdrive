package management

// Delete Hyperdrive's services, and remove the config and data folders
func (m *HyperdriveManager) TerminateService(composeFiles []string, configPath string) error {
	/*
		// Get the command to run with root privileges
		rootCmd, err := getEscalationCommand()
		if err != nil {
			return fmt.Errorf("could not get privilege escalation command: %w", err)
		}

		// Terminate the Docker containers
		cmd, err := m.compose(composeFiles, "down -v")
		if err != nil {
			return fmt.Errorf("error creating Docker artifact removal command: %w", err)
		}
		err = printOutput(cmd)
		if err != nil {
			return fmt.Errorf("error removing Docker artifacts: %w", err)
		}

		// Delete the Hyperdrive directory
		path, err := homedir.Expand(configPath)
		if err != nil {
			return fmt.Errorf("error loading Hyperdrive directory: %w", err)
		}
		fmt.Printf("Deleting Hyperdrive directory (%s)...\n", path)
		cmd = fmt.Sprintf("%s rm -rf %s", rootCmd, path)
		_, err = readOutput(cmd)
		if err != nil {
			return fmt.Errorf("error deleting Hyperdrive directory: %w", err)
		}

		fmt.Println("Termination complete.")
	*/
	return nil
}
