package helium

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/coinbase/rosetta-sdk-go/types"
)

func PayloadGenerator(operations []*types.Operation, metadata map[string]interface{}) (*types.ConstructionPayloadsResponse, *types.Error) {

	transactionPreprocessor, err := OpsToTransaction(operations)
	if err != nil {
		return nil, err
	}

	var operationMetadata map[string]interface{}
	marshalledPreprocessor, _ := json.Marshal(transactionPreprocessor)
	json.Unmarshal(marshalledPreprocessor, &operationMetadata)

	if !reflect.DeepEqual(metadata["options"], operationMetadata) {
		return nil, WrapErr(ErrUnclearIntent, errors.New(`payload operations options result do not match provided metadata options (metadata["options"])`))
	}

	jsonValue, jErr := json.Marshal(metadata)
	if jErr != nil {
		fmt.Print(jErr)
	}

	var payload map[string]interface{}
	resp, ctErr := http.Post("http://localhost:3000/create-tx", "application/json", bytes.NewBuffer(jsonValue))
	if ctErr != nil {
		return nil, WrapErr(ErrUnclearIntent, ctErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&payload)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}

	decodedByteArray, hErr := hex.DecodeString(payload["payload"].(string))
	if hErr != nil {
		return nil, WrapErr(ErrUnableToParseTxn, hErr)
	}

	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: payload["unsigned_txn"].(string),
		Payloads: []*types.SigningPayload{
			{
				Bytes: decodedByteArray,
			},
		},
	}, nil
}
