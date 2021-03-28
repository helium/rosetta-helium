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

package helium

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/coinbase/rosetta-sdk-go/types"
)

func readLBSfile() *int64 {
	file, err := os.Open("lbs.txt")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.TrimSuffix(buf.String(), "\n")

	lbs, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	return &lbs
}

const (
	// NodeVersion is the version of helium we are using.
	NodeVersion = "1.1.23"

	// Blockchain is Helium.
	Blockchain string = "Helium"

	// MainnetNetwork is the value of the network
	// in MainnetNetworkIdentifier.
	MainnetNetwork string = "Mainnet"

	// TestnetNetwork is the value of the network
	// in TestnetNetworkIdentifier.
	TestnetNetwork string = "Testnet"

	// Symbol is the symbol value
	// used in Currency.
	Symbol = "HNT"

	// Decimals is the decimals value
	// used in Currency.
	Decimals = 8

	// AddGatewayOpType is used to describe
	// adding a gateway.
	AddGatewayOpType = "add_gateway_v1"

	// AssertLocationOpType is used to describe
	// asserting a gateway's location.
	AssertLocationOpType = "assert_location_v1"

	// ChainVarsOpType is used to describe
	// changing a chain variable.
	ChainVarsOpType = "vars_v1"

	// CoinbaseOpType is used to describe
	// a coinbase transaction.
	CoinbaseOpType = "COINBASE"

	// CoinbaseDataCreditsOpType is used to describe
	// the initial 10k credits to bring initial miners online.
	CoinbaseDataCreditsOpType = "dc_coinbase_v1"

	// ConsensusGroupOpType is used to describe
	// the election of a new consensus group
	ConsensusGroupOpType = "consensus_group_v1"

	// CreateHashedTimelockOpType is used to describe
	// creating a hashed timelock
	CreateHashedTimelockOpType = "create_htlc_v1"

	// CreateProofOfCoverageRequestOpType is used to describe
	// a proof of coverage request
	CreateProofOfCoverageRequestOpType = "poc_request_v1"

	// DataCreditsOpType is used to describe
	// burning HNT for DCs
	DataCreditsOpType = "token_burn_v1"

	// GenesisGatewayOpType is used to describe
	// initial group of miners that bootstrapped the blockchain
	GenesisGatewayOpType = "gen_gateway_v1"

	// MultiPaymentOpType is used to describe
	// a transaction from one wallet to multiple
	MultiPaymentOpType = "payment_v2"

	// OUIType is used to describe
	// a new OUI for a new router on the network
	OUIOpType = "oui_v1"

	// PaymentOpType is used to describe
	// sending HNT from one wallet to another
	PaymentV1OpType = "payment_v1"

	// ProofOfCoverageReceiptsOpType is used to describe
	// completed POC submitted to the network
	ProofOfCoverageReceiptsOpType = "poc_receipts_v1"

	// RedeemHashedTimelockOpType is used to describe
	// redeeming a hashed timelock
	RedeemHashedTimelockOpType = "redeem_htlc_v1"

	// RewardOpType is used to describe
	// a token payout for a specific event on the network
	RewardOpType = "reward_v1"

	// RewardsOpType is used to describe
	// a bundle of multiple reward transactions
	RewardsOpType = "rewards_v1"

	// RoutingOpType is used to describe
	// updating the routing information with an OUI
	RoutingOpType = "routing_v1"

	// SecurityCoinbaseOpType is used to describe
	// the distribution of security tokens in genesis block
	SecurityCoinbaseOpType = "security_coinbase_v1"

	// SecurityExchangeOpType is used to describe
	// the transfer of security tokens from one address to another
	SecurityExchangeOpType = "security_exchange_v1"

	// StateChannelOpenOpType is used to describe
	// opening a new state channel on a Helium router
	StateChannelOpenOpType = "state_channel_open_v1"

	// StateChannelCloseOpType is used to describe
	// closing a state channel on a Helium router
	StateChannelCloseOpType = "state_channel_close_v1"

	// TokenBurnExchangeRateOpType is used to describe
	// changing the exchange rate for burning HNT to DCs
	TokenBurnExchangeRateOpType = "price_oracle_v1"

	// TransferHotspotOpType is used to describe
	// transferring hotspots from one wallet to another
	TransferHotspotOpType = "transfer_hotspot_v1"

	// SuccessStatus is the status of any
	// Helium operation considered successful.
	SuccessStatus = "SUCCESS"

	// FailureStatus is the status of any
	// Helium operation considered unsuccessful.
	FailureStatus = "FAILURE"

	// HistoricalBalanceSupported is whether
	// historical balance is supported.
	HistoricalBalanceSupported = false

	// Genesis is the index of the
	// genesis block for blockchain-etl instances
	GenesisBlockIndex = int64(0)

	// IncludeMempoolCoins does not apply to rosetta-ethereum as it is not UTXO-based.
	IncludeMempoolCoins = false
)

var (

	// GenesisBlockIdentifier is the *types.BlockIdentifier
	// of the mainnet genesis block.
	GenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  "La6PuV80Ps9qTP0339Pwm64q3_deMTkv6JOo1251EJI",
		Index: 1,
	}

	// Currency is the *types.Currency for all
	// Ethereum networks.
	Currency = &types.Currency{
		Symbol:   Symbol,
		Decimals: Decimals,
	}

	// OperationTypes are all suppoorted operation types.
	OperationTypes = []string{
		AddGatewayOpType,
		AssertLocationOpType,
		ChainVarsOpType,
		CoinbaseOpType,
		CoinbaseDataCreditsOpType,
		ConsensusGroupOpType,
		CreateHashedTimelockOpType,
		CreateProofOfCoverageRequestOpType,
		DataCreditsOpType,
		GenesisGatewayOpType,
		MultiPaymentOpType,
		OUIOpType,
		PaymentV1OpType,
		ProofOfCoverageReceiptsOpType,
		RedeemHashedTimelockOpType,
		RewardOpType,
		RewardsOpType,
		RoutingOpType,
		SecurityCoinbaseOpType,
		SecurityExchangeOpType,
		StateChannelOpenOpType,
		StateChannelCloseOpType,
		TokenBurnExchangeRateOpType,
		TransferHotspotOpType,
	}

	// OperationStatuses are all supported operation statuses.
	OperationStatuses = []*types.OperationStatus{
		{
			Status:     SuccessStatus,
			Successful: true,
		},
		{
			Status:     FailureStatus,
			Successful: false,
		},
	}

	//LastBlessedSnapshotIndex
	LBS = readLBSfile()
)

type Block struct {
	Hash         string   `json:"hash"`
	Height       int64    `json:"height"`
	PrevHash     string   `json:"prev_hash"`
	Time         int64    `json:"time"`
	Transactions []string `json:"transactions"`
}

type Transaction struct {
	Hash string `json:"hash"`
	Type string `json:"type"`
}

type Payment struct {
	Payee  string `json:"payee"`
	Amount int64  `json:"amount"`
}

type PaymentV2Transaction struct {
	Transaction
	Fee      int64     `json:"fee"`
	Nonce    int64     `json:"nonce"`
	Payer    string    `json:"payer"`
	Payments []Payment `json:"payments"`
}
