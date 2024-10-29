package response

import "errors"

type errorCode int

const (
	ErrCodeDefault errorCode = iota + 1000
	ErrCodeEmptyData
	ErrAlreadyExistUserId
	ErrCodeAccountNotYet
	ErrFaucetLimit
	ErrInvalidAddress
	ErrTxnPendingState
	ErrTxnNonceTooLow
	ErrTxnInvalidSig
	ErrTxnInsufficientToken
)

var (
	errCodeMessages = map[errorCode]string{
		ErrCodeDefault:          "undefined error",
		ErrCodeEmptyData:        "empty response data",
		ErrAlreadyExistUserId:   "alreay exist userId",
		ErrCodeAccountNotYet:    "account is not yet created. wait for sync or request creation",
		ErrFaucetLimit:          "faucet limit has not been reset yet",
		ErrInvalidAddress:       "invalid address",
		ErrTxnPendingState:      "transaction failed, still in pending state",
		ErrTxnNonceTooLow:       "transaction failed, nonce too low",
		ErrTxnInvalidSig:        "transaction failed, invalid signature",
		ErrTxnInsufficientToken: "transaction failed, insufficient token",
	}
	errCodeErrors map[errorCode]error
)

func init() {
	errCodeErrors = make(map[errorCode]error, len(errCodeMessages))
	for code, msg := range errCodeMessages {
		errCodeErrors[code] = errors.New(msg)
	}
}

func ErrByCode(code errorCode) error {
	return errCodeErrors[code]
}
