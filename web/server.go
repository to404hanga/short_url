package web

import (
	"short_url/service"

	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	svc service.ShortUrlService
}

var _ Handler = (*ServerHandler)(nil)

func NewServerHandler(svc service.ShortUrlService) *ServerHandler {
	return &ServerHandler{
		svc: svc,
	}
}

func (h *ServerHandler) RegisterRoutes(srv *gin.Engine) {
	srv.GET("/:short_url", h.Redirect)
}

func (h *ServerHandler) Redirect(ctx *gin.Context) {
	shortUrl := ctx.Param("short_url")
	if ok := h.svc.CheckShortUrl(ctx.Request.Context(), shortUrl); !ok {
		ctx.JSON(404, gin.H{"error": "Short URL not found"})
		return
	}
	originUrl, err := h.svc.Redirect(ctx.Request.Context(), shortUrl)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Redirect(301, originUrl)
}
