package helium

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/helium/rosetta-helium/utils"
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

func TransactionToOps(txn map[string]interface{}, status string, block *types.BlockIdentifier) ([]*types.Operation, *types.Error) {
	hash := fmt.Sprint(txn["hash"])
	switch txn["type"] {

	case AddGatewayV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"])+utils.JsonNumberToInt64(txn["staking_fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return AddGatewayV1(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["owner"]),
			feeDetails,
			txn,
		)

	case AssertLocationV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"])+utils.JsonNumberToInt64(txn["staking_fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return AssertLocationV1(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["owner"]),
			feeDetails,
			txn,
		)

	case AssertLocationV2Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"])+utils.JsonNumberToInt64(txn["staking_fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return AssertLocationV2(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["owner"]),
			feeDetails,
			txn,
		)

	case PaymentV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return PaymentV1(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["payee"]),
			utils.JsonNumberToInt64(txn["amount"]),
			feeDetails)

	case PaymentV2Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		var payments []*Payment
		for _, p := range txn["payments"].([]interface{}) {
			payments = append(payments, &Payment{
				Payee:  fmt.Sprint(p.(map[string]interface{})["payee"]),
				Amount: utils.JsonNumberToInt64(p.(map[string]interface{})["amount"]),
			})
		}
		return PaymentV2(
			fmt.Sprint(txn["payer"]),
			payments,
			feeDetails,
			status,
		)

	case RewardsV1Txn, RewardsV2Txn:
		// rewards_v1 and rewards_v2 have the same structure
		return RewardsV1(
			txn["rewards"].([]interface{}),
		)

	case CoinbaseV1Txn:
		return CoinbaseV1(
			fmt.Sprint(txn["payee"]),
			utils.JsonNumberToInt64(txn["amount"]),
		)

	case SecurityExchangeV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return SecurityExchangeV1(
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["payee"]),
			feeDetails,
			utils.JsonNumberToInt64(txn["amount"]),
		)

	case TokenBurnV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return TokenBurnV1(
			fmt.Sprint(txn["payer"]),
			utils.JsonNumberToInt64(txn["amount"]),
			feeDetails,
			txn,
		)

	case TransferHotspotV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return TransferHotspotV1(
			fmt.Sprint(txn["buyer"]),
			fmt.Sprint(txn["seller"]),
			utils.JsonNumberToInt64(txn["amount_to_seller"]),
			feeDetails,
			txn,
		)

	case StakeValidatorV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return StakeValidatorV1(
			fmt.Sprint(txn["owner"]),
			utils.JsonNumberToInt64(txn["stake"]),
			feeDetails,
			txn,
		)

	case UnstakeValidatorV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return UnstakeValidatorV1(
			fmt.Sprint(txn["owner"]),
			utils.JsonNumberToInt64(txn["stake_amount"]),
			utils.JsonNumberToInt64(txn["stake_release_height"]),
			feeDetails,
			txn,
		)

	case TransferValidatorStakeV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return TransferValidatorStakeV1(
			fmt.Sprint(txn["new_owner"]),
			fmt.Sprint(txn["old_owner"]),
			utils.JsonNumberToInt64(txn["payment_amount"]),
			feeDetails,
			txn,
		)

	case OUIV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return FeeOnlyTxn(
			OUIOp,
			fmt.Sprint(txn["payer"]),
			fmt.Sprint(txn["owner"]),
			feeDetails,
			txn,
		)

	case UpdateGatewayOUIV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}

		owner, oErr := GetGatewayOwner(fmt.Sprint(txn["address"]), block.Index)
		if oErr != nil {
			return nil, oErr
		}

		return FeeOnlyTxn(
			OUIOp,
			*owner,
			*owner,
			feeDetails,
			txn,
		)

	case RoutingV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return FeeOnlyTxn(
			RoutingOp,
			"",
			fmt.Sprint(txn["owner"]),
			feeDetails,
			txn,
		)

	case StateChannelOpenV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return FeeOnlyTxn(
			StateChannelOpenOp,
			"",
			fmt.Sprint(txn["owner"]),
			feeDetails,
			txn,
		)

	case CreateHTLCV1Txn:
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		return CreateHTLCV1(
			fmt.Sprint(txn["payer"]),
			utils.JsonNumberToInt64(txn["amount"]),
			feeDetails,
			map[string]interface{}{
				"payee":    fmt.Sprint(txn["payee"]),
				"address":  fmt.Sprint(txn["address"]),
				"hashlock": fmt.Sprint(txn["hashlock"]),
				"timelock": fmt.Sprint(txn["timelock"]),
			},
		)

	case RedeemHTLCV1Txn:
		address := fmt.Sprint(txn["address"])
		feeDetails, feeErr := GetFee(&hash, utils.JsonNumberToInt64(txn["fee"]))
		if feeErr != nil {
			return nil, feeErr
		}
		htlcDetails, htlcErr := GetHTLCReceipt(address)
		if htlcErr != nil {
			return nil, htlcErr
		}

		return RedeemHTLCV1(
			fmt.Sprint(txn["payee"]),
			htlcDetails.Balance,
			feeDetails,
			map[string]interface{}{
				"address":     address,
				"payer":       htlcDetails.Payer,
				"hashlock":    htlcDetails.Hashlock,
				"timelock":    htlcDetails.Timelock,
				"redeemed_at": htlcDetails.RedeemedAt,
				"preimage":    fmt.Sprint(txn["preimage"]),
			},
		)

	default:
		for _, types := range TransactionTypes {
			if fmt.Sprint(txn["type"]) == types {
				return PassthroughTxn(txn)
			}
		}
		return nil, WrapErr(ErrNotFound, errors.New("txn type not found: "+fmt.Sprint(txn["type"])))
	}
}
