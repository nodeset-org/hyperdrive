package config

const mevDisabledID string = "step-mev-disabled"

func createMevDisabledStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	helperText := mevDisabled

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0)
	}

	done := func(buttonIndex int, buttonLabel string) {
		wiz.md.Config.Hyperdrive.MevBoost.Enable.Value = false
		wiz.finishedModal.show()
	}

	back := func() {
		wiz.metricsModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		[]string{"Continue"},
		[]string{},
		76,
		"MEV-Boost Mode",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		mevDisabledID,
	)
}
