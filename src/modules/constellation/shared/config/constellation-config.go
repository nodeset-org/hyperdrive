package constconfig

import (
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/config/validator"
)

const (
// Param IDs
// StakewiseEnableID      string = "enable"
// OperatorContainerTagID string = "operatorContainerTag"
// AdditionalOpFlagsID    string = "additionalOpFlags"
// VerifyDepositRootsID   string = "verifyDepositRoots"
)

// Configuration for Stakewise
type ConstellationConfig struct {
	hdCfg *config.HyperdriveConfig

	// Custom command line flags
	AdditionalOpFlags config.Parameter[string]

	// Validator client configs
	VcCommon   *validator.ValidatorClientCommonConfig
	Lighthouse *validator.LighthouseVcConfig
	Lodestar   *validator.LodestarVcConfig
	Nimbus     *validator.NimbusVcConfig
	Prysm      *validator.PrysmVcConfig
	Teku       *validator.TekuVcConfig
}

// The title for the config
func (cfg *ConstellationConfig) GetTitle() string {
	return "Constellation"
}
