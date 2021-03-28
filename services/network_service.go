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
	"fmt"
	"strconv"

	"github.com/syuan100/rosetta-helium/helium"

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
	fmt.Println("Getting latest block: ")
	currentBlock := GetBlock(CurrentBlockHeight())
	lastBlessedBlockHeight, err := strconv.ParseInt(*helium.LBS, 10, 64)
	if err != nil {
		return nil, wrapErr(
			ErrEnvVariableMissing,
			err,
		)
	}
	fmt.Println("Getting last blessed block: ")
	lastBlessedBlock := GetBlock(lastBlessedBlockHeight)
	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: currentBlock.BlockIdentifier.Index,
			Hash:  currentBlock.BlockIdentifier.Hash,
		},
		CurrentBlockTimestamp: currentBlock.Timestamp * 1000,
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: lastBlessedBlock.BlockIdentifier.Index,
			Hash:  lastBlessedBlock.BlockIdentifier.Hash,
		},
		OldestBlockIdentifier: &types.BlockIdentifier{
			Index: lastBlessedBlock.BlockIdentifier.Index,
			Hash:  lastBlessedBlock.BlockIdentifier.Hash,
		},
		Peers: nil,
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
			Errors:                  Errors,
			OperationTypes:          helium.OperationTypes,
			OperationStatuses:       helium.OperationStatuses,
			HistoricalBalanceLookup: helium.HistoricalBalanceSupported,
		},
	}, nil
}
