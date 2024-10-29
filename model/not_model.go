package model

import (
	"math/big"
	"unicode"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chujang/demo-aa/common/utils/conv"
)

func toHexBsonValue(value interface{}) interface{} {
	switch v := value.(type) {
	case common.Address:
		return v.Hex()
	case []byte:
		return hexutil.Encode(v)
	case *big.Int:
		return hexutil.EncodeBig(v)
	case uint64:
		return hexutil.EncodeUint64(v)
	default:
		return v
	}
}

// convert [common.Address | []byte | *big.Int | uint64]
func toHexBsonMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	hm := make(map[string]interface{}, len(m))
	for key, value := range m {
		hm[key] = toHexBsonValue(value)
	}
	return hm
}

func toHintString(m RawDataHelper) *string {
	var (
		typ   = conv.TypeOf(m, true)
		size  = len(typ)
		bytes = make([]byte, 0, size-1)
	)
	if size > 1 {
		bytes = append(bytes, byte(unicode.ToLower(rune(typ[1]))))
	}
	if size > 2 {
		bytes = append(bytes, typ[2:]...)
	}
	hint := conv.B2S(bytes)
	return &hint
}
