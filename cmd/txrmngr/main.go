package main

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/go-chujang/demo-aa/common/logx"
	"github.com/go-chujang/demo-aa/common/sig"
	"github.com/go-chujang/demo-aa/common/utils/id"
	"github.com/go-chujang/demo-aa/common/utils/slice"
	"github.com/go-chujang/demo-aa/config"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx/aa"
	"github.com/go-chujang/demo-aa/platform/kafka"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		groupId        = config.Get(config.CONSUMER_GROUP_ID)
		addrs, version = kafka.EnvAddrsVersion()
		maxInterval, _ = config.GetDuration(config.CONSUME_MAX_INTERVAL)
		batchSize, _   = config.GetInt(config.CONSUME_BATCHSIZE)
		topics         = []kafka.Topic{kafka.TopicMembership, kafka.TopicLiquidity}
		rpcUri         = config.Get(config.RPC_ENDPOINT)
	)
	err := kafka.UseConsumer(addrs, version, groupId, maxInterval, batchSize)
	if err != nil {
		panic(err)
	}
	db, err := mongox.New(mongox.EnvUri())
	if err != nil {
		panic(err)
	}
	mngr, err := aa.NewMngrBackend(db, rpcUri)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	closefn := func() error {
		cancel()
		return errors.Join(db.Stop(), kafka.Unuse())
	}

	eg := new(errgroup.Group)
	eg.Go(func() error {
		return kafka.Consume(ctx, topics, func(msgs []*sarama.ConsumerMessage) ([]bool, []error) {
			var (
				size           = len(msgs)
				payloadList    = make([]aa.BundlePayload, size)
				rawDataHelpers = make([]model.RawDataHelper, size)
				markFlags      = make([]bool, size)
				errs           = make([]error, size)
				reOrdering     = false
			)
			debugkey := id.Uuid()
			logx.Debug(config.AppTag(), "input messages %d | %s", size, debugkey)

			for idx, message := range msgs {
				var parseErr error
				switch message.Topic {
				case kafka.TopicMembership.String():
					param := &model.CreateAccount{}
					if parseErr = param.Parse(message); parseErr != nil {
						markFlags[idx] = true
					} else {
						payloadList[idx] = mngr.PrepareCreateAccount(param)
						rawDataHelpers[idx] = param
					}
				case kafka.TopicLiquidity.String():
					param := &model.Faucet{}
					if parseErr = param.Parse(message); parseErr != nil {
						markFlags[idx] = true
					} else {
						payloadList[idx] = mngr.PrepareFaucet(param)
						rawDataHelpers[idx] = param
					}
				default:
					parseErr = errors.New("mismatced topic")
				}

				errs[idx] = parseErr
				if !reOrdering && parseErr != nil {
					reOrdering = true
				}
			}
			if reOrdering {
				logx.Debug(config.AppTag(), "reordered | %s", debugkey)
				return markFlags, errs
			}

			bundleErrs, err := mngr.BundleExec(ctx, payloadList, rawDataHelpers)
			logx.Debug(config.AppTag(), "bundleErrs: %v, err: %v | %s", bundleErrs, err, debugkey)
			if err != nil {
				return slice.WithDefault(false, size), bundleErrs
			}
			return slice.WithDefault(true, size), bundleErrs
		})
	})

	select {
	case <-sig.Chan(sig.DefaultSigs, logx.GetLogWriter(), closefn):
	case err = <-func() <-chan error {
		ch := make(chan error, 1)
		go func() {
			ch <- eg.Wait()
			close(ch)
		}()
		return ch
	}():
		if closeErr := closefn(); closeErr != nil {
			err = errors.Join(err, closefn())
		}
		logx.Critical(groupId, err.Error())
	}
}