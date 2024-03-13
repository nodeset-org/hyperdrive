package service

// import (
// 	"fmt"

// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/contract"
// )

// const (
// 	UserSettingPath string = ".rocketpool/user-settings.yml"
// )

// // type ValidatorAccountCreator interface {
// // 	CreateNewValidatorAccount(nodeOperator common.Address) (*common.Hash, error)
// // }

// type RewardsClaimService struct {
// 	ContractService *contract.Contract
// }

// func NewRewardsClaimService(contractService *contract.Contract) *RewardsClaimService {
// 	return &RewardsClaimService{
// 		ContractService: contractService,
// 	}
// }

// func (vs *RewardsClaimService) ClaimRewards(address common.Address) (*common.Hash, error) {
// 	tx, err := vs.ContractService.Transact(&bind.TransactOpts{}, "harvest", address)

// 	if err != nil {
// 		return nil, err
// 	}

// 	tx_hash := tx.Hash()
// 	fmt.Printf("Transaction hash: %s\n", tx_hash.String())
// 	return &tx_hash, nil
// }
