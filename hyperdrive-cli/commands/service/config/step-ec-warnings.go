package config

import (
	"fmt"
	"strings"

	"github.com/rocket-pool/node-manager-core/config"
)

const (
	localRethWarning    string = "[orange]WARNING: Reth is still in beta and has been shown to have some incompatibilities with the StakeWise Operator service, preventing you from submitting validator deposits properly. We strongly recommend you pick a different client instead until the incompatibility is fixed."
	externalRethWarning string = "[orange]WARNING: Reth is still in beta and has been shown to have some incompatibilities with the StakeWise Operator service, preventing you from submitting validator deposits properly. We strongly recommend avoiding it as your external client until the incompatibility is fixed."
)

// Get a more verbose client description, including warnings
func getAugmentedLocalEcDescription(client config.ExecutionClient, originalDescription string) string {
	switch client {
	case config.ExecutionClient_Reth:
		if !strings.HasSuffix(originalDescription, localRethWarning) {
			return fmt.Sprintf("%s\n\n%s", originalDescription, localRethWarning)
		}
	}

	return originalDescription
}

// Get a more verbose client description, including warnings
func getAugmentedExternalEcDescription(client config.ExecutionClient, originalDescription string) string {
	switch client {
	case config.ExecutionClient_Reth:
		if !strings.HasSuffix(originalDescription, externalRethWarning) {
			return fmt.Sprintf("%s\n\n%s", originalDescription, externalRethWarning)
		}
	}

	return originalDescription
}

func createLocalRethWarningStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	helperText := localRethWarning

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0)
	}

	done := func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 0 {
			wiz.localEcModal.show()
		} else {
			wiz.localBnModal.show()
		}
	}

	back := func() {
		wiz.localEcModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		[]string{"Choose Again", "Keep Reth"},
		[]string{},
		76,
		"Execution Client > Selection",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		"step-local-reth-warning",
	)
}

func createExternalRethWarningStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	helperText := externalRethWarning

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0)
	}

	done := func(buttonIndex int, buttonLabel string) {
		wiz.externalEcSettingsModal.show()
	}

	back := func() {
		wiz.externalEcSelectModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		[]string{"I Understand"},
		[]string{},
		76,
		"Execution Client (External) > Selection",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		"step-external-reth-warning",
	)
}
