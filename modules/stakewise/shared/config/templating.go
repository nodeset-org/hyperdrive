package swconfig

func (c *StakewiseConfig) DaemonRoute() string {
	return DaemonRoute
}

func (c *StakewiseConfig) WalletFilename() string {
	return WalletFilename
}

func (c *StakewiseConfig) PasswordFilename() string {
	return PasswordFilename
}

func (c *StakewiseConfig) DaemonContainerSuffix() string {
	return DaemonContainerSuffix
}

func (c *StakewiseConfig) OperatorContainerSuffix() string {
	return OperatorContainerSuffix
}

func (c *StakewiseConfig) VcContainerSuffix() string {
	return VcContainerSuffix
}
