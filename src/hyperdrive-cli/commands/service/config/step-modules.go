package config

import "fmt"

func createModulesStep(wiz *wizard, currentStep int, totalSteps int) *checkBoxWizardStep {
	// Create the labels
	fmt.Printf("!!!createModulesStep")
	stakewiseCfg := wiz.md.Config.Stakewise
	stakewiseLabel := stakewiseCfg.GetTitle()
	fmt.Printf("!!! stakewiseLabel: %s\n", stakewiseLabel)
	constellationCfg := wiz.md.Config.Constellation
	constellationLabel := constellationCfg.GetTitle()

	fmt.Printf("!!! constellationLabel: %s\n", constellationLabel)
	helperText := "Select the NodeSet modules you would like to enable below."

	show := func(modal *checkBoxModalLayout) {
		fmt.Printf("!!! inside show")
		wiz.md.setPage(modal.page)
		fmt.Printf("!!! setPage")
		modal.focus()
		for label, box := range modal.checkboxes {
			fmt.Printf("!!! label: %s, box: %v\n", label, box)
			switch label {
			case stakewiseLabel:
				box.SetChecked(wiz.md.Config.Stakewise.Enabled.Value)
			case constellationLabel:
				box.SetChecked(wiz.md.Config.Constellation.Enabled.Value)
			}
		}
	}
	fmt.Printf("!!! show")
	done := func(choices map[string]bool) {
		fmt.Printf("!!! inside done")
		stakewiseCfg.Enabled.Value = choices[stakewiseLabel]
		constellationCfg.Enabled.Value = choices[constellationLabel]
		wiz.metricsModal.show()
	}
	fmt.Printf("!!! back")
	back := func() {
		wiz.useFallbackModal.show()
	}
	fmt.Printf("!!! return")
	return newCheckBoxStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		90,
		"Modules",
		[]string{stakewiseLabel, constellationLabel},
		[]string{stakewiseCfg.Enabled.Description, constellationCfg.Enabled.Description},
		[]bool{stakewiseCfg.Enabled.Value, constellationCfg.Enabled.Value},
		show,
		done,
		back,
		"step-modules-local",
	)
}
