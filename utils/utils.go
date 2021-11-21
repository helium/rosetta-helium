package utils

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ybbus/jsonrpc"
)

var (
	NodeClient = jsonrpc.NewClient("http://localhost:4467")
)

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
