package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetQRCodes(c *gin.Context) {
	id, _ := c.Get("userId")
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
