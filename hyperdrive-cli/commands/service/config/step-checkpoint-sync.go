package config

func createCheckpointSyncStep(wiz *wizard, currentStep int, totalSteps int) *textBoxWizardStep {

	// Create the labels and args
	checkpointSyncLabel := wiz.md.Config.LocalBeaconConfig.CheckpointSyncProvider.Name

	helperText := "Your client supports Checkpoint Sync. This powerful feature allows it to copy the most recent state from a separate Consensus client that you trust, so you don't have to wait for it to sync from scratch - you can start using it instantly!\n\nTake a look at our documentation for an example of how to use it:\nhttps://docs.rocketpool.net/guides/node/config-docker.html#beacon-chain-checkpoint-syncing\n\nIf you would like to use Checkpoint Sync, please provide the provider URL here. If you don't want to use it, leave it blank."

	show := func(modal *textBoxModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus()
		for label, box := range modal.textboxes {
			for _, param := range wiz.md.Config.LocalBeaconConfig.GetParameters() {
				if param.GetCommon().Name == label {
					box.SetText(param.GetValueAsString())
				}
			}
		}
	}

	done := func(text map[string]string) {
		wiz.md.Config.LocalBeaconConfig.CheckpointSyncProvider.Value = text[checkpointSyncLabel]
		wiz.useFallbackModal.show()
	}

	back := func() {
		wiz.bnLocalModal.show()
	}

	return newTextBoxWizardStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		76,
		"Beacon Node > Checkpoint Sync",
		[]string{checkpointSyncLabel},
		[]int{wiz.md.Config.LocalBeaconConfig.CheckpointSyncProvider.MaxLength},
		[]string{wiz.md.Config.LocalBeaconConfig.CheckpointSyncProvider.Regex},
		show,
		done,
		back,
		"step-checkpoint-sync",
	)

}
