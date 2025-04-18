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
	externalRethWarning     *choiceWizardStep
	externalEcSettingsModal *textBoxWizardStep

	// Step 5 - BN settings
	localBnModal                *choiceWizardStep
	localBnRandomModal          *choiceWizardStep
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
	modulesModal         *checkBoxWizardStep
	modulesDisabledModal *choiceWizardStep

	// Step 8 - Metrics
	metricsModal *choiceWizardStep

	// Step 9 - MEV Boost
	mevModeModal     *choiceWizardStep
	localMevModal    *checkBoxWizardStep
	externalMevModal *textBoxWizardStep
	mevWarningModal  *choiceWizardStep
	mevDisabledModal *choiceWizardStep

	// Done
	finishedModal *choiceWizardStep
}

// Create a new Wizard display
func newWizard(md *mainDisplay) *wizard {
	wiz := &wizard{
		md: md,
	}

	totalSteps := 10
	stepCount := 1

	// Step 1 - Welcome
	wiz.welcomeModal = createWelcomeStep(wiz, stepCount, totalSteps)
	stepCount++

	// Step 2 - Network
	wiz.networkModal = createNetworkStep(wiz, stepCount, totalSteps)
	stepCount++

	// Step 3 - Client mode
	wiz.modeModal = createModeStep(wiz, stepCount, totalSteps)
	stepCount++

	// Step 4 - EC settings
	wiz.localEcModal = createLocalEcStep(wiz, stepCount, totalSteps)
	wiz.externalEcSelectModal = createExternalEcSelectStep(wiz, stepCount, totalSteps)
	wiz.externalRethWarning = createExternalRethWarningStep(wiz, stepCount, totalSteps)
	wiz.externalEcSettingsModal = createExternalEcSettingsStep(wiz, stepCount, totalSteps)
	stepCount++

	// Step 5 - BN settings
	wiz.localBnModal = createLocalBnStep(wiz, stepCount, totalSteps)
	wiz.localBnPrysmWarning = createPrysmWarningStep(wiz, stepCount, totalSteps)
	wiz.localBnTekuWarning = createTekuWarningStep(wiz, stepCount, totalSteps)
	wiz.checkpointSyncProviderModal = createCheckpointSyncStep(wiz, stepCount, totalSteps)
	wiz.externalBnSelectModal = createExternalBnSelectStep(wiz, stepCount, totalSteps)
	wiz.externalBnSettingsModal = createExternalBnSettingsStep(wiz, stepCount, totalSteps)
	wiz.externalPrysmSettingsModal = createExternalPrysmSettingsStep(wiz, stepCount, totalSteps)
	stepCount++

	// Step 6 - Fallback clients
	wiz.useFallbackModal = createUseFallbackStep(wiz, stepCount, totalSteps)
	wiz.fallbackNormalModal = createFallbackNormalStep(wiz, stepCount, totalSteps)
	wiz.fallbackPrysmModal = createFallbackPrysmStep(wiz, stepCount, totalSteps)
	stepCount++

	// Step 7 - Modules
	wiz.modulesModal = createModulesStep(wiz, stepCount, totalSteps)
	wiz.modulesDisabledModal = createModulesDisabledStep(wiz, stepCount, totalSteps)
	stepCount++

	// Step 8 - Metrics
	wiz.metricsModal = createMetricsStep(wiz, stepCount, totalSteps)
	stepCount++

	// Step 9 - MEV Boost
	wiz.mevModeModal = createMevModeStep(wiz, stepCount, totalSteps)
	wiz.localMevModal = createLocalMevStep(wiz, stepCount, totalSteps)
	wiz.externalMevModal = createExternalMevStep(wiz, stepCount, totalSteps)
	wiz.mevWarningModal = createMevWarningStep(wiz, stepCount, totalSteps)
	wiz.mevDisabledModal = createMevDisabledStep(wiz, stepCount, totalSteps)
	stepCount++

	// Done
	wiz.finishedModal = createFinishedStep(wiz, stepCount, totalSteps)

	return wiz
}
