package main

import (
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
)

func progressTransfer() (msg string, ok bool) {
	var err error
	switch SEQ_PROGRESS {
	case 0:
		return "receiver address", true
	case 1:
		if !ethutil.IsAddress(LAST_INPUT) {
			return "invalid to address", true
		}
		STORED_PARAM["to"] = LAST_INPUT
		return "value for transfer", true
	case 2:
		value := LAST_INPUT
		err = tokenTransfer(STORED_PARAM["to"], value)
		DONE_PROGRESS = true
	}
	if err != nil {
		return err.Error(), false
	}
	return "success", true
}

func tokenTransfer(to, value string) error {
	return postUserOp("/svc/v1/users/operations/tokens/transfer", map[string]interface{}{
		"to":    to,
		"value": value,
	})
}
