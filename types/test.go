package types

type TestOperator struct {
	Operator       OperatorAddr
	StakePerQuorum map[QuorumNum]StakeAmount
}
