package daemon

import (
	"github.com/docker/docker/client"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/hyperdrive-daemon/common/services"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/config"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/types"
	"github.com/nodeset-org/hyperdrive-stakewise-daemon/shared/utils/log"
)

// Update deposit data task
type UpdateDepositData struct {
	sp  *services.ServiceProvider
	cfg *config.HyperdriveConfig
	log log.ColorLogger
	ec  types.IExecutionClient
	bc  types.IBeaconClient
	d   *client.Client
}

// Create update deposit data task
func NewUpdateDepositData(sp *services.ServiceProvider, logger log.ColorLogger) *UpdateDepositData {
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
