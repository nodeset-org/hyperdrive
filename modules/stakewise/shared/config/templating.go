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

func (c *StakewiseConfig) KeystorePasswordFile() string {
	return KeystorePasswordFile
}

func (c *StakewiseConfig) DaemonContainerSuffix() string {
	return string(ContainerID_StakewiseDaemon)
}

func (c *StakewiseConfig) OperatorContainerSuffix() string {
	return string(ContainerID_StakewiseOperator)
}

func (c *StakewiseConfig) VcContainerSuffix() string {
	return string(ContainerID_StakewiseValidator)
}

func (c *StakewiseConfig) DepositDataFile() string {
	return DepositDataFile
}
