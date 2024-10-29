package auth

import (
	"encoding/base64"
	"strings"

	"github.com/go-chujang/demo-aa/common/utils/comp"
	"github.com/go-chujang/demo-aa/common/utils/conv"
)

func parseBasic(auth string) (string, string, bool) {
	if len(auth) <= 6 || !comp.EqualFold(auth[:6], "basic ") {
		return "", "", false
	}
	raw, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return "", "", false
	}
	creds := conv.B2S(raw)
	index := strings.Index(creds, ":")
	if index == -1 {
		return "", "", false
	}
	return creds[:index], creds[index+1:], true
}
