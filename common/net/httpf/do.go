package httpf

import (
	"context"

	"github.com/gofiber/fiber/v3/client"
	"github.com/valyala/fasthttp"
)

func Post(uri string, body interface{}, headers ...Header) ([]byte, error) {
	res, _, err := Do(context.Background(), uri, fasthttp.MethodPost, body, headers...)
	return res, err
}

func Get(uri string, body interface{}, headers ...Header) ([]byte, error) {
	res, _, err := Do(context.Background(), uri, fasthttp.MethodGet, body, headers...)
	return res, err
}

func Do(ctx context.Context, uri, method string, body interface{}, headers ...Header) ([]byte, int, error) {
	res, code, err := do(ctx, uri, method, body, headers...)
	if err != nil {
		return nil, 0, err
	}
	return res, code, nil
}

func Json(v interface{}, uri, method string, body interface{}, headers ...Header) error {
	res, _, err := do(context.Background(), uri, method, body, headers...)
	if err != nil {
		return err
	}
	return defaultClient.JSONUnmarshal()(res, v)
}

func do(ctx context.Context, uri, method string, body interface{}, headers ...Header) ([]byte, int, error) {
	req := client.AcquireRequest().SetClient(defaultClient).
		SetContext(ctx).
		SetURL(uri).
		SetMethod(method).
		SetHeaders(toHeaderMap(headers...))
	defer client.ReleaseRequest(req)

	switch b := body.(type) {
	case nil:
		// skip
	case []byte:
		req.SetRawBody(b)
	default:
		encoded, err := defaultClient.JSONMarshal()(b)
		if err != nil {
			return nil, 0, err
		}
		req.SetRawBody(encoded)
	}
	res, err := req.Send()
	if err != nil {
		return nil, 0, err
	}
	defer client.ReleaseResponse(res)
	return res.Body(), res.StatusCode(), nil
}
