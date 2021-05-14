package helium

import (
	"errors"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
)

func CreateDebitOp(
	payer string,
	amount int64,
	currency *types.Currency,
	opIndex int64,
	metadata map[string]interface{},
) (*types.Operation, *types.Error) {
	if amount < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative payment amount not allowed"))
	}

	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type:   DebitOp,
		Status: &SuccessStatus,
		Account: &types.AccountIdentifier{
			Address: payer,
		},
		Amount: &types.Amount{
			Value:    "-" + fmt.Sprint(amount),
			Currency: currency,
		},
		Metadata: metadata,
	}, nil
}

func CreateCreditOp(
	payee string,
	amount int64,
	currency *types.Currency,
	opIndex int64,
	metadata map[string]interface{},
) (*types.Operation, *types.Error) {
	if amount < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative payment amount not allowed"))
	}
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type:   CreditOp,
		Status: &SuccessStatus,
		Account: &types.AccountIdentifier{
			Address: payee,
		},
		Amount: &types.Amount{
			Value:    fmt.Sprint(amount),
			Currency: currency,
		},
		Metadata: metadata,
	}, nil
}

func CreateFeeOp(payer string, fee int64, feeType string, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	var FeeOp *types.Operation
	var FeeCurrency *types.Currency

	switch feeType {
	case "HNT":
		FeeCurrency = HNT
		metadata["implicit_burn"] = true
	case "DC":
		// No reconciliation for DC fees, this is only an FYI
		metadata["implicit_burn"] = false
		metadata["dc_fee"] = fee
		return &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: opIndex,
			},
			Type:   DebitOp,
			Status: &SuccessStatus,
			Account: &types.AccountIdentifier{
				Address: payer,
			},
			Amount: &types.Amount{
				Value:    "0",
				Currency: HNT,
			},
			Metadata: metadata,
		}, nil
	default:
		return nil, WrapErr(ErrNotFound, errors.New("incorrect or missing feeType"))
	}

	if fee < 0 {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New("negative fee amount not allowed"))
	}

	FeeOp = &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type:   DebitOp,
		Status: &SuccessStatus,
		Account: &types.AccountIdentifier{
			Address: payer,
		},
		Amount: &types.Amount{
			Value:    "-" + fmt.Sprint(fee),
			Currency: FeeCurrency,
		},
		Metadata: metadata,
	}
	return FeeOp, nil
}

func CreateAddGatewayOp(gateway, owner string, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	metadata["gateway"] = gateway
	metadata["owner"] = owner
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type:     AddGatewayOp,
		Status:   &SuccessStatus,
		Metadata: metadata,
	}, nil
}

func CreateAssertLocationOp(gateway, owner, location string, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	metadata["gateway"] = gateway
	metadata["owner"] = owner
	metadata["location"] = location
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type:     AddGatewayOp,
		Status:   &SuccessStatus,
		Metadata: metadata,
	}, nil
}
