package rpcx

const (
	rpcVersion   = "2.0"
	defaultRpcId = 1

	// methods
	MethodBlockNumber           = "eth_blockNumber"
	MethodCall                  = "eth_call"
	MethodChainId               = "eth_chainId"
	MethodEstimateGas           = "eth_estimateGas"
	MethodGetBlockByNumber      = "eth_getBlockByNumber"
	MethodGasPrice              = "eth_gasPrice"
	MethodGetBalance            = "eth_getBalance"
	MethodGetLogs               = "eth_getLogs"
	MethodNonceAt               = "eth_getTransactionCount"
	MethodGetTransactionByHash  = "eth_getTransactionByHash"
	MethodGetTransactionReceipt = "eth_getTransactionReceipt"
	MethodSendRawTransaction    = "eth_sendRawTransaction"
)
