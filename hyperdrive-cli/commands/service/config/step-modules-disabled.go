package config

const modulesDisabledID string = "step-modules-disabled"

func createModulesDisabledStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0)
	}

	done := func(buttonIndex int, buttonLabel string) {
		wiz.md.Config.StakeWise.Enabled.Value = false
		wiz.md.Config.Constellation.Enabled.Value = false
		wiz.metricsModal.show()
	}

	back := func() {
		wiz.useFallbackModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		modulesDisabled,
		[]string{"Continue"},
		[]string{},
		76,
		"Modules",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		modulesDisabledID,
	)
}
