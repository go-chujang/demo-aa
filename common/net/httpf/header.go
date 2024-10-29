package httpf

import "github.com/valyala/fasthttp"

type Header struct {
	Key   string
	Value string
}

func NewHeader(key string, value string) Header {
	return Header{Key: key, Value: value}
}

func NewHeaderAuthorization(auth string) Header {
	return Header{Key: fasthttp.HeaderAuthorization, Value: auth}
}

func NewHeaderContentType(mimeType string) Header {
	return Header{Key: fasthttp.HeaderContentType, Value: mimeType}
}

func toHeaderMap(headers ...Header) map[string]string {
	headerMap := make(map[string]string, len(headers))
	for _, v := range headers {
		headerMap[v.Key] = v.Value
	}
	return headerMap
}
