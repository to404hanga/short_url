package routes

import (
	"short_url/pkg/generator"
	short_url_v1 "short_url/proto/short_url/v1"

	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	svc     short_url_v1.ShortUrlServiceClient
	weights []int
}

var _ Handler = (*ServerHandler)(nil)

func NewServerHandler(svc short_url_v1.ShortUrlServiceClient) *ServerHandler {
	return &ServerHandler{
		svc:     svc,
		weights: []int{5, 67, 23, 71, 73, 79},
	}
}

func (h *ServerHandler) RegisterRoutes(srv *gin.Engine) {
	srv.GET("/:short_url", h.Redirect)
}

func (h *ServerHandler) Redirect(ctx *gin.Context) {
	shortUrl := ctx.Param("short_url")
	if ok := generator.CheckShortUrl(shortUrl, h.weights); !ok {
		ctx.JSON(404, gin.H{"error": "Short URL not found"})
		return
	}
	resp, err := h.svc.GetOriginUrl(ctx, &short_url_v1.GetOriginUrlRequest{
		ShortUrl: shortUrl,
	})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Redirect(301, resp.GetOriginUrl())
}
