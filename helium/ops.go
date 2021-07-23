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
		Type:     AssertLocationOp,
		Status:   &SuccessStatus,
		Metadata: metadata,
	}, nil
}

func CreateTransferHotspotOp(buyer, seller, gateway string, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	metadata["buyer"] = buyer
	metadata["seller"] = seller
	metadata["gateway"] = gateway
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type:     TransferHotspotOp,
		Status:   &SuccessStatus,
		Metadata: metadata,
	}, nil
}

func CreateTokenBurnOp(payer, payee string, amount int64, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	metadata["payee"] = payee
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Account: &types.AccountIdentifier{
			Address: payer,
		},
		Amount: &types.Amount{
			Value:    "-" + fmt.Sprint(amount),
			Currency: HNT,
		},
		Type:     TokenBurnOp,
		Status:   &SuccessStatus,
		Metadata: metadata,
	}, nil
}

func CreateStakeValidatorOp(owner, ownerSignature, address string, stake int64, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	metadata["owner"] = owner
	metadata["owner_signature"] = ownerSignature
	metadata["address"] = address
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Account: &types.AccountIdentifier{
			Address: owner,
		},
		Amount: &types.Amount{
			Value:    "-" + fmt.Sprint(stake),
			Currency: HNT,
		},
		Type:     StakeValidatorOp,
		Status:   &SuccessStatus,
		Metadata: metadata,
	}, nil
}

func CreateUnstakeValidatorOp(owner, ownerSignature, address string, stake int64, stakeStatus string, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	metadata["owner"] = owner
	metadata["owner_signature"] = ownerSignature
	metadata["address"] = address
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Account: &types.AccountIdentifier{
			Address: owner,
		},
		Amount: &types.Amount{
			Value:    fmt.Sprint(stake),
			Currency: HNT,
		},
		Type:     UnstakeValidatorOp,
		Status:   &stakeStatus,
		Metadata: metadata,
	}, nil
}

func CreateTransferValidatorOp(newOwner, oldOwner, newAddress, oldAddress, newOwnerSignature, oldOwnerSignature string, stakeAmount, opIndex int64, metadata map[string]interface{}) (*types.Operation, *types.Error) {
	metadata["new_owner"] = newOwner
	metadata["old_owner"] = oldOwner
	metadata["new_address"] = newAddress
	metadata["old_address"] = oldAddress
	metadata["new_owner_signature"] = newOwnerSignature
	metadata["old_owner_signature"] = oldOwnerSignature
	metadata["stake_amount"] = stakeAmount
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		Type:     TransferValidatorOp,
		Status:   &SuccessStatus,
		Metadata: metadata,
	}, nil
}
