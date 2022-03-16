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

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/helium/rosetta-helium/helium"
	"github.com/helium/rosetta-helium/services"
	"github.com/helium/rosetta-helium/utils"
	"go.uber.org/zap"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"

	badger "github.com/dgraph-io/badger/v3"

	rocksdb "github.com/linxGnu/grocksdb"
)

const (
	serverPort = 8080
)

type HeliumRocksDB struct {
	BalancesDB            *rocksdb.DB
	TransactionsDB        *rocksdb.DB
	BlockchainDB          *rocksdb.DB
	EntriesCF             *rocksdb.ColumnFamilyHandle
	HeightsCF             *rocksdb.ColumnFamilyHandle
	TransactionsDefaultCF *rocksdb.ColumnFamilyHandle
}

// NewBlockchainRouter creates a Mux http.Handler from a collection
// of server controllers.
func NewBlockchainRouter(
	network *types.NetworkIdentifier,
	a *asserter.Asserter,
) http.Handler {
	networkAPIService := services.NewNetworkAPIService(network)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		a,
	)

	blockAPIService := services.NewBlockAPIService(network)
	blockAPIController := server.NewBlockAPIController(
		blockAPIService,
		a,
	)

	accountAPIService := services.NewAccountAPIService(network)
	accountAPIController := server.NewAccountAPIController(
		accountAPIService,
		a,
	)

	constructionAPIService := services.NewConstructionAPIService(network)
	constructionAPIController := server.NewConstructionAPIController(
		constructionAPIService,
		a,
	)

	return server.NewRouter(networkAPIController, blockAPIController, accountAPIController, constructionAPIController)
}

func LoadGhostTxns(network *types.NetworkIdentifier, db *badger.DB) error {
	// Assume mainnet
	networkDir := "mainnet"

	if network.Network == helium.TestnetNetwork {
		networkDir = "testnet"
	}

	if werr := filepath.Walk("ghost-transactions/"+networkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			zap.S().Error(err.Error())
		}
		if info.IsDir() {
			return nil
		}

		jsonFile, jerr := os.Open("ghost-transactions/" + networkDir + "/" + info.Name())
		if jerr != nil {
			return jerr
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var result []map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)

		for _, txn := range result {
			var txnMetadata helium.UnstakeTransaction
			var txnMetadataMap map[string]interface{}
			json.Unmarshal([]byte(fmt.Sprint(txn["Fields"])), &txnMetadata)
			json.Unmarshal([]byte(fmt.Sprint(txn["Fields"])), &txnMetadataMap)

			dbKey := &utils.GhostTxnKey{
				Network: network,
				Block: &types.BlockIdentifier{
					Index: txnMetadata.StakeReleaseHeight,
				},
				Transaction: &types.TransactionIdentifier{
					Hash: txnMetadata.Hash,
				},
			}

			Unstake, uErr := helium.CreateCreditOp(helium.UnstakeValidatorOp, txnMetadata.Address, txnMetadata.StakeAmount, helium.HNT, helium.SuccessStatus, 0, txnMetadataMap)
			if uErr != nil {
				return errors.New(fmt.Sprint(uErr))
			}

			cerr := utils.CreateGhostTxn(dbKey, &utils.GhostTxnMetadata{
				Operations: []*types.Operation{
					Unstake,
				},
				Metadata: txnMetadataMap,
			})

			if cerr != nil && cerr != badger.ErrBannedKey {
				return cerr
			}
		}

		zap.S().Info("Loaded " + fmt.Sprint(len(result)) + " ghost txns from file " + info.Name())
		return nil

	}); werr != nil {
		return werr
	}

	return nil
}

func openRocksDB(dataDir string) (heliumDB *HeliumRocksDB, err error) {
	opts := rocksdb.NewDefaultOptions()
	columnOpts := rocksdb.NewDefaultOptions()

	prefixExtractor := rocksdb.NewFixedPrefixTransform(33)
	columnOpts.SetPrefixExtractor(prefixExtractor)

	balancesCfNames := []string{"default", "entries"}
	blockchainCfNames := []string{"default", "heights"}
	transactionsCfNames := []string{"default", "transactions"}

	dbBal, balancesCfHandles, dbBalErr := rocksdb.OpenDbAsSecondaryColumnFamilies(
		opts,
		dataDir+"/balances.db",
		"rocksdb/balances.db",
		balancesCfNames,
		[]*rocksdb.Options{opts, columnOpts},
	)
	if dbBalErr != nil {
		return nil, dbBalErr
	}

	dbBlock, blockchainCfHandles, dbBlockErr := rocksdb.OpenDbAsSecondaryColumnFamilies(
		opts,
		dataDir+"/blockchain.db",
		"rocksdb/blockchain.db",
		blockchainCfNames,
		[]*rocksdb.Options{opts, opts},
	)
	if dbBlockErr != nil {
		return nil, dbBalErr
	}

	dbTxn, txnCfHandles, dbTxnErr := rocksdb.OpenDbAsSecondaryColumnFamilies(
		opts,
		dataDir+"/transactions.db",
		"rocksdb/transactions.db",
		transactionsCfNames,
		[]*rocksdb.Options{opts, opts},
	)
	if dbTxnErr != nil {
		return nil, dbTxnErr
	}

	return &HeliumRocksDB{
		BalancesDB:            dbBal,
		TransactionsDB:        dbTxn,
		BlockchainDB:          dbBlock,
		EntriesCF:             balancesCfHandles[1],
		HeightsCF:             blockchainCfHandles[1],
		TransactionsDefaultCF: txnCfHandles[0],
	}, nil
}

func main() {
	// Logging tool
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	globalLogger := zap.ReplaceGlobals(logger)
	defer globalLogger()

	// Ghost transaction DB setup
	bdb, err := badger.Open(badger.DefaultOptions("badger"))
	if err != nil {
		zap.S().Fatal(err)
	}
	defer bdb.Close()
	utils.DB = bdb

	// CLI Flag Parsing
	// ****************
	//
	// blockchain-node data dir direct connection CLI arg
	var blockchainNodeDataDir string
	flag.StringVar(&blockchainNodeDataDir, "data", "", "path to blockchain-node data dir")

	// Testnet CLI arg
	var testnet bool
	flag.BoolVar(&testnet, "testnet", false, "run testnet version of rosetta-helium")

	// Parse flags
	flag.Parse()

	// Network setup
	var network *types.NetworkIdentifier
	if !testnet {
		zap.S().Info("Initilizing mainnet node...")
		network = &types.NetworkIdentifier{
			Blockchain: "Helium",
			Network:    helium.MainnetNetwork,
		}

		if lerr := LoadGhostTxns(network, bdb); lerr != nil {
			zap.S().Error("Cannot load mainnet ghost transactions: " + lerr.Error())
			os.Exit(1)
		}

	} else {
		zap.S().Info("Initilizing testnet node...")
		network = &types.NetworkIdentifier{
			Blockchain: "Helium",
			Network:    helium.TestnetNetwork,
		}

		if lerr := LoadGhostTxns(network, bdb); lerr != nil {
			zap.S().Error("Cannot load testnet ghost transactions: " + lerr.Error())
			os.Exit(1)
		}
	}

	helium.CurrentNetwork = network

	// Blockchain-node rocksdb direct connection setup
	if blockchainNodeDataDir != "" {
		// Set total number of retries before canceling DB open
		retries := 10

		zap.S().Info("Attempting to load rocksdb directly at dir '" + blockchainNodeDataDir + "'")

		for retries > 0 {
			heliumDB, dbOpenErr := openRocksDB(blockchainNodeDataDir)
			if dbOpenErr != nil {
				zap.S().Warn(dbOpenErr.Error() + ": Unable to open rocksdb at this time. Retrying - (" + fmt.Sprint(retries) + " attemps left)")
				time.Sleep(5 * time.Second)
			} else {
				zap.S().Info("Loaded rocksdb directly at dir '" + blockchainNodeDataDir + "'")
				helium.NodeBalancesDB = heliumDB.BalancesDB
				helium.NodeBlocksDB = heliumDB.BlockchainDB
				helium.NodeTransactionsDB = heliumDB.TransactionsDB
				helium.NodeBalancesDBEntriesHandle = heliumDB.EntriesCF
				helium.NodeBlockchainDBHeightsHandle = heliumDB.HeightsCF
				helium.NodeTransactionsDBDefaultHandle = heliumDB.TransactionsDefaultCF
				break
			}
		}
	}

	// The asserter automatically rejects incorrectly formatted
	// requests.
	a, err := asserter.NewServer(
		helium.OperationTypes,
		helium.HistoricalBalanceSupported,
		[]*types.NetworkIdentifier{network},
		nil,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create the main router handler then apply the logger and Cors
	// middlewares in sequence.
	router := NewBlockchainRouter(network, a)
	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)
	zap.S().Info("Listening on port ", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), corsRouter))
}
