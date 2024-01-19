package daemon

import (
	"time"

	"github.com/fatih/color"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/common/services"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/utils/log"
)

// Config
var tasksInterval, _ = time.ParseDuration("5m")
var taskCooldown, _ = time.ParseDuration("10s")

const (
	ErrorColor             = color.FgRed
	WarningColor           = color.FgYellow
	UpdateDepositDataColor = color.FgHiWhite
)

// Run daemon
func Run(sp *services.ServiceProvider) error {
	// Initialize loggers
	errorLog := log.NewColorLogger(ErrorColor)

	// Initialize tasks
	updateDepositData := NewUpdateDepositData(sp, log.NewColorLogger(UpdateDepositDataColor))

	for {
		// Check the EC status
		err := sp.WaitEthClientSynced(false) // Force refresh the primary / fallback EC status
		if err != nil {
			errorLog.Println(err)
			time.Sleep(taskCooldown)
			continue
		}

		// Check the BC status
		err = sp.WaitBeaconClientSynced(false) // Force refresh the primary / fallback BC status
		if err != nil {
			errorLog.Println(err)
			time.Sleep(taskCooldown)
			continue
		}

		// Update deposit data from the NodeSet server
		if err := updateDepositData.Run(); err != nil {
			errorLog.Println(err)
		}
		// time.Sleep(taskCooldown)

		time.Sleep(tasksInterval)
	}

	/*
		// Run metrics loop
		go func() {
			err := runMetricsServer(sp, log.NewColorLogger(MetricsColor), stateLocker)
			if err != nil {
				errorLog.Println(err)
			}
			wg.Done()
		}()
	*/
}
