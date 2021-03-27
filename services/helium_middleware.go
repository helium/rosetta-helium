package services

import (
	"fmt"
	"log"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/helium"
	"github.com/ybbus/jsonrpc"
)

var (
	NodeClient = jsonrpc.NewClient("http://localhost:4467")
)

func CurrentBlockHeight() int64 {

	var result int64

	fmt.Print("Getting current block height: ")

	if err := NodeClient.CallFor(&result, "block_height", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	return result
}

func GetBlock(index int64) *types.Block {

	fmt.Print("Getting block at height ")
	fmt.Println(index)
	type request struct {
		Height int64 `json:"height"`
	}

	var result helium.Block

	req := request{Height: index}
	if err := NodeClient.CallFor(&result, "block_get", req); err != nil {
		log.Fatal(err)
	}

	var processedTxs []*types.Transaction
	for _, tx := range result.Transactions {
		ptx := GetTransaction(tx)
		processedTxs = append(processedTxs, ptx)
	}

	currentBlock := &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: result.Height,
			Hash:  result.Hash,
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: result.Height,
			Hash:  result.Hash,
		},
		Timestamp:    result.Time,
		Transactions: processedTxs,
		Metadata:     nil,
	}

	return currentBlock
}

func GetTransaction(txHash string) *types.Transaction {

	type request struct {
		Hash string `json:"hash"`
	}

	var result map[string]interface{}

	req := request{Hash: txHash}
	if err := NodeClient.CallFor(&result, "transaction_get", req); err != nil {
		log.Fatal(err)
	}

	transaction := &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: fmt.Sprint(result["hash"]),
		},
		Operations: []*types.Operation{
			ParseOperationFromTx(result, 0),
		},
		RelatedTransactions: nil,
		Metadata:            nil,
	}

	return transaction

}

func ParseOperationFromTx(tx map[string]interface{}, index int64) *types.Operation {
	txType := tx["type"]
	status := helium.SuccessStatus

	operation := &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: index,
		},
		RelatedOperations: nil,
		Type:              fmt.Sprint(txType),
		Status:            &status,
		Account:           nil,
		Amount:            nil,
		CoinChange:        nil,
		Metadata:          nil,
	}

	return operation
}

func GetAmount(address string) *types.Amount {

	type request struct {
		Address string `json:"address"`
	}

	var result map[string]interface{}

	req := request{Address: address}
	if err := NodeClient.CallFor(&result, "account_get", req); err != nil {
		log.Fatal(err)
	}

	amount := &types.Amount{
		Value: fmt.Sprint(result["balance"]),
		Currency: &types.Currency{
			Symbol:   helium.Currency.Symbol,
			Decimals: helium.Currency.Decimals,
		},
	}

	return amount

}
