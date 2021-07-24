package helium

import (
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/utils"
)

func PaymentV1(payer, payee string, amount int64, fee *Fee) ([]*types.Operation, *types.Error) {
	PaymentDebit, pErr := CreateDebitOp(payer, amount, HNT, 0, map[string]interface{}{"debit_category": "payment"})
	if pErr != nil {
		return nil, pErr
	}

	PaymentCredit, pcErr := CreateCreditOp(payee, amount, HNT, 1, map[string]interface{}{"credit_category": "payment"})
	if pcErr != nil {
		return nil, pcErr
	}

	Fee, fErr := CreateFeeOp(payer, fee, 2, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		PaymentDebit,
		PaymentCredit,
		Fee,
	}, nil
}

func PaymentV2(payer string, payments []*Payment, fee *Fee) ([]*types.Operation, *types.Error) {
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

	Fee, fErr := CreateFeeOp(payer, fee, int64(len(paymentV2Operations)), map[string]interface{}{})
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
	fee *Fee,
	gateway,
	owner string,
	baseFee,
	stakingFee int64,
) ([]*types.Operation, *types.Error) {
	var AddGatewayOps []*types.Operation

	agwOp, agwErr := CreateAddGatewayOp(gateway, owner, 1, map[string]interface{}{})
	if agwErr != nil {
		return nil, agwErr
	}

	feePayer := owner
	if (payer != owner) && (payer != "1Wh4bh") && (payer != "") {
		feePayer = payer
	}

	feeOp, feeErr := CreateFeeOp(feePayer, fee, 0, map[string]interface{}{"fee": baseFee, "staking_fee": stakingFee})
	if feeErr != nil {
		return nil, feeErr
	}

	AddGatewayOps = append(AddGatewayOps, agwOp, feeOp)

	return AddGatewayOps, nil
}

func AssertLocationV1(
	baseFee int64,
	gateway,
	location string,
	owner,
	payer string,
	stakingFee int64,
	fee *Fee,
) ([]*types.Operation, *types.Error) {
	var AssertLocationOps []*types.Operation

	alOp, alErr := CreateAssertLocationOp(gateway, owner, location, 1, map[string]interface{}{})
	if alErr != nil {
		return nil, alErr
	}

	feePayer := owner
	if (payer != owner) && (payer != "1Wh4bh") && (payer != "") {
		feePayer = payer
	}

	feeOp, feeErr := CreateFeeOp(feePayer, fee, 0, map[string]interface{}{"fee": baseFee, "staking_fee": stakingFee})
	if feeErr != nil {
		return nil, feeErr
	}

	AssertLocationOps = append(AssertLocationOps, alOp, feeOp)
	return AssertLocationOps, nil
}

func AssertLocationV2(
	elevation,
	baseFee,
	gain int64,
	gateway,
	location string,
	owner,
	payer string,
	stakingFee int64,
	fee *Fee,
) ([]*types.Operation, *types.Error) {
	var AssertLocationOps []*types.Operation

	alOp, alErr := CreateAssertLocationOp(gateway, owner, location, 0, map[string]interface{}{
		"elevation": elevation,
		"gain":      gain,
	})
	if alErr != nil {
		return nil, alErr
	}

	feePayer := owner
	if (payer != owner) && (payer != "1Wh4bh") && (payer != "") {
		feePayer = payer
	}

	feeOp, feeErr := CreateFeeOp(feePayer, fee, 1, map[string]interface{}{"fee": baseFee, "staking_fee": stakingFee})
	if feeErr != nil {
		return nil, feeErr
	}

	AssertLocationOps = append(AssertLocationOps, alOp, feeOp)
	return AssertLocationOps, nil
}

func SecurityExchangeV1(payer, payee string, fee *Fee, amount int64) ([]*types.Operation, *types.Error) {
	PaymentDebit, pErr := CreateDebitOp(payer, amount, HST, 0, map[string]interface{}{"credit_category": "payment"})
	if pErr != nil {
		return nil, pErr
	}

	PaymentCredit, pcErr := CreateCreditOp(payee, amount, HST, 1, map[string]interface{}{"debit_category": "payment"})
	if pcErr != nil {
		return nil, pcErr
	}

	Fee, fErr := CreateFeeOp(payer, fee, 2, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		PaymentDebit,
		PaymentCredit,
		Fee,
	}, nil
}

func TransferHotspotV1(
	amountToSeller int64,
	buyer string,
	fee *Fee,
	gateway string,
	seller string,
) ([]*types.Operation, *types.Error) {
	ops := []*types.Operation{}
	index := int64(0)

	TransferHotspot, tErr := CreateTransferHotspotOp(buyer, seller, gateway, index, map[string]interface{}{})
	if tErr != nil {
		return nil, tErr
	}
	index++
	ops = append(ops, TransferHotspot)

	if amountToSeller > 0 {
		Debit, dErr := CreateDebitOp(buyer, amountToSeller, HNT, index, map[string]interface{}{})
		if dErr != nil {
			return nil, dErr
		}
		index++

		Credit, cErr := CreateCreditOp(seller, amountToSeller, HNT, index, map[string]interface{}{})
		if cErr != nil {
			return nil, cErr
		}
		index++

		ops = append(ops, Debit, Credit)
	}

	Fee, fErr := CreateFeeOp(buyer, fee, index, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	ops = append(ops, Fee)

	return ops, nil
}

func TokenBurnV1(
	payer string,
	payee string,
	memo string,
	amount int64,
	fee *Fee,
) ([]*types.Operation, *types.Error) {
	TokenBurn, tErr := CreateTokenBurnOp(payer, payee, amount, 0, map[string]interface{}{})
	if tErr != nil {
		return nil, tErr
	}

	Fee, fErr := CreateFeeOp(payer, fee, 1, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		TokenBurn,
		Fee,
	}, nil
}

func StakeValidatorV1(
	owner string,
	ownerSignature string,
	address string,
	stake int64,
	fee *Fee,
) ([]*types.Operation, *types.Error) {
	StakeValidator, sErr := CreateStakeValidatorOp(owner, ownerSignature, address, stake, 0, map[string]interface{}{})
	if sErr != nil {
		return nil, sErr
	}

	Fee, fErr := CreateFeeOp(owner, fee, 1, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		StakeValidator,
		Fee,
	}, nil
}

func UnstakeValidatorV1(
	owner string,
	ownerSignature string,
	address string,
	stake int64,
	releaseHeight int64,
	fee *Fee,
) ([]*types.Operation, *types.Error) {
	stakeStatus := PendingStatus
	currentHeight, cErr := GetCurrentHeight()
	if cErr != nil {
		return nil, cErr
	}

	if releaseHeight >= *currentHeight {
		stakeStatus = SuccessStatus
	}

	Unstake, uErr := CreateUnstakeValidatorOp(owner, ownerSignature, address, stake, stakeStatus, 0, map[string]interface{}{})
	if uErr != nil {
		return nil, uErr
	}

	Fee, fErr := CreateFeeOp(owner, fee, 1, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		Unstake,
		Fee,
	}, nil
}

func TransferValidatorStakeV1(
	newOwner string,
	oldOwner string,
	newAddress string,
	oldAddress string,
	newOwnerSignature string,
	oldOwnerSignature string,
	stakeAmount int64,
	paymentAmount int64,
	fee *Fee,
) ([]*types.Operation, *types.Error) {
	ops := []*types.Operation{}
	index := int64(0)
	TransferValidator, tErr := CreateTransferValidatorOp(newOwner, oldOwner, newAddress, oldAddress, newOwnerSignature, oldOwnerSignature, stakeAmount, index, map[string]interface{}{})
	if tErr != nil {
		return nil, tErr
	}
	index++
	ops = append(ops, TransferValidator)

	if paymentAmount > int64(0) {
		Debit, dErr := CreateDebitOp(newOwner, paymentAmount, HNT, index, map[string]interface{}{})
		if dErr != nil {
			return nil, dErr
		}
		index++

		Credit, cErr := CreateCreditOp(oldOwner, paymentAmount, HNT, index, map[string]interface{}{})
		if cErr != nil {
			return nil, cErr
		}
		index++

		ops = append(ops, Debit, Credit)
	}

	Fee, fErr := CreateFeeOp(oldOwner, fee, index, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	ops = append(ops, Fee)

	return ops, nil
}

func OUIV1(
	oui int64,
	owner string,
	payer string,
	filter string,
	addresses []string,
	baseFee int64,
	stakingFee int64,
	requestedSubnetSize int64,
	fee *Fee,
) ([]*types.Operation, *types.Error) {
	ops := []*types.Operation{}
	OUI, oErr := CreateOUIOp(oui, owner, payer, filter, addresses, requestedSubnetSize, 0, map[string]interface{}{})
	if oErr != nil {
		return nil, oErr
	}

	feePayer := owner
	if (payer != owner) && (payer != "1Wh4bh") && (payer != "") {
		feePayer = payer
	}

	Fee, fErr := CreateFeeOp(feePayer, fee, 1, map[string]interface{}{"fee": baseFee, "staking_fee": stakingFee})
	if fErr != nil {
		return nil, fErr
	}

	ops = append(ops, OUI, Fee)

	return ops, nil
}
