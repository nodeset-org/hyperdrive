package config

func createFallbackNormalStep(wiz *wizard, currentStep int, totalSteps int) *textBoxWizardStep {
	// Create the labels
	ecHttpLabel := wiz.md.Config.Hyperdrive.Fallback.EcHttpUrl.Name
	ccHttpLabel := wiz.md.Config.Hyperdrive.Fallback.BnHttpUrl.Name

	helperText := "You can use any Execution Client and Beacon Node pair as a fallback.\n\nPlease enter the URLs of the HTTP APIs for your fallback clients.\n\nFor example: `http://192.168.1.45:8545` for your Execution client and `http://192.168.1.45:5052` for your Beacon Node."

	show := func(modal *textBoxModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus()
		for label, box := range modal.textboxes {
			for _, param := range wiz.md.Config.Hyperdrive.Fallback.GetParameters() {
				if param.GetCommon().Name == label {
					box.SetText(param.String())
				}
			}
		}
	}

	done := func(text map[string]string) {
		wiz.md.Config.Hyperdrive.Fallback.EcHttpUrl.Value = text[ecHttpLabel]
		wiz.md.Config.Hyperdrive.Fallback.BnHttpUrl.Value = text[ccHttpLabel]
		wiz.modulesModal.show()
	}

	back := func() {
		wiz.useFallbackModal.show()
	}

	return newTextBoxWizardStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		96,
		"Fallback Client URLs",
		[]string{ecHttpLabel, ccHttpLabel},
		[]int{wiz.md.Config.Hyperdrive.Fallback.EcHttpUrl.MaxLength, wiz.md.Config.Hyperdrive.Fallback.BnHttpUrl.MaxLength},
		[]string{wiz.md.Config.Hyperdrive.Fallback.EcHttpUrl.Regex, wiz.md.Config.Hyperdrive.Fallback.BnHttpUrl.Regex},
		show,
		done,
		back,
		"step-fallback-normal",
	)
}
