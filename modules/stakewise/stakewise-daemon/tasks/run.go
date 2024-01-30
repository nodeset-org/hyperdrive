package tasks

import (
	"context"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

// Config
var tasksInterval, _ = time.ParseDuration("5m")
var taskCooldown, _ = time.ParseDuration("10s")

const (
	ErrorColor             = color.FgRed
	WarningColor           = color.FgYellow
	UpdateDepositDataColor = color.FgHiWhite
)

type TaskLoop struct {
	ctx    context.Context
	cancel context.CancelFunc
	sp     *common.StakewiseServiceProvider
	wg     *sync.WaitGroup
}

func NewTaskLoop(sp *common.StakewiseServiceProvider, wg *sync.WaitGroup) *TaskLoop {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskLoop{
		ctx:    ctx,
		cancel: cancel,
		sp:     sp,
		wg:     wg,
	}
}

// Run daemon
func (t *TaskLoop) Run() error {
	// Initialize loggers
	errorLog := log.NewColorLogger(ErrorColor)

	// Initialize tasks
	updateDepositData := NewUpdateDepositData(t.sp, log.NewColorLogger(UpdateDepositDataColor))

	// Run the loop
	go func() {
		for {
			// Check the EC status
			err := t.sp.WaitEthClientSynced(false) // Force refresh the primary / fallback EC status
			if err != nil {
				errorLog.Println(err)
				if t.sleepAndCheckIfCancelled(taskCooldown) {
					break
				}
				continue
			}

			// Check the BC status
			err = t.sp.WaitBeaconClientSynced(false) // Force refresh the primary / fallback BC status
			if err != nil {
				errorLog.Println(err)
				if t.sleepAndCheckIfCancelled(taskCooldown) {
					break
				}
				continue
			}

			// Update deposit data from the NodeSet server
			if err := updateDepositData.Run(); err != nil {
				errorLog.Println(err)
			}
			// time.Sleep(taskCooldown)

			if t.sleepAndCheckIfCancelled(tasksInterval) {
				break
			}
		}

		// Signal the task loop is done
		t.wg.Done()
	}()

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
	return nil
}

func (t *TaskLoop) Stop() {
	t.cancel()
}

func (t *TaskLoop) sleepAndCheckIfCancelled(duration time.Duration) bool {
	timer := time.NewTimer(duration)
	select {
	case <-t.ctx.Done():
		// Cancel occurred
		timer.Stop()
		return true

	case <-timer.C:
		// Duration has passed without a cancel
		return false
	}
}
