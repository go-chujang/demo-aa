package main

import (
	"github.com/go-chujang/demo-aa/api/middleware/routechart"
	"github.com/go-chujang/demo-aa/api/router"
	"github.com/go-chujang/demo-aa/common/logx"
	"github.com/go-chujang/demo-aa/common/sig"
	"github.com/go-chujang/demo-aa/config"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/kafka"
	"github.com/go-chujang/demo-aa/platform/mongox"
)

func main() {
	addrs, version := kafka.EnvAddrsVersion()
	if err := kafka.UseProducer(addrs, version); err != nil {
		panic(err)
	}

	db, err := mongox.New(mongox.EnvUri())
	if err != nil {
		panic(err)
	}
	mongox.SetDefault(db)

	if err = aa.UseApi(db, config.Get(config.RPC_ENDPOINT)); err != nil {
		panic(err)
	}

	timeout, _ := config.GetDuration(config.API_TIMEOUT)
	app, err := router.New(router.Config{
		AppName:        config.AppTag(),
		Port:           config.Get(config.API_PORT),
		CtxTimeout:     timeout,
		ReleaseMode:    config.IsRelease(),
		RequestLogger:  logx.GetLogWriter(),
		ResponseLogger: logx.GetLogWriter(),
		RouteChart:     []routechart.Chart{routechart.ServiceV1},
	})
	if err != nil {
		panic(err)
	}

	go app.Start()

	sig.Wait(sig.DefaultSigs, logx.GetLogWriter(),
		app.Stop,
		db.Stop,
	)
}
