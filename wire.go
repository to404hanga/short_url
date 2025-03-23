//go:build wireinject

package main

import (
	"short_url/ioc"
	"short_url/repository/dao"
	"short_url/web"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func Init() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitRedis,

		dao.NewGormShortUrlDAO,

		ioc.InitRedisCache,
		ioc.InitCachedRepository,
		ioc.InitService,

		web.NewApiHandler,
		web.NewServerHandler,

		ioc.InitWebServer,
	)
	return gin.Default()
}
