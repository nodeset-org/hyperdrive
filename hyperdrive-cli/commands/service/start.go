package service

import (
	"fmt"
	"regexp"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	"github.com/urfave/cli/v2"
)

// Start the Hyperdrive service
func startService(c *cli.Context, ignoreConfigSuggestion bool) error {
	// Get RP client
	rp := client.NewClientFromCtx(c)

	// Update the Prometheus template with the assigned ports
	cfg, isNew, err := rp.LoadConfig()
	if err != nil {
		return fmt.Errorf("Error loading user settings: %w", err)
	}

	if isNew {
		return fmt.Errorf("No configuration detected. Please run `hyperdrive service config` to set up Hyperdrive before running it.")
	}

	// Check if this is a new install
	isUpdate, err := rp.IsFirstRun()
	if err != nil {
		return fmt.Errorf("error checking for first-run status: %w", err)
	}
	if isUpdate && !ignoreConfigSuggestion {
		if c.Bool(utils.YesFlag.Name) || utils.Confirm("Hyperdrive upgrade detected - starting will overwrite certain settings with the latest defaults (such as container versions).\nYou may want to run `hyperdrive service config` first to see what's changed.\n\nWould you like to continue starting the service?") {
			cfg.UpdateDefaults()
			rp.SaveConfig(cfg)
			fmt.Printf("%sUpdated settings successfully.%s\n", terminal.ColorGreen, terminal.ColorReset)
		} else {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Update the Prometheus template with the assigned ports
	metricsEnabled := cfg.Metrics.EnableMetrics.Value
	if metricsEnabled {
		err := rp.UpdatePrometheusConfiguration(cfg)
		if err != nil {
			return err
		}
	}

	// Validate the config
	errors := cfg.Validate()
	if len(errors) > 0 {
		fmt.Printf("%sYour configuration encountered errors. You must correct the following in order to start Hyperdrive:\n\n", terminal.ColorRed)
		for _, err := range errors {
			fmt.Printf("%s\n\n", err)
		}
		fmt.Println(terminal.ColorReset)
		return nil
	}

	// Start service
	err = rp.StartService(getComposeFiles(c))
	if err != nil {
		return err
	}

	// Remove the upgrade flag if it's there
	return rp.RemoveUpgradeFlagFile()
}

// Extract the image name from a Docker image string
func getDockerImageName(imageString string) (string, error) {
	// Return the empty string if the validator didn't exist (probably because this is the first time starting it up)
	if imageString == "" {
		return "", nil
	}

	reg := regexp.MustCompile(dockerImageRegex)
	matches := reg.FindStringSubmatch(imageString)
	if matches == nil {
		return "", fmt.Errorf("Couldn't parse the Docker image string [%s]", imageString)
	}
	imageIndex := reg.SubexpIndex("image")
	if imageIndex == -1 {
		return "", fmt.Errorf("Image name not found in Docker image [%s]", imageString)
	}

	imageName := matches[imageIndex]
	return imageName, nil
}
