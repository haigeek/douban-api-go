package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/haigeek/douban-api-go/internal/api/movie"
	"github.com/haigeek/douban-api-go/internal/config"
)

type Handlers struct {
	movie *movie.Service
	cfg   config.Config
}

func NewHandlers(movieService *movie.Service, cfg config.Config) *Handlers {
	return &Handlers{movie: movieService, cfg: cfg}
}

func (h *Handlers) Index(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
       接口列表：<br/>
       /movies?q={movie_name}<br/>
       /movies?q={movie_name}&type=full<br/>
       /movies/{sid}<br/>
       /movies/{sid}/celebrities<br/>
       /celebrities/{cid}<br/>
       /photo/{sid}<br/>
       /v2/book/search?q={book_name}<br/>
       /v2/book/id/{sid}<br/>
       /v2/book/isbn/{isbn}<br/>
       /v2/media/hot/tv?start=0&limit=20<br/>
       /v2/media/hot/movie?start=0&limit=20<br/>
       /v2/media/latest/movie?start=0&limit=20<br/>
       /v2/media/high-rating/movie?start=0&limit=20<br/>
    `))
}

func (h *Handlers) Movies(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusOK, []movie.Movie{})
		return
	}

	count := 0
	if rawCount, ok := c.GetQuery("count"); ok {
		v, err := strconv.Atoi(rawCount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid count"})
			return
		}
		count = v
	}

	fromJellyfin := c.GetHeader("User-Agent") == ""
	if count == 0 && fromJellyfin {
		count = h.cfg.Limit
	}

	searchType := c.DefaultQuery("type", "")
	imageSize := c.DefaultQuery("s", "")
	if searchType == "full" {
		result, err := h.movie.SearchFull(c.Request.Context(), q, count, imageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	result, err := h.movie.Search(c.Request.Context(), q, count, imageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) Movie(c *gin.Context) {
	sid := c.Param("sid")
	imageSize := c.DefaultQuery("s", "")
	result, err := h.movie.GetMovieInfo(c.Request.Context(), sid, imageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) Celebrities(c *gin.Context) {
	sid := c.Param("sid")
	result, err := h.movie.GetCelebrities(c.Request.Context(), sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) Celebrity(c *gin.Context) {
	id := c.Param("id")
	result, err := h.movie.GetCelebrity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) Photo(c *gin.Context) {
	sid := c.Param("sid")
	result, err := h.movie.GetWallpaper(c.Request.Context(), sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) Proxy(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "url is required"})
		return
	}

	resp, body, err := h.movie.ProxyImage(c.Request.Context(), rawURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Data(resp.StatusCode, contentType, body)
}
