package config

import (
	"fmt"
	"strings"

	"github.com/rocket-pool/node-manager-core/config"
)

const (
	externalRethWarning string = "[orange]WARNING: Ensure that your external Reth client is configured to preserve event logs from **all contracts** (at least as far back as the Beacon deposit contract's deployment), not just the deposit contract logs! Event logs are required for many modules, and if your Reth has pruned them, the modules will not work."
)

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
