package ethutil

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/common/utils/conv"
)

func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func isHex(str string) bool {
	if !Has0xPrefix(str) {
		return false
	}
	for _, c := range conv.S2B(Trim0x(str)) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

func Has0xPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

func Trim0x(hex string) string {
	if Has0xPrefix(hex) {
		hex = hex[2:]
	}
	return hex
}

func IsHexHash(s string) bool {
	s = Trim0x(s)
	return len(s) == 64 && isHex(s)
}

func HexOrString2X[T *big.Int | uint64](hexOrString string) (x T, err error) {
	var (
		tmp  any
		base = 10
	)
	if isHex(hexOrString) {
		hexOrString = Trim0x(hexOrString)
		base = 16
	}
	switch v := any(x).(type) {
	case *big.Int:
		var ok bool
		v, ok = new(big.Int).SetString(hexOrString, base)
		if !ok {
			err = errors.New("invalid hex string")
			return
		}
		tmp = v
	case uint64:
		v, err = strconv.ParseUint(hexOrString, base, 64)
		if err != nil {
			return
		}
		tmp = v
	}
	return tmp.(T), nil
}

func Hex2X[T *big.Int | uint64 | []byte](hex string) (x T, err error) {
	if !isHex(hex) {
		return x, errors.New("missing 0x prefix")
	}
	var tmp any
	switch any(x).(type) {
	case *big.Int:
		tmp, err = HexOrString2X[*big.Int](hex)
	case uint64:
		tmp, err = HexOrString2X[uint64](hex)
	case []byte:
		tmp = common.Hex2Bytes(Trim0x(hex))
	}
	if err != nil {
		return x, err
	}
	return tmp.(T), nil
}
