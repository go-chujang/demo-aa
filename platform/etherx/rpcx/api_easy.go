package rpcx

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chujang/demo-aa/common/net/httpf"
)

func EasyEstimateGas(uri string, from, to common.Address, packed []byte) (*big.Int, error) {
	return EasyBig(uri, MethodEstimateGas, ethereum.CallMsg{
		From: from,
		To:   &to,
		Data: packed,
	})
}

func EasyBig(uri string, method string, params ...interface{}) (*big.Int, error) {
	hex, err := Do[string](uri, method, params...)
	if err != nil {
		return nil, err
	}
	return hexutil.DecodeBig(hex)
}

func EasyUint(uri string, method string, params ...interface{}) (uint64, error) {
	hex, err := Do[string](uri, method, params...)
	if err != nil {
		return 0, err
	}
	return hexutil.DecodeUint64(hex)
	// return strconv.ParseUint(ethutil.Trim0x(hex), 16, 64)
}

func EasyBatch(uri string, reqs []Request) (ids []int, results []string, rpcErrs []error, err error) {
	req, err := json.Marshal(reqs)
	if err != nil {
		return nil, nil, nil, err
	}
	res, err := httpf.Post(uri, req)
	if err != nil {
		return nil, nil, nil, err
	}

	cap := len(reqs)
	batchRes := make([]Response[string], 0, cap)
	if err = json.Unmarshal(res, &batchRes); err != nil {
		return nil, nil, nil, err
	}

	ids = make([]int, cap)
	results = make([]string, cap)
	rpcErrs = make([]error, cap)

	for i, elem := range batchRes {
		ids[i] = elem.Id
		results[i] = elem.Result
		if elem.Error != nil {
			rpcErrs[i] = elem.Error
		}
	}
	return ids, results, rpcErrs, nil
}
