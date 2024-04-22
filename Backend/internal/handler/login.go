package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tgloginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ChatId   int64  `json:"chatId"`
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
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) TgLogin(c *gin.Context) {

	var input tgloginInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input.ChatId)
	err := h.services.Login.TgLogin(input.Username, input.Password, input.ChatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "")
}
