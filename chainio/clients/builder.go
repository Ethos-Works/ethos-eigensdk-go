package clients

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethos-works/ethos-eigensdk-go/chainio/clients/avsregistry"
	"github.com/ethos-works/ethos-eigensdk-go/chainio/clients/elcontracts"
	"github.com/ethos-works/ethos-eigensdk-go/chainio/clients/eth"
	"github.com/ethos-works/ethos-eigensdk-go/chainio/txmgr"
	chainioutils "github.com/ethos-works/ethos-eigensdk-go/chainio/utils"
	"github.com/ethos-works/ethos-eigensdk-go/logging"
	"github.com/ethos-works/ethos-eigensdk-go/signerv2"
	"github.com/ethos-works/ethos-eigensdk-go/types"
)

type BuildAllConfig struct {
	EthHttpUrl                 string
	RegistryCoordinatorAddr    string
	OperatorStateRetrieverAddr string
	AvsName                    string
}

// TODO: this is confusing right now because clients are not instrumented clients, but
// we return metrics and prometheus reg, so user has to build instrumented clients at the call
// site if they need them. We should probably separate into two separate constructors, one
// for non-instrumented clients that doesn't return metrics/reg, and another instrumented-constructor
// that returns instrumented clients and the metrics/reg.
type Clients struct {
	AvsRegistryChainReader *avsregistry.AvsRegistryChainReader
	AvsRegistryChainWriter *avsregistry.AvsRegistryChainWriter
	ElChainReader          *elcontracts.ELChainReader
	ElChainWriter          *elcontracts.ELChainWriter
	EthHttpClient          *eth.Client
}

func BuildAll(
	config BuildAllConfig,
	signerAddr gethcommon.Address,
	signerFn signerv2.SignerFn,
	logger logging.Logger,
) (*Clients, error) {
	config.validate(logger)

	ethHttpClient, err := eth.NewClient(config.EthHttpUrl)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to create Eth Http client"), err)
	}

	txMgr := txmgr.NewSimpleTxManager(ethHttpClient, logger, signerFn, signerAddr)
	// creating EL clients: Reader, Writer and Subscriber
	elChainReader, elChainWriter, err := config.buildElClients(
		ethHttpClient,
		txMgr,
		logger,
	)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to create EL Reader, Writer and Subscriber"), err)
	}

	// creating AVS clients: Reader and Writer
	avsRegistryChainReader, avsRegistryChainWriter, err := config.buildAvsClients(
		elChainReader,
		ethHttpClient,
		txMgr,
		logger,
	)
	if err != nil {
		return nil, types.WrapError(errors.New("Failed to create AVS Registry Reader and Writer"), err)
	}

	return &Clients{
		ElChainReader:          elChainReader,
		ElChainWriter:          elChainWriter,
		AvsRegistryChainReader: avsRegistryChainReader,
		AvsRegistryChainWriter: avsRegistryChainWriter,
		EthHttpClient:          ethHttpClient,
	}, nil

}

func (config *BuildAllConfig) buildElClients(
	ethHttpClient eth.EthClient,
	txMgr txmgr.TxManager,
	logger logging.Logger,
) (*elcontracts.ELChainReader, *elcontracts.ELChainWriter, error) {

	avsRegistryContractBindings, err := chainioutils.NewAVSRegistryContractBindings(
		gethcommon.HexToAddress(config.RegistryCoordinatorAddr),
		gethcommon.HexToAddress(config.OperatorStateRetrieverAddr),
		ethHttpClient,
		logger,
	)
	if err != nil {
		return nil, nil, types.WrapError(errors.New("Failed to create AVSRegistryContractBindings"), err)
	}

	delegationManagerAddr, err := avsRegistryContractBindings.StakeRegistry.Delegation(&bind.CallOpts{})
	if err != nil {
		logger.Fatal("Failed to fetch Slasher contract", "err", err)
	}
	avsDirectoryAddr, err := avsRegistryContractBindings.ServiceManager.AvsDirectory(&bind.CallOpts{})
	if err != nil {
		logger.Fatal("Failed to fetch Slasher contract", "err", err)
	}

	elContractBindings, err := chainioutils.NewEigenlayerContractBindings(
		delegationManagerAddr,
		avsDirectoryAddr,
		ethHttpClient,
		logger,
	)
	if err != nil {
		return nil, nil, types.WrapError(errors.New("Failed to create EigenlayerContractBindings"), err)
	}

	// get the Reader for the EL contracts
	elChainReader := elcontracts.NewELChainReader(
		elContractBindings.Slasher,
		elContractBindings.DelegationManager,
		elContractBindings.StrategyManager,
		elContractBindings.AvsDirectory,
		logger,
		ethHttpClient,
	)

	elChainWriter := elcontracts.NewELChainWriter(
		elContractBindings.Slasher,
		elContractBindings.DelegationManager,
		elContractBindings.StrategyManager,
		elContractBindings.StrategyManagerAddr,
		elChainReader,
		ethHttpClient,
		logger,
		txMgr,
	)
	if err != nil {
		return nil, nil, types.WrapError(errors.New("Failed to create ELChainWriter"), err)
	}

	return elChainReader, elChainWriter, nil
}

func (config *BuildAllConfig) buildAvsClients(
	elReader elcontracts.ELReader,
	ethHttpClient eth.EthClient,
	txMgr txmgr.TxManager,
	logger logging.Logger,
) (*avsregistry.AvsRegistryChainReader, *avsregistry.AvsRegistryChainWriter, error) {

	avsRegistryContractBindings, err := chainioutils.NewAVSRegistryContractBindings(
		gethcommon.HexToAddress(config.RegistryCoordinatorAddr),
		gethcommon.HexToAddress(config.OperatorStateRetrieverAddr),
		ethHttpClient,
		logger,
	)
	if err != nil {
		return nil, nil, types.WrapError(errors.New("Failed to create AVSRegistryContractBindings"), err)
	}

	avsRegistryChainReader := avsregistry.NewAvsRegistryChainReader(
		avsRegistryContractBindings.RegistryCoordinatorAddr,
		avsRegistryContractBindings.RegistryCoordinator,
		avsRegistryContractBindings.OperatorStateRetriever,
		avsRegistryContractBindings.StakeRegistry,
		logger,
		ethHttpClient,
	)

	avsRegistryChainWriter, err := avsregistry.NewAvsRegistryChainWriter(
		avsRegistryContractBindings.ServiceManagerAddr,
		avsRegistryContractBindings.RegistryCoordinator,
		avsRegistryContractBindings.OperatorStateRetriever,
		avsRegistryContractBindings.StakeRegistry,
		elReader,
		logger,
		ethHttpClient,
		txMgr,
	)
	if err != nil {
		return nil, nil, types.WrapError(errors.New("Failed to create AVSRegistryChainWriter"), err)
	}

	if err != nil {
		return nil, nil, types.WrapError(errors.New("Failed to create ELChainSubscriber"), err)
	}

	return avsRegistryChainReader, avsRegistryChainWriter, nil
}

// Very basic validation that makes sure all fields are nonempty
// we might eventually want more sophisticated validation, based on regexp,
// or use something like https://json-schema.org/ (?)
func (config *BuildAllConfig) validate(logger logging.Logger) {
	if config.EthHttpUrl == "" {
		logger.Fatalf("BuildAllConfig.validate: Missing eth http url")
	}
	if config.RegistryCoordinatorAddr == "" {
		logger.Fatalf("BuildAllConfig.validate: Missing bls registry coordinator address")
	}
	if config.OperatorStateRetrieverAddr == "" {
		logger.Fatalf("BuildAllConfig.validate: Missing bls operator state retriever address")
	}
	if config.AvsName == "" {
		logger.Fatalf("BuildAllConfig.validate: Missing avs name")
	}
}
