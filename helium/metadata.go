package helium

import (
	"encoding/json"
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
			fmt.Println(v)
		default:
			return nil, nil
		}
	}

	return nil, nil
}
