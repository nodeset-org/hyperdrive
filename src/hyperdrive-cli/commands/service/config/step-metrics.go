package config

func createMetricsStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	helperText := "Would you like to enable Hyperdrive's metrics monitoring system? This will monitor things such as hardware stats (CPU usage, RAM usage, free disk space), your validator stats, stats about your node such as total ETH rewards, and much more. It also enables the Grafana dashboard to quickly and easily view these metrics.\n\nNone of this information will be sent to any remote servers for collection an analysis; this is purely for your own usage on your node."

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		if !wiz.md.Config.Hyperdrive.Metrics.EnableMetrics.Value {
			modal.focus(0)
		} else {
			modal.focus(1)
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 1 {
			wiz.md.Config.Hyperdrive.Metrics.EnableMetrics.Value = true
		} else {
			wiz.md.Config.Hyperdrive.Metrics.EnableMetrics.Value = false
		}
		wiz.finishedModal.show()
	}

	back := func() {
		wiz.modulesModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		[]string{"No", "Yes"},
		[]string{},
		76,
		"Metrics",
		DirectionalModalHorizontal,
		show,
		done,
		back,
		"step-metrics",
	)
}
