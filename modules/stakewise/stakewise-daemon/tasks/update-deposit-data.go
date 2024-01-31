package swtasks

import (
	"fmt"

	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
)

// Update deposit data task
type UpdateDepositData struct {
	sp  *swcommon.StakewiseServiceProvider
	cfg *config.HyperdriveConfig
	log log.ColorLogger
}

// Create update deposit data task
func NewUpdateDepositData(sp *swcommon.StakewiseServiceProvider, logger log.ColorLogger) *UpdateDepositData {
	return &UpdateDepositData{
		sp:  sp,
		log: logger,
	}
}

// Update deposit data
func (t *UpdateDepositData) Run() error {
	t.log.Println("Checking version of NodeSet data on disk...")

	// Get services
	w := t.sp.GetWallet()
	hd := t.sp.GetHyperdriveClient()
	ns := t.sp.GetNodesetClient()
	ddMgr := t.sp.GetDepositDataManager()

	// Check the version of the NodeSet

	// Get the version on the server
	remoteVersion, err := ns.GetServerDepositDataVersion()
	if err != nil {
		return fmt.Errorf("error getting latest deposit data version: %w", err)
	}

	// Compare versions
	localVersion := w.GetLatestDepositDataVersion()
	if remoteVersion == localVersion {
		t.log.Printlnf("Local data is up to date (version %d).", localVersion)
		return nil
	}

	// Get the new data
	t.log.Printlnf("Latest data version is %d but we have %d, retrieving latest data...", remoteVersion, localVersion)
	_, depositData, err := ns.GetServerDepositData()
	if err != nil {
		return fmt.Errorf("error getting latest deposit data: %w", err)
	}

	// Save it
	err = ddMgr.UpdateDepositData(depositData)
	if err != nil {
		return fmt.Errorf("error saving deposit data: %w", err)
	}
	err = w.SetLatestDepositDataVersion(remoteVersion)
	if err != nil {
		return fmt.Errorf("error updating latest saved version number: %w", err)
	}

	// Restart the Stakewise op container
	t.log.Printlnf("Restarting Stakewise operator...")
	_, err = hd.Service.RestartContainer(string(swconfig.ContainerID_StakewiseOperator))
	if err != nil {
		return fmt.Errorf("error restarting %s container: %w", swconfig.ContainerID_StakewiseOperator, err)
	}

	t.log.Println("Done! Your deposit data is now up to date.")
	return nil
}
