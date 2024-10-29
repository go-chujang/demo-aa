package main

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chujang/demo-aa/common/utils/conv"
)

func state() (string, bool) {
	res, _, err := get("/svc/v1/users/accounts/state", nil)
	if err != nil {
		return "please retry", false
	}
	parsed := res.(map[string]interface{})
	parsed["ownerBalance"] = hexutil.MustDecodeBig(parsed["ownerBalance"].(string))
	parsed["accountBalance"] = hexutil.MustDecodeBig(parsed["accountBalance"].(string))
	parsed["accountNonce"] = hexutil.MustDecodeBig(parsed["accountNonce"].(string))
	bytes, _ := json.MarshalIndent(parsed, "", "    ")
	return conv.B2S(bytes), false
}
