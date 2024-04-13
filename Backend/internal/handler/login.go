package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {

	var input loginInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.services.Login.GenerateToken(input.Username, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) checkAuth(c *gin.Context) {

	user, err := h.services.GetUserById(c.GetInt("userId"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		//c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}
