package config

func createExternalEcSelectStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	// Create the button names and descriptions from the config
	clients := wiz.md.Config.Hyperdrive.ExternalExecutionClient.ExecutionClient.Options
	clientNames := []string{}
	for _, client := range clients {
		clientNames = append(clientNames, client.Name)
	}

	helperText := "Which Execution Client are you externally managing? Each of them has small behavioral differences, so we'll need to know which one you're using in order to connect to it properly."

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0) // Catch-all for safety

		for i, option := range wiz.md.Config.Hyperdrive.ExternalExecutionClient.ExecutionClient.Options {
			if option.Value == wiz.md.Config.Hyperdrive.ExternalExecutionClient.ExecutionClient.Value {
				modal.focus(i)
				break
			}
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		selectedClient := clients[buttonIndex].Value
		wiz.md.Config.Hyperdrive.ExternalExecutionClient.ExecutionClient.Value = selectedClient
		wiz.externalEcSettingsModal.show()
	}

	back := func() {
		wiz.modeModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		clientNames,
		[]string{},
		70,
		"Execution Client (External) > Selection",
		DirectionalModalVertical,
		show,
		done,
		back,
		"step-external-ec-select",
	)
}
