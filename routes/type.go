package routes

import "github.com/gin-gonic/gin"

type Route interface {
	RegisterRoutes(srv *gin.Engine)
}
