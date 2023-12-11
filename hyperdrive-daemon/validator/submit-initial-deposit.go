package validator

import (
	"fmt"
	"math/big"

	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/services/beacon"
	"github.com/rocket-pool/rocketpool-go/node"
)

type DepositContext struct {
	// cfg *config.RocketPoolConfig
	// rp *rocketpool.RocketPool
	Bc beacon.Client
	// w   *wallet.LocalWallet

	// amount     *big.Int
	// minNodeFee float64
	Salt *big.Int
	Node *node.Node
	// pSettings  *protocol.ProtocolDaoSettings
	// oSettings  *oracle.OracleDaoSettings
	// mpMgr      *minipool.MinipoolManager
}

func (c *DepositContext) SubmitInitDeposit() error {
	// Initial population
	// Where do I get credit balance? Can I assume 0?
	// data.CreditBalance = c.node.Credit.Get()
	// Same with these fields
	// data.DepositDisabled = !c.pSettings.Node.IsDepositingEnabled.Get()
	// data.DepositBalance = c.depositPool.Balance.Get()

	// Can prob delete this
	// data.ScrubPeriod = c.oSettings.Minipool.ScrubPeriod.Formatted()

	// Get Beacon config
	eth2Config, err := c.Bc.GetEth2Config()
	fmt.Printf("eth2Config: %v\n", eth2Config)

	if err != nil {
		return fmt.Errorf("error getting Beacon config: %w", err)
	}

	// Adjust the salt
	if c.Salt.Cmp(big.NewInt(0)) == 0 {
		fmt.Printf("salt is 0\n")
		// nonce, err := c.rp.Client.NonceAt(context.Background(), c.Node.Address, nil)
		// 	if err != nil {
		// 		return fmt.Errorf("error getting node's latest nonce: %w", err)
		// 	}
		// 	c.salt.SetUint64(nonce)
	}

	// // Check node balance
	// data.NodeBalance, err = c.rp.Client.BalanceAt(context.Background(), c.node.Address, nil)
	// if err != nil {
	// 	return fmt.Errorf("error getting node's ETH balance: %w", err)
	// }

	// // Check the node's collateral
	// collateral, err := collateral.CheckCollateral(c.rp, c.node.Address, nil)
	// if err != nil {
	// 	return fmt.Errorf("error checking node collateral: %w", err)
	// }
	// ethMatched := collateral.EthMatched
	// ethMatchedLimit := collateral.EthMatchedLimit
	// pendingMatchAmount := collateral.PendingMatchAmount

	// // Check for insufficient balance
	// totalBalance := big.NewInt(0).Add(data.NodeBalance, data.CreditBalance)
	// data.InsufficientBalance = (c.amount.Cmp(totalBalance) > 0)

	// // Check if the credit balance can be used
	// data.CanUseCredit = (data.DepositBalance.Cmp(eth.EthToWei(1)) >= 0)

	// // Check data
	// validatorEthWei := eth.EthToWei(ValidatorEth)
	// matchRequest := big.NewInt(0).Sub(validatorEthWei, c.amount)
	// availableToMatch := big.NewInt(0).Sub(ethMatchedLimit, ethMatched)
	// availableToMatch.Sub(availableToMatch, pendingMatchAmount)
	// data.InsufficientRplStake = (availableToMatch.Cmp(matchRequest) == -1)

	// // Update response
	// data.CanDeposit = !(data.InsufficientBalance || data.InsufficientRplStake || data.InvalidAmount || data.DepositDisabled)
	// if data.CanDeposit && !data.CanUseCredit && data.NodeBalance.Cmp(c.amount) < 0 {
	// 	// Can't use credit and there's not enough ETH in the node wallet to deposit so error out
	// 	data.InsufficientBalanceWithoutCredit = true
	// 	data.CanDeposit = false
	// }

	// // Return if depositing won't work
	// if !data.CanDeposit {
	// 	return nil
	// }

	// // Make sure ETH2 is on the correct chain
	// depositContractInfo, err := rputils.GetDepositContractInfo(c.rp, c.cfg, c.bc)
	// if err != nil {
	// 	return fmt.Errorf("error verifying the EL and BC are on the same chain: %w", err)
	// }
	// if depositContractInfo.RPNetwork != depositContractInfo.BeaconNetwork ||
	// 	depositContractInfo.RPDepositContract != depositContractInfo.BeaconDepositContract {
	// 	return fmt.Errorf("FATAL: Beacon network mismatch! Expected %s on chain %d, but beacon is using %s on chain %d.",
	// 		depositContractInfo.RPDepositContract.Hex(),
	// 		depositContractInfo.RPNetwork,
	// 		depositContractInfo.BeaconDepositContract.Hex(),
	// 		depositContractInfo.BeaconNetwork)
	// }

	// // Get how much credit to use
	// if data.CanUseCredit {
	// 	remainingAmount := big.NewInt(0).Sub(c.amount, data.CreditBalance)
	// 	if remainingAmount.Cmp(big.NewInt(0)) > 0 {
	// 		// Send the remaining amount if the credit isn't enough to cover the whole deposit
	// 		opts.Value = remainingAmount
	// 	}
	// } else {
	// 	opts.Value = c.amount
	// }

	// // Get the next available validator key without saving it
	// validatorKey, index, err := c.w.GetNextValidatorKey()
	// if err != nil {
	// 	return fmt.Errorf("error getting next available validator key: %w", err)
	// }
	// data.Index = index

	// // Get the next minipool address
	// var minipoolAddress common.Address
	// err = c.rp.Query(func(mc *batch.MultiCaller) error {
	// 	c.node.GetExpectedMinipoolAddress(mc, &minipoolAddress, c.salt)
	// 	return nil
	// }, nil)
	// if err != nil {
	// 	return fmt.Errorf("error getting expected minipool address: %w", err)
	// }
	// data.MinipoolAddress = minipoolAddress

	// // Get the withdrawal credentials
	// var withdrawalCredentials common.Hash
	// err = c.rp.Query(func(mc *batch.MultiCaller) error {
	// 	c.mpMgr.GetMinipoolWithdrawalCredentials(mc, &withdrawalCredentials, minipoolAddress)
	// 	return nil
	// }, nil)
	// if err != nil {
	// 	return fmt.Errorf("error getting minipool withdrawal credentials: %w", err)
	// }

	// // Get validator deposit data and associated parameters
	// depositAmount := uint64(1e9) // 1 ETH in gwei
	// depositData, depositDataRoot, err := validator.GetDepositData(validatorKey, withdrawalCredentials, eth2Config, depositAmount)
	// if err != nil {
	// 	return fmt.Errorf("error getting deposit data: %w", err)
	// }
	// pubkey := rptypes.BytesToValidatorPubkey(depositData.PublicKey)
	// signature := rptypes.BytesToValidatorSignature(depositData.Signature)
	// data.ValidatorPubkey = pubkey

	// // Make sure a validator with this pubkey doesn't already exist
	// status, err := c.bc.GetValidatorStatus(pubkey, nil)
	// if err != nil {
	// 	return fmt.Errorf("Error checking for existing validator status: %w\nYour funds have not been deposited for your own safety.", err)
	// }
	// if status.Exists {
	// 	return fmt.Errorf("**** ALERT ****\n"+
	// 		"Your minipool %s has the following as a validator pubkey:\n\t%s\n"+
	// 		"This key is already in use by validator %d on the Beacon chain!\n"+
	// 		"Rocket Pool will not allow you to deposit this validator for your own safety so you do not get slashed.\n"+
	// 		"PLEASE REPORT THIS TO THE ROCKET POOL DEVELOPERS.\n"+
	// 		"***************\n", minipoolAddress.Hex(), pubkey.Hex(), status.Index)
	// }

	// // Do a final sanity check
	// err = validateDepositInfo(eth2Config, uint64(depositAmount), pubkey, withdrawalCredentials, signature)
	// if err != nil {
	// 	return fmt.Errorf("FATAL: Your deposit failed the validation safety check: %w\n"+
	// 		"For your safety, this deposit will not be submitted and your ETH will not be staked.\n"+
	// 		"PLEASE REPORT THIS TO THE ROCKET POOL DEVELOPERS and include the following information:\n"+
	// 		"\tDomain Type: 0x%s\n"+
	// 		"\tGenesis Fork Version: 0x%s\n"+
	// 		"\tGenesis Validator Root: 0x%s\n"+
	// 		"\tDeposit Amount: %d gwei\n"+
	// 		"\tValidator Pubkey: %s\n"+
	// 		"\tWithdrawal Credentials: %s\n"+
	// 		"\tSignature: %s\n",
	// 		err,
	// 		hex.EncodeToString(eth2types.DomainDeposit[:]),
	// 		hex.EncodeToString(eth2Config.GenesisForkVersion),
	// 		hex.EncodeToString(eth2types.ZeroGenesisValidatorsRoot),
	// 		depositAmount,
	// 		pubkey.Hex(),
	// 		withdrawalCredentials.Hex(),
	// 		signature.Hex(),
	// 	)
	// }

	// // Get tx info
	// var txInfo *core.TransactionInfo
	// var funcName string
	// if data.CanUseCredit {
	// 	txInfo, err = c.node.DepositWithCredit(c.amount, c.minNodeFee, pubkey, signature, depositDataRoot, c.salt, minipoolAddress, opts)
	// 	funcName = "DepositWithCredit"
	// } else {
	// 	txInfo, err = c.node.Deposit(c.amount, c.minNodeFee, pubkey, signature, depositDataRoot, c.salt, minipoolAddress, opts)
	// 	funcName = "Deposit"
	// }
	// if err != nil {
	// 	return fmt.Errorf("error getting TX info for %s: %w", funcName, err)
	// }
	// data.TxInfo = txInfo

	return nil
}
