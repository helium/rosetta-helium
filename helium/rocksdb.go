package helium

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	b64 "encoding/base64"

	"github.com/btcsuite/btcutil/base58"
	"github.com/helium/rosetta-helium/utils"
	rocksdb "github.com/linxGnu/grocksdb"
	"github.com/okeuday/erlang_go/v2/erlang"
	"go.uber.org/zap"
)

type Entry struct {
	Name   string
	Nonce  int64
	Amount int64
}

type AccountEntry struct {
	Entry    Entry
	DCEntry  Entry
	SecEntry Entry
}

var heightKeyBinary []byte = []byte("height")

func extractEntryFromAccountTuple(account interface{}, position int64) (*Entry, error) {
	entryBin := account.(erlang.OtpErlangTuple)[position].(erlang.OtpErlangBinary).Value
	entry, eErr := erlang.BinaryToTerm(entryBin[1:])
	if eErr != nil {
		return nil, eErr
	}

	entryName := fmt.Sprint(entry.(erlang.OtpErlangTuple)[0].(erlang.OtpErlangAtom))
	entryNonce := utils.Num64(entry.(erlang.OtpErlangTuple)[1])
	entryAmount := utils.Num64(entry.(erlang.OtpErlangTuple)[2])

	return &Entry{
		entryName,
		entryNonce,
		entryAmount,
	}, nil
}

func addressToBinary(address string) ([]byte, error) {
	addressBinary, _, err := base58.CheckDecode("14XASpMbTophdTzmEwc2hTTyv78YmWgu2ckgUmEZeXX1poXGuHZ")
	if err != nil {
		return nil, err
	}

	return addressBinary, nil
}

func heightToBinary(height int64) ([]byte, error) {
	heightBinary := new(bytes.Buffer)
	err := binary.Write(heightBinary, binary.BigEndian, height)
	if err != nil {
		return nil, err
	}

	return heightBinary.Bytes(), nil
}

func RocksDBBlockHashGet(height int64) (*string, error) {
	heightBin, hbErr := heightToBinary(height)
	if hbErr != nil {
		return nil, hbErr
	}

	ro := rocksdb.NewDefaultReadOptions()
	hashBin, hErr := NodeBlocksDB.GetCF(ro, NodeBlockchainDBHeightsHandle, heightBin)
	if hErr != nil {
		return nil, hErr
	}

	hash := b64.RawURLEncoding.EncodeToString(hashBin.Data())

	return &hash, nil
}

func RocksDBTransactionsHeightGet() (*int64, error) {
	ro := rocksdb.NewDefaultReadOptions()

	heightBin, hErr := NodeTransactionsDB.GetCF(ro, NodeTransactionsDBDefaultHandle, heightKeyBinary)
	if hErr != nil {
		return nil, hErr
	}

	height := int64(binary.LittleEndian.Uint64(heightBin.Data()))

	return &height, nil
}

func RocksDBBalancesHeightGet() (*int64, error) {
	ro := rocksdb.NewDefaultReadOptions()

	heightBin, hErr := NodeBalancesDB.GetCF(ro, NodeBalancesDBDefaultHandle, heightKeyBinary)
	if hErr != nil {
		return nil, hErr
	}

	zap.S().Info(fmt.Sprint(heightBin))

	height := int64(binary.LittleEndian.Uint64(heightBin.Data()))

	return &height, nil
}

func RocksDBAccountGet(address string, height int64) (*AccountEntry, error) {
	addressBin, abErr := addressToBinary(address)
	if abErr != nil {
		return nil, abErr
	}

	heightBin, hbErr := heightToBinary(height)
	if hbErr != nil {
		return nil, hbErr
	}

	key := append(addressBin, heightBin...)

	readOptions := rocksdb.NewDefaultReadOptions()
	readOptions.SetFillCache(false)
	readOptions.SetTotalOrderSeek(true)
	iterator := NodeBalancesDB.NewIteratorCF(readOptions, NodeBalancesDBEntriesHandle)
	defer iterator.Close()
	iterator.SeekForPrev(key)

	if iterator.ValidForPrefix(addressBin) {
		accountEntryBin := iterator.Value()
		accountEntryTuple, bErr := erlang.BinaryToTerm(accountEntryBin.Data())
		if bErr != nil {
			return nil, bErr
		}
		accountEntryBin.Free()

		entry, entryErr := extractEntryFromAccountTuple(accountEntryTuple, 0)
		if entryErr != nil {
			return nil, entryErr
		}

		dcEntry, entryErr := extractEntryFromAccountTuple(accountEntryTuple, 1)
		if entryErr != nil {
			return nil, entryErr
		}

		secEntry, entryErr := extractEntryFromAccountTuple(accountEntryTuple, 2)
		if entryErr != nil {
			return nil, entryErr
		}

		accountEntry := &AccountEntry{
			*entry,
			*dcEntry,
			*secEntry,
		}
		return accountEntry, nil
	} else {
		return nil, errors.New("invalid iterator")
	}
}
