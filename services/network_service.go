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

	"github.com/helium/rosetta-helium/helium"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// NetworkAPIService implements the server.NetworkAPIServicer interface.
type NetworkAPIService struct {
	network *types.NetworkIdentifier
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(network *types.NetworkIdentifier) server.NetworkAPIServicer {
	return &NetworkAPIService{
		network: network,
	}
}

// NetworkList implements the /network/list endpoint
func (s *NetworkAPIService) NetworkList(
	ctx context.Context,
	request *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			s.network,
		},
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (s *NetworkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {

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

	currentHeight, chErr := helium.GetCurrentHeight()
	if chErr != nil {
		return nil, chErr
	}

	currentBlock, cbErr := helium.GetBlockMeta(&types.PartialBlockIdentifier{
		Index: currentHeight,
	})

	if cbErr != nil {
		return nil, cbErr
	}

	currentBlockID := &types.BlockIdentifier{
		Index: currentBlock.Height,
		Hash:  currentBlock.Hash,
	}

	currentBlockTimestamp := currentBlock.Time

	if currentBlockID.Index < *helium.LBS {
		return nil, helium.WrapErr(helium.ErrNodeSync, errors.New("node is catching up to snapshot height"))
	}

	lastBlessedBlock, lbErr := helium.GetBlockIdentifier(&types.PartialBlockIdentifier{
		Index: helium.LBS,
	})

	if lbErr != nil {
		return nil, lbErr
	}

	peers, pErr := helium.GetPeers()
	if pErr != nil {
		return nil, pErr
	}

	// Removing syncStatus for now since currentBlock is
	// retrieved from external API which is not ideal.
	//
	// syncStatus, sErr := helium.GetSyncStatus()
	// if sErr != nil {
	// 	return nil, sErr
	// }

	genesisIndex := helium.MainnetGenesisBlockIndex
	genesisHash := helium.MainnetGenesisBlockHash

	if request.NetworkIdentifier.Network == helium.TestnetNetwork {
		genesisIndex = helium.TestnetGenesisBlockIndex
		genesisHash = helium.TestnetGenesisBlockHash
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: currentBlockID,
		CurrentBlockTimestamp:  currentBlockTimestamp,
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: genesisIndex,
			Hash:  genesisHash,
		},
		OldestBlockIdentifier: lastBlessedBlock,
		Peers:                 peers,
		// SyncStatus:            syncStatus,
	}, nil
}

// NetworkOptions implements the /network/options endpoint.
func (s *NetworkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: "1.4.10",
			NodeVersion:    helium.NodeVersion,
		},
		Allow: &types.Allow{
			Errors:                  helium.Errors,
			OperationTypes:          helium.OperationTypes,
			OperationStatuses:       helium.OperationStatuses,
			HistoricalBalanceLookup: helium.HistoricalBalanceSupported,
		},
	}, nil
}
