package avsregistry

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethos-works/ethos-eigensdk-go/types"
)

type FakeAvsRegistryService struct {
	operators map[types.BlockNum]map[types.OperatorAddr]types.OperatorAvsState
}

func NewFakeAvsRegistryService(blockNum types.BlockNum, operators []types.TestOperator) *FakeAvsRegistryService {
	fakeAvsRegistryService := &FakeAvsRegistryService{
		operators: map[types.BlockNum]map[types.OperatorAddr]types.OperatorAvsState{
			blockNum: {},
		},
	}
	for _, operator := range operators {
		fakeAvsRegistryService.operators[blockNum][operator.Operator] = types.OperatorAvsState{
			Operator:       operator.Operator,
			StakePerQuorum: operator.StakePerQuorum,
			BlockNumber:    blockNum,
		}
	}
	return fakeAvsRegistryService
}

var _ AvsRegistryService = (*FakeAvsRegistryService)(nil)

func (f *FakeAvsRegistryService) GetOperatorsAvsStateAtBlock(ctx context.Context, quorumNumbers []types.QuorumNum, blockNumber types.BlockNum) (map[types.OperatorAddr]types.OperatorAvsState, error) {
	operatorsAvsState, ok := f.operators[blockNumber]
	if !ok {
		return nil, errors.New("block number not found")
	}
	return operatorsAvsState, nil
}

func (f *FakeAvsRegistryService) GetQuorumsAvsStateAtBlock(ctx context.Context, quorumNumbers []types.QuorumNum, blockNumber types.BlockNum) (map[types.QuorumNum]types.QuorumAvsState, error) {
	operatorsAvsState, ok := f.operators[blockNumber]
	if !ok {
		return nil, errors.New("block number not found")
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
