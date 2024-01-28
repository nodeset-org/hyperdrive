package wallet

const (
	StakewiseValidatorPath string = "m/12381/3600/%d/1/0" // Stakewise keys are generated with `use` index 1
)

// Data relating to Stakewise's wallet
type StakewiseWalletData struct {
	// The next account to generate the key for
	NextAccount uint64 `json:"nextAccount`
}
