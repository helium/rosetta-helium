package helium

import (
	"errors"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
)

func CreatePaymentDebitOp(payer *string, amount *int64, opIndex int64) (*types.Operation, *types.Error) {
	if *amount < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative payment amount not allowed"))
	} else {
		return &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: opIndex,
			},
			Type:   PaymentDebitOp,
			Status: &SuccessStatus,
			Account: &types.AccountIdentifier{
				Address: *payer,
			},
			Amount: &types.Amount{
				Value:    "-" + fmt.Sprint(*amount),
				Currency: HNT,
			},
		}, nil
	}
}

func CreatePaymentCreditOp(payee *string, amount *int64, opIndex int64) (*types.Operation, *types.Error) {
	if *amount < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative payment amount not allowed"))
	} else {
		return &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: opIndex,
			},
			Type:   PaymentCreditOp,
			Status: &SuccessStatus,
			Account: &types.AccountIdentifier{
				Address: *payee,
			},
			Amount: &types.Amount{
				Value:    fmt.Sprint(*amount),
				Currency: HNT,
			},
		}, nil
	}
}

func CreateFeeOp(payer *string, fee *int64, feeType *string, opIndex int64) (*types.Operation, *types.Error) {

	var FeeOp *types.Operation
	var FeeCurrency *types.Currency
	var metadata map[string]interface{}

	switch *feeType {
	case "HNT":
		FeeCurrency = HNT
		metadata = map[string]interface{}{
			"implicit_burn": true,
		}
	case "DC":
		FeeCurrency = DC
		metadata = map[string]interface{}{
			"implicit_burn": false,
		}
	default:
		return nil, WrapErr(ErrNotFound, errors.New("incorrect or missing feeType"))
	}

	if *fee < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative fee amount not allowed"))
	} else {
		FeeOp = &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: opIndex,
			},
			Type:   PaymentCreditOp,
			Status: &SuccessStatus,
			Account: &types.AccountIdentifier{
				Address: *payer,
			},
			Amount: &types.Amount{
				Value:    "-" + fmt.Sprint(*fee),
				Currency: FeeCurrency,
			},
			Metadata: metadata,
		}
		return FeeOp, nil
	}
}

func CreateRewardOp(payee *string, amount *int64, opIndex int64) (*types.Operation, *types.Error) {
	if *amount < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative reward amount not allowed"))
	} else {
		return &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: opIndex,
			},
			Type:   RewardOp,
			Status: &SuccessStatus,
			Account: &types.AccountIdentifier{
				Address: *payee,
			},
			Amount: &types.Amount{
				Value:    fmt.Sprint(*amount),
				Currency: HNT,
			},
		}, nil
	}
}
