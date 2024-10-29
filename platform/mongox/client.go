package mongox

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connStringForm = "mongodb://%s:%s@%s/?%s"
	connTimeout    = time.Second * 3
)

type Client struct {
	cfg    Config
	client *mongo.Client
}

func New(uri string, config ...Config) (*Client, error) {
	cfg := configDefault(config...)
	opts := options.Client().ApplyURI(uri).
		SetConnectTimeout(connTimeout).
		SetTimeout(cfg.Timeout).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetRetryWrites(*cfg.RetryWrites).
		SetRetryReads(*cfg.RetryReads).
		SetHTTPClient(cfg.HttpClient).
		SetReplicaSet(*cfg.ReplicaSet)
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	cli, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	return &Client{cfg: cfg, client: cli}, cli.Ping(context.Background(), readpref.Primary())
}

func (c *Client) Stop() error {
	return c.client.Disconnect(context.Background())
}

func (c *Client) Sdk() *mongo.Client {
	return c.client
}

func (c *Client) ctx(ctxorigin ...context.Context) context.Context {
	if ctxorigin != nil {
		return ctxorigin[0]
	}
	return context.Background()
}

func BaseUri(user, password, host string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin",
		user,
		password,
		host,
	)
}

func EnvUri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin",
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
	)
}
