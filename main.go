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

	"github.com/helium/rosetta-helium/helium"
	"github.com/helium/rosetta-helium/services"
	"github.com/helium/rosetta-helium/utils"
	"go.uber.org/zap"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"

	badger "github.com/dgraph-io/badger/v3"
)

const (
	serverPort = 8080
)

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
	if werr := filepath.Walk("ghost-transactions", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			zap.S().Error(err.Error())
		}
		if info.IsDir() {
			return nil
		}

		jsonFile, jerr := os.Open("ghost-transactions/" + info.Name())
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

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any

	globalLogger := zap.ReplaceGlobals(logger)
	defer globalLogger()

	bdb, err := badger.Open(badger.DefaultOptions("badger"))
	if err != nil {
		zap.S().Fatal(err)
	}
	defer bdb.Close()

	utils.DB = bdb

	var testnet bool
	var network *types.NetworkIdentifier

	flag.BoolVar(&testnet, "testnet", false, "run testnet version of rosetta-helium")
	flag.Parse()

	if !testnet {
		zap.S().Info("Initilizing mainnet node...")
		network = &types.NetworkIdentifier{
			Blockchain: "Helium",
			Network:    helium.MainnetNetwork,
		}

		if lerr := LoadGhostTxns(network, bdb); lerr != nil {
			zap.S().Error("Cannot load ghost transactions: " + lerr.Error())
			os.Exit(1)
		}

	} else {
		zap.S().Info("Initilizing testnet node...")
		network = &types.NetworkIdentifier{
			Blockchain: "Helium",
			Network:    helium.TestnetNetwork,
		}
	}

	helium.CurrentNetwork = network

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
