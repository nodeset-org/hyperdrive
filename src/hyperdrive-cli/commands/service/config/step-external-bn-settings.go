package config

func createExternalBnSettingsStep(wiz *wizard, currentStep int, totalSteps int) *textBoxWizardStep {
	// Create the labels
	httpUrlLabel := wiz.md.Config.Hyperdrive.ExternalBeaconConfig.HttpUrl.Name

	helperText := "Please provide the URL of your Beacon Node's client's HTTP API (for example: `http://192.168.1.40:5052`).\n\nNote that if you're running it on the same machine as Hyperdrive, you cannot use `localhost` or `127.0.0.1`; you must use your machine's LAN IP address."

	show := func(modal *textBoxModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus()
		for label, box := range modal.textboxes {
			for _, param := range wiz.md.Config.Hyperdrive.ExternalBeaconConfig.GetParameters() {
				if param.GetCommon().Name == label {
					box.SetText(param.String())
				}
			}
		}
	}

	done := func(text map[string]string) {
		wiz.md.Config.Hyperdrive.ExternalBeaconConfig.HttpUrl.Value = text[httpUrlLabel]
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
		[]string{httpUrlLabel},
		[]int{wiz.md.Config.Hyperdrive.ExternalBeaconConfig.HttpUrl.MaxLength},
		[]string{wiz.md.Config.Hyperdrive.ExternalBeaconConfig.HttpUrl.Regex},
		show,
		done,
		back,
		"step-external-bn-settings",
	)
}
