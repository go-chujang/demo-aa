package ethutil

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func PackLowHighByBig(low128, high128 *big.Int) [32]byte {
	if low128 == nil {
		low128 = big.NewInt(0)
	}
	if high128 == nil {
		high128 = big.NewInt(0)
	}

	var buf [32]byte
	copy(buf[:16], common.LeftPadBytes(low128.Bytes(), 16))
	copy(buf[16:], common.LeftPadBytes(high128.Bytes(), 16))
	return buf
}
