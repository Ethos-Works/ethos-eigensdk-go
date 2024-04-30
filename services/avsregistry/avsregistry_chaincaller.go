package avsregistry

import (
	"context"
	"errors"
	"math/big"

	avsregistry "github.com/ethos-works/ethos-eigensdk-go/chainio/clients/avsregistry"
	"github.com/ethos-works/ethos-eigensdk-go/logging"
	"github.com/ethos-works/ethos-eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// AvsRegistryServiceChainCaller is a wrapper around AvsRegistryReader that transforms the data into
// nicer golang types that are easier to work with
type AvsRegistryServiceChainCaller struct {
	avsregistry.AvsRegistryReader
	logger logging.Logger
}

var _ AvsRegistryService = (*AvsRegistryServiceChainCaller)(nil)

func NewAvsRegistryServiceChainCaller(avsRegistryReader avsregistry.AvsRegistryReader, logger logging.Logger) *AvsRegistryServiceChainCaller {
	return &AvsRegistryServiceChainCaller{
		AvsRegistryReader: avsRegistryReader,
		logger:            logger,
	}
}

func (ar *AvsRegistryServiceChainCaller) GetOperatorsAvsStateAtBlock(ctx context.Context, quorumNumbers []types.QuorumNum, blockNumber types.BlockNum) (map[types.OperatorAddr]types.OperatorAvsState, error) {
	operatorsAvsState := make(map[types.OperatorAddr]types.OperatorAvsState)
	// Get operator state for each quorum by querying BLSOperatorStateRetriever (this call is why this service implementation is called ChainCaller)
	operatorsStakesInQuorums, err := ar.AvsRegistryReader.GetOperatorsStakeInQuorumsAtBlock(&bind.CallOpts{Context: ctx}, quorumNumbers, blockNumber)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to get operator state"), err)
	}
	numquorums := len(quorumNumbers)
	if len(operatorsStakesInQuorums) != numquorums {
		ar.logger.Fatal("Number of quorums returned from GetOperatorsStakeInQuorumsAtBlock does not match number of quorums requested. Probably pointing to old contract or wrong implementation.", "service", "AvsRegistryServiceChainCaller")
	}

	for quorumIdx, quorumNum := range quorumNumbers {
		for _, operator := range operatorsStakesInQuorums[quorumIdx] {
			if err != nil {
				return nil, types.WrapError(errors.New("Failed to find pubkeys for operator while building operatorsAvsState"), err)
			}
			if operatorAvsState, ok := operatorsAvsState[operator.Operator]; ok {
				operatorAvsState.StakePerQuorum[quorumNum] = operator.Stake
			} else {
				stakePerQuorum := make(map[types.QuorumNum]types.StakeAmount)
				stakePerQuorum[quorumNum] = operator.Stake
				operatorsAvsState[operator.Operator] = types.OperatorAvsState{
					Operator:       operator.Operator,
					StakePerQuorum: stakePerQuorum,
					BlockNumber:    blockNumber,
				}
				operatorsAvsState[operator.Operator].StakePerQuorum[quorumNum] = operator.Stake
			}
		}
	}

	return operatorsAvsState, nil
}

func (ar *AvsRegistryServiceChainCaller) GetQuorumsAvsStateAtBlock(ctx context.Context, quorumNumbers []types.QuorumNum, blockNumber types.BlockNum) (map[types.QuorumNum]types.QuorumAvsState, error) {
	operatorsAvsState, err := ar.GetOperatorsAvsStateAtBlock(ctx, quorumNumbers, blockNumber)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to get quorum state"), err)
	}
	quorumsAvsState := make(map[types.QuorumNum]types.QuorumAvsState)
	for _, quorumNum := range quorumNumbers {
		totalStake := big.NewInt(0)
		for _, operator := range operatorsAvsState {
			// only include operators that have a stake in this quorum
			if stake, ok := operator.StakePerQuorum[quorumNum]; ok {
				totalStake.Add(totalStake, stake)
			}
		}
		quorumsAvsState[quorumNum] = types.QuorumAvsState{
			QuorumNumber: quorumNum,
			TotalStake:   totalStake,
			BlockNumber:  blockNumber,
		}
	}
	return quorumsAvsState, nil
}
