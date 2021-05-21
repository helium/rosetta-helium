package helium

import (
	"errors"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/utils"
)

type Preprocessor struct {
	TransactionType   string                 `json:"transaction_type"`
	RequestedMetadata map[string]interface{} `json:"requested_metadata"`
}

func OpsToTransaction(operations []*types.Operation) (*Preprocessor, *types.Error) {
	var preprocessedTransaction Preprocessor

	switch operations[0].Type {
	case DebitOp:
		// Set txn type
		preprocessedTransaction.TransactionType = PaymentV2Txn

		// Parse payer
		if operations[0].Account == nil {
			return nil, WrapErr(ErrNotFound, errors.New("payment_v2 ops require Accounts"))
		} else {
			preprocessedTransaction.RequestedMetadata["payer"] = operations[0].Account.Address
		}

		// Parse payments
		var payments []map[string]interface{}
		for i, operation := range operations {
			if operation.Account == nil {
				return nil, WrapErr(ErrNotFound, errors.New("payment_v2 ops require Accounts"))
			}

			// Even Ops must be debits, odd Ops must be credits
			if i%2 == 0 {
				// Confirm payer is the same
				if preprocessedTransaction.RequestedMetadata["payer"] != operations[0].Account.Address {
					return nil, WrapErr(ErrNotFound, errors.New("cannot exceed more than one payer for payment_v2 txn"))
				}
			} else {
				if operations[i].Amount == nil {
					return nil, WrapErr(ErrNotFound, errors.New(CreditOp+"s require Amounts"))
				}
				if operations[i].Amount.Value != utils.TrimLeftChar(operations[i-1].Amount.Value) {
					return nil, WrapErr(ErrNotFound, errors.New("debit value does not match credit value"))
				}
				payments = append(payments, map[string]interface{}{"payee": operations[i].Account.Address, "amount": operations[i].Amount.Value})
			}
		}
		preprocessedTransaction.RequestedMetadata["payments"] = payments
		return &preprocessedTransaction, nil
	default:
		return nil, WrapErr(ErrUnclearIntent, errors.New("supported transactions cannot start with "+operations[0].Type+" ops"))
	}
}
