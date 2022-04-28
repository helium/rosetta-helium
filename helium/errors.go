// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helium

import (
	"github.com/coinbase/rosetta-sdk-go/types"
)

var (
	// Errors contains all errors that could be returned
	// by this Rosetta implementation.
	Errors = []*types.Error{
		ErrUnimplemented,
		ErrNotFound,
		ErrFailed,
		ErrInvalidParameter,
		ErrInvalidPassword,
		ErrUnableToDerive,
		ErrUnclearIntent,
		ErrUnableToParseIntermediateResult,
		ErrUnableToParseTxn,
		ErrUnableToDecodeAddress,
		ErrSignatureInvalid,
		ErrEnvVariableMissing,
		ErrNodeSync,
		ErrDBCatchup,
	}

	// ErrUnimplemented is returned when an endpoint
	// is called that is not implemented.
	ErrUnimplemented = &types.Error{
		Code:    0,
		Message: "Endpoint not implemented",
	}

	// ErrNotFound is returned when
	// the requested object is not found
	ErrNotFound = &types.Error{
		Code:    1,
		Message: "Object not found",
	}

	// ErrFailed is returned when
	// an endpoint fails
	ErrFailed = &types.Error{
		Code:    2,
		Message: "Endpoint failed",
	}

	// ErrInvalidParameter is returned when
	// an invalid parameter is passed to an endpoint
	ErrInvalidParameter = &types.Error{
		Code:    3,
		Message: "Invalid parameter",
	}

	// ErrInvalidPassword is returned when
	// an invalid password is passed to a wallet
	ErrInvalidPassword = &types.Error{
		Code:    4,
		Message: "InvalidPassword",
	}

	// ErrUnableToDerive is returned when an address
	// cannot be derived from a provided public key.
	ErrUnableToDerive = &types.Error{
		Code:    5,
		Message: "Unable to derive address",
	}

	// ErrUnclearIntent is returned when operations
	// provided in /construction/preprocess or /construction/payloads
	// are not valid.
	ErrUnclearIntent = &types.Error{
		Code:    6,
		Message: "Unable to parse intent",
	}

	// ErrUnableToParseIntermediateResult is returned
	// when a data structure passed between Construction
	// API calls is not valid.
	ErrUnableToParseIntermediateResult = &types.Error{
		Code:    7,
		Message: "Unable to parse intermediate result",
	}

	// ErrUnableToDecodeAddress is returned when an address
	// cannot be parsed during construction.
	ErrUnableToDecodeAddress = &types.Error{
		Code:    8,
		Message: "Unable to decode address",
	}

	// ErrSignatureInvalid is returned when a signature
	// cannot be parsed.
	ErrSignatureInvalid = &types.Error{
		Code:    9,
		Message: "Signature invalid",
	}

	// ErrEnvVariableMissing is returned when an env variable
	// cannot be found.
	ErrEnvVariableMissing = &types.Error{
		Code:    10,
		Message: "Environment variable missing",
	}

	// ErrUnableToParseTxn is returned when a txn
	// cannot be parsed into valid operations
	ErrUnableToParseTxn = &types.Error{
		Code:    11,
		Message: "Unable to parse transaction into valid operations",
	}

	ErrNodeSync = &types.Error{
		Code:      12,
		Message:   "Node is not ready",
		Retriable: true,
	}
	ErrDBCatchup = &types.Error{
		Code:      13,
		Message:   "RocksDB secondary needs to catchup",
		Retriable: true,
	}
)

// WrapErr adds details to the types.Error provided. We use a function
// to do this so that we don't accidentially overrwrite the standard
// errors.
func WrapErr(rErr *types.Error, err error) *types.Error {
	newErr := &types.Error{
		Code:      rErr.Code,
		Message:   rErr.Message,
		Retriable: rErr.Retriable,
	}
	if err != nil {
		newErr.Details = map[string]interface{}{
			"context": err.Error(),
		}
	}

	return newErr
}
