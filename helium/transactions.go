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
	case PaymentV2Txn:
		feeDetails, bErr := GetImplicitBurn(fmt.Sprint(txn["hash"]))
		if bErr != nil {
			feeDetails = &Fee{
				Amount:   utils.MapToInt64(txn["fee"]),
				Payer:    fmt.Sprint(txn["payer"]),
				Currency: DC,
			}
		}
		var payments []*Payment
		for _, p := range txn["payments"].([]map[string]interface{}) {
			payments = append(payments, &Payment{
				Payee:  fmt.Sprint(p["payee"]),
				Amount: utils.MapToInt64(p["amount"]),
			})
		}
		return PaymentV2(
			fmt.Sprint(txn["payer"]),
			payments,
			feeDetails.Amount,
			feeDetails.Currency.Symbol,
		)
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

func PaymentV2(payer string, payments []*Payment, fee int64, feeType string) ([]*types.Operation, *types.Error) {

	var paymentV2Operations []*types.Operation

	for i, p := range payments {
		PaymentDebit, pErr := CreatePaymentDebitOp(&p.Payee, &p.Amount, int64(2*i))
		if pErr != nil {
			return nil, pErr
		}

		PaymentCredit, pcErr := CreatePaymentCreditOp(&p.Payee, &p.Amount, int64((2*i)+1))
		if pcErr != nil {
			return nil, pcErr
		}

		paymentV2Operations = append(paymentV2Operations, PaymentDebit, PaymentCredit)
	}

	Fee, fErr := CreateFeeOp(&payer, &fee, &feeType, int64(len(paymentV2Operations)))
	if fErr != nil {
		return nil, fErr
	}

	paymentV2Operations = append(paymentV2Operations, Fee)

	return paymentV2Operations, nil

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
