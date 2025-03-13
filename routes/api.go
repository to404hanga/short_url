package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiRoute struct {
}

var _ Route = (*ApiRoute)(nil)

func NewApiRoute() *ApiRoute {
	return &ApiRoute{}
}

func (r *ApiRoute) RegisterRoutes(srv *gin.Engine) {
	api := srv.Group("/api")
	{
		api.GET("/ping", r.Ping)
		api.POST("/create", r.Create)
	}
}

func (r *ApiRoute) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong")
}

func (r *ApiRoute) Create(ctx *gin.Context) {
	type CreateRequest struct {
		OriginURL string `json:"origin_url"`
	}
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.OriginURL == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid params",
		})
		return
	}

}
