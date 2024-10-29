package mongox

import (
	"net/http"
	"time"

	"github.com/go-chujang/demo-aa/common/utils/conv"
)

type Config struct {
	Timeout         time.Duration
	MaxConnIdleTime time.Duration
	MaxPoolSize     uint64
	RetryWrites     *bool
	RetryReads      *bool
	HttpClient      *http.Client
	ReplicaSet      *string
}

var ConfigDefault = Config{
	Timeout:         time.Second * 3,
	MaxConnIdleTime: time.Minute * 30,
	MaxPoolSize:     100,
	RetryWrites:     conv.ToPtr(false),
	RetryReads:      conv.ToPtr(true),
	HttpClient:      nil,
	ReplicaSet:      conv.ToPtr(""),
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}
	cfg := config[0]
	if cfg.MaxPoolSize == 0 {
		cfg.MaxPoolSize = ConfigDefault.MaxPoolSize
	}
	if cfg.RetryWrites == nil {
		cfg.RetryWrites = ConfigDefault.RetryWrites
	}
	if cfg.RetryReads == nil {
		cfg.RetryReads = ConfigDefault.RetryReads
	}
	if cfg.HttpClient == nil {
		cfg.HttpClient = ConfigDefault.HttpClient
	}
	if cfg.ReplicaSet == nil {
		cfg.ReplicaSet = ConfigDefault.ReplicaSet
	}
	return cfg
}
