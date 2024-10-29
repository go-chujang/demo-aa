package main

import (
	"github.com/go-chujang/demo-aa/api/middleware/routechart"
	"github.com/go-chujang/demo-aa/api/router"
	"github.com/go-chujang/demo-aa/common/logx"
	"github.com/go-chujang/demo-aa/common/sig"
	"github.com/go-chujang/demo-aa/config"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/mongox"
)

func main() {
	var (
		rpcUri        = config.Get(config.RPC_ENDPOINT)
		interval, _   = config.GetDuration(config.WATCHDOG_INTERVAL)
		blockRange, _ = config.GetUint64(config.WATCHDOG_BLOCK_RANGE)
	)
	db, err := mongox.New(mongox.EnvUri())
	if err != nil {
		panic(err)
	}
	watchdog, err := aa.NewWatchDog(db, rpcUri, interval, blockRange)
	if err != nil {
		panic(err)
	}
	app, err := router.New(router.Config{
		AppName: config.AppTag(),
		// Port:           config.Get(config.API_PORT),
		Port:           "5100",
		ReleaseMode:    config.IsRelease(),
		RequestLogger:  logx.GetLogWriter(),
		ResponseLogger: logx.GetLogWriter(),
		RouteChart: []routechart.Chart{
			chart(db, watchdog, interval, blockRange),
		},
	})
	if err != nil {
		panic(err)
	}

	go watchdog.Start()
	go app.Start()

	sig.Wait(sig.DefaultSigs, logx.GetLogWriter(),
		app.Stop,
		watchdog.Stop,
		db.Stop,
	)
}
