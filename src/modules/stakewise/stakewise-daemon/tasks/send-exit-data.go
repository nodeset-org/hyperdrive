package swtasks

import (
	"context"
	"fmt"

	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/daemon-utils/validator/utils"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
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
	w := t.sp.GetWallet()
	ns := t.sp.GetNodesetClient()
	bc := t.sp.GetBeaconClient()

	resp, err := ns.GetRegisteredValidators()
	if err != nil {
		return fmt.Errorf("error getting registered validators: %w", err)
	}
	pubkeys := []beacon.ValidatorPubkey{}
	for _, v := range resp {
		pubkeys = append(pubkeys, v.Pubkey)
	}
	statuses, err := bc.GetValidatorStatuses(context.Background(), pubkeys, nil)
	if err != nil {
		return fmt.Errorf("error getting validator statuses: %w", err)
	}
	head, err := bc.GetBeaconHead(context.Background())

	epoch := head.Epoch
	signatureDomain, err := bc.GetDomainData(context.Background(), eth2types.DomainVoluntaryExit[:], epoch, false)
	exitData := []swcommon.ExitData{}
	for _, v := range resp {
		fmt.Printf("Validator: %v\n", v)
		if !v.Uploaded {
			fmt.Printf("Validator %v has not been uploaded\n", v.Pubkey)
			fmt.Printf("Attempting to generate exit data for validator %v\n", v.Pubkey)
			key, err := w.GetPrivateKeyForPubkey(v.Pubkey)
			if err != nil {
				return fmt.Errorf("error getting private key for pubkey %v: %w", v.Pubkey, err)
			}
			index := statuses[v.Pubkey].Index
			signature, err := utils.GetSignedExitMessage(key, index, epoch, signatureDomain)
			if err != nil {
				return fmt.Errorf("error getting signed exit message: %w", err)
			}
			exitData = append(exitData, swcommon.ExitData{
				Pubkey: v.Pubkey.HexWithPrefix(),
				ExitMessage: swcommon.ExitMessage{
					Message: swcommon.ExitMessageDetails{
						Epoch:          string(epoch),
						ValidatorIndex: index,
					},
					Signature: signature.HexWithPrefix(),
				},
			})
		}
	}
	ns.PostExitData(exitData)
	fmt.Printf("Registered validators: %v\n", resp)
	return nil
}
