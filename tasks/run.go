package tasks

import (
	"context"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/utils"
)

// Config
const (
	tasksInterval time.Duration = time.Minute * 5
	taskCooldown  time.Duration = time.Second * 10

	ErrorColor             = color.FgRed
	WarningColor           = color.FgYellow
	UpdateDepositDataColor = color.FgHiWhite
)

type TaskLoop struct {
	ctx    context.Context
	logger *log.Logger
	sp     *common.ServiceProvider
	wg     *sync.WaitGroup
}

func NewTaskLoop(sp *common.ServiceProvider, wg *sync.WaitGroup) *TaskLoop {
	taskLoop := &TaskLoop{
		sp:     sp,
		logger: sp.GetTasksLogger(),
		wg:     wg,
	}
	taskLoop.ctx = taskLoop.logger.CreateContextWithLogger(sp.GetBaseContext())
	return taskLoop
}

// Run daemon
func (t *TaskLoop) Run() error {
	// Initialize tasks

	// Run the loop
	t.wg.Add(1)
	go func() {
		for {
			// Check the EC status
			err := t.sp.WaitEthClientSynced(t.ctx, false) // Force refresh the primary / fallback EC status
			if err != nil {
				t.logger.Error(err.Error())
				if utils.SleepWithCancel(t.ctx, taskCooldown) {
					break
				}
				continue
			}

			// Check the BC status
			err = t.sp.WaitBeaconClientSynced(t.ctx, false) // Force refresh the primary / fallback BC status
			if err != nil {
				t.logger.Error(err.Error())
				if utils.SleepWithCancel(t.ctx, taskCooldown) {
					break
				}
				continue
			}

			// Tasks go here

			if utils.SleepWithCancel(t.ctx, tasksInterval) {
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
