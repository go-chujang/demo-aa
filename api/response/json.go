package response

import (
	"encoding/json"
)

var _ json.Marshaler = (*body)(nil)

func (b body) MarshalJSON() ([]byte, error) {
	type body struct {
		Success bool        `json:"success"`
		Code    errorCode   `json:"code"`
		Message *string     `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
	var enc body

	enc.Success = !(b.data == nil || b.code > 0 || b.message != nil || b.err != nil)
	if enc.Success {
		enc.Data = b.data
	} else {
		var (
			code       = ErrCodeDefault
			msg, exist = errCodeMessages[b.code]
		)
		switch {
		case exist:
			code = b.code
		case b.message != nil:
			msg = *b.message
		case b.err != nil:
			msg = b.err.Error()
		case b.data == nil:
			code = ErrCodeEmptyData
		}
		if msg == "" {
			msg = errCodeMessages[code]
		}
		enc.Code = code
		enc.Message = &msg
		enc.Data = nil
	}
	return json.Marshal(&enc)
}
