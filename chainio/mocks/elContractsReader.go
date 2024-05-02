// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ethos-works/ethos-eigensdk-go/chainio/clients/elcontracts (interfaces: ELReader)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/elContractsReader.go -package=mocks github.com/ethos-works/ethos-eigensdk-go/chainio/clients/elcontracts ELReader
//
// Package mocks is a generated GoMock package.
package mocks

import (
	big "math/big"
	reflect "reflect"

	contractIERC20 "github.com/ethos-works/ethos-eigensdk-go/contracts/bindings/IERC20"
	contractIStrategy "github.com/ethos-works/ethos-eigensdk-go/contracts/bindings/IStrategy"
	types "github.com/ethos-works/ethos-eigensdk-go/types"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	common "github.com/ethereum/go-ethereum/common"
	gomock "go.uber.org/mock/gomock"
)

// MockELReader is a mock of ELReader interface.
type MockELReader struct {
	ctrl     *gomock.Controller
	recorder *MockELReaderMockRecorder
}

// MockELReaderMockRecorder is the mock recorder for MockELReader.
type MockELReaderMockRecorder struct {
	mock *MockELReader
}

// NewMockELReader creates a new mock instance.
func NewMockELReader(ctrl *gomock.Controller) *MockELReader {
	mock := &MockELReader{ctrl: ctrl}
	mock.recorder = &MockELReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockELReader) EXPECT() *MockELReaderMockRecorder {
	return m.recorder
}

// CalculateDelegationApprovalDigestHash mocks base method.
func (m *MockELReader) CalculateDelegationApprovalDigestHash(arg0 *bind.CallOpts, arg1, arg2, arg3 common.Address, arg4 [32]byte, arg5 *big.Int) ([32]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateDelegationApprovalDigestHash", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([32]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateDelegationApprovalDigestHash indicates an expected call of CalculateDelegationApprovalDigestHash.
func (mr *MockELReaderMockRecorder) CalculateDelegationApprovalDigestHash(arg0, arg1, arg2, arg3, arg4, arg5 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateDelegationApprovalDigestHash", reflect.TypeOf((*MockELReader)(nil).CalculateDelegationApprovalDigestHash), arg0, arg1, arg2, arg3, arg4, arg5)
}

// CalculateOperatorAVSRegistrationDigestHash mocks base method.
func (m *MockELReader) CalculateOperatorAVSRegistrationDigestHash(arg0 *bind.CallOpts, arg1, arg2 common.Address, arg3 [32]byte, arg4 *big.Int) ([32]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateOperatorAVSRegistrationDigestHash", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([32]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateOperatorAVSRegistrationDigestHash indicates an expected call of CalculateOperatorAVSRegistrationDigestHash.
func (mr *MockELReaderMockRecorder) CalculateOperatorAVSRegistrationDigestHash(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateOperatorAVSRegistrationDigestHash", reflect.TypeOf((*MockELReader)(nil).CalculateOperatorAVSRegistrationDigestHash), arg0, arg1, arg2, arg3, arg4)
}

// GetOperatorDetails mocks base method.
func (m *MockELReader) GetOperatorDetails(arg0 *bind.CallOpts, arg1 types.Operator) (types.Operator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOperatorDetails", arg0, arg1)
	ret0, _ := ret[0].(types.Operator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOperatorDetails indicates an expected call of GetOperatorDetails.
func (mr *MockELReaderMockRecorder) GetOperatorDetails(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOperatorDetails", reflect.TypeOf((*MockELReader)(nil).GetOperatorDetails), arg0, arg1)
}

// GetOperatorSharesInStrategy mocks base method.
func (m *MockELReader) GetOperatorSharesInStrategy(arg0 *bind.CallOpts, arg1, arg2 common.Address) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOperatorSharesInStrategy", arg0, arg1, arg2)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOperatorSharesInStrategy indicates an expected call of GetOperatorSharesInStrategy.
func (mr *MockELReaderMockRecorder) GetOperatorSharesInStrategy(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOperatorSharesInStrategy", reflect.TypeOf((*MockELReader)(nil).GetOperatorSharesInStrategy), arg0, arg1, arg2)
}

// GetStrategyAndUnderlyingERC20Token mocks base method.
func (m *MockELReader) GetStrategyAndUnderlyingERC20Token(arg0 *bind.CallOpts, arg1 common.Address) (*contractIStrategy.ContractIStrategy, contractIERC20.ContractIERC20Methods, common.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStrategyAndUnderlyingERC20Token", arg0, arg1)
	ret0, _ := ret[0].(*contractIStrategy.ContractIStrategy)
	ret1, _ := ret[1].(contractIERC20.ContractIERC20Methods)
	ret2, _ := ret[2].(common.Address)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetStrategyAndUnderlyingERC20Token indicates an expected call of GetStrategyAndUnderlyingERC20Token.
func (mr *MockELReaderMockRecorder) GetStrategyAndUnderlyingERC20Token(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStrategyAndUnderlyingERC20Token", reflect.TypeOf((*MockELReader)(nil).GetStrategyAndUnderlyingERC20Token), arg0, arg1)
}

// GetStrategyAndUnderlyingToken mocks base method.
func (m *MockELReader) GetStrategyAndUnderlyingToken(arg0 *bind.CallOpts, arg1 common.Address) (*contractIStrategy.ContractIStrategy, common.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStrategyAndUnderlyingToken", arg0, arg1)
	ret0, _ := ret[0].(*contractIStrategy.ContractIStrategy)
	ret1, _ := ret[1].(common.Address)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetStrategyAndUnderlyingToken indicates an expected call of GetStrategyAndUnderlyingToken.
func (mr *MockELReaderMockRecorder) GetStrategyAndUnderlyingToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStrategyAndUnderlyingToken", reflect.TypeOf((*MockELReader)(nil).GetStrategyAndUnderlyingToken), arg0, arg1)
}

// IsOperatorRegistered mocks base method.
func (m *MockELReader) IsOperatorRegistered(arg0 *bind.CallOpts, arg1 types.Operator) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOperatorRegistered", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsOperatorRegistered indicates an expected call of IsOperatorRegistered.
func (mr *MockELReaderMockRecorder) IsOperatorRegistered(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOperatorRegistered", reflect.TypeOf((*MockELReader)(nil).IsOperatorRegistered), arg0, arg1)
}

// OperatorIsFrozen mocks base method.
func (m *MockELReader) OperatorIsFrozen(arg0 *bind.CallOpts, arg1 common.Address) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OperatorIsFrozen", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OperatorIsFrozen indicates an expected call of OperatorIsFrozen.
func (mr *MockELReaderMockRecorder) OperatorIsFrozen(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OperatorIsFrozen", reflect.TypeOf((*MockELReader)(nil).OperatorIsFrozen), arg0, arg1)
}

// ServiceManagerCanSlashOperatorUntilBlock mocks base method.
func (m *MockELReader) ServiceManagerCanSlashOperatorUntilBlock(arg0 *bind.CallOpts, arg1, arg2 common.Address) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceManagerCanSlashOperatorUntilBlock", arg0, arg1, arg2)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ServiceManagerCanSlashOperatorUntilBlock indicates an expected call of ServiceManagerCanSlashOperatorUntilBlock.
func (mr *MockELReaderMockRecorder) ServiceManagerCanSlashOperatorUntilBlock(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceManagerCanSlashOperatorUntilBlock", reflect.TypeOf((*MockELReader)(nil).ServiceManagerCanSlashOperatorUntilBlock), arg0, arg1, arg2)
}
