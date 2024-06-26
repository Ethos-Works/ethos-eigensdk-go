package types

import "github.com/ethos-works/ethos-eigensdk-go/crypto/bls"

type TaskIndex = uint32
type TaskResponseDigest [32]byte

type SignedTaskResponseDigest struct {
	TaskResponseDigest          TaskResponseDigest
	BlsSignature                *bls.Signature
	OperatorId                  bls.OperatorId
	SignatureVerificationErrorC chan error
}
