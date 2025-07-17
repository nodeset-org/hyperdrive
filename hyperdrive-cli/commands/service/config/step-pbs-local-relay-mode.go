package config

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive-daemon/shared/config/pbs"
)

func createPbsLocalRelayModeStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	// Create the button names and descriptions from the config
	modeNames := []string{}
	modeDescriptions := []string{}
	for _, mode := range wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Options {
		modeNames = append(modeNames, mode.Name)
		modeDescriptions = append(modeDescriptions, mode.Description)
	}

	helperText := "Your PBS client will connect to different relays, which submit your block building requests to third-party builders. How would you like to select which relays you want to connect to?"

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0) // Catch-all for safety

		for i, mode := range wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Options {
			if mode.Value == wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Value {
				modal.focus(i)
				return
			}
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		mode := wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Options[buttonIndex].Value
		wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Value = mode
		switch mode {
		case pbs.PbsRelaySelectionMode_All:
			wiz.finishedModal.show()
		case pbs.PbsRelaySelectionMode_Manual:
			wiz.localPbsRelaySelectionModal.show()
		default:
			panic(fmt.Sprintf("Unhandled value %d", buttonIndex))
		}
	}

	back := func() {
		wiz.metricsModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		modeNames,
		modeDescriptions,
		76,
		"PBS Relay Mode",
		DirectionalModalVertical,
		show,
		done,
		back,
		"step-pbs-local-relay-mode",
	)
}
