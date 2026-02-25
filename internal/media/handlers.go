package media

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 20
	maxLimit     = 50
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) HotTV(c *gin.Context) {
	start, limit, ok := parsePageParams(c)
	if !ok {
		return
	}
	result, err := h.service.HotTV(c.Request.Context(), start, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) HotMovie(c *gin.Context) {
	start, limit, ok := parsePageParams(c)
	if !ok {
		return
	}
	result, err := h.service.HotMovie(c.Request.Context(), start, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) LatestMovie(c *gin.Context) {
	start, limit, ok := parsePageParams(c)
	if !ok {
		return
	}
	result, err := h.service.LatestMovie(c.Request.Context(), start, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) HighRatingMovie(c *gin.Context) {
	start, limit, ok := parsePageParams(c)
	if !ok {
		return
	}
	result, err := h.service.HighRatingMovie(c.Request.Context(), start, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func parsePageParams(c *gin.Context) (start int, limit int, ok bool) {
	start = 0
	limit = defaultLimit

	if raw, exists := c.GetQuery("start"); exists {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid start"})
			return 0, 0, false
		}
		start = v
	}

	if raw, exists := c.GetQuery("limit"); exists {
		v, err := strconv.Atoi(raw)
		if err != nil || v <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid limit"})
			return 0, 0, false
		}
		if v > maxLimit {
			c.JSON(http.StatusBadRequest, gin.H{"message": "limit不能大于50"})
			return 0, 0, false
		}
		limit = v
	}

	return start, limit, true
}
