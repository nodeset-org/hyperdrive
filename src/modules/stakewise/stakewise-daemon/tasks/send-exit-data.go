package swtasks

import (
	"fmt"

	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

// Send exit data task
type SendExitData struct {
	sp *swcommon.StakewiseServiceProvider

	log log.ColorLogger
}

// Create Exit data task
func NewSendExitData(sp *swcommon.StakewiseServiceProvider, logger log.ColorLogger) *SendExitData {
	return &SendExitData{
		sp:  sp,
		log: logger,
	}
}

// Update Exit data
func (t *SendExitData) Run() error {
	t.log.Println("Checking Nodeset API...")
	// w := t.sp.GetWallet()
	// hd := t.sp.GetHyperdriveClient()
	ns := t.sp.GetNodesetClient()
	// ddMgr := t.sp.GetDepositDataManager()
	// cfg := t.sp.GetModuleConfig()

	resp, err := ns.GetRegisteredValidators()
	if err != nil {
		return fmt.Errorf("error getting registered validators: %w", err)
	}

	fmt.Printf("Registered validators: %v\n", resp)
	return nil
}
