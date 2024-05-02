package handler

import (
	"classrooms/internal/models"
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

	user, err := h.services.Login.WebAuth(input.Username, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 60*60*24*60, "/", "", true, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": accessToken,
		"user":  user,
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

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("refresh_token", "", 1, "/", "", true, true)

	c.JSON(http.StatusOK, "User logged out successfully")
}

func (h *Handler) AuthRefresh(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusBadRequest, "куки не найдены")
		return
	}

	userId, err := h.services.Login.ParseToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	user, err := h.services.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 60*60*24*60, "/", "localhost", true, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": accessToken,
	})
}

func (h *Handler) generateTokens(user models.User) (accessToken, refreshToken string, err error) {
	accToken, err := h.services.Login.GenerateAccessToken(user)

	if err != nil {
		return "", "", err
	}

	refToken, err := h.services.Login.GenerateRefreshToken(user.Id)

	if err != nil {
		return "", "", err
	}
	return accToken, refToken, nil

}

func (h *Handler) TgLogin(c *gin.Context) {

	var input tgloginInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input.ChatId)
	err := h.services.Login.TgAuth(input.Username, input.Password, input.ChatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "")
}
