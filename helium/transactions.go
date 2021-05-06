package helium

import (
	"errors"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/utils"
)

func OperationsFromTx(txn map[string]interface{}) ([]*types.Operation, *types.Error) {
	switch txn["type"] {
	case AddGatewayTxn:
		feeDetails := GetFee(fmt.Sprint(txn["hash"]), utils.MapToInt64(txn["fee"])+utils.MapToInt64(txn["staking_fee"]), fmt.Sprint(txn["payer"]))
		return AddGatewayV1(
			fmt.Sprint(txn["payer"]),
			feeDetails.Amount,
			feeDetails.Currency.Symbol,
			fmt.Sprint(txn["gateway"]),
			fmt.Sprint(txn["owner"]),
			utils.MapToInt64(txn["fee"]),
			utils.MapToInt64(txn["staking_fee"]))
	case PaymentV1Txn:
		feeDetails := GetFee(fmt.Sprint(txn["hash"]), utils.MapToInt64(txn["fee"]), fmt.Sprint(txn["payer"]))
		return PaymentV1(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["payee"]),
			int64(txn["amount"].(float64)),
			feeDetails.Amount,
			feeDetails.Currency.Symbol)
	case PaymentV2Txn:
		feeDetails := GetFee(fmt.Sprint(txn["hash"]), utils.MapToInt64(txn["fee"]), fmt.Sprint(txn["payer"]))
		var payments []*Payment
		for _, p := range txn["payments"].([]interface{}) {
			payments = append(payments, &Payment{
				Payee:  fmt.Sprint(p.(map[string]interface{})["payee"]),
				Amount: utils.MapToInt64(int64(p.(map[string]interface{})["amount"].(float64))),
			})
		}
		return PaymentV2(
			fmt.Sprint(txn["payer"]),
			payments,
			feeDetails.Amount,
			feeDetails.Currency.Symbol,
		)
	case RewardsTxnV1:
		return RewardsV1(txn["rewards"].([]interface{}))
	case SecurityCoinbaseTxn:
		return SecurityCoinbaseV1(fmt.Sprint(txn["payee"]), int64(txn["amount"].(float64)))
	case CoinbaseDataCreditsTxn:
		return DCCoinbaseV1(fmt.Sprint(txn["payee"]), int64(txn["amount"].(float64)))
	default:
		return nil, WrapErr(ErrNotFound, errors.New("txn type not found"))
	}
}

func PaymentV1(payer string, payee string, amount int64, fee int64, feeType string) ([]*types.Operation, *types.Error) {
	PaymentDebit, pErr := CreateDebitOp(payer, amount, HNT, 0, map[string]interface{}{"credit_category": "payment"})
	if pErr != nil {
		return nil, pErr
	}

	PaymentCredit, pcErr := CreateCreditOp(payee, amount, HNT, 1, map[string]interface{}{"debit_category": "payment"})
	if pcErr != nil {
		return nil, pcErr
	}

	Fee, fErr := CreateFeeOp(payer, fee, feeType, 2, map[string]interface{}{})
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
		PaymentDebit, pErr := CreateDebitOp(p.Payee, p.Amount, HNT, int64(2*i), map[string]interface{}{"credit_category": "payment"})
		if pErr != nil {
			return nil, pErr
		}

		PaymentCredit, pcErr := CreateCreditOp(p.Payee, p.Amount, HNT, int64((2*i)+1), map[string]interface{}{"debit_category": "payment"})
		if pcErr != nil {
			return nil, pcErr
		}

		paymentV2Operations = append(paymentV2Operations, PaymentDebit, PaymentCredit)
	}

	Fee, fErr := CreateFeeOp(payer, fee, feeType, int64(len(paymentV2Operations)), map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	paymentV2Operations = append(paymentV2Operations, Fee)

	return paymentV2Operations, nil

}

func RewardsV1(rewards []interface{}) ([]*types.Operation, *types.Error) {

	var rewardOps []*types.Operation

	for i, reward := range rewards {
		rewardOp, rErr := CreateCreditOp(
			fmt.Sprint(reward.(map[string]interface{})["account"]),
			utils.MapToInt64(int64(reward.(map[string]interface{})["amount"].(float64))),
			HNT,
			int64(i),
			map[string]interface{}{
				"credit_category": "reward",
				"gateway":         fmt.Sprint(reward.(map[string]interface{})["gateway"]),
				"type":            fmt.Sprint(reward.(map[string]interface{})["type"]),
			})
		if rErr != nil {
			return nil, rErr
		}

		rewardOps = append(rewardOps, rewardOp)
	}

	return rewardOps, nil

}

func SecurityCoinbaseV1(payee string, amount int64) ([]*types.Operation, *types.Error) {
	var securityCoinbaseOps []*types.Operation

	secOps, secErr := CreateCreditOp(payee, amount, HST, 0, map[string]interface{}{"credit_category": "security_coinbase"})
	if secErr != nil {
		return nil, secErr
	}

	securityCoinbaseOps = append(securityCoinbaseOps, secOps)

	return securityCoinbaseOps, nil

}

func DCCoinbaseV1(payee string, amount int64) ([]*types.Operation, *types.Error) {
	var DCCoinbaseOps []*types.Operation

	dccOps, dccErr := CreateCreditOp(payee, amount, DC, 0, map[string]interface{}{"credit_category": "dc_coinbase"})
	if dccErr != nil {
		return nil, dccErr
	}

	DCCoinbaseOps = append(DCCoinbaseOps, dccOps)

	return DCCoinbaseOps, nil

}

func AddGatewayV1(payer string, feeTotal int64, feeType string, gateway string, owner string, metaBaseFee int64, metaStakingFee int64) ([]*types.Operation, *types.Error) {
	var AddGatewayOps []*types.Operation

	feeOp, feeErr := CreateFeeOp(payer, feeTotal, feeType, 0, map[string]interface{}{"base_fee": metaBaseFee, "staking_fee": metaStakingFee})
	if feeErr != nil {
		return nil, feeErr
	}

	agwOp, agwErr := CreateAddGatewayOp(gateway, owner, 1, map[string]interface{}{})
	if agwErr != nil {
		return nil, agwErr
	}

	AddGatewayOps = append(AddGatewayOps, feeOp, agwOp)

	return AddGatewayOps, nil

}
