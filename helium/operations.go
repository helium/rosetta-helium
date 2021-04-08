package helium

import (
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/utils"
)

func ParsePaymentV1Txn(tx map[string]interface{}) (*types.Operation, *types.Error) {
	status := SuccessStatus
	tx_block := utils.MapToInt64(tx["block"])
	hnt_price, hErr := GetOraclePrice(tx_block)
	if hErr != nil {
		return nil, hErr
	}

	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 0,
		},
		Type:   PaymentV1Txn,
		Status: &status,
		Account: &types.AccountIdentifier{
			Address: fmt.Sprint(tx["payer"]),
		},
	}, nil
}
