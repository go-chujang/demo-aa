package httpf

import (
	"time"

	"github.com/go-chujang/demo-aa/internal/json"
	"github.com/gofiber/fiber/v3/client"
	"github.com/valyala/fasthttp"
)

const defaultTimeout = time.Second * 3

var defaultClient *client.Client

func init() {
	defaultClient = client.New().
		SetTimeout(defaultTimeout).
		SetHeader(fasthttp.HeaderAccept, "application/json").
		SetHeader(fasthttp.HeaderContentType, "application/json").
		SetJSONMarshal(json.Marshal).
		SetJSONUnmarshal(json.Unmarshal)
}

func SetTimeout(t time.Duration) {
	defaultClient.SetTimeout(t)
}

func SetHeader(header Header) {
	defaultClient.SetHeader(header.Key, header.Value)
}
