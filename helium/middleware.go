package helium

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/helium/rosetta-helium/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/ybbus/jsonrpc"
)

var (
	NodeClient = jsonrpc.NewClient("http://localhost:4467")
)

type MetadataOptions struct {
	RequestedMetadata map[string]interface{} `json:"requested_metadata"`
	HeliumMetadata    map[string]interface{} `json:"helium_metadata"`
	TransactionType   string                 `json:"transaction_type"`
}

type combination struct {
	UnsignedTransaction string             `json:"unsigned_transaction"`
	Signatures          []*types.Signature `json:"signatures"`
}

type hashRequest struct {
	Transaction string `json:"txn"`
}

type deriveRequest struct {
	CurveType string `json:"curve_type"`
	PublicKey string `json:"public_key"`
}

func GetCurrentHeight() (*int64, *types.Error) {
	var result int64

	if err := NodeClient.CallFor(&result, "block_height", nil); err != nil {
		return nil, WrapErr(ErrUnclearIntent, errors.New("error getting block_height"))
	}

	return &result, nil
}

func GetBlock(blockIdentifier *types.PartialBlockIdentifier) (*types.Block, *types.Error) {
	type request struct {
		Height int64  `json:"height,omitempty"`
		Hash   string `json:"hash,omitempty"`
	}

	var result Block
	var req request

	if blockIdentifier.Index != nil && blockIdentifier.Hash != nil {
		req = request{
			Height: *blockIdentifier.Index,
		}
	} else if blockIdentifier.Index == nil && blockIdentifier.Hash != nil {
		req = request{
			Hash: *blockIdentifier.Hash,
		}
	} else if blockIdentifier.Index != nil && blockIdentifier.Hash == nil {
		req = request{
			Height: *blockIdentifier.Index,
		}
	}

	if err := NodeClient.CallFor(&result, "block_get", req); err != nil {
		return nil, WrapErr(ErrNotFound, err)
	}

	var processedTxs []*types.Transaction
	for _, tx := range result.Transactions {
		ptx, txErr := GetTransaction(fmt.Sprint(tx["hash"]))
		if txErr != nil {
			return nil, txErr
		}

		processedTxs = append(processedTxs, ptx)
	}

	currentBlock := &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: result.Height,
			Hash:  result.Hash,
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: result.Height,
			Hash:  result.Hash,
		},
		Timestamp:    result.Time * 1000,
		Transactions: processedTxs,
		Metadata:     nil,
	}

	return currentBlock, nil
}

func GetTransaction(txHash string) (*types.Transaction, *types.Error) {
	type request struct {
		Hash string `json:"hash"`
	}

	var result map[string]interface{}

	req := request{Hash: txHash}
	if err := NodeClient.CallFor(&result, "transaction_get", req); err != nil {
		return nil, WrapErr(
			ErrNotFound,
			err,
		)
	}

	operations, oErr := TransactionToOps(result, SuccessStatus)
	if oErr != nil {
		return nil, oErr
	}

	transaction := &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: fmt.Sprint(result["hash"]),
		},
		Operations:          operations,
		RelatedTransactions: nil,
		Metadata:            nil,
	}

	return transaction, nil
}

func GetAddress(curveType types.CurveType, publicKey []byte) (*string, *types.Error) {
	jsonObject, jErr := json.Marshal(deriveRequest{
		CurveType: string(curveType),
		PublicKey: hex.EncodeToString(publicKey),
	})
	if jErr != nil {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New(`unable to decode transaction object into json`))
	}

	var payload map[string]interface{}
	resp, ctErr := http.Post("http://localhost:3000/derive", "application/json", bytes.NewBuffer(jsonObject))
	if ctErr != nil {
		return nil, WrapErr(ErrUnclearIntent, ctErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&payload)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}

	if payload["address"] == nil {
		return nil, WrapErr(ErrUnableToDerive, errors.New("constructor unable to process"))
	}
	response := payload["address"].(string)
	return &response, nil
}

func GetBalance(address string) ([]*types.Amount, *types.Error) {
	var balances []*types.Amount

	type request struct {
		Address string `json:"address"`
	}

	var result map[string]interface{}

	req := request{Address: address}
	if err := NodeClient.CallFor(&result, "account_get", req); err != nil {
		return nil, WrapErr(
			ErrNotFound,
			err,
		)
	}

	amountHNT := &types.Amount{
		Value:    fmt.Sprint(utils.MapToInt64(result["balance"])),
		Currency: HNT,
	}

	amountHST := &types.Amount{
		Value:    fmt.Sprint(utils.MapToInt64(result["sec_balance"])),
		Currency: HST,
	}

	balances = append(balances, amountHNT, amountHST)

	return balances, nil
}

func GetHash(signedTransaction string) (*string, *types.Error) {
	jsonObject, jErr := json.Marshal(hashRequest{
		Transaction: signedTransaction,
	})
	if jErr != nil {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New(`unable to decode transaction object into json`))
	}

	var payload map[string]interface{}
	resp, ctErr := http.Post("http://localhost:3000/hash", "application/json", bytes.NewBuffer(jsonObject))
	if ctErr != nil {
		return nil, WrapErr(ErrUnclearIntent, ctErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&payload)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}

	hash := payload["hash"].(string)
	return &hash, nil
}

func GetNonce(address string) (*int64, *types.Error) {
	var nonce int64

	type request struct {
		Address string `json:"address"`
	}

	var result map[string]interface{}

	req := request{Address: address}
	if err := NodeClient.CallFor(&result, "account_get", req); err != nil {
		return nil, WrapErr(
			ErrNotFound,
			err,
		)
	}

	nonce = utils.MapToInt64(result["nonce"])

	return &nonce, nil
}

func GetOraclePrice(height int64) (*int64, *types.Error) {
	type request struct {
		Height int64 `json:"height"`
	}

	var result map[string]interface{}

	req := request{Height: height}
	if err := NodeClient.CallFor(&result, "oracle_price_get", req); err != nil {
		return nil, WrapErr(
			ErrNotFound,
			err,
		)
	}

	price := utils.MapToInt64(result["price"])

	return &price, nil
}

func GetFee(hash *string, DCFeeAmount int64) *Fee {
	if hash == nil {
		return &Fee{
			Amount:      DCFeeAmount,
			Currency:    DC,
			Estimate:    true,
			DCFeeAmount: DCFeeAmount,
		}
	}

	type request struct {
		Hash string `json:"hash"`
	}

	var result map[string]interface{}

	req := request{Hash: *hash}
	if err := NodeClient.CallFor(&result, "implicit_burn_get", req); err != nil {
		return &Fee{
			Amount:      DCFeeAmount,
			Currency:    DC,
			Estimate:    false,
			DCFeeAmount: DCFeeAmount,
		}
	}

	feeResult := &Fee{
		Amount:      utils.MapToInt64(result["fee"]),
		Currency:    HNT,
		Estimate:    false,
		DCFeeAmount: DCFeeAmount,
	}

	return feeResult
}

func GetMetadata(request *types.ConstructionMetadataRequest) (*types.ConstructionMetadataResponse, *types.Error) {
	metadataResponse := types.ConstructionMetadataResponse{
		Metadata: map[string]interface{}{},
	}

	// Get chain_vars (default metadata)
	var chainVars map[string]interface{}
	resp, vErr := http.Get("http://localhost:3000/chain-vars")
	if vErr != nil {
		return nil, WrapErr(ErrUnclearIntent, vErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&chainVars)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}
	metadataResponse.Metadata["chain_vars"] = chainVars

	// Get raw options
	jsonString, _ := json.Marshal(request.Options)
	options := MetadataOptions{}
	err := json.Unmarshal(jsonString, &options)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("syntax error at byte offset %d", e.Offset)
		}
		return nil, WrapErr(ErrUnclearIntent, err)
	}
	metadataResponse.Metadata["options"] = options

	for k, v := range options.RequestedMetadata {
		switch k {
		case "get_nonce_for":
			switch t := v.(type) {
			case map[string]interface{}:
				if v.(map[string]interface{})["address"] == nil {
					return nil, WrapErr(ErrUnclearIntent, errors.New("get_nonce_for requires `address` to be present in JSON object"))
				}

				nonce, nErr := GetNonce(fmt.Sprint(v.(map[string]interface{})["address"]))
				if nErr != nil {
					return nil, nErr
				}

				metadataResponse.Metadata["get_nonce_for"] = map[string]interface{}{
					"nonce": nonce,
				}
			default:
				return nil, WrapErr(ErrUnclearIntent, errors.New("unexpected object "+fmt.Sprint(t)+" in get_nonce_for"))
			}
		default:
			return nil, WrapErr(ErrUnclearIntent, errors.New("metadata request `"+fmt.Sprint(k)+"` not recognized"))
		}
	}

	return &metadataResponse, nil
}

func GetPeers() ([]*types.Peer, *types.Error) {
	var result []map[string]interface{}
	if err := NodeClient.CallFor(&result, "peer_book_self"); err != nil {
		return nil, WrapErr(
			ErrNotFound,
			err,
		)
	}

	peers := result[0]["sessions"].([]interface{})
	var convertedPeers []*types.Peer

	for _, peer := range peers {
		var cPeer *Peer
		cErr := mapstructure.Decode(peer, &cPeer)
		if cErr != nil {
			return nil, WrapErr(ErrUnclearIntent, cErr)
		}
		convertedPeers = append(convertedPeers, &types.Peer{
			PeerID: cPeer.P2P,
		})
	}

	return convertedPeers, nil
}

func GetSyncStatus() (*types.SyncStatus, *types.Error) {
	currentHeight, chErr := GetCurrentHeight()
	if chErr != nil {
		return nil, chErr
	}

	targetHeight, thErr := GetTargetHeight()
	if thErr != nil {
		return nil, thErr
	}

	// Consider blockchain syced if the local head is < 100 blocks from remote head
	synced := (int(*currentHeight) > int(*targetHeight)-100)

	syncStatus := &types.SyncStatus{
		CurrentIndex: currentHeight,
		TargetIndex:  targetHeight,
		Synced:       &synced,
	}
	return syncStatus, nil
}

func GetTargetHeight() (*int64, *types.Error) {
	var response map[string]interface{}
	resp, rErr := http.Get("http://localhost:3000/current-height")
	if rErr != nil {
		return nil, WrapErr(ErrUnclearIntent, rErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&response)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}

	currentHeight := utils.MapToInt64(response["current_height"])
	return &currentHeight, nil
}

func CombineTransaction(unsignedTxn string, signatures []*types.Signature) (*types.ConstructionCombineResponse, *types.Error) {
	jsonObject, jErr := json.Marshal(combination{
		UnsignedTransaction: unsignedTxn,
		Signatures:          signatures,
	})
	if jErr != nil {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New(`unable to decode combination object into json`))
	}

	var payload map[string]interface{}
	resp, ctErr := http.Post("http://localhost:3000/combine-tx", "application/json", bytes.NewBuffer(jsonObject))
	if ctErr != nil {
		return nil, WrapErr(ErrUnclearIntent, ctErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&payload)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}

	signedTransaction := payload["signed_transaction"].(string)

	return &types.ConstructionCombineResponse{
		SignedTransaction: signedTransaction,
	}, nil
}

func ParseTransaction(rawTxn string, signed bool) ([]*types.Operation, *types.AccountIdentifier, *types.Error) {
	var jsonData = []byte(fmt.Sprintf(`{ "raw_transaction": "%s", "signed": %t }`, rawTxn, signed))

	var payload map[string]interface{}
	resp, ctErr := http.Post("http://localhost:3000/parse-tx", "application/json", bytes.NewBuffer(jsonData))
	if ctErr != nil {
		return nil, nil, WrapErr(ErrUnclearIntent, ctErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&payload)
	if dErr != nil {
		return nil, nil, WrapErr(ErrUnclearIntent, dErr)
	}

	operations, oErr := TransactionToOps(payload["payload"].(map[string]interface{}), "")
	if oErr != nil {
		return nil, nil, oErr
	}

	if signed {
		return operations, &types.AccountIdentifier{
			Address: fmt.Sprint(payload["signer"]),
		}, nil
	}

	return operations, nil, nil
}

func PayloadGenerator(operations []*types.Operation, metadata map[string]interface{}) (*types.ConstructionPayloadsResponse, *types.Error) {

	transactionPreprocessor, err := OpsToTransaction(operations)
	if err != nil {
		return nil, err
	}

	var operationMetadata map[string]interface{}
	marshalledPreprocessor, _ := json.Marshal(transactionPreprocessor)
	json.Unmarshal(marshalledPreprocessor, &operationMetadata)

	if !reflect.DeepEqual(metadata["options"], operationMetadata) {
		return nil, WrapErr(ErrUnclearIntent, errors.New(`payload operations options result do not match provided metadata options (metadata["options"])`))
	}

	jsonValue, jErr := json.Marshal(metadata)
	if jErr != nil {
		fmt.Print(jErr)
	}

	var payload map[string]interface{}
	resp, ctErr := http.Post("http://localhost:3000/create-tx", "application/json", bytes.NewBuffer(jsonValue))
	if ctErr != nil {
		return nil, WrapErr(ErrUnclearIntent, ctErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&payload)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}

	decodedByteArray, hErr := hex.DecodeString(payload["payload"].(string))
	if hErr != nil {
		return nil, WrapErr(ErrUnableToParseTxn, hErr)
	}

	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: payload["unsigned_txn"].(string),
		Payloads: []*types.SigningPayload{
			{
				AccountIdentifier: &types.AccountIdentifier{
					Address: fmt.Sprint(metadata["options"].(map[string]interface{})["helium_metadata"].(map[string]interface{})["payer"]),
				},
				Bytes:         decodedByteArray,
				SignatureType: types.Ed25519,
			},
		},
	}, nil
}

func SubmitTransaction(signedTransaction string) (*string, *types.Error) {
	type request struct {
		SignedTransaction string `json:"signed_transaction,omitempty"`
	}

	jsonObject, jErr := json.Marshal(request{
		SignedTransaction: signedTransaction,
	})
	if jErr != nil {
		return nil, WrapErr(ErrUnableToParseTxn, errors.New(`unable to decode combination object into json`))
	}

	var payload map[string]interface{}
	resp, rErr := http.Post("http://localhost:3000/submit-tx", "application/json", bytes.NewBuffer(jsonObject))
	if rErr != nil {
		return nil, WrapErr(ErrUnclearIntent, rErr)
	}
	defer resp.Body.Close()
	dErr := json.NewDecoder(resp.Body).Decode(&payload)
	if dErr != nil {
		return nil, WrapErr(ErrUnclearIntent, dErr)
	}

	hash := payload["hash"].(string)

	return &hash, nil
}
