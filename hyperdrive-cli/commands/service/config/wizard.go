package config

type wizard struct {
	md                              *mainDisplay
	welcomeModal                    *choiceWizardStep
	networkModal                    *choiceWizardStep
	modeModal                       *choiceWizardStep
	executionLocalModal             *choiceWizardStep
	executionExternalModal          *textBoxWizardStep
	consensusLocalModal             *choiceWizardStep
	consensusExternalSelectModal    *choiceWizardStep
	graffitiModal                   *textBoxWizardStep
	checkpointSyncProviderModal     *textBoxWizardStep
	doppelgangerDetectionModal      *choiceWizardStep
	lighthouseExternalSettingsModal *textBoxWizardStep
	nimbusExternalSettingsModal     *textBoxWizardStep
	lodestarExternalSettingsModal   *textBoxWizardStep
	prysmExternalSettingsModal      *textBoxWizardStep
	tekuExternalSettingsModal       *textBoxWizardStep
	externalGraffitiModal           *textBoxWizardStep
	metricsModal                    *choiceWizardStep
	mevModeModal                    *choiceWizardStep
	localMevModal                   *checkBoxWizardStep
	externalMevModal                *textBoxWizardStep
	finishedModal                   *choiceWizardStep
	consensusLocalRandomModal       *choiceWizardStep
	consensusLocalRandomPrysmModal  *choiceWizardStep
	consensusLocalPrysmWarning      *choiceWizardStep
	consensusLocalTekuWarning       *choiceWizardStep
	externalDoppelgangerModal       *choiceWizardStep
	executionLocalRandomModal       *choiceWizardStep
	useFallbackModal                *choiceWizardStep
	fallbackNormalModal             *textBoxWizardStep
	fallbackPrysmModal              *textBoxWizardStep
}

func newWizard(md *mainDisplay) *wizard {
	wiz := &wizard{
		md: md,
	}

	totalSteps := 9

	// Docker mode
	wiz.welcomeModal = createWelcomeStep(wiz, 1, totalSteps)
	wiz.networkModal = createNetworkStep(wiz, 2, totalSteps)
	wiz.modeModal = createModeStep(wiz, 3, totalSteps)
	wiz.executionLocalModal = createLocalEcStep(wiz, 4, totalSteps)
	wiz.executionExternalModal = createExternalEcStep(wiz, 4, totalSteps)
	wiz.consensusLocalModal = createLocalCcStep(wiz, 5, totalSteps)
	wiz.consensusExternalSelectModal = createExternalCcStep(wiz, 5, totalSteps)
	wiz.consensusLocalPrysmWarning = createPrysmWarningStep(wiz, 5, totalSteps)
	wiz.consensusLocalTekuWarning = createTekuWarningStep(wiz, 5, totalSteps)
	wiz.graffitiModal = createGraffitiStep(wiz, 5, totalSteps)
	wiz.checkpointSyncProviderModal = createCheckpointSyncStep(wiz, 5, totalSteps)
	wiz.doppelgangerDetectionModal = createDoppelgangerStep(wiz, 5, totalSteps)
	wiz.lighthouseExternalSettingsModal = createExternalLhStep(wiz, 5, totalSteps)
	wiz.nimbusExternalSettingsModal = createExternalNimbusStep(wiz, 5, totalSteps)
	wiz.lodestarExternalSettingsModal = createExternalLodestarStep(wiz, 5, totalSteps)
	wiz.prysmExternalSettingsModal = createExternalPrysmStep(wiz, 5, totalSteps)
	wiz.tekuExternalSettingsModal = createExternalTekuStep(wiz, 5, totalSteps)
	wiz.externalGraffitiModal = createExternalGraffitiStep(wiz, 5, totalSteps)
	wiz.externalDoppelgangerModal = createExternalDoppelgangerStep(wiz, 5, totalSteps)
	wiz.useFallbackModal = createUseFallbackStep(wiz, 6, totalSteps)
	wiz.fallbackNormalModal = createFallbackNormalStep(wiz, 6, totalSteps)
	wiz.fallbackPrysmModal = createFallbackPrysmStep(wiz, 6, totalSteps)
	wiz.metricsModal = createMetricsStep(wiz, 7, totalSteps)
	wiz.mevModeModal = createMevModeStep(wiz, 8, totalSteps)
	wiz.localMevModal = createLocalMevStep(wiz, 8, totalSteps)
	wiz.externalMevModal = createExternalMevStep(wiz, 8, totalSteps)
	wiz.finishedModal = createFinishedStep(wiz, 9, totalSteps)

	return wiz
}
