package config

import (
	"github.com/nodeset-org/hyperdrive-daemon/shared/config/pbs"
	"github.com/rocket-pool/node-manager-core/config"
)

const pbsWarningID string = "step-pbs-warning"

func createPbsWarningStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	helperText := mevWarning

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0)
	}

	done := func(buttonIndex int, buttonLabel string) {
		if wiz.md.Config.Hyperdrive.Pbs.Mode.Value == config.ClientMode_Local {
			if wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.RelaySelectionMode.Value == pbs.PbsRelaySelectionMode_All {
				wiz.finishedModal.show()
			} else {
				wiz.localPbsRelaySelectionModal.show()
			}
		} else {
			if wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_Local {
				wiz.externalPbsModal.show()
			} else {
				wiz.finishedModal.show()
			}
		}
	}

	back := func() {
		wiz.pbsModeModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		[]string{"Continue"},
		[]string{},
		76,
		"PBS Client Mode",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		pbsWarningID,
	)
}
