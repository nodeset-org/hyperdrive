package config

const pbsDisabledID string = "step-pbs-disabled"

func createPbsDisabledStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	helperText := pbsDisabled

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0)
	}

	done := func(buttonIndex int, buttonLabel string) {
		wiz.md.Config.Hyperdrive.Pbs.Enable.Value = false
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
		"PBS Client Mode",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		pbsDisabledID,
	)
}
