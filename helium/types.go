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

	// Symbol for Data Credits
	DCSymbol = "DC"

	// Decimals for Data Credits
	DCDecimals = 0

	// Symbol for Security Tokens
	HSTSymbol = "HST"

	// Decimals for Security Tokens
	HSTDecimals = 8

	// AddGatewayTxn is used to describe
	// adding a gateway.
	AddGatewayTxn = "add_gateway_v1"

	// AssertLocationTxn is used to describe
	// asserting a gateway's location.
	AssertLocationTxn = "assert_location_v1"

	// ChainVarsTxn is used to describe
	// changing a chain variable.
	ChainVarsTxn = "vars_v1"

	// CoinbaseTxn is used to describe
	// a coinbase transaction.
	CoinbaseTxn = "COINBASE"

	// CoinbaseDataCreditsTxn is used to describe
	// the initial 10k credits to bring initial miners online.
	CoinbaseDataCreditsTxn = "dc_coinbase_v1"

	// ConsensusGroupTxn is used to describe
	// the election of a new consensus group
	ConsensusGroupTxn = "consensus_group_v1"

	// CreateHashedTimelockTxn is used to describe
	// creating a hashed timelock
	CreateHashedTimelockTxn = "create_htlc_v1"

	// CreateProofOfCoverageRequestTxn is used to describe
	// a proof of coverage request
	CreateProofOfCoverageRequestTxn = "poc_request_v1"

	// DataCreditsTxn is used to describe
	// burning HNT for DCs
	DataCreditsTxn = "token_burn_v1"

	// GenesisGatewayTxn is used to describe
	// initial group of miners that bootstrapped the blockchain
	GenesisGatewayTxn = "gen_gateway_v1"

	// PaymentV2Txn is used to describe
	// a transaction from one wallet to multiple
	PaymentV2Txn = "payment_v2"

	// OUIType is used to describe
	// a new OUI for a new router on the network
	OUITxn = "oui_v1"

	// PaymentTxn is used to describe
	// sending HNT from one wallet to another
	PaymentV1Txn = "payment_v1"

	// ProofOfCoverageReceiptsTxn is used to describe
	// completed POC submitted to the network
	ProofOfCoverageReceiptsTxn = "poc_receipts_v1"

	// RedeemHashedTimelockTxn is used to describe
	// redeeming a hashed timelock
	RedeemHashedTimelockTxn = "redeem_htlc_v1"

	// RewardTxn is used to describe
	// a token payout for a specific event on the network
	RewardTxnV1 = "reward_v1"

	// RewardsTxn is used to describe
	// a bundle of multiple reward transactions
	RewardsTxnV1 = "rewards_v1"

	// RewardTxn is used to describe
	// a token payout for a specific event on the network
	RewardTxnV2 = "reward_v2"

	// RewardsTxn is used to describe
	// a bundle of multiple reward transactions
	RewardsTxnV2 = "rewards_v2"

	// RoutingTxn is used to describe
	// updating the routing information with an OUI
	RoutingTxn = "routing_v1"

	// SecurityCoinbaseTxn is used to describe
	// the distribution of security tokens in genesis block
	SecurityCoinbaseTxn = "security_coinbase_v1"

	// SecurityExchangeTxn is used to describe
	// the transfer of security tokens from one address to another
	SecurityExchangeTxn = "security_exchange_v1"

	// StateChannelOpenTxn is used to describe
	// opening a new state channel on a Helium router
	StateChannelOpenTxn = "state_channel_open_v1"

	// StateChannelCloseTxn is used to describe
	// closing a state channel on a Helium router
	StateChannelCloseTxn = "state_channel_close_v1"

	// TokenBurnExchangeRateTxn is used to describe
	// changing the exchange rate for burning HNT to DCs
	TokenBurnExchangeRateTxn = "price_oracle_v1"

	// TransferHotspotTxn is used to describe
	// transferring hotspots from one wallet to another
	TransferHotspotTxn = "transfer_hotspot_v1"

	// CreditOp is used to describe
	// a credit to an account (HNT, HST, or DC)
	CreditOp = "credit_op"

	// DebitOp is used to describe
	// a debit from an account (HNT, HST, or DC)
	DebitOp = "debit_op"

	// FeeOp is used to describe
	// a transaction fee (Negative HNT or DC)
	FeeOp = "fee_op"

	// AddGatewayOp is used to describe
	// associating a new gateway with an owner
	AddGatewayOp = "add_gateway_op"

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

	// SuccessStatus is the status of any
	// Helium operation considered successful.
	SuccessStatus = "SUCCESS"

	// FailureStatus is the status of any
	// Helium operation considered unsuccessful.
	FailureStatus = "FAILURE"

	// GenesisBlockIdentifier is the *types.BlockIdentifier
	// of the mainnet genesis block.
	GenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  "La6PuV80Ps9qTP0339Pwm64q3_deMTkv6JOo1251EJI",
		Index: 1,
	}

	// HNT is the *types.Currency for HNT.
	HNT = &types.Currency{
		Symbol:   Symbol,
		Decimals: Decimals,
	}

	// DC is the *types.Currency for DC.
	DC = &types.Currency{
		Symbol:   DCSymbol,
		Decimals: DCDecimals,
	}

	// HST is the *types.Currency for HST.
	HST = &types.Currency{
		Symbol:   HSTSymbol,
		Decimals: HSTDecimals,
	}

	// TransactionTypes are all suppoorted operation types.
	TransactionTypes = []string{
		AddGatewayTxn,
		AssertLocationTxn,
		ChainVarsTxn,
		CoinbaseTxn,
		CoinbaseDataCreditsTxn,
		ConsensusGroupTxn,
		CreateHashedTimelockTxn,
		CreateProofOfCoverageRequestTxn,
		DataCreditsTxn,
		GenesisGatewayTxn,
		OUITxn,
		PaymentV1Txn,
		PaymentV2Txn,
		ProofOfCoverageReceiptsTxn,
		RedeemHashedTimelockTxn,
		RewardTxnV1,
		RewardsTxnV1,
		RewardTxnV2,
		RewardsTxnV2,
		RoutingTxn,
		SecurityCoinbaseTxn,
		SecurityExchangeTxn,
		StateChannelOpenTxn,
		StateChannelCloseTxn,
		TokenBurnExchangeRateTxn,
		TransferHotspotTxn,
	}

	// OperationTypes are all supported base operations
	// that make up a transaction
	OperationTypes = []string{
		CreditOp,
		DebitOp,
		FeeOp,
		AddGatewayOp,
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

	// LBS is the LastBlessedBlock height as an int64
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

type Fee struct {
	Amount   int64
	Payer    string
	Currency *types.Currency
}

type Reward struct {
	Account string `json:"account"`
	Amount  int64  `json:"amount"`
	Gateway string `json:"gateway"`
	Type    string `json:"type"`
}

type PaymentV2Transaction struct {
	Transaction
	Fee      int64     `json:"fee"`
	Nonce    int64     `json:"nonce"`
	Payer    string    `json:"payer"`
	Payments []Payment `json:"payments"`
}
