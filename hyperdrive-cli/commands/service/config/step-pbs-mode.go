package config

import (
	"github.com/nodeset-org/hyperdrive-daemon/shared/config/pbs"
	"github.com/rocket-pool/node-manager-core/config"
)

func createPbsModeStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	// Create the button names and descriptions from the config
	options := wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.Client.Options
	names := []string{}
	descriptions := []string{}
	for _, option := range options {
		names = append(names, option.Name)
		descriptions = append(descriptions, option.Description)
	}
	names = append(names, "Externally Managed", "Disable PBS Support")
	descriptions = append(descriptions,
		"Connect to an external PBS client that you manage yourself.",
		"Disable PBS support. When your validators propose a block, they will use a block that your clients create by themselves.")

	helperText := "Hyperdrive supports PBS clients such as MEV-Boost and Commit-Boost, which allow you to capture extra profits from your validator's block proposals. If you'd like to use it, Hyperdrive can either manage a client for you or connect to an instance you manage yourself. How would you like to use PBS?\n\n[lime]To learn more about PBS and MEV, please visit:\nhttps://docs.flashbots.net/new-to-mev\n"

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0) // Catch-all for safety

		if !wiz.md.Config.Hyperdrive.Pbs.Enable.Value {
			modal.focus(len(names) - 1)
			return
		}

		switch wiz.md.Config.Hyperdrive.Pbs.Mode.Value {
		case config.ClientMode_Local:
			for i, option := range options {
				if option.Value == pbs.PbsClient(wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.Client.Value) {
					modal.focus(i)
				}
			}
		case config.ClientMode_External:
			modal.focus(len(names) - 2)
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		switch buttonIndex {
		case len(names) - 1:
			// Disable PBS support
			wiz.md.Config.Hyperdrive.Pbs.Enable.Value = false
			wiz.finishedModal.show()
		case len(names) - 2:
			// Externally Managed
			wiz.md.Config.Hyperdrive.Pbs.Enable.Value = true
			wiz.md.Config.Hyperdrive.Pbs.Mode.Value = config.ClientMode_External
			wiz.externalPbsModal.show()
		default:
			// Local PBS Client
			wiz.md.Config.Hyperdrive.Pbs.Enable.Value = true
			wiz.md.Config.Hyperdrive.Pbs.Mode.Value = config.ClientMode_Local
			wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.Client.Value = options[buttonIndex].Value
			if wiz.md.Config.Hyperdrive.ClientMode.Value == config.ClientMode_External {
				wiz.pbsWarningModal.show()
			} else {
				wiz.localPbsRelayModeModal.show()
			}
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
		names,
		descriptions,
		76,
		"PBS Mode",
		DirectionalModalVertical,
		show,
		done,
		back,
		"step-pbs-mode",
	)
}
