package rpcx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func EasyReqBig(uri string, req Request) (*big.Int, error) {
	hex, err := DoReq[string](uri, req)
	if err != nil {
		return nil, err
	}
	return hexutil.DecodeBig(hex)
}

func EasyReqUint(uri string, req Request) (uint64, error) {
	hex, err := DoReq[string](uri, req)
	if err != nil {
		return 0, err
	}
	return hexutil.DecodeUint64(hex)
}
