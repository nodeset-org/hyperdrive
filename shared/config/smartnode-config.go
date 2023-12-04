package config

// Configuration for the Smartnode
type SmartnodeConfig struct {
	Title string `yaml:"-"`

	// The parent config
	// parent *RocketPoolConfig

	////////////////////////////
	// User-editable settings //
	////////////////////////////

	// Docker container prefix
	ProjectName Parameter `yaml:"projectName,omitempty"`

	// The path of the data folder where everything is stored
	DataPath Parameter `yaml:"dataPath,omitempty"`

	// The path of the watchtower's persistent state storage
	WatchtowerStatePath Parameter `yaml:"watchtowerStatePath"`

	// Which network we're on
	Network Parameter `yaml:"network,omitempty"`

	// Manual max fee override
	ManualMaxFee Parameter `yaml:"manualMaxFee,omitempty"`

	// Manual priority fee override
	PriorityFee Parameter `yaml:"priorityFee,omitempty"`

	// Threshold for automatic transactions
	AutoTxGasThreshold Parameter `yaml:"minipoolStakeGasThreshold,omitempty"`

	// The amount of ETH in a minipool's balance before auto-distribute kicks in
	DistributeThreshold Parameter `yaml:"distributeThreshold,omitempty"`

	// Mode for acquiring Merkle rewards trees
	RewardsTreeMode Parameter `yaml:"rewardsTreeMode,omitempty"`

	// URL for an EC with archive mode, for manual rewards tree generation
	ArchiveECUrl Parameter `yaml:"archiveEcUrl,omitempty"`

	// Token for Oracle DAO members to use when uploading Merkle trees to Web3.Storage
	Web3StorageApiToken Parameter `yaml:"web3StorageApiToken,omitempty"`

	// Manual override for the watchtower's max fee
	WatchtowerMaxFeeOverride Parameter `yaml:"watchtowerMaxFeeOverride,omitempty"`

	// Manual override for the watchtower's priority fee
	WatchtowerPrioFeeOverride Parameter `yaml:"watchtowerPrioFeeOverride,omitempty"`

	// The toggle for rolling records
	UseRollingRecords Parameter `yaml:"useRollingRecords,omitempty"`

	// The rolling record checkpoint interval
	RecordCheckpointInterval Parameter `yaml:"recordCheckpointInterval,omitempty"`

	// The checkpoint retention limit
	CheckpointRetentionLimit Parameter `yaml:"checkpointRetentionLimit,omitempty"`

	// The path of the records folder where snapshots of rolling record info is stored during a rewards interval
	RecordsPath Parameter `yaml:"recordsPath,omitempty"`

	///////////////////////////
	// Non-editable settings //
	///////////////////////////

	// TODO: Uncomment as needed
	// // The URL to provide the user so they can follow pending transactions
	// txWatchUrl map[Network]string `yaml:"-"`

	// // The URL to use for staking rETH
	// stakeUrl map[Network]string `yaml:"-"`

	// // The map of networks to execution chain IDs
	// chainID map[Network]uint `yaml:"-"`

	// // The contract address of RocketStorage
	// storageAddress map[Network]string `yaml:"-"`

	// // The contract address of the RPL token
	// rplTokenAddress map[Network]string `yaml:"-"`

	// // The contract address of the RPL faucet
	// rplFaucetAddress map[Network]string `yaml:"-"`

	// // The contract address for Snapshot delegation
	// snapshotDelegationAddress map[Network]string `yaml:"-"`

	// // The Snapshot API domain
	// snapshotApiDomain map[Network]string `yaml:"-"`

	// // The contract address of rETH
	// rethAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketRewardsPool from v1.0.0
	// v1_0_0_RewardsPoolAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketClaimNode from v1.0.0
	// v1_0_0_ClaimNodeAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketClaimTrustedNode from v1.0.0
	// v1_0_0_ClaimTrustedNodeAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketMinipoolManager from v1.0.0
	// v1_0_0_MinipoolManagerAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketNetworkPrices from v1.1.0
	// v1_1_0_NetworkPricesAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketNodeStaking from v1.1.0
	// v1_1_0_NodeStakingAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketNodeDeposit from v1.1.0
	// v1_1_0_NodeDepositAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketMinipoolQueue from v1.1.0
	// v1_1_0_MinipoolQueueAddress map[Network]string `yaml:"-"`

	// // The contract address of rocketMinipoolFactory from v1.1.0
	// v1_1_0_MinipoolFactoryAddress map[Network]string `yaml:"-"`

	// // Addresses for RocketRewardsPool that have been upgraded during development
	// previousRewardsPoolAddresses map[Network][]common.Address `yaml:"-"`

	// // The RocketOvmPriceMessenger Optimism address for each network
	// optimismPriceMessengerAddress map[Network]string `yaml:"-"`

	// // The RocketPolygonPriceMessenger Polygon address for each network
	// polygonPriceMessengerAddress map[Network]string `yaml:"-"`

	// // The RocketArbitumPriceMessenger Arbitrum address for each network
	// arbitrumPriceMessengerAddress map[Network]string `yaml:"-"`

	// // The RocketZkSyncPriceMessenger zkSyncEra address for each network
	// zkSyncEraPriceMessengerAddress map[Network]string `yaml:"-"`

	// // The RocketBasePriceMessenger Base address for each network
	// basePriceMessengerAddress map[Network]string `yaml:"-"`

	// // The UniswapV3 pool address for each network (used for RPL price TWAP info)
	// rplTwapPoolAddress map[Network]string `yaml:"-"`

	// // The multicall contract address
	// multicallAddress map[Network]string `yaml:"-"`

	// // The BalanceChecker contract address
	// balancebatcherAddress map[Network]string `yaml:"-"`

	// // The FlashBots Protect RPC endpoint
	// flashbotsProtectUrl map[Network]string `yaml:"-"`
}
