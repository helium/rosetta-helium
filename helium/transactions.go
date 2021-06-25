package helium

import (
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/utils"
)

func PaymentV1(payer, payee string, amount, fee int64, feeType string) ([]*types.Operation, *types.Error) {
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
	indexIncrementer := 2
	for i, p := range payments {
		PaymentDebit, pErr := CreateDebitOp(
			payer,
			p.Amount,
			HNT,
			int64(indexIncrementer*i),
			map[string]interface{}{"debit_category": "payment"},
		)
		if pErr != nil {
			return nil, pErr
		}

		PaymentCredit, pcErr := CreateCreditOp(
			p.Payee,
			p.Amount,
			HNT,
			int64((indexIncrementer*i)+1),
			map[string]interface{}{"credit_category": "payment"},
		)
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

func AddGatewayV1(
	payer string,
	feeTotal int64,
	feeType,
	gateway,
	owner string,
	metaBaseFee,
	metaStakingFee int64,
) ([]*types.Operation, *types.Error) {
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

func AssertLocationV1(
	metaBaseFee int64,
	gateway,
	location string,
	owner,
	payer string,
	metaStakingFee,
	feeTotal int64,
	feeType string,
) ([]*types.Operation, *types.Error) {
	var AssertLocationOps []*types.Operation

	feeOp, feeErr := CreateFeeOp(payer, feeTotal, feeType, 0, map[string]interface{}{"base_fee": metaBaseFee, "staking_fee": metaStakingFee})
	if feeErr != nil {
		return nil, feeErr
	}

	alOp, alErr := CreateAssertLocationOp(gateway, owner, location, 1, map[string]interface{}{})
	if alErr != nil {
		return nil, alErr
	}

	AssertLocationOps = append(AssertLocationOps, feeOp, alOp)
	return AssertLocationOps, nil
}

func AssertLocationV2(
	elevation,
	metaBaseFee,
	gain int64,
	gateway,
	location string,
	owner,
	payer string,
	metaStakingFee,
	feeTotal int64,
	feeType string,
) ([]*types.Operation, *types.Error) {
	var AssertLocationOps []*types.Operation

	feeOp, feeErr := CreateFeeOp(payer, feeTotal, feeType, 0, map[string]interface{}{"base_fee": metaBaseFee, "staking_fee": metaStakingFee})
	if feeErr != nil {
		return nil, feeErr
	}

	alOp, alErr := CreateAssertLocationOp(gateway, owner, location, 1, map[string]interface{}{
		"elevation": elevation,
		"gain":      gain,
	})
	if alErr != nil {
		return nil, alErr
	}

	AssertLocationOps = append(AssertLocationOps, feeOp, alOp)
	return AssertLocationOps, nil
}

func SecurityExchangeV1(payer, payee string, fee int64, feeType string, amount int64) ([]*types.Operation, *types.Error) {
	PaymentDebit, pErr := CreateDebitOp(payer, amount, HST, 0, map[string]interface{}{"credit_category": "payment"})
	if pErr != nil {
		return nil, pErr
	}

	PaymentCredit, pcErr := CreateCreditOp(payee, amount, HST, 1, map[string]interface{}{"debit_category": "payment"})
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
