package data

import (
	"bus-geo-service/internal/conf"
	"context"
	"crypto/tls"

	"github.com/Nerzal/gocloak/v13"
	gocent "github.com/centrifugal/gocent/v3"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedisClient, NewCent, NewBusRepo, NewKeyCloakAPI, NewKeycloak)

// Data .
type Data struct {
	rdb  *redis.Client
	cent *gocent.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, rdb *redis.Client, cent *gocent.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{rdb: rdb, cent: cent}, cleanup, nil
}

func NewRedisClient(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		DB:           0,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Errorf("failed opening connection to redis: %v", err)
		panic("failed to connect redis")
	}
	return rdb
}

func NewKeycloak(c *conf.Data) *gocloak.GoCloak {
	client := gocloak.NewClient(c.Keycloak.Hostname)
	restyClient := client.RestyClient()
	// restyClient.SetDebug(true)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return client
}
func NewCent(conf *conf.Data) *gocent.Client {
	c := gocent.New(gocent.Config{
		Addr: conf.AddressMessage,
		Key:  conf.ApiKey,
	})
	return c
}
