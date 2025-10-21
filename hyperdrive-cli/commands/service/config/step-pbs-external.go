package config

import "github.com/rocket-pool/node-manager-core/config"

func createPbsExternalStep(wiz *wizard, currentStep int, totalSteps int) *textBoxWizardStep {
	// Create the labels
	urlLabel := wiz.md.Config.Hyperdrive.Pbs.ExternalPbsClientConfig.ExternalUrl.Name

	helperText := "Please enter the URL of your external PBS client.\n\nFor example: `http://192.168.1.46:18550`"

	show := func(modal *textBoxModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus()
		for label, box := range modal.textboxes {
			for _, param := range wiz.md.Config.Hyperdrive.Pbs.ExternalPbsClientConfig.GetParameters() {
				if param.GetCommon().Name == label {
					box.SetText(param.String())
				}
			}
		}
	}

	done := func(text map[string]string) {
		wiz.md.Config.Hyperdrive.Pbs.Enable.Value = true
		wiz.md.Config.Hyperdrive.Pbs.Mode.Value = config.ClientMode_External
		wiz.md.Config.Hyperdrive.Pbs.ExternalPbsClientConfig.ExternalUrl.Value = text[urlLabel]
		wiz.finishedModal.show()
	}

	back := func() {
		wiz.pbsModeModal.show()
	}

	return newTextBoxWizardStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		70,
		"PBS Client (External)",
		[]string{urlLabel},
		[]int{wiz.md.Config.Hyperdrive.Pbs.ExternalPbsClientConfig.ExternalUrl.MaxLength},
		[]string{wiz.md.Config.Hyperdrive.Pbs.ExternalPbsClientConfig.ExternalUrl.Regex},
		show,
		done,
		back,
		"step-external-pbs",
	)
}
