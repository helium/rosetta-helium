package helium

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/coinbase/rosetta-sdk-go/types"
	rocksdb "github.com/linxGnu/grocksdb"
	"go.uber.org/zap"
)

func readLBSfile() *int64 {
	file, err := os.Open("lbs.txt")
	if err != nil {
		zap.S().Info("No lbs.txt found: attempting to set last blessed snapshot as genesis (1)...")
		tmpLBS := int64(1)
		return &tmpLBS
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.TrimSuffix(buf.String(), "\n")

	lbs, err := strconv.ParseInt(s, 10, 64)

	// Increment last blessed snapshot block in order to account
	// for previous_block query in /network/status
	lbs = lbs + 1

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

	// Mainnet geneis block index
	MainnetGenesisBlockIndex int64 = 1

	// Mainnet geneis block hash
	MainnetGenesisBlockHash string = "La6PuV80Ps9qTP0339Pwm64q3_deMTkv6JOo1251EJI"

	// TestnetNetwork is the value of the network
	// in TestnetNetworkIdentifier.
	TestnetNetwork string = "Testnet"

	// Testnet geneis block index
	TestnetGenesisBlockIndex int64 = 1

	// Testnet geneis block hash
	TestnetGenesisBlockHash string = "edKekSj8vMmVPJ4I9vAPGVCisn7ypJ22RE8RLj6LiWk"

	// Symbol is the symbol value
	// used in Currency.
	Symbol = "HNT"

	// Decimals is the decimals value
	// used in Currency.
	Decimals = 8

	// Symbol is the symbol value
	// used in Currency.
	DCSymbol = "DC"

	// Decimals is the decimals value
	// used in Currency.
	DCDecimals = 1

	// Symbol for Security Tokens
	HSTSymbol = "HST"

	// Decimals for Security Tokens
	HSTDecimals = 8

	// AddGatewayV1Txn is used to describe
	// adding a gateway.
	AddGatewayV1Txn = "add_gateway_v1"

	// AssertLocationV1Txn is used to describe
	// asserting a gateway's location.
	AssertLocationV1Txn = "assert_location_v1"

	// AssertLocationV2Txn is used to describe
	// asserting a gateway's location with extra info
	AssertLocationV2Txn = "assert_location_v2"

	// VarsV1Txn is used to describe
	// changing a chain variable.
	VarsV1Txn = "vars_v1"

	// CoinbaseTxn is used to describe
	// a coinbase transaction.
	CoinbaseV1Txn = "coinbase_v1"

	// DCCoinbaseV1Txn is used to describe
	// the initial 10k credits to bring initial miners online.
	DCCoinbaseV1Txn = "dc_coinbase_v1"

	// ConsensusGroupV1Txn is used to describe
	// the election of a new consensus group
	ConsensusGroupV1Txn = "consensus_group_v1"

	// ConsensusGroupFailiureV1Txn is used to describe
	// the failure to elect a new consensus group
	ConsensusGroupFailureV1Txn = "consensus_group_failure_v1"

	// CreateHTLCV1Txn is used to describe
	// creating a hashed timelock
	CreateHTLCV1Txn = "create_htlc_v1"

	// POCRequestV1Txn is used to describe
	// a proof of coverage request
	POCRequestV1Txn = "poc_request_v1"

	// TokenBurnV1Txn is used to describe
	// burning HNT for DCs
	TokenBurnV1Txn = "token_burn_v1"

	// TokenBurnExchangeRateV1Txn is used to describe
	// the exchange rate for burning HNT for DCs
	TokenBurnExchangeRateV1Txn = "token_burn_exchange_rate_v1"

	// GenGatewayV1Txn is used to describe
	// initial group of miners that bootstrapped the blockchain
	GenGatewayV1Txn = "gen_gateway_v1"

	// PaymentV2Txn is used to describe
	// a transaction from one wallet to multiple
	PaymentV2Txn = "payment_v2"

	// OUIV1Txn is used to describe
	// a new OUI for a new router on the network
	OUIV1Txn = "oui_v1"

	// UpdateGatewayOUIV1Txn is used to describe
	// updating a gateway OUI
	UpdateGatewayOUIV1Txn = "update_gateway_oui_v1"

	// PaymentTxn is used to describe
	// sending HNT from one wallet to another
	PaymentV1Txn = "payment_v1"

	// POCReceiptsV1 is used to describe
	// completed POC submitted to the network
	POCReceiptsV1 = "poc_receipts_v1"

	// RedeemHTLCV1Txn is used to describe
	// redeeming a hashed timelock
	RedeemHTLCV1Txn = "redeem_htlc_v1"

	// RewardsTxn is used to describe
	// a bundle of multiple reward transactions
	RewardsV1Txn = "rewards_v1"

	// RewardsTxn is used to describe
	// a bundle of multiple reward transactions
	RewardsV2Txn = "rewards_v2"

	// RoutingV1Txn is used to describe
	// updating the routing information with an OUI
	RoutingV1Txn = "routing_v1"

	// SecurityCoinbaseV1Txn is used to describe
	// the distribution of security tokens in genesis block
	SecurityCoinbaseV1Txn = "security_coinbase_v1"

	// SecurityExchangeV1Txn is used to describe
	// the transfer of security tokens from one address to another
	SecurityExchangeV1Txn = "security_exchange_v1"

	// StateChannelOpenV1Txn is used to describe
	// opening a new state channel on a Helium router
	StateChannelOpenV1Txn = "state_channel_open_v1"

	// StateChannelCloseV1Txn is used to describe
	// closing a state channel on a Helium router
	StateChannelCloseV1Txn = "state_channel_close_v1"

	// GenValidatorV1Txn is used to describe
	// genesis validators
	GenValidatorV1Txn = "gen_validator_v1"

	// StakeValidatorV1Txn is used to describe
	// staking a new validator
	StakeValidatorV1Txn = "stake_validator_v1"

	// UnstakeValidatorV1Txn is used to describe
	// unstaking a validator
	UnstakeValidatorV1Txn = "unstake_validator_v1"

	// TransferValidatorStakeV1Txn is used to describe
	// transferring a validator to a new owner and/or address
	TransferValidatorStakeV1Txn = "transfer_validator_stake_v1"

	// ValidatorHeartbeatV1Txn is used to describe
	// when a validator provides proof of liveness
	ValidatorHeartbeatV1Txn = "validator_heartbeat_v1"

	// GenPriceOracleV1Txn is used to describe
	// a genesis price oracle
	GenPriceOracleV1Txn = "gen_price_oracle_v1"

	// PriceOracleV1Txn is used to describe
	// changing the exchange rate for burning HNT to DCs
	PriceOracleV1Txn = "price_oracle_v1"

	// TransferHotspotV1Txn is used to describe
	// transferring hotspots from one wallet to another
	TransferHotspotV1Txn = "transfer_hotspot_v1"

	// TransferHotspotV1Txn is used to describe
	// transferring hotspots from one wallet to another
	TransferHotspotV2Txn = "transfer_hotspot_v2"

	// GhostTxn is used to describe
	// a placeholder transaction that does not actually
	// exist on the Helium blockchain but is required
	// for account reconciliation
	GhostTxn = "ghost_txn"

	// CoinbaseOp is used to describe
	// the coinbase transaction at genesis (testnet only)
	CoinbaseOp = "coinbase_op"

	// CreditOp is used to describe
	// a credit to an account
	CreditOp = "credit_op"

	// DebitOp is used to describe
	// a debit from an account
	DebitOp = "debit_op"

	// RewardOp is used to describe
	// a blockchain rewarded to an account
	RewardOp = "reward_op"

	// RoutingOp is used to describe
	// routing packets on the network
	RoutingOp = "routing_op"

	// FeeOp is used to describe
	// a transaction fee
	FeeOp = "fee_op"

	// AddGatewayOp is used to describe
	// associating a new gateway with an owner
	AddGatewayOp = "add_gateway_op"

	// AssertLocationOp is used to describe
	// asserting the location of a hotspot
	AssertLocationOp = "assert_location_op"

	// TransferHotspotOp is used to describe
	// transferring a hotspot to a new owner
	TransferHotspotOp = "transfer_hotspot_op"

	// TokenBurnOp is used to describe
	// a burn of HNT for DC
	TokenBurnOp = "token_burn_op"

	// StakeValidatorOp is used to describe
	// staking a new validator
	StakeValidatorOp = "stake_validator_op"

	// UnstakeValidatorOp is used to describe
	// unstaking a validator
	UnstakeValidatorOp = "unstake_validator_op"

	// TransferValidatorStakeOp is used to describe
	// transferring a validator to a new owner and/or address
	TransferValidatorStakeOp = "transfer_validator_op"

	// StateChannelOpenOp is used to describe
	// opening a state channel
	StateChannelOpenOp = "state_channel_open_op"

	// OUIOp is used to describe
	// creating a new OUI
	OUIOp = "oui_op"

	// CreateHTLCOp is used to describe
	// creating an HTLC transaction
	CreateHTLCOp = "create_htlc_op"

	// CreateHTLCOp is used to describe
	// creating an HTLC transaction
	RedeemHTLCOp = "redeem_htlc_op"

	// UpdateGatewayOUIOp is used to describe
	// updating a gateway's OUI
	UpdateGatewayOUIOp = "update_gateway_oui_op"

	// PassthroughOp is used to describe
	// passthrough transactions
	PassthroughOp = "passthrough_op"

	// HistoricalBalanceSupported is whether
	// historical balance is supported.
	HistoricalBalanceSupported = true

	// Genesis is the index of the
	// genesis block for blockchain-etl instances
	GenesisBlockIndex = int64(1)

	// IncludeMempoolCoins does not apply to rosetta-ethereum as it is not UTXO-based.
	IncludeMempoolCoins = false
)

var (

	// MainnetNetworkBytes is the value of the
	// mainnet network in bytes
	MainnetNetworkBytes = []byte{0}

	// TestnetNetworkBytes is the value of the
	// testnet network in bytes
	TestnetNetworkBytes = []byte{1}

	// SyncedRocksDBHeight is the height of the
	// latest iterator
	SyncedRocksDBHeight = int64(0)

	// SuccessStatus is the status of any
	// Helium operation considered successful.
	SuccessStatus = "SUCCESS"

	// PendingStatus is the status of any
	// Helium operation considered unsuccessful.
	PendingStatus = "PENDING"

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
		AddGatewayV1Txn,
		AssertLocationV1Txn,
		AssertLocationV2Txn,
		VarsV1Txn,
		CoinbaseV1Txn,
		DCCoinbaseV1Txn,
		ConsensusGroupV1Txn,
		ConsensusGroupFailureV1Txn,
		CreateHTLCV1Txn,
		POCRequestV1Txn,
		TokenBurnV1Txn,
		TokenBurnExchangeRateV1Txn,
		GenGatewayV1Txn,
		OUIV1Txn,
		UpdateGatewayOUIV1Txn,
		PaymentV1Txn,
		PaymentV2Txn,
		POCReceiptsV1,
		RedeemHTLCV1Txn,
		RewardsV1Txn,
		RewardsV2Txn,
		RoutingV1Txn,
		SecurityCoinbaseV1Txn,
		SecurityExchangeV1Txn,
		StateChannelOpenV1Txn,
		StateChannelCloseV1Txn,
		StakeValidatorV1Txn,
		GenPriceOracleV1Txn,
		PriceOracleV1Txn,
		GenValidatorV1Txn,
		TransferHotspotV1Txn,
		StakeValidatorV1Txn,
		UnstakeValidatorV1Txn,
		TransferValidatorStakeV1Txn,
		ValidatorHeartbeatV1Txn,
		TransferHotspotV2Txn,
	}

	// OperationTypes are all supported base operations
	// that make up a transaction
	OperationTypes = []string{
		AddGatewayOp,
		AssertLocationOp,
		CoinbaseOp,
		CreateHTLCOp,
		CreditOp,
		DebitOp,
		RedeemHTLCOp,
		RewardOp,
		RoutingOp,
		FeeOp,
		TransferHotspotOp,
		TokenBurnOp,
		StakeValidatorOp,
		UnstakeValidatorOp,
		TransferValidatorStakeOp,
		StateChannelOpenOp,
		OUIOp,
		PassthroughOp,
		UpdateGatewayOUIOp,
	}

	// OperationStatuses are all supported operation statuses.
	OperationStatuses = []*types.OperationStatus{
		{
			Status:     SuccessStatus,
			Successful: true,
		},
		{
			Status:     PendingStatus,
			Successful: false,
		},
		{
			Status:     FailureStatus,
			Successful: false,
		},
	}

	// LBS is the LastBlessedBlock height as an int64
	LBS = readLBSfile()

	// Optional RocksDB vars for node db
	NodeBalancesDB                  *rocksdb.DB
	NodeBlocksDB                    *rocksdb.DB
	NodeTransactionsDB              *rocksdb.DB
	NodeBalancesDBEntriesHandle     *rocksdb.ColumnFamilyHandle
	NodeBalancesDBDefaultHandle     *rocksdb.ColumnFamilyHandle
	NodeBlockchainDBHeightsHandle   *rocksdb.ColumnFamilyHandle
	NodeTransactionsDBDefaultHandle *rocksdb.ColumnFamilyHandle
)

type Block struct {
	Hash         string                   `json:"hash"`
	Height       int64                    `json:"height"`
	PrevHash     string                   `json:"prev_hash"`
	Time         int64                    `json:"time"`
	Transactions []map[string]interface{} `json:"transactions"`
}

type Peer struct {
	Local  string `json:"local"`
	Name   string `json:"name"`
	P2P    string `json:"p2p"`
	Remote string `json:"remote"`
}

type Transaction struct {
	Hash string `json:"hash"`
	Type string `json:"type"`
}

type UnstakeTransaction struct {
	Address            string `json:"address"`
	Fee                int64  `json:"fee"`
	Hash               string `json:"hash"`
	Owner              string `json:"owner"`
	OwnerSignature     string `json:"owner_signature"`
	StakeAmount        int64  `json:"stake_amount"`
	StakeReleaseHeight int64  `json:"stake_release_height"`
}

type Payment struct {
	Payee  string `json:"payee"`
	Amount int64  `json:"amount"`
}

type Fee struct {
	Amount      int64
	Currency    *types.Currency
	Estimate    bool
	DCFeeAmount int64
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
