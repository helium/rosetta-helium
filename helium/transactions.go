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
	case RewardTxnV1:
		return RewardV1(
			fmt.Sprint(txn["payee"]),
			utils.MapToInt64(txn["amount"]))
	default:
		return nil, WrapErr(ErrNotFound, errors.New("txn type not found"))
	}
}

func PaymentV1(payer string, payee string, amount int64, fee int64, feeType string) ([]*types.Operation, *types.Error) {

	PaymentDebit, pErr := CreatePaymentDebitOp(&payer, &amount, 0)
	if pErr != nil {
		return nil, pErr
	}

	PaymentCredit, pcErr := CreatePaymentCreditOp(&payee, &amount, 1)
	if pcErr != nil {
		return nil, pcErr
	}

	Fee, fErr := CreateFeeOp(&payer, &fee, &feeType, 2)
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		PaymentDebit,
		PaymentCredit,
		Fee,
	}, nil

}

func RewardV1(payee string, amount int64) ([]*types.Operation, *types.Error) {

	Reward, rErr := CreateRewardOp(&payee, &amount, 0)
	if rErr != nil {
		return nil, rErr
	}

	return []*types.Operation{
		Reward,
	}, nil

}
