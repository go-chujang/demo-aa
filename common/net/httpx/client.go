package httpx

import (
	"net"
	"net/http"
	"time"
)

var (
	timeout       = time.Second * 3
	defaultClient = &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout:     timeout,
			TLSHandshakeTimeout: timeout,
			Dial: (&net.Dialer{
				Timeout:   75 * time.Second,
				KeepAlive: 75 * time.Second,
			}).Dial,
		},
	}
)
