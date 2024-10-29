package rpcx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("json-rpc error %d %s", err.Code, err.Message)
}

func (err *Error) ErrorCode() int {
	return err.Code
}

func (err *Error) ErrorData() interface{} {
	return err.Data
}

func (b BlockParam) String() string {
	return string(b)
}

var (
	BlockParamLatest  = BlockParam(rpc.LatestBlockNumber.String())
	BlockParamPending = BlockParam(rpc.PendingBlockNumber.String())
)

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	return fmt.Sprintf("<invalid %d>", number)
}

func (q FilterQuery) toArg() map[string]interface{} {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	} else {
		arg["fromBlock"] = toBlockNumArg(q.FromBlock)
	}
	arg["toBlock"] = toBlockNumArg(q.ToBlock)
	return arg
}
