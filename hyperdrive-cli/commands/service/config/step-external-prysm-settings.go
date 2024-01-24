package config

func createExternalPrysmSettingsStep(wiz *wizard, currentStep int, totalSteps int) *textBoxWizardStep {
	// Create the labels
	httpUrlLabel := wiz.md.Config.ExternalBeaconConfig.HttpUrl.Name
	jsonRpcUrlLabel := wiz.md.Config.ExternalBeaconConfig.PrysmRpcUrl.Name

	helperText := "Please provide the URL of your Prysm client's HTTP API (for example: `http://192.168.1.40:5052`) and the URL of its JSON RPC API (e.g., `192.168.1.40:5053`) too.\n\nNote that if you're running it on the same machine as the Smartnode, you cannot use `localhost` or `127.0.0.1`; you must use your machine's LAN IP address."

	show := func(modal *textBoxModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus()
		for label, box := range modal.textboxes {
			for _, param := range wiz.md.Config.ExternalBeaconConfig.GetParameters() {
				if param.GetCommon().Name == label {
					box.SetText(param.GetValueAsString())
				}
			}
		}
	}

	done := func(text map[string]string) {
		wiz.md.Config.ExternalBeaconConfig.HttpUrl.Value = text[httpUrlLabel]
		wiz.md.Config.ExternalBeaconConfig.PrysmRpcUrl.Value = text[jsonRpcUrlLabel]
		wiz.useFallbackModal.show()
	}

	back := func() {
		wiz.externalBnSelectModal.show()
	}

	return newTextBoxWizardStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		70,
		"Beacon Node (External) > Settings",
		[]string{httpUrlLabel, jsonRpcUrlLabel},
		[]int{wiz.md.Config.ExternalBeaconConfig.HttpUrl.MaxLength, wiz.md.Config.ExternalBeaconConfig.PrysmRpcUrl.MaxLength},
		[]string{wiz.md.Config.ExternalBeaconConfig.HttpUrl.Regex, wiz.md.Config.ExternalBeaconConfig.PrysmRpcUrl.Regex},
		show,
		done,
		back,
		"step-external-prysm-settings",
	)
}
