// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/helium/rosetta-helium/helium"
	"go.uber.org/zap"
)

// BlockAPIService implements the server.BlockAPIServicer interface.
type BlockAPIService struct {
	network *types.NetworkIdentifier
}

// NewBlockAPIService creates a new instance of a BlockAPIService.
func NewBlockAPIService(network *types.NetworkIdentifier) server.BlockAPIServicer {
	return &BlockAPIService{
		network: network,
	}
}

// Block implements the /block endpoint.
func (s *BlockAPIService) Block(
	ctx context.Context,
	request *types.BlockRequest,
) (*types.BlockResponse, *types.Error) {

	zap.S().Info("/block request: " + fmt.Sprint(*request.BlockIdentifier.Index))

	// Update all secondary rocksdb references
	if tErr := helium.NodeBalancesDB.TryCatchUpWithPrimary(); tErr != nil {
		return nil, helium.WrapErr(helium.ErrFailed, tErr)
	}

	if tErr := helium.NodeBlocksDB.TryCatchUpWithPrimary(); tErr != nil {
		return nil, helium.WrapErr(helium.ErrFailed, tErr)
	}

	if tErr := helium.NodeTransactionsDB.TryCatchUpWithPrimary(); tErr != nil {
		return nil, helium.WrapErr(helium.ErrFailed, tErr)
	}

	previousBlockIndex := *request.BlockIdentifier.Index - 1
	if previousBlockIndex == 0 {
		previousBlockIndex = 1
	}

	requestedBlock, rErr := helium.GetBlock(request.BlockIdentifier)
	if rErr != nil {
		return nil, rErr
	}

	previousBlock, pErr := helium.GetBlock(&types.PartialBlockIdentifier{
		Index: &previousBlockIndex,
	})
	if pErr != nil {
		return nil, pErr
	}

	if request.BlockIdentifier.Hash != nil {
		if requestedBlock.BlockIdentifier.Hash != *request.BlockIdentifier.Hash {
			return nil, helium.WrapErr(
				helium.ErrNotFound,
				errors.New("ambiguous request: requested block height ("+
					strconv.FormatInt(*request.BlockIdentifier.Index, 10)+
					") does not match returned block hash ("+
					requestedBlock.BlockIdentifier.Hash+
					")"),
			)
		}
	}

	return &types.BlockResponse{
		Block: &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Index: *request.BlockIdentifier.Index,
				Hash:  requestedBlock.BlockIdentifier.Hash,
			},
			ParentBlockIdentifier: &types.BlockIdentifier{
				Index: previousBlockIndex,
				Hash:  previousBlock.BlockIdentifier.Hash,
			},
			Timestamp:    requestedBlock.Timestamp,
			Transactions: requestedBlock.Transactions,
		},
	}, nil
}

// BlockTransaction implements the /block/transaction endpoint.
func (s *BlockAPIService) BlockTransaction(
	ctx context.Context,
	request *types.BlockTransactionRequest,
) (*types.BlockTransactionResponse, *types.Error) {
	txn, txErr := helium.GetTransaction(request.TransactionIdentifier.Hash, request.BlockIdentifier)
	if txErr != nil {
		return nil, txErr
	}

	return &types.BlockTransactionResponse{
		Transaction: txn,
	}, nil
}
