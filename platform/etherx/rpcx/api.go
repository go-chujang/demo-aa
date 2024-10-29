package rpcx

import (
	"encoding/json"
	"reflect"

	"github.com/go-chujang/demo-aa/common/net/httpf"
)

func DoReq[T any](uri string, req Request) (result T, err error) {
	res, err := httpf.Post(uri, req)
	if err != nil {
		return
	}

	var parsed Response[T]
	switch err = json.Unmarshal(res, &parsed); {
	case err != nil:
		return result, err
	case parsed.Error != nil:
		return result, parsed.Error
	default:
		return parsed.Result, nil
	}
}

func Do[T any](uri string, method string, params ...interface{}) (result T, err error) {
	return DoReq[T](uri, Request{JsonRpc: rpcVersion, Id: defaultRpcId, Method: method, Params: params})
}

func ReqId(method string, id int, params ...interface{}) Request {
	return Request{JsonRpc: rpcVersion, Id: id, Method: method, Params: params}
}
func Req(method string, params ...interface{}) Request {
	return ReqId(method, defaultRpcId, params...)
}
func (q *Request) SetId(id int) *Request {
	q.Id = id
	return q
}

type (
	BatchElem     Response[json.RawMessage]
	BatchResponse []BatchElem
)

func Batch(uri string, reqs []Request) (BatchResponse, error) {
	req, err := json.Marshal(reqs)
	if err != nil {
		return nil, err
	}
	res, err := httpf.Post(uri, req)
	if err != nil {
		return nil, err
	}
	result := make(BatchResponse, 0, len(reqs))
	return result, json.Unmarshal(res, &result)
}

func ParseBatchElem[T any](elem BatchElem) (id int, parsed T, err error) {
	switch {
	case elem.Error != nil:
		err = elem.Error
	case reflect.TypeOf(parsed).Kind() == reflect.Ptr:
		err = json.Unmarshal(elem.Result, parsed)
	default:
		err = json.Unmarshal(elem.Result, &parsed)
	}
	return elem.Id, parsed, err
}
