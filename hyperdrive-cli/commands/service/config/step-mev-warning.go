package config

import (
	hdconfig "github.com/nodeset-org/hyperdrive-daemon/shared/config"
	"github.com/rocket-pool/node-manager-core/config"
)

const mevWarningID string = "step-mev-warning"

func createMevWarningStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	helperText := mevWarning

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0)
	}

	done := func(buttonIndex int, buttonLabel string) {
		if wiz.md.Config.Hyperdrive.MevBoost.Mode.Value == config.ClientMode_Local {
			if wiz.md.Config.Hyperdrive.MevBoost.SelectionMode.Value == hdconfig.MevSelectionMode_All {
				wiz.finishedModal.show()
			} else {
				wiz.localMevModal.show()
			}
		} else {
			if wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_Local {
				wiz.externalMevModal.show()
			} else {
				wiz.finishedModal.show()
			}
		}
	}

	back := func() {
		wiz.mevModeModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		[]string{"Continue"},
		[]string{},
		76,
		"MEV-Boost Mode",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		mevWarningID,
	)
}
