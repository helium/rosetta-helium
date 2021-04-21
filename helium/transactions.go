package helium

import (
	"errors"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/utils"
)

type Fee struct {
	Amount   int64
	Payer    string
	Currency *types.Currency
}

func OperationsFromTx(txn map[string]interface{}) ([]*types.Operation, *types.Error) {
	switch txn["type"] {
	case PaymentV1Txn:
		feeDetails, bErr := GetImplicitBurn(fmt.Sprint(txn["hash"]))
		if bErr != nil {
			feeDetails = &Fee{
				Amount:   utils.MapToInt64(txn["fee"]),
				Payer:    fmt.Sprint(txn["payer"]),
				Currency: DC,
			}
		}
		return PaymentV1(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["payee"]),
			utils.MapToInt64(txn["amount"]),
			feeDetails.Amount,
			feeDetails.Currency.Symbol)
	default:
		return nil, WrapErr(ErrNotFound, errors.New("txn type not found"))
	}
}

func PaymentV1(payer string, payee string, amount int64, fee int64, feeType string) ([]*types.Operation, *types.Error) {

	PaymentDebit, pErr := CreatePaymentDebitOp(&payer, &amount, 0)
	if pErr != nil {
		return nil, pErr
	}

	Fee, fErr := CreateFeeOp(&payer, &fee, &feeType, 1)
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		PaymentDebit,
		Fee,
	}, nil

}
