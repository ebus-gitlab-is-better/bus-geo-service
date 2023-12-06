package server

import (
	"bus-geo-service/internal/biz"
	"bus-geo-service/internal/conf"
	"bus-geo-service/internal/route"
	"bus-geo-service/pkg/customhttp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewCustomHttp(
	c *conf.Server,
	uc *biz.BusUseCase,
	logger log.Logger) *customhttp.CustomHTTP {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Custom.Network != "" {
		opts = append(opts, http.Network(c.Custom.Network))
	}
	if c.Custom.Addr != "" {
		opts = append(opts, http.Address(c.Custom.Addr))
	}
	if c.Custom.Timeout != nil {
		opts = append(opts, http.Timeout(c.Custom.Timeout.AsDuration()))
	}
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	config.AllowCredentials = true
	r.Use(cors.New(config))
	srv := http.NewServer(opts...)
	bus := route.NewBusRouter(uc)
	chatsRoute := r.Group("/geo")
	bus.Register(chatsRoute)
	srv.HandlePrefix("/", r)
	return &customhttp.CustomHTTP{Http: srv}
}
