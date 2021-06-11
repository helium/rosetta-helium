package helium

import (
	"errors"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/utils"
)

type Preprocessor struct {
	TransactionType   string                 `json:"transaction_type"`
	HeliumMetadata    map[string]interface{} `json:"helium_metadata"`
	RequestedMetadata map[string]interface{} `json:"requested_metadata"`
}

func OpsToTransaction(operations []*types.Operation) (*Preprocessor, *types.Error) {
	var preprocessedTransaction Preprocessor

	switch operations[0].Type {
	case DebitOp:
		// Set txn type
		preprocessedTransaction.TransactionType = PaymentV2Txn

		if len(operations) <= 1 {
			return nil, WrapErr(ErrNotFound, errors.New("payment_v2 require at least two ops (debit and credit)"))
		}

		// Parse payer
		if operations[0].Account == nil {
			return nil, WrapErr(ErrNotFound, errors.New("payment_v2 ops require Accounts"))
		} else {
			preprocessedTransaction.RequestedMetadata = map[string]interface{}{"get_nonce_for": map[string]interface{}{"address": operations[0].Account.Address}}
		}

		// Create payments helium_metadata object
		paymentMap := []Payment{}

		for i, operation := range operations {
			if operation.Account == nil {
				return nil, WrapErr(ErrNotFound, errors.New("payment_v2 ops require Accounts"))
			}

			// Even Ops must be debits, odd Ops must be credits
			if i%2 == 0 {
				// Confirm payer is the same
				if preprocessedTransaction.RequestedMetadata["get_nonce_for"].(map[string]interface{})["address"] != operations[i].Account.Address {
					return nil, WrapErr(ErrUnclearIntent, errors.New("cannot exceed more than one payer for payment_v2 txn"))
				}
				if operations[i].Amount.Value[0:1] != "-" {
					return nil, WrapErr(ErrUnclearIntent, errors.New(DebitOp+"s cannot be positive"))
				}
			} else {
				if operations[i].Amount == nil {
					return nil, WrapErr(ErrNotFound, errors.New(CreditOp+"s require Amounts"))
				}
				if operations[i].Amount.Value[0:1] == "-" {
					return nil, WrapErr(ErrUnclearIntent, errors.New(CreditOp+"s cannot be negative"))
				}
				if operations[i].Amount.Value != utils.TrimLeftChar(operations[i-1].Amount.Value) {
					return nil, WrapErr(ErrUnclearIntent, errors.New("debit value does not match credit value"))
				}
				if operations[i].Account.Address == preprocessedTransaction.RequestedMetadata["payer"] {
					return nil, WrapErr(ErrUnclearIntent, errors.New("payee and payer cannot be the same address"))
				}

				paymentAmount, err := strconv.ParseInt(operations[i].Amount.Value, 10, 64)
				if err != nil {
					return nil, WrapErr(ErrUnableToParseTxn, err)
				}

				paymentMap = append(paymentMap, Payment{
					Payee:  operations[i].Account.Address,
					Amount: paymentAmount,
				})
			}
		}

		preprocessedTransaction.HeliumMetadata = map[string]interface{}{
			"payer":    operations[0].Account.Address,
			"payments": paymentMap,
		}
		return &preprocessedTransaction, nil
	default:
		return nil, WrapErr(ErrUnclearIntent, errors.New("supported transactions cannot start with "+operations[0].Type))
	}
}
