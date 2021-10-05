package helium

import (
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/helium/rosetta-helium/utils"
)

func PaymentV1(payer, payee string, amount int64, fee *Fee) ([]*types.Operation, *types.Error) {
	PaymentDebit, pErr := CreateDebitOp(DebitOp, payer, amount, HNT, SuccessStatus, 0, map[string]interface{}{"debit_category": "payment"})
	if pErr != nil {
		return nil, pErr
	}

	PaymentCredit, pcErr := CreateCreditOp(CreditOp, payee, amount, HNT, SuccessStatus, 1, map[string]interface{}{"credit_category": "payment"})
	if pcErr != nil {
		return nil, pcErr
	}

	Fee, fErr := CreateFeeOp(payer, fee, SuccessStatus, 2, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		PaymentDebit,
		PaymentCredit,
		Fee,
	}, nil
}

func PaymentV2(payer string, payments []*Payment, fee *Fee, statusString string) ([]*types.Operation, *types.Error) {
	var paymentV2Operations []*types.Operation
	indexIncrementer := 2

	for i, p := range payments {
		PaymentDebit, pErr := CreateDebitOp(
			DebitOp,
			payer,
			p.Amount,
			HNT,
			statusString,
			int64(indexIncrementer*i),
			map[string]interface{}{"debit_category": "payment"},
		)
		if pErr != nil {
			return nil, pErr
		}

		PaymentCredit, pcErr := CreateCreditOp(
			CreditOp,
			p.Payee,
			p.Amount,
			HNT,
			statusString,
			int64((indexIncrementer*i)+1),
			map[string]interface{}{"credit_category": "payment"},
		)
		if pcErr != nil {
			return nil, pcErr
		}

		paymentV2Operations = append(paymentV2Operations, PaymentDebit, PaymentCredit)
	}

	Fee, fErr := CreateFeeOp(payer, fee, statusString, int64(len(paymentV2Operations)), map[string]interface{}{})
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
			RewardOp,
			fmt.Sprint(reward.(map[string]interface{})["account"]),
			utils.MapToInt64(reward.(map[string]interface{})["amount"]),
			HNT,
			SuccessStatus,
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

func CreateHTLCV1(payer string, amount int64, fee *Fee, metadata map[string]interface{}) ([]*types.Operation, *types.Error) {
	var CreateHTLCOps []*types.Operation

	createHTLCOps, chErr := CreateDebitOp(CreateHTLCOp, payer, amount, HNT, SuccessStatus, 0, metadata)
	if chErr != nil {
		return nil, chErr
	}

	CreateHTLCOps = append(CreateHTLCOps, createHTLCOps)

	return CreateHTLCOps, nil
}

func CoinbaseV1(payee string, amount int64) ([]*types.Operation, *types.Error) {
	var CoinbaseOps []*types.Operation

	coinbaseOps, cbErr := CreateCreditOp(CoinbaseOp, payee, amount, HNT, SuccessStatus, 0, map[string]interface{}{"credit_category": "coinbase"})
	if cbErr != nil {
		return nil, cbErr
	}

	CoinbaseOps = append(CoinbaseOps, coinbaseOps)

	return CoinbaseOps, nil
}

func AddGatewayV1(
	payer,
	owner string,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	var AddGatewayOps []*types.Operation

	agwOp, agwErr := CreateGenericOp(AddGatewayOp, SuccessStatus, 0, metadata)
	if agwErr != nil {
		return nil, agwErr
	}

	feePayer := owner
	if (payer != owner) && (payer != "1Wh4bh") && (payer != "") {
		feePayer = payer
	}

	feeOp, feeErr := CreateFeeOp(feePayer, fee, SuccessStatus, 1, map[string]interface{}{})
	if feeErr != nil {
		return nil, feeErr
	}

	AddGatewayOps = append(AddGatewayOps, agwOp, feeOp)

	return AddGatewayOps, nil
}

func AssertLocationV1(
	payer,
	owner string,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	var AssertLocationOps []*types.Operation

	alOp, alErr := CreateGenericOp(AssertLocationOp, SuccessStatus, 0, metadata)
	if alErr != nil {
		return nil, alErr
	}

	feePayer := owner
	if (payer != owner) && (payer != "1Wh4bh") && (payer != "") {
		feePayer = payer
	}

	feeOp, feeErr := CreateFeeOp(feePayer, fee, SuccessStatus, 1, map[string]interface{}{})
	if feeErr != nil {
		return nil, feeErr
	}

	AssertLocationOps = append(AssertLocationOps, alOp, feeOp)
	return AssertLocationOps, nil
}

func AssertLocationV2(
	payer,
	owner string,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	var AssertLocationOps []*types.Operation

	alOp, alErr := CreateGenericOp(AssertLocationOp, SuccessStatus, 0, metadata)
	if alErr != nil {
		return nil, alErr
	}

	feePayer := owner
	if (payer != owner) && (payer != "1Wh4bh") && (payer != "") {
		feePayer = payer
	}

	feeOp, feeErr := CreateFeeOp(feePayer, fee, SuccessStatus, 1, map[string]interface{}{})
	if feeErr != nil {
		return nil, feeErr
	}

	AssertLocationOps = append(AssertLocationOps, alOp, feeOp)
	return AssertLocationOps, nil
}

func SecurityExchangeV1(payer, payee string, fee *Fee, amount int64) ([]*types.Operation, *types.Error) {
	PaymentDebit, pErr := CreateDebitOp(DebitOp, payer, amount, HST, SuccessStatus, 0, map[string]interface{}{"credit_category": "payment"})
	if pErr != nil {
		return nil, pErr
	}

	PaymentCredit, pcErr := CreateCreditOp(CreditOp, payee, amount, HST, SuccessStatus, 1, map[string]interface{}{"debit_category": "payment"})
	if pcErr != nil {
		return nil, pcErr
	}

	Fee, fErr := CreateFeeOp(payer, fee, SuccessStatus, 2, map[string]interface{}{})
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
	buyer,
	seller string,
	amountToSeller int64,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	ops := []*types.Operation{}
	index := int64(0)

	TransferHotspot, tErr := CreateGenericOp(TransferHotspotOp, SuccessStatus, index, metadata)
	if tErr != nil {
		return nil, tErr
	}
	index++
	ops = append(ops, TransferHotspot)

	if amountToSeller > 0 {
		Debit, dErr := CreateDebitOp(DebitOp, buyer, amountToSeller, HNT, SuccessStatus, index, map[string]interface{}{})
		if dErr != nil {
			return nil, dErr
		}
		index++

		Credit, cErr := CreateCreditOp(CreditOp, seller, amountToSeller, HNT, SuccessStatus, index, map[string]interface{}{})
		if cErr != nil {
			return nil, cErr
		}
		index++

		ops = append(ops, Debit, Credit)
	}

	Fee, fErr := CreateFeeOp(buyer, fee, SuccessStatus, index, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	ops = append(ops, Fee)

	return ops, nil
}

func TokenBurnV1(
	payer string,
	amount int64,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	TokenBurn, tErr := CreateDebitOp(TokenBurnOp, payer, amount, HNT, SuccessStatus, 0, map[string]interface{}{})
	if tErr != nil {
		return nil, tErr
	}

	Fee, fErr := CreateFeeOp(payer, fee, SuccessStatus, 1, map[string]interface{}{})
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
	stake int64,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	StakeValidator, sErr := CreateDebitOp(StakeValidatorOp, owner, stake, HNT, SuccessStatus, 0, metadata)
	if sErr != nil {
		return nil, sErr
	}

	Fee, fErr := CreateFeeOp(owner, fee, SuccessStatus, 1, map[string]interface{}{})
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
	stake int64,
	stakeReleaseHeight int64,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	stakeStatus := PendingStatus
	currentHeight, cErr := GetCurrentHeight()
	if cErr != nil {
		return nil, cErr
	}

	if stakeReleaseHeight >= *currentHeight {
		stakeStatus = SuccessStatus
	}

	Unstake, uErr := CreateCreditOp(UnstakeValidatorOp, owner, stake, HNT, stakeStatus, 0, metadata)
	if uErr != nil {
		return nil, uErr
	}

	Fee, fErr := CreateFeeOp(owner, fee, SuccessStatus, 1, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	return []*types.Operation{
		Unstake,
		Fee,
	}, nil
}

func TransferValidatorStakeV1(
	newOwner,
	oldOwner string,
	paymentAmount int64,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	ops := []*types.Operation{}
	index := int64(0)
	TransferValidator, tErr := CreateGenericOp(TransferValidatorStakeOp, SuccessStatus, index, metadata)
	if tErr != nil {
		return nil, tErr
	}
	index++
	ops = append(ops, TransferValidator)

	if paymentAmount > int64(0) {
		Debit, dErr := CreateDebitOp(DebitOp, newOwner, paymentAmount, HNT, SuccessStatus, index, map[string]interface{}{})
		if dErr != nil {
			return nil, dErr
		}
		index++

		Credit, cErr := CreateCreditOp(CreditOp, oldOwner, paymentAmount, HNT, SuccessStatus, index, map[string]interface{}{})
		if cErr != nil {
			return nil, cErr
		}
		index++

		ops = append(ops, Debit, Credit)
	}

	Fee, fErr := CreateFeeOp(oldOwner, fee, SuccessStatus, index, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	ops = append(ops, Fee)

	return ops, nil
}

func FeeOnlyTxn(
	opType,
	payer,
	owner string,
	fee *Fee,
	metadata map[string]interface{},
) ([]*types.Operation, *types.Error) {
	ops := []*types.Operation{}
	MainOp, oErr := CreateGenericOp(opType, SuccessStatus, 0, metadata)
	if oErr != nil {
		return nil, oErr
	}

	feePayer := owner
	if (payer != owner) && (payer != "1Wh4bh") && (payer != "") {
		feePayer = payer
	}

	Fee, fErr := CreateFeeOp(feePayer, fee, SuccessStatus, 1, map[string]interface{}{})
	if fErr != nil {
		return nil, fErr
	}

	ops = append(ops, MainOp, Fee)

	return ops, nil
}

func PassthroughTxn(metadata map[string]interface{}) ([]*types.Operation, *types.Error) {
	MainOp, mErr := CreateGenericOp(PassthroughOp, SuccessStatus, 0, metadata)
	if mErr != nil {
		return nil, mErr
	}
	return []*types.Operation{
		MainOp,
	}, nil
}
