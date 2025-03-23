package web

import (
	"short_url/service"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	svc service.ShortUrlService
}

var _ Handler = (*ApiHandler)(nil)

func NewApiHandler(svc service.ShortUrlService) *ApiHandler {
	return &ApiHandler{
		svc: svc,
	}
}

func (h *ApiHandler) RegisterRoutes(srv *gin.Engine) {
	api := srv.Group("/api")
	{
		api.POST("/create", h.Create)
	}
}

func (h *ApiHandler) Create(ctx *gin.Context) {
	type CreateRequest struct {
		OriginUrl string `json:"origin_url"`
	}
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := h.svc.Create(ctx.Request.Context(), req.OriginUrl)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"short_url": shortUrl,
	})
}
