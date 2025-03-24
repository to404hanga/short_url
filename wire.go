//go:build wireinject

package main

import (
	"short_url/ioc"
	"short_url/repository/dao"
	"short_url/web"

	"github.com/google/wire"
)

func Init() *App {
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

		ioc.InitCleanerJob,
		ioc.InitJobs,
		ioc.InitGinMiddleware,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
