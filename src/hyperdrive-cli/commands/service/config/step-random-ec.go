package config

import (
	"fmt"

	nmc_config "github.com/rocket-pool/node-manager-core/config"
)

const randomEcID string = "step-random-ec"

func createRandomEcStep(wiz *wizard, currentStep int, totalSteps int, goodOptions []*nmc_config.ParameterOption[nmc_config.ExecutionClient]) *choiceWizardStep {
	var selectedClientName string
	selectedClient := wiz.md.Config.Hyperdrive.LocalExecutionConfig.ExecutionClient.Value
	for _, clientOption := range goodOptions {
		if clientOption.Value == selectedClient {
			selectedClientName = clientOption.Name
			break
		}
	}

	helperText := fmt.Sprintf("You have been randomly assigned to %s for your Execution client.", selectedClientName)

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0)
	}

	done := func(buttonIndex int, buttonLabel string) {
		wiz.bnLocalModal.show()
	}

	back := func() {
		wiz.ecLocalModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		[]string{"Ok"},
		[]string{},
		76,
		"Execution Client > Selection",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		randomEcID,
	)
}
