package services

import (
	"context"
	"encoding/json"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/helium/rosetta-helium/helium"
)

// ConstructionAPIService implements the server.ConstructionAPIServicer interface.
type ConstructionAPIService struct {
	network *types.NetworkIdentifier
}

// NewConstructionAPIService creates a new instance of a ConstructionAPIService.
func NewConstructionAPIService(network *types.NetworkIdentifier) server.ConstructionAPIServicer {
	return &ConstructionAPIService{
		network: network,
	}
}

func (s *ConstructionAPIService) ConstructionCombine(
	ctx context.Context,
	request *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {
	combineResponse, cErr := helium.CombineTransaction(request.UnsignedTransaction, request.Signatures)
	if cErr != nil {
		return nil, cErr
	}

	return combineResponse, nil
}

func (s *ConstructionAPIService) ConstructionDerive(
	ctx context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	derivedAddress, dErr := helium.GetAddress(request.PublicKey.CurveType, request.PublicKey.Bytes)
	if dErr != nil {
		return nil, dErr
	}

	deriveResponse := &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address: *derivedAddress,
		},
	}

	return deriveResponse, nil
}

func (s *ConstructionAPIService) ConstructionHash(
	ctx context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	transaction := request.SignedTransaction
	hash, hErr := helium.GetHash(transaction)
	if hErr != nil {
		return nil, hErr
	}

	response := &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: *hash,
		},
	}

	return response, nil
}

func (s *ConstructionAPIService) ConstructionMetadata(
	ctx context.Context,
	request *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {

	metadata, err := helium.GetMetadata(request)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *ConstructionAPIService) ConstructionParse(
	ctx context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {
	operations, signer, err := helium.ParseTransaction(request.Transaction, request.Signed)
	if err != nil {
		return nil, err
	}

	parseResponse := &types.ConstructionParseResponse{
		Operations: operations,
	}

	if signer != nil {
		parseResponse.AccountIdentifierSigners = []*types.AccountIdentifier{
			signer,
		}
	}

	return parseResponse, nil
}

func (s *ConstructionAPIService) ConstructionPayloads(
	ctx context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {

	operations := request.Operations

	payload, err := helium.PayloadGenerator(operations, request.Metadata)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (s *ConstructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {

	operations := request.Operations
	transactionPreprocessor, err := helium.OpsToTransaction(operations)
	if err != nil {
		return nil, err
	}

	// Convert transactionPreprocessor into map to satisfy Option requirement
	var options map[string]interface{}
	marshalledPreprocessor, _ := json.Marshal(transactionPreprocessor)
	json.Unmarshal(marshalledPreprocessor, &options)

	return &types.ConstructionPreprocessResponse{
		Options: options,
	}, nil

}

func (s *ConstructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	submittedTxnHash, sErr := helium.SubmitTransaction(request.SignedTransaction)
	if sErr != nil {
		return nil, sErr
	}

	submitResponse := &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: *submittedTxnHash,
		},
	}

	return submitResponse, nil
}
