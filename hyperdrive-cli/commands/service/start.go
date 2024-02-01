package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/commands/wallet"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils/terminal"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	"github.com/nodeset-org/hyperdrive/shared"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	"github.com/urfave/cli/v2"
)

// Start the Hyperdrive service
func startService(c *cli.Context, ignoreConfigSuggestion bool) error {
	// Get Hyperdrive client
	hd := client.NewHyperdriveClientFromCtx(c)

	// Update the Prometheus template with the assigned ports
	cfg, isNew, err := hd.LoadConfig()
	if err != nil {
		return fmt.Errorf("Error loading user settings: %w", err)
	}

	if isNew {
		return fmt.Errorf("No configuration detected. Please run `hyperdrive service config` to set up Hyperdrive before running it.")
	}

	// Check if this is a new install
	oldVersion := strings.TrimPrefix(cfg.Version, "v")
	currentVersion := strings.TrimPrefix(shared.HyperdriveVersion, "v")
	isUpdate := oldVersion != currentVersion
	if isUpdate && !ignoreConfigSuggestion {
		if c.Bool(utils.YesFlag.Name) || utils.Confirm("Hyperdrive upgrade detected - starting will overwrite certain settings with the latest defaults (such as container versions).\nYou may want to run `hyperdrive service config` first to see what's changed.\n\nWould you like to continue starting the service?") {
			cfg.UpdateDefaults()
			hd.SaveConfig(cfg)
			fmt.Printf("%sUpdated settings successfully.%s\n", terminal.ColorGreen, terminal.ColorReset)
		} else {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Update the Prometheus and Grafana config templates with the assigned ports
	metricsEnabled := cfg.Metrics.EnableMetrics.Value
	if metricsEnabled {
		err := hd.UpdatePrometheusConfiguration(cfg)
		if err != nil {
			return err
		}
		err = hd.UpdateGrafanaDatabaseConfiguration(cfg)
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

	// TODO: SLASHING DELAY
	if !c.Bool(ignoreSlashTimerFlag.Name) {
		// Do the client swap check
		firstRun, err := checkForValidatorChange(hd, cfg)
		if err != nil {
			fmt.Printf("%sWARNING: couldn't verify that the validator container can be safely restarted:\n\t%s\n", terminal.ColorYellow, err.Error())
			fmt.Println("If you are changing to a different client, it may resubmit an attestation you have already submitted.")
			fmt.Println("This will slash your validator!")
			fmt.Println("To prevent slashing, you must wait 15 minutes from the time you stopped the clients before starting them again.")
			fmt.Println()
			fmt.Println("**If you did NOT change clients, you can safely ignore this warning.**")
			fmt.Println()
			if !utils.Confirm(fmt.Sprintf("Press y when you understand the above warning, have waited, and are ready to start Hyperdrive:%s", terminal.ColorReset)) {
				fmt.Println("Cancelled.")
				return nil
			}
		} else if firstRun {
			fmt.Println("It looks like this is your first time starting Hyperdrive.")
			existingNode := utils.Confirm("Just to be sure, are you rebuilding a node that has existing, active validators attesting on the Beacon Chain?")
			if !existingNode {
				fmt.Println("Okay, great! You're safe to start. Have fun!")
			} else {
				fmt.Printf("%sSince you're rebuilding a node, Hyperdrive can't determine if you attested in the last 15 minutes.\n", terminal.ColorYellow)
				fmt.Println("If you did, it may resubmit an attestation you have already submitted.")
				fmt.Println("This will slash your validator!")
				fmt.Println("To prevent slashing, you must wait 15 minutes from the time you stopped the clients before starting them again.")
				fmt.Println()
				if !utils.Confirm(fmt.Sprintf("Press y when you understand the above warning, have waited, and are ready to start Hyperdrive:%s", terminal.ColorReset)) {
					fmt.Println("Cancelled.")
					return nil
				}
			}
		}
	} else {
		fmt.Printf("%sIgnoring anti-slashing safety delay.%s\n", terminal.ColorYellow, terminal.ColorReset)
	}

	// Write a note on doppelganger protection
	if cfg.IsDoppelgangerEnabled() {
		fmt.Printf("%sNOTE: You currently have Doppelganger Protection enabled on at least one module.\nYour Validator Client will miss up to 3 attestations when it starts.\nThis is *intentional* and does not indicate a problem with your node.%s\n\n", terminal.ColorBold, terminal.ColorReset)
	}

	// Start service
	err = hd.StartService(getComposeFiles(c))
	if err != nil {
		return fmt.Errorf("error starting service: %w", err)
	}

	// Check wallet status
	fmt.Println()
	fmt.Println("Checking node wallet status...")
	var status *types.WalletStatus
	retries := 5
	for i := 0; i < retries; i++ {
		response, err := hd.Api.Wallet.Status()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		status = &response.Data.WalletStatus
		break
	}

	// Handle errors
	if status == nil {
		fmt.Println("Hyperdrive couldn't check your node wallet status yet. Check on it again later with `hyperdrive wallet status`. If you haven't madea wallet yet, you can do so now with `hyperdrive wallet init`.")
		return nil
	}

	// All set
	if status.Wallet.IsLoaded {
		fmt.Printf("Your node wallet with address %s%s%s is loaded and ready to use.\n", terminal.ColorBlue, status.Wallet.WalletAddress.Hex(), terminal.ColorReset)
		return nil
	}

	// Prompt for password
	if status.Wallet.IsOnDisk {
		return promptForPassword(c, hd)
	}

	// Init
	fmt.Println("You don't have a node wallet yet.")
	if c.Bool(utils.YesFlag.Name) || !utils.Confirm("Would you like to create one now?") {
		fmt.Println("Please create one using `hyperdrive wallet init` when you're ready.")
		return nil
	}
	err = wallet.InitWallet(c, hd)
	if err != nil {
		return fmt.Errorf("error initializing node wallet: %w", err)
	}

	// Get the wallet status
	return nil
}

// Prompt for the wallet password upon startup if it isn't available, but a wallet is on disk
func promptForPassword(c *cli.Context, hd *client.HyperdriveClient) error {
	fmt.Println("Your node wallet is saved, but the password is not stored on disk so it cannot be loaded automatically.")
	// Get the password
	passwordString := c.String(wallet.PasswordFlag.Name)
	if passwordString == "" {
		passwordString = wallet.PromptExistingPassword()
	}
	password, err := input.ValidateNodePassword("password", passwordString)
	if err != nil {
		return fmt.Errorf("error validating password: %w", err)
	}

	// Get the save flag
	savePassword := c.Bool(wallet.SavePasswordFlag.Name) || utils.Confirm("Would you like to save the password to disk? If you do, your node will be able to handle transactions automatically after a client restart; otherwise, you will have to repeat this command to manually enter the password after each restart.")

	// Run it
	_, err = hd.Api.Wallet.SetPassword(password, savePassword)
	if err != nil {
		fmt.Printf("%sError setting password: %s%s\n", terminal.ColorYellow, err.Error(), terminal.ColorReset)
		fmt.Println("Your service has started, but you'll need to provide the node wallet password later with `hyperdrive wallet set-password`.")
		return nil
	}

	// Refresh the status
	response, err := hd.Api.Wallet.Status()
	if err != nil {
		fmt.Printf("Wallet password set.\n%sError checking node wallet: %s%s\n", terminal.ColorYellow, err.Error(), terminal.ColorReset)
		fmt.Println("Please check the service logs with `hyperdrive service logs daemon` for more information.")
		return nil
	}
	status := response.Data.WalletStatus
	if !status.Wallet.IsLoaded {
		fmt.Println("Wallet password set, but the node wallet could not be loaded.")
		fmt.Println("Please check the service logs with `hyperdrive service logs daemon` for more information.")
		return nil
	}
	fmt.Printf("Your node wallet with address %s%s%s is now loaded and ready to use.\n", terminal.ColorBlue, status.Wallet.WalletAddress.Hex(), terminal.ColorReset)
	return nil
}

// Check if any of the VCs has changed and force a wait for slashing protection, since all VCs are tied to the BN selection
func checkForValidatorChange(hd *client.HyperdriveClient, cfg *config.HyperdriveConfig) (bool, error) {
	// Get all of the VCs belonging to the project
	prefix := cfg.ProjectName.Value
	vcs, err := hd.GetValidatorContainers(prefix + "_")
	if err != nil {
		return false, fmt.Errorf("error getting validator client containers: %w", err)
	}

	// Break if there aren't any
	if len(vcs) == 0 {
		return true, nil
	}

	/*
		// TODO: DEBUG
			fmt.Println("Found the following Validator Clients:")
			for _, vc := range vcs {
				fmt.Println(vc)
			}
			fmt.Println()
	*/

	// Get the map of VCs to their new tags in the config
	newTagMap, err := getVcContainerTagParamMap(cfg, vcs)
	if err != nil {
		return false, err
	}

	// Get the list of any VCs that can't be safely started yet
	longestRemainingTime := time.Duration(0)
	for _, vc := range vcs {
		remainingTime, err := checkValidatorClient(hd, vc, newTagMap)
		if err != nil {
			return false, err
		}

		// If this VC has remaining time before it can be safely started, see if it's more than the current max
		if remainingTime > longestRemainingTime {
			longestRemainingTime = remainingTime
		}
	}

	// Show the slashing prevention dialog
	if longestRemainingTime > 0 {
		showSlashingDelay(longestRemainingTime)
	}
	return false, nil
}

func checkValidatorClient(hd *client.HyperdriveClient, vcName string, newTagMap map[string]string) (time.Duration, error) {
	// Get the current and pending VC images
	currentTag, err := hd.GetDockerImage(vcName)
	if err != nil {
		return 0, fmt.Errorf("error getting Docker image tag for [%s]: %w", vcName, err)
	}
	currentVcType, err := getDockerImageName(hd, currentTag)
	if err != nil {
		return 0, fmt.Errorf("error parsing current Docker image tag [%s] for [%s]: %w", currentTag, vcName, err)
	}
	pendingTag := newTagMap[vcName]
	pendingVcType, err := getDockerImageName(hd, pendingTag)
	if err != nil {
		return 0, fmt.Errorf("error parsing pending Docker image tag [%s] for [%s]: %w", pendingTag, vcName, err)
	}

	// Compare the clients and warn if necessary
	if currentVcType == pendingVcType {
		fmt.Printf("Validator Client [%s] is still [%s] - no slashing prevention delay necessary.\n", vcName, currentVcType)
		return 0, nil
	} else {
		validatorFinishTime, err := hd.GetDockerContainerShutdownTime(vcName)
		if err != nil {
			return 0, fmt.Errorf("error getting VC [%s] shutdown time: %w", vcName, err)
		}

		// If it hasn't exited yet, shut it down
		zeroTime := time.Time{}
		status, err := hd.GetDockerStatus(vcName)
		if err != nil {
			return 0, fmt.Errorf("error getting VC [%s] status: %w", vcName, err)
		}
		if validatorFinishTime == zeroTime || status == "running" {
			fmt.Printf("%sValidator Client [%s] is currently running, stopping it...%s\n", terminal.ColorYellow, vcName, terminal.ColorReset)
			err := hd.StopContainer(vcName)
			if err != nil {
				return 0, fmt.Errorf("error stopping VC [%s]: %w", vcName, err)
			}
			validatorFinishTime = time.Now()
		}

		// Print the warning and start the time lockout
		safeStartTime := validatorFinishTime.Add(15 * time.Minute)
		remainingTime := time.Until(safeStartTime)
		if remainingTime <= 0 {
			fmt.Printf("Validator Client [%s] has been offline for %s, which is long enough to prevent slashing.\n", vcName, time.Since(validatorFinishTime))
			return 0, nil
		}

		// If this VC has remaining time before it can be safely started, add it to the list
		if remainingTime > 0 {
			fmt.Printf("Validator Client [%s] has changed types from [%s] to [%s].\n", vcName, currentVcType, pendingVcType)
			fmt.Printf("Only %s has elapsed since you stopped it.\n", time.Since(validatorFinishTime))
		}

		// This can't be safely started, return its info
		return remainingTime, nil
	}
}

func showSlashingDelay(remainingTime time.Duration) {
	fmt.Printf("%s=== WARNING ===\n", terminal.ColorRed)
	fmt.Println("You have changed validator clients. You must wait at least 15 minutes before safely starting them to prevent attesting to the same block twice, which would result in slashing your ETH.")
	fmt.Println("To prevent slashing, Hyperdrive will delay activating the new client until it is safe.")
	fmt.Println("See the documentation for a more detailed explanation: https://docs.nodeset.io")
	fmt.Printf("If you have read the documentation, understand the risks, and want to bypass this cooldown, run `hyperdrive service start --%s`.%s\n\n", ignoreSlashTimerFlag.Name, terminal.ColorReset)

	// Wait for 15 minutes
	safeStartTime := time.Now().Add(remainingTime)
	for remainingTime > 0 {
		fmt.Printf("Remaining time: %s", remainingTime)
		time.Sleep(1 * time.Second)
		remainingTime = time.Until(safeStartTime)
		fmt.Printf("%s\r", terminal.ClearLine)
	}

	fmt.Println(terminal.ColorReset)
	fmt.Println("You may now safely start Hyperdrive without fear of being slashed.")
}

// !!!!
// THIS IS A TOTAL HACK DESIGNED TO PROTECT YOU FROM FORGETTING TO ADD A MODULE, DO THIS BETTER WHEN MODULES ARE PROPERLY SPLIT OUT
// !!!!
func getVcContainerTagParamMap(cfg *config.HyperdriveConfig, vcs []string) (map[string]string, error) {
	containerTagMap := map[string]string{}

	// Get the VC names
	stakewiseVc := cfg.GetDockerArtifactName(string(swconfig.ContainerID_StakewiseValidator))

	// Get the BN being used
	var bn types.BeaconNode
	if cfg.IsLocalMode() {
		bn = cfg.LocalBeaconConfig.BeaconNode.Value
	} else {
		bn = cfg.ExternalBeaconConfig.BeaconNode.Value
	}

	switch bn {
	case types.BeaconNode_Lighthouse:
		containerTagMap[stakewiseVc] = cfg.Modules.Stakewise.Lighthouse.ContainerTag.Value
	case types.BeaconNode_Lodestar:
		containerTagMap[stakewiseVc] = cfg.Modules.Stakewise.Lodestar.ContainerTag.Value
	case types.BeaconNode_Nimbus:
		containerTagMap[stakewiseVc] = cfg.Modules.Stakewise.Nimbus.ContainerTag.Value
	case types.BeaconNode_Prysm:
		containerTagMap[stakewiseVc] = cfg.Modules.Stakewise.Prysm.ContainerTag.Value
	case types.BeaconNode_Teku:
		containerTagMap[stakewiseVc] = cfg.Modules.Stakewise.Teku.ContainerTag.Value
	default:
		panic(fmt.Sprintf("unknown Beacon Node type: [%v]", bn))
	}

	// SANITY CHECK
	for _, vc := range vcs {
		_, exists := containerTagMap[vc]
		if !exists {
			return nil, fmt.Errorf("validator client [%s] was missing from the slashing prevention check", vc)
		}
	}

	return containerTagMap, nil
}

// Extract the image name from a Docker image string
func getDockerImageName(hd *client.HyperdriveClient, image string) (string, error) {
	// Return the empty string if the validator didn't exist (probably because this is the first time starting it up)
	if image == "" {
		return "", nil
	}

	reg := regexp.MustCompile(dockerImageRegex)
	matches := reg.FindStringSubmatch(image)
	if matches == nil {
		return "", fmt.Errorf("Couldn't parse the Docker image string [%s]", image)
	}
	imageIndex := reg.SubexpIndex("image")
	if imageIndex == -1 {
		return "", fmt.Errorf("Image name not found in Docker image [%s]", image)
	}

	imageName := matches[imageIndex]
	return imageName, nil
}
