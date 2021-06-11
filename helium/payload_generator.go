package helium

import (
	"errors"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
)

func PayloadGenerator(operations []*types.Operation, metadata MetadataOptions) (*types.ConstructionPayloadsResponse, *types.Error) {

	switch metadata.TransactionType {
	case PaymentV2Txn:
		fmt.Print(PaymentV2Txn)
	default:
		return nil, WrapErr(ErrNotFound, errors.New(`invalid tranasction_type in metadata`))
	}

	return nil, nil
}
