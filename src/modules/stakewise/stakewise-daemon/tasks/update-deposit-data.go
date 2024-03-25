package swtasks

import (
	"fmt"

	swcontracts "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common/contracts"

	"github.com/ethereum/go-ethereum/common"
	swconfig "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/config"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
	batch "github.com/rocket-pool/batch-query"
)

// Update deposit data task
type UpdateDepositData struct {
	sp  *swcommon.StakewiseServiceProvider
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
	cfg := t.sp.GetModuleConfig()

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

	// Verify the merkle roots if enabled
	if cfg.VerifyDepositsRoot.Value {
		isMatch, err := t.verifyDepositsRoot(depositData)
		if err != nil {
			return err
		}
		if !isMatch {
			return nil
		}
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

// Verify the Merkle root from the deposits data matches what's on chain before saving
func (t *UpdateDepositData) verifyDepositsRoot(depositData []types.ExtendedDepositData) (bool, error) {
	// Get services
	ec := t.sp.GetEthClient()
	res := t.sp.GetResources()
	txMgr := t.sp.GetTransactionManager()
	q := t.sp.GetQueryManager()
	ddMgr := t.sp.GetDepositDataManager()

	// Get the Merkle root from it
	localRoot, err := ddMgr.ComputeMerkleRoot(depositData)
	if err != nil {
		return false, fmt.Errorf("error computing Merkle root from deposit data: %w", err)
	}
	t.log.Printlnf("Computed Merkle root:   %s", localRoot.Hex())

	// Get the Merkle root from the vault
	vault, err := swcontracts.NewStakewiseVault(res.Vault, ec, txMgr)
	if err != nil {
		return false, fmt.Errorf("error creating Stakewise Vault binding: %w", err)
	}
	var contractRoot common.Hash
	err = q.Query(func(mc *batch.MultiCaller) error {
		vault.GetValidatorsRoot(mc, &contractRoot)
		return nil
	}, nil)
	if err != nil {
		return false, fmt.Errorf("error getting canonical deposit root from the Stakewise Vault: %w", err)
	}
	t.log.Printlnf("Contract's Merkle root: %s", contractRoot.Hex())

	// Compare them
	if localRoot != contractRoot {
		t.log.Printlnf("WARNING: Locally computed deposits data root does not match the value stored on chain, refusing to save for safety!")
		return false, nil
	} else {
		t.log.Println("Locally computed deposits data root matches the root stored on-chain, updating may proceed.")
	}
	return true, nil
}
