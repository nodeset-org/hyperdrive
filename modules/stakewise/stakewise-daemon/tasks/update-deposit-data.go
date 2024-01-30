package tasks

import (
	"github.com/docker/docker/client"
	"github.com/nodeset-org/eth-utils/eth"
	"github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

// Update deposit data task
type UpdateDepositData struct {
	sp  *common.StakewiseServiceProvider
	cfg *config.HyperdriveConfig
	log log.ColorLogger
	ec  eth.IExecutionClient
	bc  types.IBeaconClient
	d   *client.Client
}

// Create update deposit data task
func NewUpdateDepositData(sp *common.StakewiseServiceProvider, logger log.ColorLogger) *UpdateDepositData {
	return &UpdateDepositData{
		sp:  sp,
		log: logger,
	}
}

// Update deposit data
func (t *UpdateDepositData) Run() error {
	// Get services
	t.cfg = t.sp.GetConfig()
	t.ec = t.sp.GetEthClient()
	t.bc = t.sp.GetBeaconClient()
	t.d = t.sp.GetDocker()

	// TODO

	return nil
}
