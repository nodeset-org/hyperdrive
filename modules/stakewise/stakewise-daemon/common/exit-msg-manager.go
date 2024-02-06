package swcommon

import (
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

type ExitMessageManager struct {
	sp *common.ServiceProvider
}

func (e *ExitMessageManager) GenerateExitMessage(key *eth2types.BLSPrivateKey, index uint64) {

}
