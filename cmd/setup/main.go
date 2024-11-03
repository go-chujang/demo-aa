package main

import (
	"log"
	"os"
	"regexp"

	"github.com/go-chujang/demo-aa/config"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/kafka"
	"github.com/go-chujang/demo-aa/platform/mongox"
)

var (
	rpcuri        = config.Get(config.RPC_ENDPOINT)
	mnemonic      = config.Get(config.MNEMONIC_STRING)
	accountCnt, _ = config.GetInt(config.MNEMONIC_ACCOUNTS)
	db            *mongox.Client
	err           error
)

func main() {
	switch addr, version := kafka.EnvAddrsVersion(); {
	case addr == nil:
		panic("KAFKA_ADDRS must not be empty")
	case !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(version.String()):
		panic("KAFKA_VERSION must not be empty")
	default:
		if db, err = mongox.New(mongox.EnvUri()); err != nil {
			panic(err)
		}
		defer db.Stop()
		log.Println("start setup")
	}
	// deploy contract & insert mngd_accounts
	if cnt, err := db.Count(mongox.NewQuery().SetColl(model.CollectionMngdAccounts)); err != nil {
		panic(err)
	} else if cnt == 0 {
		setup_aa()
	}
	log.Println("done setup_aa")

	// kafka topics
	setup_kafka()
	log.Println("done setup_kafka")
	os.Exit(0)
}
