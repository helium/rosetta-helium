package helium

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/coinbase/rosetta-sdk-go/types"
)

type combination struct {
	UnsignedTransaction string             `json:"unsigned_transaction"`
	Signatures          []*types.Signature `json:"signatures"`
}

func CombineTransaction(unsignedTxn string, signatures []*types.Signature) (*types.ConstructionCombineResponse, *types.Error) {

	jsonObject, jErr := json.Marshal(combination{
		UnsignedTransaction: unsignedTxn,
		Signatures:          signatures,
	})
	if jErr != nil {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New(`unable to decode combination object into json`))
	}

	// fmt.Println(jsonObject)

	var payload map[string]interface{}
	resp, ctErr := http.Post("http://localhost:3000/combine-tx", "application/json", bytes.NewBuffer(jsonObject))
	if ctErr != nil {
		return nil, WrapErr(ErrUnclearIntent, ctErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&payload)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}

	signedTransaction := payload["signed_transaction"].(string)

	return &types.ConstructionCombineResponse{
		SignedTransaction: signedTransaction,
	}, nil
}
