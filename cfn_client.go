package rediscloud_api

import (
	"net/http"
	"os"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/service/account"
	"github.com/RedisLabs/rediscloud-go-api/service/cloud_accounts"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
)

func NewClientV2(configs ...Option) (*Client, error) {
	config := &Options{
		baseUrl:   "https://api.redislabs.com/v1",
		userAgent: userAgent,
		apiKey:    os.Getenv(AccessKeyEnvVar),
		secretKey: os.Getenv(SecretKeyEnvVar),
		logger:    &defaultLogger{},
		transport: http.DefaultTransport,
	}

	for _, option := range configs {
		option(config)
	}

	httpClient := &http.Client{
		Transport: config.roundTripper(),
	}

	client, err := internal.NewHttpClient(httpClient, config.baseUrl)
	if err != nil {
		return nil, err
	}

	t := internal.NewAPI(client, config.logger)

	a := account.NewAPI(client)
	c := cloud_accounts.NewAPIV2(client, t, config.logger)
	d := databases.NewAPI(client, t, config.logger)
	s := subscriptions.NewAPI(client, t, config.logger)

	return &Client{
		Account:      a,
		CloudAccount: c,
		Database:     d,
		Subscription: s,
	}, nil
}
