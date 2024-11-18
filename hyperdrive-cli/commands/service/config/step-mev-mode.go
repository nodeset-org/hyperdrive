package config

import (
	"fmt"

	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
)

func createMevModeStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	// Create the button names and descriptions from the config
	modeNames := []string{
		"Use All Relays",
		"Choose Relays",
		"Externally Managed",
		"Disable MEV-Boost",
	}
	modeDescriptions := []string{
		"Allow Hyperdrive to manage MEV-Boost for you, and subscribe to all of the built-in MEV relays. You can view the built-in relays during the manual configuration process once this wizard is finished, and add your own custom relays as well.",
		"Allow Hyperdrive to manage MEV-Boost for you, but manually select which built-in relays to use. You can add your own custom relays as well.",
		"Connect to an external MEV-Boost client that you manage yourself.",
		"Disable MEV-Boost support. When your validators propose a block, they will use a block that your clients create by themselves.",
	}

	helperText := "Hyperdrive supports MEV-Boost, which allows you to capture extra profits from your validator's block proposals. If you'd like to use it, Hyperdrive can either manage MEV-Boost for you or connect to an instance you manage yourself. How would you like to use MEV-Boost?\n\n[lime]To learn more about MEV, please visit:\nhttps://docs.flashbots.net/new-to-mev\n"

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0) // Catch-all for safety

		if !wiz.md.Config.Hyperdrive.MevBoost.Enable.Value {
			modal.focus(3)
			return
		}

		switch wiz.md.Config.Hyperdrive.MevBoost.Mode.Value {
		case config.ClientMode_Local:
			switch wiz.md.Config.Hyperdrive.MevBoost.SelectionMode.Value {
			case hdconfig.MevSelectionMode_All:
				modal.focus(0)
			case hdconfig.MevSelectionMode_Manual:
				modal.focus(1)
			}
		case config.ClientMode_External:
			modal.focus(2)
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		switch buttonIndex {
		case 0:
			wiz.md.Config.Hyperdrive.MevBoost.Enable.Value = true
			wiz.md.Config.Hyperdrive.MevBoost.Mode.Value = config.ClientMode_Local
			wiz.md.Config.Hyperdrive.MevBoost.SelectionMode.Value = hdconfig.MevSelectionMode_All
			if wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_External {
				wiz.mevWarningModal.show()
			} else {
				wiz.finishedModal.show()
			}
		case 1:
			wiz.md.Config.Hyperdrive.MevBoost.Enable.Value = true
			wiz.md.Config.Hyperdrive.MevBoost.Mode.Value = config.ClientMode_Local
			wiz.md.Config.Hyperdrive.MevBoost.SelectionMode.Value = hdconfig.MevSelectionMode_Manual
			if wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_External {
				wiz.mevWarningModal.show()
			} else {
				wiz.localMevModal.show()
			}
		case 2:
			wiz.md.Config.Hyperdrive.MevBoost.Enable.Value = true
			wiz.md.Config.Hyperdrive.MevBoost.Mode.Value = config.ClientMode_External
			if wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_External {
				wiz.mevWarningModal.show()
			} else {
				wiz.externalMevModal.show()
			}
		case 3:
			wiz.md.Config.Hyperdrive.MevBoost.Enable.Value = false
			wiz.finishedModal.show()
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
		"MEV-Boost Mode",
		DirectionalModalVertical,
		show,
		done,
		back,
		"step-mev-mode",
	)
}
