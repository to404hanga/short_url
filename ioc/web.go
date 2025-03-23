package ioc

import (
	"short_url/web"

	"github.com/gin-gonic/gin"
)

func InitWebServer(apiSrv *web.ApiHandler, serverSrv *web.ServerHandler) *gin.Engine {
	router := gin.Default()

	apiSrv.RegisterRoutes(router)
	serverSrv.RegisterRoutes(router)

	return router
}
