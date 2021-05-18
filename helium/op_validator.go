package helium

import (
	"github.com/coinbase/rosetta-sdk-go/types"
)

type Preprocessor struct {
	TransactionType   string                 `json:"transaction_type"`
	RequestedMetadata map[string]interface{} `json:"requested_metadata"`
}

func OpsToTransaction(operation []*types.Operation) (*Preprocessor, *types.Error) {
	return nil, nil
}
