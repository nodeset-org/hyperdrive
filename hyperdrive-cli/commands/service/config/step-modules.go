package config

func createModulesStep(wiz *wizard, currentStep int, totalSteps int) *checkBoxWizardStep {
	// Create the labels
	stakewiseCfg := wiz.md.Config.StakeWise
	stakewiseLabel := stakewiseCfg.GetTitle()
	constellationCfg := wiz.md.Config.Constellation
	constellationLabel := constellationCfg.GetTitle()

	helperText := "Select the NodeSet modules you would like to enable below."

	show := func(modal *checkBoxModalLayout) {
		modal.generateCheckboxes(
			[]string{stakewiseLabel, constellationLabel},
			[]string{stakewiseCfg.Enabled.Description, constellationCfg.Enabled.Description},
			[]bool{stakewiseCfg.Enabled.Value, constellationCfg.Enabled.Value},
		)

		wiz.md.setPage(modal.page)
		modal.focus()
		for label, box := range modal.checkboxes {
			switch label {
			case stakewiseLabel:
				box.SetChecked(wiz.md.Config.StakeWise.Enabled.Value)
			case constellationLabel:
				box.SetChecked(wiz.md.Config.Constellation.Enabled.Value)
			}
		}
	}

	done := func(choices map[string]bool) {
		stakewiseCfg.Enabled.Value = choices[stakewiseLabel]
		constellationCfg.Enabled.Value = choices[constellationLabel]
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
		show,
		done,
		back,
		"step-modules-local",
	)
}
