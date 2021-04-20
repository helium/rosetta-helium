package helium

import (
	"github.com/coinbase/rosetta-sdk-go/types"
)

func OperationsFromTx(txn map[string]interface{}) ([]*types.Operation, *types.Error) {

}

func PaymentV1(payer *string, payee *string, amount *int64, fee *int64, feeType *string) ([]*types.Operation, *types.Error) {

	PaymentDebit, pErr := CreatePaymentDebitOp(payer, amount, 0)
	if pErr != nil {
		return nil, pErr
	}

	Fee, fErr := CreateFeeOp(payer, fee, feeType, 1)
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		PaymentDebit,
		Fee,
	}, nil

}
