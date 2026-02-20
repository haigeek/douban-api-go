package book

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusOK, []any{})
		return
	}

	count := 2
	if raw, ok := c.GetQuery("count"); ok {
		v, err := strconv.Atoi(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid count"})
			return
		}
		count = v
	}

	if count > 20 {
		c.Data(http.StatusBadRequest, "application/json; charset=utf-8", []byte(`{"message":"count不能大于20"}`))
		return
	}

	result, err := h.service.Search(c.Request.Context(), q, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handlers) ByID(c *gin.Context) {
	sid := c.Param("sid")
	info, err := h.service.GetBookInfo(c.Request.Context(), sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handlers) ByISBN(c *gin.Context) {
	isbn := c.Param("isbn")
	info, err := h.service.GetBookInfoByISBN(c.Request.Context(), isbn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}
