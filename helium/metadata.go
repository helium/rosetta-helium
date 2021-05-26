package helium

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
)

type MetadataOptions struct {
	RequestedMetadata map[string]interface{} `json:"requested_metadata"`
	TransactionType   string                 `json:"transaction_type"`
}

func GetMetadata(request *types.ConstructionMetadataRequest) (*types.ConstructionMetadataResponse, *types.Error) {
	jsonString, _ := json.Marshal(request.Options)
	options := MetadataOptions{}
	err := json.Unmarshal(jsonString, &options)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("syntax error at byte offset %d", e.Offset)
		}
		return nil, WrapErr(ErrUnclearIntent, err)
	}

	for k, v := range options.RequestedMetadata {
		switch k {
		case "get_nonce_for":
			switch t := v.(type) {
			case map[string]interface{}:
				if v.(map[string]interface{})["address"] == nil {
					return nil, WrapErr(ErrUnclearIntent, errors.New("get_nonce_for requires `address` to be present in JSON object"))
				}

				nonce, nErr := GetNonce(fmt.Sprint(v.(map[string]interface{})["address"]))
				if nErr != nil {
					return nil, nErr
				}

				return &types.ConstructionMetadataResponse{
					Metadata: map[string]interface{}{
						"nonce": nonce,
					},
				}, nil
			default:
				return nil, WrapErr(ErrUnclearIntent, errors.New("unexpected object "+fmt.Sprint(t)+" in get_nonce_for"))
			}
		default:
			return nil, WrapErr(ErrUnclearIntent, errors.New("metadata request `"+fmt.Sprint(k)+"` not recognized"))
		}
	}

	return nil, nil
}
