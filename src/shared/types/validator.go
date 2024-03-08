package types

import (
	"github.com/rocket-pool/node-manager-core/beacon"
)

// Extended deposit data beyond what is required in an actual deposit message to Beacon, emulating what the deposit CLI produces
type ExtendedDepositData struct {
	beacon.ExtendedDepositData
	HyperdriveVersion string `json:"hyperdrive_version,omitempty"`
}
