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

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/helium/rosetta-helium/helium"
	"go.uber.org/zap"
)

// AccountAPIService implements the server.AccountAPIServicer interface.
type AccountAPIService struct {
	network *types.NetworkIdentifier
}

// NewBlockAPIService creates a new instance of a BlockAPIService.
func NewAccountAPIService(network *types.NetworkIdentifier) server.AccountAPIServicer {
	return &AccountAPIService{
		network: network,
	}
}

// AccountBalance implements the /account/balance endpoint.
func (s *AccountAPIService) AccountBalance(
	ctx context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {

	balanceRequest := helium.GetBalanceRequest{
		Address: request.AccountIdentifier.Address,
	}

	zap.S().Info(request.AccountIdentifier.Address + " " + fmt.Sprint(*request.BlockIdentifier.Index))

	if request.BlockIdentifier != nil {
		if request.BlockIdentifier.Index == nil {
			return nil, helium.WrapErr(
				helium.ErrInvalidParameter,
				errors.New("request.BlockIdentifier requires an Index"),
			)
		}
		balanceRequest.Height = *request.BlockIdentifier.Index
	}

	if helium.NodeBalancesDB != nil {
		var accountBalances []*types.Amount
		accountEntry, aeErr := helium.RocksDBAccountGet(request.AccountIdentifier.Address, balanceRequest.Height)
		if aeErr != nil {
			zap.S().Info("no balance found for " + balanceRequest.Address + " at height " + fmt.Sprint(balanceRequest.Height) + ". Returning balanaces of 0.")
			accountBalances = []*types.Amount{
				{
					Value:    "0",
					Currency: helium.HNT,
				},
				{
					Value:    "0",
					Currency: helium.HST,
				},
			}
		} else {
			accountBalances = []*types.Amount{
				{
					Value:    fmt.Sprint(accountEntry.Entry.Amount),
					Currency: helium.HNT,
				},
				{
					Value:    fmt.Sprint(accountEntry.SecEntry.Amount),
					Currency: helium.HST,
				},
			}
		}

		blockHash, bhErr := helium.RocksDBBlockHashGet(*request.BlockIdentifier.Index)
		if bhErr != nil {
			return nil, helium.WrapErr(helium.ErrFailed, bhErr)
		}

		blockIdentifier := &types.BlockIdentifier{
			Index: *request.BlockIdentifier.Index,
			Hash:  *blockHash,
		}

		return &types.AccountBalanceResponse{
			BlockIdentifier: blockIdentifier,
			Balances:        accountBalances,
		}, nil
	} else {
		zap.S().Info("Old path")
		accountBalances, aErr := helium.GetBalance(balanceRequest)
		if aErr != nil {
			zap.S().Info("no balance found for " + balanceRequest.Address + " at height " + fmt.Sprint(balanceRequest.Height) + ". Returning balanaces of 0.")
			accountBalances = []*types.Amount{
				{
					Value:    "0",
					Currency: helium.HNT,
				},
				{
					Value:    "0",
					Currency: helium.HST,
				},
			}
		}

		var blockId types.BlockIdentifier

		if request.BlockIdentifier == nil {
			currentHeight, chErr := helium.GetCurrentHeight()
			if chErr != nil {
				return nil, chErr
			}

			currentBlock, cErr := helium.GetBlockIdentifier(&types.PartialBlockIdentifier{
				Index: currentHeight,
			})
			if cErr != nil {
				return nil, cErr
			}

			blockId = *currentBlock
		} else {
			requestedBlock, rErr := helium.GetBlockIdentifier(&types.PartialBlockIdentifier{
				Index: request.BlockIdentifier.Index,
			})
			if rErr != nil {
				return nil, rErr
			}

			blockId = *requestedBlock
		}

		return &types.AccountBalanceResponse{
			BlockIdentifier: &blockId,
			Balances:        accountBalances,
		}, nil
	}
}

// AccountCoins implements the /account/coins endpoint.
func (s *AccountAPIService) AccountCoins(
	ctx context.Context,
	request *types.AccountCoinsRequest,
) (*types.AccountCoinsResponse, *types.Error) {
	return nil, nil
}
