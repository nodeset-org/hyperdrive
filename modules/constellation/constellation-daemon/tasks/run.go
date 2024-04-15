package consttasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/utils"

	constcommon "github.com/nodeset-org/hyperdrive/modules/constellation/constellation-daemon/common"
)

// Config
const (
	tasksInterval time.Duration = time.Minute * 5
	taskCooldown  time.Duration = time.Second * 10

	ErrorColor             = color.FgRed
	WarningColor           = color.FgYellow
	UpdateDepositDataColor = color.FgHiWhite
	SendExitDataColor      = color.FgGreen
)

type TaskLoop struct {
	ctx    context.Context
	logger *log.Logger
	sp     *constcommon.ConstellationServiceProvider
	wg     *sync.WaitGroup
	// stakePrelaunchMinipools *node.StakePrelaunchMinipools
	// ec                      eth.IExecutionClient
	// bc                      beacon.IBeaconClient
	// stateLocker             *collectors.StateLocker
	// cfg                     *config.SmartNodeConfig
	// rp                      *rocketpool.RocketPool
}

func NewTaskLoop(sp *constcommon.ConstellationServiceProvider, wg *sync.WaitGroup) *TaskLoop {
	fmt.Printf("!!! executing task loop\n")
	fmt.Printf("!!! user dir: %s\n", sp.ServiceProvider.GetUserDir())
	taskLoop := &TaskLoop{
		sp:     sp,
		logger: sp.ServiceProvider.GetTasksLogger(),
		wg:     wg,
		// stakePrelaunchMinipools: node.NewStakePrelaunchMinipools(sp.RpServiceProvider, sp.ServiceProvider.GetTasksLogger()),
	}
	taskLoop.ctx = taskLoop.logger.CreateContextWithLogger(sp.ServiceProvider.GetBaseContext())
	fmt.Printf("!!! executing task loop done\n")
	return taskLoop
}

// Run daemon
func (t *TaskLoop) Run() error {
	// Initialize tasks
	// 	updateDepositData := NewUpdateDepositDataTask(t.ctx, t.sp, t.logger)
	// 	sendExitData := NewSendExitData(t.ctx, t.sp, t.logger)
	fmt.Printf("!!! executing task loop\n")
	t.logger.Warn("!!! executing task loop\n")

	// Run the loop
	t.wg.Add(1)
	go func() {
		fmt.Printf("!!! executing task loop 1\n")
		t.logger.Warn("!!! executing task loop 1\n")

		for {
			err := t.sp.ServiceProvider.WaitEthClientSynced(t.ctx, false) // Force refresh the primary / fallback EC status
			if err != nil {
				t.logger.Error(err.Error())
				if utils.SleepWithCancel(t.ctx, taskCooldown) {
					break
				}
				continue
			}

			// Check the BC status
			err = t.sp.ServiceProvider.WaitBeaconClientSynced(t.ctx, false) // Force refresh the primary / fallback BC status
			if err != nil {
				t.logger.Error(err.Error())
				if utils.SleepWithCancel(t.ctx, taskCooldown) {
					break
				}
				continue
			}

			// Tasks start here

			// // Update deposit data from the NodeSet server
			// if err := updateDepositData.Run(); err != nil {
			// 	t.logger.Error(err.Error())
			// }
			// if utils.SleepWithCancel(t.ctx, taskCooldown) {
			// 	break
			// }

			// // Submit missing exit messages to the NodeSet server
			// if err := sendExitData.Run(); err != nil {
			// 	t.logger.Error(err.Error())
			// }

			// Tasks end here

			if utils.SleepWithCancel(t.ctx, tasksInterval) {
				break
			}
		}
		fmt.Printf("!!! executing task loop 2\n")
		t.logger.Warn("!!! executing task loop 2\n")

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
	fmt.Printf("!!! executing task loop done\n")
	t.logger.Warn("!!! executing task loop done\n")

	return nil
}
