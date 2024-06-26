package avsregistry

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	chainiomocks "github.com/ethos-works/ethos-eigensdk-go/chainio/mocks"
	opstateretrievar "github.com/ethos-works/ethos-eigensdk-go/contracts/bindings/OperatorStateRetriever"
	"github.com/ethos-works/ethos-eigensdk-go/logging"
	"github.com/ethos-works/ethos-eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/mock/gomock"
)

type testOperator struct {
	operatorAddr common.Address
}

func TestAvsRegistryServiceChainCaller_GetOperatorsAvsState(t *testing.T) {
	logger := logging.NewNoopLogger()
	testOperator := testOperator{
		operatorAddr: common.HexToAddress("0x1"),
	}

	var tests = []struct {
		name                      string
		mocksInitializationFunc   func(*chainiomocks.MockAvsRegistryReader)
		queryQuorumNumbers        []types.QuorumNum
		queryBlockNum             types.BlockNum
		wantErr                   error
		wantOperatorsAvsStateDict map[types.OperatorAddr]types.OperatorAvsState
	}{
		{
			name: "should return operatorsAvsState",
			mocksInitializationFunc: func(mockAvsRegistryReader *chainiomocks.MockAvsRegistryReader) {
				mockAvsRegistryReader.EXPECT().GetOperatorsStakeInQuorumsAtBlock(gomock.Any(), []types.QuorumNum{1}, types.BlockNum(1)).Return([][]opstateretrievar.OperatorStateRetrieverOperator{
					{
						{
							Operator: testOperator.operatorAddr,
							Stake:    big.NewInt(123),
						},
					},
				}, nil)
			},
			queryQuorumNumbers: []types.QuorumNum{1},
			queryBlockNum:      1,
			wantErr:            nil,
			wantOperatorsAvsStateDict: map[types.OperatorAddr]types.OperatorAvsState{
				testOperator.operatorAddr: {
					Operator:       testOperator.operatorAddr,
					StakePerQuorum: map[types.QuorumNum]types.StakeAmount{1: big.NewInt(123)},
					BlockNumber:    1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockCtrl := gomock.NewController(t)
			mockAvsRegistryReader := chainiomocks.NewMockAvsRegistryReader(mockCtrl)

			if tt.mocksInitializationFunc != nil {
				tt.mocksInitializationFunc(mockAvsRegistryReader)
			}
			// Create a new instance of the avsregistry service
			service := NewAvsRegistryServiceChainCaller(mockAvsRegistryReader, logger)

			// Call the GetOperatorPubkeys method with the test operator address
			gotOperatorsAvsStateDict, gotErr := service.GetOperatorsAvsStateAtBlock(context.Background(), tt.queryQuorumNumbers, tt.queryBlockNum)
			if tt.wantErr != gotErr {
				t.Fatalf("GetOperatorsAvsState returned wrong error. Got: %v, want: %v.", gotErr, tt.wantErr)
			}
			if tt.wantErr == nil && !reflect.DeepEqual(tt.wantOperatorsAvsStateDict, gotOperatorsAvsStateDict) {
				t.Fatalf("GetOperatorsAvsState returned wrong operatorsAvsStateDict. Got: %v, want: %v.", gotOperatorsAvsStateDict, tt.wantOperatorsAvsStateDict)
			}
		})
	}
}

func TestAvsRegistryServiceChainCaller_GetQuorumsAvsState(t *testing.T) {
	logger := logging.NewNoopLogger()
	testOperator := testOperator{
		operatorAddr: common.HexToAddress("0x1"),
	}

	var tests = []struct {
		name                    string
		mocksInitializationFunc func(*chainiomocks.MockAvsRegistryReader)
		queryQuorumNumbers      []types.QuorumNum
		queryBlockNum           types.BlockNum
		wantErr                 error
		wantQuorumsAvsStateDict map[types.QuorumNum]types.QuorumAvsState
	}{
		{
			name: "should return operatorsAvsState",
			mocksInitializationFunc: func(mockAvsRegistryReader *chainiomocks.MockAvsRegistryReader) {
				mockAvsRegistryReader.EXPECT().GetOperatorsStakeInQuorumsAtBlock(gomock.Any(), []types.QuorumNum{1}, types.BlockNum(1)).Return([][]opstateretrievar.OperatorStateRetrieverOperator{
					{
						{
							Operator: testOperator.operatorAddr,
							Stake:    big.NewInt(123),
						},
					},
				}, nil)
			},
			queryQuorumNumbers: []types.QuorumNum{1},
			queryBlockNum:      1,
			wantErr:            nil,
			wantQuorumsAvsStateDict: map[types.QuorumNum]types.QuorumAvsState{
				1: types.QuorumAvsState{
					QuorumNumber: types.QuorumNum(1),
					TotalStake:   big.NewInt(123),
					BlockNumber:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockCtrl := gomock.NewController(t)
			mockAvsRegistryReader := chainiomocks.NewMockAvsRegistryReader(mockCtrl)

			if tt.mocksInitializationFunc != nil {
				tt.mocksInitializationFunc(mockAvsRegistryReader)
			}
			// Create a new instance of the avsregistry service
			service := NewAvsRegistryServiceChainCaller(mockAvsRegistryReader, logger)

			// Call the GetOperatorPubkeys method with the test operator address
			aggG1PubkeyPerQuorum, gotErr := service.GetQuorumsAvsStateAtBlock(context.Background(), tt.queryQuorumNumbers, tt.queryBlockNum)
			if tt.wantErr != gotErr {
				t.Fatalf("GetOperatorsAvsState returned wrong error. Got: %v, want: %v.", gotErr, tt.wantErr)
			}
			if tt.wantErr == nil && !reflect.DeepEqual(tt.wantQuorumsAvsStateDict, aggG1PubkeyPerQuorum) {
				t.Fatalf("GetOperatorsAvsState returned wrong aggG1PubkeyPerQuorum. Got: %v, want: %v.", aggG1PubkeyPerQuorum, tt.wantQuorumsAvsStateDict)
			}
		})
	}
}
