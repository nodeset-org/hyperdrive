package config

func createModulesStep(wiz *wizard, currentStep int, totalSteps int) *checkBoxWizardStep {
	// Create the labels
	stakewiseCfg := wiz.md.Config.Stakewise
	stakewiseLabel := stakewiseCfg.GetTitle()

	helperText := "Select the NodeSet modules you would like to enable below."

	show := func(modal *checkBoxModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus()
		for label, box := range modal.checkboxes {
			switch label {
			case stakewiseLabel:
				box.SetChecked(wiz.md.Config.Stakewise.Enabled.Value)
			}
		}
	}

	done := func(choices map[string]bool) {
		stakewiseCfg.Enabled.Value = choices[stakewiseLabel]
		wiz.metricsModal.show()
	}

	back := func() {
		wiz.useFallbackModal.show()
	}

	return newCheckBoxStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		90,
		"Modules",
		[]string{stakewiseLabel},
		[]string{stakewiseCfg.Enabled.Description},
		[]bool{stakewiseCfg.Enabled.Value},
		show,
		done,
		back,
		"step-modules-local",
	)
}
