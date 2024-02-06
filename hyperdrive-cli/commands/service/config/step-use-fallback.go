package config

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
)

func createUseFallbackStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	helperText := "If you have an extra externally-managed Execution Client and Beacon Node pair that you trust, you can use them as \"fallback\" clients.\nHyperrive and any modules' Validator Clients will connect to these if your primary clients go offline for any reason, so your node will continue functioning properly until they're back online.\n\nWould you like to use a fallback client pair?"

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		if wiz.md.Config.Hyperdrive.Fallback.UseFallbackClients.Value == false {
			modal.focus(0)
		} else {
			modal.focus(1)
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 1 {
			wiz.md.Config.Hyperdrive.Fallback.UseFallbackClients.Value = true
			if (wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_Local && wiz.md.Config.Hyperdrive.LocalBeaconConfig.BeaconNode.Value == config.BeaconNode_Prysm) ||
				(wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_External && wiz.md.Config.Hyperdrive.ExternalBeaconConfig.BeaconNode.Value == config.BeaconNode_Prysm) {
				wiz.fallbackPrysmModal.show()
			} else {
				wiz.fallbackNormalModal.show()
			}
		} else {
			wiz.md.Config.Hyperdrive.Fallback.UseFallbackClients.Value = false
			wiz.modulesModal.show()
		}
	}

	back := func() {
		if wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_Local {
			wiz.checkpointSyncProviderModal.show()
		} else {
			wiz.externalBnSelectModal.show()
		}
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		[]string{"No", "Yes"},
		[]string{},
		76,
		"Use Fallback Clients",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		"step-use-fallback",
	)
}
