package rpcx

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-chujang/demo-aa/common/net/httpf"
)

func FilterLogs(uri string, query FilterQuery) ([]types.Log, error) {
	args := query.toArg()
	return Do[[]types.Log](uri, MethodGetLogs, args)
}

func FilterLogs2[E any](uri string, query FilterQuery) ([]E, error) {
	args := query.toArg()
	return Do[[]E](uri, MethodGetLogs, args)
}

func GetTxReceipts(uri string, txHashes []string) ([]types.Receipt, error) {
	return GetTxReceipts2[types.Receipt](uri, txHashes)
}

func GetTxReceipts2[E any](uri string, txHashes []string) ([]E, error) {
	reqs := make([]Request, 0, len(txHashes))
	for i, v := range txHashes {
		reqs = append(reqs, ReqId(MethodGetTransactionReceipt, i, v))
	}
	res, err := httpf.Post(uri, reqs)
	if err != nil {
		return nil, err
	}
	var batchRes []Response[E]
	if err = json.Unmarshal(res, &batchRes); err != nil {
		return nil, err
	}

	receipts := make([]E, 0, len(batchRes))
	for _, elem := range batchRes {
		if elem.Error != nil {
			return nil, elem.Error
		}
		receipts = append(receipts, elem.Result)
	}
	return receipts, nil
}
