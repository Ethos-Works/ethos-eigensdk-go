package avsregistry

import (
	"context"
	"errors"
	"math"
	"math/big"

	"github.com/ethos-works/ethos-eigensdk-go/chainio/clients/eth"
	"github.com/ethos-works/ethos-eigensdk-go/logging"
	"github.com/ethos-works/ethos-eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethcommon "github.com/ethereum/go-ethereum/common"

	contractOperatorStateRetriever "github.com/ethos-works/ethos-eigensdk-go/contracts/bindings/OperatorStateRetriever"
	opstateretriever "github.com/ethos-works/ethos-eigensdk-go/contracts/bindings/OperatorStateRetriever"
	regcoord "github.com/ethos-works/ethos-eigensdk-go/contracts/bindings/RegistryCoordinator"
	stakeregistry "github.com/ethos-works/ethos-eigensdk-go/contracts/bindings/StakeRegistry"
)

type AvsRegistryReader interface {
	GetQuorumCount(opts *bind.CallOpts) (uint8, error)

	GetOperatorsStakeInQuorumsAtCurrentBlock(
		opts *bind.CallOpts,
		quorumNumbers []byte,
	) ([][]opstateretriever.OperatorStateRetrieverOperator, error)

	GetOperatorsStakeInQuorumsAtBlock(
		opts *bind.CallOpts,
		quorumNumbers []byte,
		blockNumber uint32,
	) ([][]opstateretriever.OperatorStateRetrieverOperator, error)

	GetOperatorAddrsInQuorumsAtCurrentBlock(
		opts *bind.CallOpts,
		quorumNumbers []byte,
	) ([][]common.Address, error)

	GetOperatorsStakeInQuorumsOfOperatorAtBlock(
		opts *bind.CallOpts,
		operator common.Address,
		blockNumber uint32,
	) ([]types.QuorumNum, [][]opstateretriever.OperatorStateRetrieverOperator, error)

	GetOperatorsStakeInQuorumsOfOperatorAtCurrentBlock(
		opts *bind.CallOpts,
		operator common.Address,
	) ([]types.QuorumNum, [][]opstateretriever.OperatorStateRetrieverOperator, error)

	GetOperatorStakeInQuorumsOfOperatorAtCurrentBlock(
		opts *bind.CallOpts,
		operator common.Address,
	) (map[types.QuorumNum]types.StakeAmount, error)

	IsOperatorRegistered(opts *bind.CallOpts, operatorAddress gethcommon.Address) (bool, error)
}

type AvsRegistryChainReader struct {
	logger                  logging.Logger
	registryCoordinatorAddr gethcommon.Address
	registryCoordinator     *regcoord.ContractRegistryCoordinator
	operatorStateRetriever  *opstateretriever.ContractOperatorStateRetriever
	stakeRegistry           *stakeregistry.ContractStakeRegistry
	ethClient               eth.EthClient
}

// forces AvsReader to implement the clients.ReaderInterface interface
var _ AvsRegistryReader = (*AvsRegistryChainReader)(nil)

func NewAvsRegistryChainReader(
	registryCoordinatorAddr gethcommon.Address,
	registryCoordinator *regcoord.ContractRegistryCoordinator,
	operatorStateRetriever *opstateretriever.ContractOperatorStateRetriever,
	stakeRegistry *stakeregistry.ContractStakeRegistry,
	logger logging.Logger,
	ethClient eth.EthClient,
) *AvsRegistryChainReader {
	return &AvsRegistryChainReader{
		registryCoordinatorAddr: registryCoordinatorAddr,
		registryCoordinator:     registryCoordinator,
		operatorStateRetriever:  operatorStateRetriever,
		stakeRegistry:           stakeRegistry,
		logger:                  logger,
		ethClient:               ethClient,
	}
}

func BuildAvsRegistryChainReader(
	registryCoordinatorAddr gethcommon.Address,
	operatorStateRetrieverAddr gethcommon.Address,
	ethClient eth.EthClient,
	logger logging.Logger,
) (*AvsRegistryChainReader, error) {
	contractRegistryCoordinator, err := regcoord.NewContractRegistryCoordinator(registryCoordinatorAddr, ethClient)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to create contractRegistryCoordinator"), err)
	}
	stakeRegistryAddr, err := contractRegistryCoordinator.StakeRegistry(&bind.CallOpts{})
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to get stakeRegistryAddr"), err)
	}
	contractStakeRegistry, err := stakeregistry.NewContractStakeRegistry(stakeRegistryAddr, ethClient)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to create contractStakeRegistry"), err)
	}
	contractOperatorStateRetriever, err := contractOperatorStateRetriever.NewContractOperatorStateRetriever(
		operatorStateRetrieverAddr,
		ethClient,
	)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to create contractOperatorStateRetriever"), err)
	}
	return NewAvsRegistryChainReader(
		registryCoordinatorAddr,
		contractRegistryCoordinator,
		contractOperatorStateRetriever,
		contractStakeRegistry,
		logger,
		ethClient,
	), nil
}

func (r *AvsRegistryChainReader) GetQuorumCount(opts *bind.CallOpts) (uint8, error) {
	return r.registryCoordinator.QuorumCount(opts)
}

func (r *AvsRegistryChainReader) GetOperatorsStakeInQuorumsAtCurrentBlock(
	opts *bind.CallOpts,
	quorumNumbers []byte,
) ([][]opstateretriever.OperatorStateRetrieverOperator, error) {
	if opts.Context == nil {
		opts.Context = context.Background()
	}
	curBlock, err := r.ethClient.BlockNumber(opts.Context)
	if err != nil {
		return nil, types.WrapError(errors.New("Cannot get current block number"), err)
	}
	if curBlock > math.MaxUint32 {
		return nil, types.WrapError(errors.New("Current block number is too large to be converted to uint32"), err)
	}
	return r.GetOperatorsStakeInQuorumsAtBlock(opts, quorumNumbers, uint32(curBlock))
}

// the contract stores historical state, so blockNumber should be the block number of the state you want to query
// and the blockNumber in opts should be the block number of the latest block (or set to nil, which is equivalent)
func (r *AvsRegistryChainReader) GetOperatorsStakeInQuorumsAtBlock(
	opts *bind.CallOpts,
	quorumNumbers []byte,
	blockNumber uint32,
) ([][]opstateretriever.OperatorStateRetrieverOperator, error) {
	operatorStakes, err := r.operatorStateRetriever.GetOperatorState0(
		opts,
		r.registryCoordinatorAddr,
		quorumNumbers,
		blockNumber)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to get operators state"), err)
	}
	return operatorStakes, nil
}

func (r *AvsRegistryChainReader) GetOperatorAddrsInQuorumsAtCurrentBlock(
	opts *bind.CallOpts,
	quorumNumbers []byte,
) ([][]common.Address, error) {
	if opts.Context == nil {
		opts.Context = context.Background()
	}
	curBlock, err := r.ethClient.BlockNumber(opts.Context)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to get current block number"), err)
	}
	if curBlock > math.MaxUint32 {
		return nil, types.WrapError(errors.New("Current block number is too large to be converted to uint32"), err)
	}
	operatorStakes, err := r.operatorStateRetriever.GetOperatorState0(
		opts,
		r.registryCoordinatorAddr,
		quorumNumbers,
		uint32(curBlock),
	)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to get operators state"), err)
	}
	var quorumOperatorAddrs [][]common.Address
	for _, quorum := range operatorStakes {
		var operatorAddrs []common.Address
		for _, operator := range quorum {
			operatorAddrs = append(operatorAddrs, operator.Operator)
		}
		quorumOperatorAddrs = append(quorumOperatorAddrs, operatorAddrs)
	}
	return quorumOperatorAddrs, nil

}

func (r *AvsRegistryChainReader) GetOperatorsStakeInQuorumsOfOperatorAtBlock(
	opts *bind.CallOpts,
	operator common.Address,
	blockNumber uint32,
) ([]types.QuorumNum, [][]opstateretriever.OperatorStateRetrieverOperator, error) {
	quorumBitmap, operatorStakes, err := r.operatorStateRetriever.GetOperatorState(
		opts,
		r.registryCoordinatorAddr,
		operator,
		blockNumber)
	if err != nil {
		return nil, nil, types.WrapError(errors.New("Failed to get operators state"), err)
	}
	quorums := types.BitmapToQuorumIds(quorumBitmap)
	return quorums, operatorStakes, nil
}

// opts will be modified to have the latest blockNumber, so make sure not to reuse it
// blockNumber in opts will be ignored, and the chain will be queried to get the latest blockNumber
func (r *AvsRegistryChainReader) GetOperatorsStakeInQuorumsOfOperatorAtCurrentBlock(
	opts *bind.CallOpts,
	operator common.Address,
) ([]types.QuorumNum, [][]opstateretriever.OperatorStateRetrieverOperator, error) {
	if opts.Context == nil {
		opts.Context = context.Background()
	}
	curBlock, err := r.ethClient.BlockNumber(opts.Context)
	if err != nil {
		return nil, nil, types.WrapError(errors.New("Failed to get current block number"), err)
	}
	if curBlock > math.MaxUint32 {
		return nil, nil, types.WrapError(errors.New("Current block number is too large to be converted to uint32"), err)
	}
	opts.BlockNumber = big.NewInt(int64(curBlock))
	return r.GetOperatorsStakeInQuorumsOfOperatorAtBlock(opts, operator, uint32(curBlock))
}

// GetOperatorStakeInQuorumsOfOperatorAtCurrentBlock could have race conditions
// it currently makes a bunch of calls to fetch "current block" information,
// so some of them could actually return information from different blocks
func (r *AvsRegistryChainReader) GetOperatorStakeInQuorumsOfOperatorAtCurrentBlock(
	opts *bind.CallOpts,
	operator common.Address,
) (map[types.QuorumNum]types.StakeAmount, error) {
	quorumBitmap, err := r.registryCoordinator.GetCurrentQuorumBitmap(opts, operator)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to get operator quorums"), err)
	}
	quorums := types.BitmapToQuorumIds(quorumBitmap)
	quorumStakes := make(map[types.QuorumNum]types.StakeAmount)
	for _, quorum := range quorums {
		stake, err := r.stakeRegistry.GetCurrentStake(
			&bind.CallOpts{},
			operator,
			quorum,
		)
		if err != nil {
			return nil, types.WrapError(errors.New("Failed to get operator stake"), err)
		}
		quorumStakes[quorum] = stake
	}
	return quorumStakes, nil
}

func (r *AvsRegistryChainReader) IsOperatorRegistered(
	opts *bind.CallOpts,
	operatorAddress gethcommon.Address,
) (bool, error) {
	operatorStatus, err := r.registryCoordinator.GetOperatorStatus(opts, operatorAddress)
	if err != nil {
		return false, types.WrapError(errors.New("Failed to get operator status"), err)
	}

	// 0 = NEVER_REGISTERED, 1 = REGISTERED, 2 = DEREGISTERED
	registeredWithAvs := operatorStatus == 1
	return registeredWithAvs, nil
}
