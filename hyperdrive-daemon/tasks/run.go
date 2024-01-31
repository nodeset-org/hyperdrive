package tasks

import (
	"context"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
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
	sp     *common.ServiceProvider
	wg     *sync.WaitGroup
}

func NewTaskLoop(sp *common.ServiceProvider, wg *sync.WaitGroup) *TaskLoop {
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
	// Nothing here yet

	// Run the loop
	go func() {
		for {
			// Check the EC status
			err := t.sp.WaitEthClientSynced(t.ctx, false) // Force refresh the primary / fallback EC status
			if err != nil {
				errorLog.Println(err)
				if t.sleepAndCheckIfCancelled(taskCooldown) {
					break
				}
				continue
			}

			// Check the BC status
			err = t.sp.WaitBeaconClientSynced(t.ctx, false) // Force refresh the primary / fallback BC status
			if err != nil {
				errorLog.Println(err)
				if t.sleepAndCheckIfCancelled(taskCooldown) {
					break
				}
				continue
			}

			// Tasks go here

			if t.sleepAndCheckIfCancelled(tasksInterval) {
				break
			}
		}

		// Signal the task loop is done
		t.wg.Done()
	}()
	t.wg.Add(1)

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
