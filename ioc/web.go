package ioc

import (
	"short_url/web"
	"short_url/web/middlewares"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/to404hanga/pkg404/logger"
)

func InitWebServer(apiSrv *web.ApiHandler, serverSrv *web.ServerHandler, mdls []gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	router.Use(mdls...)
	apiSrv.RegisterRoutes(router)
	serverSrv.RegisterRoutes(router)

	return router
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
		middlewares.ZapLogger(l),
	}
}
