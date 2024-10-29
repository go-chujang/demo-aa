package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chujang/demo-aa/common/net/httpf"
	"github.com/go-chujang/demo-aa/common/utils/conv"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/gofiber/fiber/v3"
)

func basicauth() httpf.Header {
	id, pass, ok := getIdPassword()
	if !ok {
		return httpf.Header{}
	}
	auth := base64.StdEncoding.EncodeToString(conv.S2B(id + ":" + pass))
	return httpf.NewHeader("Authorization", fmt.Sprintf("Basic %s", auth))
}

func get(path string, req map[string]interface{}) (interface{}, int, error) {
	return do[interface{}](fiber.MethodGet, path, req)
}

func post(path string, req map[string]interface{}) (interface{}, int, error) {
	return do[interface{}](fiber.MethodPost, path, req)
}

func do[T any](method, path string, req map[string]interface{}) (data T, code int, err error) {
	var (
		uri = fmt.Sprintf("%s%s", URL, path)
		res struct {
			Success bool    `json:"success"`
			Code    int     `json:"code"`
			Message *string `json:"message,omitempty"`
			Data    T       `json:"data,omitempty"`
		}
	)
	if err = httpf.Json(&res, uri, method, req, basicauth()); err != nil {
		return data, 0, err
	}
	if res.Code != 0 {
		err = fmt.Errorf("path: %s, err: %v", path, *res.Message)
	}
	return res.Data, res.Code, err
}

type (
	userop         model.PackedUserOperation
	useropWithHash struct {
		UserOperation *userop       `json:"userOperation"`
		Hash          hexutil.Bytes `json:"hash"`
	}
)

var (
	_ json.Marshaler   = (*userop)(nil)
	_ json.Unmarshaler = (*userop)(nil)
)

func postUserOp(uri string, data map[string]interface{}) error {
	res, _, err := do[aa.UserOperationWithHash](fiber.MethodPost, uri, data)
	if err != nil {
		return err
	}
	sig, err := ethutil.Signature(res.Hash, PRIVATE_KEY)
	if err != nil {
		return err
	}
	uop := res.UserOperation.SetSignature(sig)
	_, _, err = post("/svc/v1/users/operations/execute", map[string]interface{}{
		"op": uop,
	})
	return err
}

func (o userop) MarshalJSON() ([]byte, error) { return json.Marshal(model.PackedUserOperation(o)) }
func (o *userop) UnmarshalJSON(input []byte) error {
	var dec model.PackedUserOperation
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	*o = userop(dec)
	return nil
}

/*
func postUserOp(uri string, data map[string]interface{}) error {
	res, _, err := do[aa.UserOperationWithHash](fiber.MethodPost, uri, data)
	if err != nil {
		return err
	}
	sig, err := ethutil.Signature(res.Hash, PRIVATE_KEY)
	if err != nil {
		return err
	}
	uop := res.UserOperation.SetSignature(sig)
	_, _, err = post("/svc/v1/users/operations/execute", map[string]interface{}{
		"op": uop,
	})
	return err
}
*/
