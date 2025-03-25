package ioc

import (
	"context"
	"short_url/web"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/to404hanga/pkg404/logger"
)

func InitWebServer(apiSrv *web.ApiHandler, serverSrv *web.ServerHandler, mdls []gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 1024 * 1024 * 1024 * 2

	router.Use(mdls...)
	apiSrv.RegisterRoutes(router)
	serverSrv.RegisterRoutes(router)

	return router
}

// 强制给每个请求 5 秒超时时间
func timeout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := ctx.Request.Context().Deadline(); !ok {
			newCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*5)
			defer cancel()
			ctx.Request = ctx.Request.Clone(newCtx)
		}
		ctx.Next()
	}
}

func InitGinMiddleware(l logger.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "127.0.0.1") {
					return true
				}
				return strings.Contains(origin, "your_company.com") // TODO 将 your_company.com 改为实际前端服务器域名或 ip 地址
			},
			MaxAge: 12 * time.Hour,
		}),
		timeout(),
		// middlewares.ZapLogger(l),
	}
}
