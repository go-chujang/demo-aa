package ethutil

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ZeroAddress    = common.Address{}
	ZeroAddressHex = common.Address{}.Hex()

	ErrInvalidAddress = errors.New("invalid address")
)

func IsAddress(address string) bool {
	return common.IsHexAddress(address)
}
