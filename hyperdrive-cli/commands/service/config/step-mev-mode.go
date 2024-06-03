package config

import (
	"fmt"

	"github.com/rocket-pool/node-manager-core/config"
)

func createMevModeStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	// Create the button names and descriptions from the config
	modes := wiz.md.Config.Hyperdrive.MevBoost.Mode.Options
	modeNames := []string{}
	modeDescriptions := []string{}
	for _, mode := range modes {
		modeNames = append(modeNames, mode.Name)
		modeDescriptions = append(modeDescriptions, mode.Description)
	}

	helperText := "By default, Hyperdrive has MEV-Boost enabled. This allows you to capture extra profits from block proposals. Would you like Hyperdrive to manage MEV-Boost for you, or would you like to manage it yourself?\n\n[lime]To learn more about MEV, please visit:\nhttps://docs.flashbots.net/new-to-mev\n"

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0) // Catch-all for safety

		for i, option := range wiz.md.Config.Hyperdrive.MevBoost.Mode.Options {
			if option.Value == wiz.md.Config.Hyperdrive.MevBoost.Mode.Value {
				modal.focus(i)
				break
			}
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		wiz.md.Config.Hyperdrive.MevBoost.Mode.Value = modes[buttonIndex].Value
		switch modes[buttonIndex].Value {
		case config.ClientMode_Local:
			wiz.localMevModal.show()
		case config.ClientMode_External:
			if wiz.md.Config.Hyperdrive.IsLocalMode() {
				wiz.externalMevModal.show()
			} else {
				wiz.finishedModal.show()
			}
		default:
			panic(fmt.Sprintf("Unknown MEV mode %s", modes[buttonIndex].Value))
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
		"MEV-Boost Mode",
		DirectionalModalVertical,
		show,
		done,
		back,
		"step-mev-mode",
	)
}
