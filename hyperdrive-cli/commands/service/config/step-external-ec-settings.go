package config

func createExternalEcSettingsStep(wiz *wizard, currentStep int, totalSteps int) *textBoxWizardStep {
	// Create the labels
	httpLabel := wiz.md.Config.ExternalExecutionConfig.HttpUrl.Name
	wsLabel := wiz.md.Config.ExternalExecutionConfig.WebsocketUrl.Name

	helperText := "Please enter the URL of the HTTP-based RPC API and the URL of the Websocket-based RPC API for your existing client.\n\nFor example: `http://192.168.1.45:8545` and `ws://192.168.1.45:8546`"

	show := func(modal *textBoxModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus()
		for label, box := range modal.textboxes {
			for _, param := range wiz.md.Config.ExternalExecutionConfig.GetParameters() {
				if param.GetCommon().Name == label {
					box.SetText(param.GetValueAsString())
				}
			}
		}
	}

	done := func(text map[string]string) {
		wiz.md.Config.ExternalExecutionConfig.HttpUrl.Value = text[httpLabel]
		wiz.md.Config.ExternalExecutionConfig.WebsocketUrl.Value = text[wsLabel]
		wiz.externalBnSelectModal.show()
	}

	back := func() {
		wiz.externalEcSelectModal.show()
	}

	return newTextBoxWizardStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		70,
		"Execution Client (External) > Settings",
		[]string{httpLabel, wsLabel},
		[]int{wiz.md.Config.ExternalExecutionConfig.HttpUrl.MaxLength, wiz.md.Config.ExternalExecutionConfig.WebsocketUrl.MaxLength},
		[]string{wiz.md.Config.ExternalExecutionConfig.HttpUrl.Regex, wiz.md.Config.ExternalExecutionConfig.WebsocketUrl.Regex},
		show,
		done,
		back,
		"step-ec-external-settings",
	)
}
