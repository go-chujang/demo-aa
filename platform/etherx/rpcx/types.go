package rpcx

import (
	"github.com/ethereum/go-ethereum"
)

type Request struct {
	JsonRpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Response[T any] struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  T      `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type (
	BlockParam string

	// BlockHash always ignore
	FilterQuery ethereum.FilterQuery
)
