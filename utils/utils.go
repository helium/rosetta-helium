package utils

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/coinbase/rosetta-sdk-go/types"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/ybbus/jsonrpc"
	"go.uber.org/zap"
)

var (
	NodeClient = jsonrpc.NewClient("http://localhost:4467")
)

var DB *badger.DB

type GhostTxnKey struct {
	Network     *types.NetworkIdentifier     `json:"network"`
	Block       *types.BlockIdentifier       `json:"block"`
	Transaction *types.TransactionIdentifier `json:"transaction"`
}

type GhostTxnMetadata struct {
	Operations []*types.Operation     `json:"operations"`
	Metadata   map[string]interface{} `json:"metadata"`
}

func GetKeyBytes(key *GhostTxnKey, seekKeyOnly bool) ([]byte, error) {
	networkBytes, nerr := json.Marshal(key.Network)
	if nerr != nil {
		return nil, nerr
	}

	heightBytes, herr := json.Marshal(key.Block)
	if herr != nil {
		return nil, herr
	}

	transactionBytes, terr := json.Marshal(key.Transaction)
	if terr != nil {
		return nil, terr
	}

	seekKeyBytes := append(networkBytes, heightBytes...)

	if seekKeyOnly {
		return seekKeyBytes, nil
	}

	keyBytes := append(seekKeyBytes, transactionBytes...)

	return keyBytes, nil
}

func CreateGhostTxn(key *GhostTxnKey, metadata *GhostTxnMetadata) error {
	keyBytes, kerr := GetKeyBytes(key, false)
	if kerr != nil {
		return kerr
	}

	verr := DB.View(func(txn *badger.Txn) error {
		_, err := txn.Get(keyBytes)
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		} else {
			zap.S().Info("cannot create new ghost txn, badger db entry already exists")
			return badger.ErrBannedKey
		}
	})
	if verr != nil {
		return verr
	}

	metadataBytes, merr := json.Marshal(metadata)
	if merr != nil {
		return merr
	}

	uerr := DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(keyBytes, metadataBytes)
		return err
	})
	if uerr != nil {
		return uerr
	}

	return nil
}

func GetGhostTxn(key *GhostTxnKey) (*GhostTxnMetadata, error) {
	keyBytes, kerr := GetKeyBytes(key, false)
	if kerr != nil {
		return nil, kerr
	}

	var txnMetadata GhostTxnMetadata

	verr := DB.View(func(txn *badger.Txn) error {
		item, gerr := txn.Get(keyBytes)
		if gerr != nil {
			return gerr
		}

		ierr := item.Value(func(val []byte) error {
			if terr := json.Unmarshal(val, &txnMetadata); terr != nil {
				return terr
			}
			return nil
		})
		if ierr != nil {
			return ierr
		}
		return nil
	})
	if verr != nil {
		return nil, verr
	}

	return &txnMetadata, nil
}

func SeekGhostTxnsInBlock(network *types.NetworkIdentifier, block_height int64) ([]*types.Transaction, error) {
	var transactions []*types.Transaction

	seekKeyBytes, kerr := GetKeyBytes(
		&GhostTxnKey{
			Network: network,
			Block: &types.BlockIdentifier{
				Index: block_height,
			},
			Transaction: nil,
		},
		true,
	)
	if kerr != nil {
		return nil, kerr
	}

	DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(seekKeyBytes); it.ValidForPrefix(seekKeyBytes); it.Next() {
			item := it.Item()
			var txnMetadata GhostTxnMetadata
			err := item.Value(func(v []byte) error {
				if terr := json.Unmarshal(v, &txnMetadata); terr != nil {
					return terr
				}

				transactions = append(transactions, &types.Transaction{
					TransactionIdentifier: &types.TransactionIdentifier{
						Hash: string(item.Key()),
					},
					Operations: txnMetadata.Operations,
				})
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return transactions, nil
}

func JsonNumberToInt64(m interface{}) int64 {
	convertedInt, _ := m.(json.Number).Int64()
	return convertedInt
}

func DecodeCallAsNumber(call *jsonrpc.RPCResponse, err error) (map[string]interface{}, error) {
	if err != nil {
		return nil, errors.New("unable to decode json-rpc response")
	}

	stringResult, serr := json.Marshal(call.Result)
	if serr != nil {
		return nil, errors.New("unable to marshal json-rpc response")
	}

	d := json.NewDecoder(strings.NewReader(string(stringResult)))
	d.UseNumber()
	var result map[string]interface{}
	if derr := d.Decode(&result); derr != nil {
		return nil, errors.New("unable to decode json-rpc response with json.Number")
	}

	return result, nil
}

func TrimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
