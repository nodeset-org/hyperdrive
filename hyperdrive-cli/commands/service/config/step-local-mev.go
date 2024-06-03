package config

import (
	"strings"

	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
	nconfig "github.com/rocket-pool/node-manager-core/config"
)

func createLocalMevStep(wiz *wizard, currentStep int, totalSteps int) *checkBoxWizardStep {
	// Create the labels
	regulatedAllLabel := strings.TrimPrefix(wiz.md.Config.Hyperdrive.MevBoost.EnableRegulatedAllMev.Name, "Enable ")
	unregulatedAllLabel := strings.TrimPrefix(wiz.md.Config.Hyperdrive.MevBoost.EnableUnregulatedAllMev.Name, "Enable ")

	helperText := "Select the profiles you would like to enable below. Read the descriptions carefully! Leave all options unchecked if you wish to disable MEV-Boost.\n\n[lime]To learn more about MEV, please visit:\nhttps://docs.flashbots.net/new-to-mev\n"

	show := func(modal *checkBoxModalLayout) {
		labels, descriptions, selections := getMevChoices(wiz.md.Config.Hyperdrive.MevBoost, wiz.md.Config.Hyperdrive.Network.Value)
		modal.generateCheckboxes(labels, descriptions, selections)

		wiz.md.setPage(modal.page)
		modal.focus()
	}

	done := func(choices map[string]bool) {
		wiz.md.Config.Hyperdrive.MevBoost.SelectionMode.Value = config.MevSelectionMode_Profile
		wiz.md.Config.Hyperdrive.MevBoost.Enable.Value = false

		atLeastOneEnabled := false
		enabled, exists := choices[regulatedAllLabel]
		if exists {
			wiz.md.Config.Hyperdrive.MevBoost.EnableRegulatedAllMev.Value = enabled
			atLeastOneEnabled = atLeastOneEnabled || enabled
		}
		enabled, exists = choices[unregulatedAllLabel]
		if exists {
			wiz.md.Config.Hyperdrive.MevBoost.EnableUnregulatedAllMev.Value = enabled
			atLeastOneEnabled = atLeastOneEnabled || enabled
		}

		wiz.md.Config.Hyperdrive.MevBoost.Enable.Value = atLeastOneEnabled
		wiz.finishedModal.show()
	}

	back := func() {
		wiz.mevModeModal.show()
	}

	return newCheckBoxStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		90,
		"MEV-Boost",
		show,
		done,
		back,
		"step-mev-local",
	)
}

func getMevChoices(config *config.MevBoostConfig, network nconfig.Network) ([]string, []string, []bool) {
	labels := []string{}
	descriptions := []string{}
	settings := []bool{}

	regulatedAllMev, unregulatedAllMev := config.GetAvailableProfiles()

	if unregulatedAllMev {
		label := strings.TrimPrefix(config.EnableUnregulatedAllMev.Name, "Enable ")
		labels = append(labels, label)
		description := config.EnableUnregulatedAllMev.DescriptionsByNetwork[network]
		descriptions = append(descriptions, getDescriptionBody(description))
		settings = append(settings, config.EnableUnregulatedAllMev.Value)
	}
	if regulatedAllMev {
		label := strings.TrimPrefix(config.EnableRegulatedAllMev.Name, "Enable ")
		labels = append(labels, label)
		description := config.EnableRegulatedAllMev.DescriptionsByNetwork[network]
		descriptions = append(descriptions, getDescriptionBody(description))
		settings = append(settings, config.EnableRegulatedAllMev.Value)
	}

	return labels, descriptions, settings
}

func getDescriptionBody(description string) string {
	index := strings.Index(description, "Select this")
	return description[index:]
}
