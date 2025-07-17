package config

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
)

func createPbsLocalRelaySelectionStep(wiz *wizard, currentStep int, totalSteps int) *checkBoxWizardStep {
	helperText := "Select the relays you would like to enable below. Note that all of Hyperdrive's built-in relays support regional sanction lists (such as the US OFAC list) and are compliant with regulations. If you'd like to add your own custom relays, choose \"Review All Settings\" at the end of the wizard and go to the PBS Client section.\n\n[lime]To learn more about MEV, please visit:\nhttps://docs.flashbots.net/new-to-mev\n"

	show := func(modal *checkBoxModalLayout) {
		labels, descriptions, selections := getPbsRelayChoices(wiz.md.Config)
		modal.generateCheckboxes(labels, descriptions, selections)

		wiz.md.setPage(modal.page)
		modal.focus()
	}

	done := func(choices map[string]bool) {
		atLeastOneEnabled := false
		for label, enabled := range choices {
			for _, param := range wiz.md.Config.Hyperdrive.Pbs.LocalPbsClientConfig.GetParameters() {
				if param.GetCommon().Name == "Enable "+label {
					param.SetValue(enabled)
					atLeastOneEnabled = atLeastOneEnabled || enabled
					break
				}
			}
		}

		wiz.md.Config.Hyperdrive.Pbs.Enable.Value = atLeastOneEnabled
		wiz.finishedModal.show()
	}

	back := func() {
		wiz.localPbsRelayModeModal.show()
	}

	return newCheckBoxStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		80,
		"PBS Relay Selection",
		show,
		done,
		back,
		"step-pbs-local-relay-selection",
	)
}

func getPbsRelayChoices(config *client.GlobalConfig) ([]string, []string, []bool) {
	labels := []string{}
	descriptions := []string{}
	settings := []bool{}

	localPbsConfig := config.Hyperdrive.Pbs.LocalPbsClientConfig
	relays := localPbsConfig.GetAvailableRelays(config.Hyperdrive.GetEthNetworkName())
	for _, relay := range relays {
		label := relay.Name
		labels = append(labels, label)
		description := relay.Description
		descriptions = append(descriptions, description)
		for _, parameter := range localPbsConfig.GetParameters() {
			if parameter.GetCommon().Name == "Enable "+relay.Name {
				settings = append(settings, parameter.GetValueAsAny().(bool))
				break
			}
		}
	}

	return labels, descriptions, settings
}
