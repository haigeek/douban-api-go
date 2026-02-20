package server

import (
	"github.com/gin-gonic/gin"

	"github.com/haigeek/douban-api-go/internal/book"
)

func NewRouter(h *Handlers, b *book.Handlers, debug bool) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	if h.cfg.BasicUser != "" && h.cfg.BasicPass != "" {
		r.Use(gin.BasicAuth(gin.Accounts{
			h.cfg.BasicUser: h.cfg.BasicPass,
		}))
	}

	r.GET("/", h.Index)
	r.GET("/movies", h.Movies)
	r.GET("/movies/:sid", h.Movie)
	r.GET("/movies/:sid/celebrities", h.Celebrities)
	r.GET("/celebrities/:id", h.Celebrity)
	r.GET("/photo/:sid", h.Photo)
	r.GET("/proxy", h.Proxy)

	r.GET("/v2/book/search", b.Search)
	r.GET("/v2/book/id/:sid", b.ByID)
	r.GET("/v2/book/isbn/:isbn", b.ByISBN)

	return r
}
