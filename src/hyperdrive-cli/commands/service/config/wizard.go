package config

// The wizard display when walking through the general config step-by-step
type wizard struct {
	md *mainDisplay

	// Step 1 - Welcome
	welcomeModal *choiceWizardStep

	// Step 2 - Network
	networkModal *choiceWizardStep

	// Step 3 - Client mode
	modeModal *choiceWizardStep

	// Step 4 - EC settings
	localEcModal            *choiceWizardStep
	localEcRandomModal      *choiceWizardStep
	externalEcSelectModal   *choiceWizardStep
	externalEcSettingsModal *textBoxWizardStep

	// Step 5 - BN settings
	localBnModal                *choiceWizardStep
	localBnRandomModal          *choiceWizardStep
	localBnRandomPrysmModal     *choiceWizardStep
	localBnPrysmWarning         *choiceWizardStep
	localBnTekuWarning          *choiceWizardStep
	checkpointSyncProviderModal *textBoxWizardStep
	externalBnSelectModal       *choiceWizardStep
	externalBnSettingsModal     *textBoxWizardStep
	externalPrysmSettingsModal  *textBoxWizardStep

	// Step 6 - Fallback clients
	useFallbackModal    *choiceWizardStep
	fallbackNormalModal *textBoxWizardStep
	fallbackPrysmModal  *textBoxWizardStep

	// Step 7 - Modules
	modulesModal *checkBoxWizardStep

	// Step 8 - Metrics
	metricsModal *choiceWizardStep

	// Done
	finishedModal *choiceWizardStep
}

// Create a new Wizard display
func newWizard(md *mainDisplay) *wizard {
	wiz := &wizard{
		md: md,
	}

	totalSteps := 9

	// Step 1 - Welcome
	wiz.welcomeModal = createWelcomeStep(wiz, 1, totalSteps)

	// Step 2 - Network
	wiz.networkModal = createNetworkStep(wiz, 2, totalSteps)

	// Step 3 - Client mode
	wiz.modeModal = createModeStep(wiz, 3, totalSteps)

	// Step 4 - EC settings
	wiz.localEcModal = createLocalEcStep(wiz, 4, totalSteps)
	wiz.externalEcSelectModal = createExternalEcSelectStep(wiz, 4, totalSteps)
	wiz.externalEcSettingsModal = createExternalEcSettingsStep(wiz, 4, totalSteps)

	// Step 5 - BN settings
	wiz.localBnModal = createLocalBnStep(wiz, 5, totalSteps)
	wiz.localBnPrysmWarning = createPrysmWarningStep(wiz, 5, totalSteps)
	wiz.localBnTekuWarning = createTekuWarningStep(wiz, 5, totalSteps)
	wiz.checkpointSyncProviderModal = createCheckpointSyncStep(wiz, 5, totalSteps)
	wiz.externalBnSelectModal = createExternalBnSelectStep(wiz, 5, totalSteps)
	wiz.externalBnSettingsModal = createExternalBnSettingsStep(wiz, 5, totalSteps)
	wiz.externalPrysmSettingsModal = createExternalPrysmSettingsStep(wiz, 5, totalSteps)

	// Step 6 - Fallback clients
	wiz.useFallbackModal = createUseFallbackStep(wiz, 6, totalSteps)
	wiz.fallbackNormalModal = createFallbackNormalStep(wiz, 6, totalSteps)
	wiz.fallbackPrysmModal = createFallbackPrysmStep(wiz, 6, totalSteps)

	// Step 7 - Modules
	wiz.modulesModal = createModulesStep(wiz, 7, totalSteps)

	// Step 8 - Metrics
	wiz.metricsModal = createMetricsStep(wiz, 8, totalSteps)

	// Done
	wiz.finishedModal = createFinishedStep(wiz, 9, totalSteps)

	return wiz
}
