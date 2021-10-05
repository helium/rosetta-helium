package helium

import (
	"errors"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
)

func CreateGenericOp(opType string, status string, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type:     opType,
		Status:   &status,
		Metadata: metadata,
	}, nil
}

func CreateDebitOp(
	opType,
	payer string,
	amount int64,
	currency *types.Currency,
	status string,
	opIndex int64,
	metadata map[string]interface{},
) (*types.Operation, *types.Error) {
	if amount < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative payment amount not allowed"))
	}

	debitOp := &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type: opType,
		Account: &types.AccountIdentifier{
			Address: payer,
		},
		Amount: &types.Amount{
			Value:    "-" + fmt.Sprint(amount),
			Currency: currency,
		},
		Metadata: metadata,
	}

	if status != "" {
		debitOp.Status = &status
	}

	return debitOp, nil
}

func CreateCreditOp(
	opType,
	payee string,
	amount int64,
	currency *types.Currency,
	status string,
	opIndex int64,
	metadata map[string]interface{},
) (*types.Operation, *types.Error) {
	if amount < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative payment amount not allowed"))
	}

	creditOp := &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type: opType,
		Account: &types.AccountIdentifier{
			Address: payee,
		},
		Amount: &types.Amount{
			Value:    fmt.Sprint(amount),
			Currency: currency,
		},
		Metadata: metadata,
	}

	if status != "" {
		creditOp.Status = &status
	}

	return creditOp, nil
}

func CreateFeeOp(payer string, fee *Fee, status string, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	FeeOpObject := &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type: FeeOp,
		Account: &types.AccountIdentifier{
			Address: payer,
		},
		Metadata: metadata,
	}

	if fee.Amount < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative fee amount not allowed"))
	}

	switch fee.Currency.Symbol {
	case "HNT":
		FeeOpObject.Amount = &types.Amount{
			Value:    "-" + fmt.Sprint(fee.Amount),
			Currency: HNT,
		}
		metadata["debit_category"] = "fee"
		metadata["implicit_burn"] = true
		metadata["dc_fee"] = fee.DCFeeAmount
		FeeOpObject.Metadata = metadata
	case "DC":
		metadata["debit_category"] = "fee"
		metadata["implicit_burn"] = false
		metadata["dc_fee"] = fee.DCFeeAmount
		FeeOpObject.Metadata = metadata
	default:
		return nil, WrapErr(ErrNotFound, errors.New("incorrect or missing feeType"))
	}

	if status != "" {
		FeeOpObject.Status = &status
	}

	return FeeOpObject, nil
}
