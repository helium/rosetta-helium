package helium

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/syuan100/rosetta-helium/utils"
)

func OpsToTransaction(operations []*types.Operation) (*MetadataOptions, *types.Error) {
	var preprocessedTransaction MetadataOptions

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

func TransactionToOps(txn map[string]interface{}) ([]*types.Operation, *types.Error) {
	hash := fmt.Sprint(txn["hash"])
	switch txn["type"] {

	case AddGatewayV1Txn:
		feeDetails := GetFee(&hash, utils.MapToInt64(txn["fee"])+utils.MapToInt64(txn["staking_fee"]), fmt.Sprint(txn["payer"]))
		return AddGatewayV1(
			fmt.Sprint(txn["payer"]),
			feeDetails.Amount,
			feeDetails.Currency.Symbol,
			fmt.Sprint(txn["gateway"]),
			fmt.Sprint(txn["owner"]),
			utils.MapToInt64(txn["fee"]),
			utils.MapToInt64(txn["staking_fee"]))

	case AssertLocationV1Txn:
		feeDetails := GetFee(&hash, utils.MapToInt64(txn["fee"])+utils.MapToInt64(txn["staking_fee"]), fmt.Sprint(txn["payer"]))
		return AssertLocationV1(
			utils.MapToInt64(txn["fee"]),
			fmt.Sprint(txn["gateway"]),
			fmt.Sprint(txn["location"]),
			fmt.Sprint(txn["owner"]),
			fmt.Sprint(txn["payer"]),
			utils.MapToInt64(txn["staking_fee"]),
			feeDetails.Amount,
			feeDetails.Currency.Symbol)

	case AssertLocationV2Txn:
		feeDetails := GetFee(&hash, utils.MapToInt64(txn["fee"])+utils.MapToInt64(txn["staking_fee"]), fmt.Sprint(txn["payer"]))
		return AssertLocationV2(
			utils.MapToInt64(txn["elevation"]),
			utils.MapToInt64(txn["fee"]),
			utils.MapToInt64(txn["gain"]),
			fmt.Sprint(txn["gateway"]),
			fmt.Sprint(txn["location"]),
			fmt.Sprint(txn["owner"]),
			fmt.Sprint(txn["payer"]),
			utils.MapToInt64(txn["staking_fee"]),
			feeDetails.Amount,
			feeDetails.Currency.Symbol)

	case PaymentV1Txn:
		feeDetails := GetFee(&hash, utils.MapToInt64(txn["fee"]), fmt.Sprint(txn["payer"]))
		return PaymentV1(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["payee"]),
			int64(txn["amount"].(float64)),
			feeDetails.Amount,
			feeDetails.Currency.Symbol)

	case PaymentV2Txn:
		feeDetails := GetFee(&hash, utils.MapToInt64(txn["fee"]), fmt.Sprint(txn["payer"]))
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

	case RewardsV1Txn, RewardsV2Txn:
		// rewards_v1 and rewards_v2 have the same structure
		return RewardsV1(
			txn["rewards"].([]interface{}),
		)

	case SecurityCoinbaseV1Txn:
		return SecurityCoinbaseV1(
			fmt.Sprint(txn["payee"]),
			int64(txn["amount"].(float64)),
		)

	case SecurityExchangeV1Txn:
		feeDetails := GetFee(&hash, utils.MapToInt64(txn["fee"]), fmt.Sprint(txn["payer"]))
		return SecurityExchangeV1(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["payee"]),
			feeDetails.Amount,
			feeDetails.Currency.Symbol,
			int64(txn["amount"].(float64)),
		)

	// case TokenBurnV1Txn:
	// 	feeDetails := GetFee(&hash, utils.MapToInt64(txn["fee"])+utils.MapToInt64(txn["staking_fee"]), fmt.Sprint(txn["payer"]))

	default:
		return nil, WrapErr(ErrNotFound, errors.New("txn type not found"))
	}
}
