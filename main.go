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
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/helium/rosetta-helium/helium"
	"github.com/helium/rosetta-helium/services"
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

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any

	globalLogger := zap.ReplaceGlobals(logger)
	defer globalLogger()

	db, err := badger.Open(badger.DefaultOptions("badger"))
	if err != nil {
		zap.S().Fatal(err)
	}
	defer db.Close()

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
	} else {
		zap.S().Info("Initilizing testnet node...")
		network = &types.NetworkIdentifier{
			Blockchain: "Helium",
			Network:    helium.TestnetNetwork,
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
