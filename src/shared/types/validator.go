package types

import (
	nmc_beacon "github.com/rocket-pool/node-manager-core/beacon"
)

const (
	// Hyperdrive distinguishes keys by module to prevent overlapping between modules.
	// It uses the `use` field of the path, as defined in EIP-2334, to represent each module.

	RocketPoolValidatorPath    string = "m/12381/3600/%d/0/0"
	StakewiseValidatorPath     string = "m/12381/3600/%d/1/0"
	ConstellationValidatorPath string = "m/12381/3600/%d/2/0"
	SoloValidatorPath          string = "m/12381/3600/%d/3/0"
)

// Extended deposit data beyond what is required in an actual deposit message to Beacon, emulating what the deposit CLI produces
type ExtendedDepositData struct {
	nmc_beacon.ExtendedDepositData
	HyperdriveVersion string `json:"hyperdrive_version,omitempty"`
}
