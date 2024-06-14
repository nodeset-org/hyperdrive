package shared

const (
	// Hyperdrive distinguishes keys by module to prevent overlapping between modules.
	// It uses the `use` field of the path, as defined in EIP-2334, to represent each module.

	RocketPoolValidatorPath    string = "m/12381/3600/%d/0/0"
	StakeWiseValidatorPath     string = "m/12381/3600/%d/1/0"
	ConstellationValidatorPath string = "m/12381/3600/%d/2/0"
	SoloValidatorPath          string = "m/12381/3600/%d/3/0"
)
