package handler

import (
	"classrooms/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.services.Users.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) DeleteUserById(c *gin.Context) {
	userIdStr := c.Param("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.services.Users.DeleteUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]int{
		"id": userId})

}

func (h *Handler) CreateUser(c *gin.Context) {

	var user models.User
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "json decoding error")
		return
	}

	err := h.services.Users.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "User created!")
}

func (h *Handler) UpdateUser(c *gin.Context) {
	updatedUserIdStr := c.Param("userId")

	updatedUserId, err := strconv.Atoi(updatedUserIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedUser models.User
	err = json.NewDecoder(c.Request.Body).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	existingUser, err := h.services.Login.GetUserById(updatedUserId)

	if err != nil {
		c.JSON(http.StatusBadRequest, "пользователь не найден")
		return
	}

	userId := c.GetInt("userId")

	user, err := h.services.Login.GetUserById(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, "пользователь не найден")
		return
	}

	if userId != existingUser.Id && user.Role != "admin" {
		c.JSON(http.StatusBadRequest, "недостаточно прав")
		return
	}

	if updatedUser.Username != "" {
		existingUser.Username = updatedUser.Username
	}
	if updatedUser.Fullname != "" {
		existingUser.Fullname = updatedUser.Fullname
	}
	if updatedUser.Email != "" {
		existingUser.Email = updatedUser.Email
	}
	if updatedUser.Role != "" {
		existingUser.Role = updatedUser.Role
	}

	err = h.services.Users.UpdateUser(existingUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, "не удалось обновить")
		return
	}

	c.JSON(http.StatusOK, map[string]int{
		"id": userId})
}

type passwords struct {
	OldPassword string
	NewPassword string
}

func (h *Handler) ChangeUserPassword(c *gin.Context) {
	var pass passwords
	err := json.NewDecoder(c.Request.Body).Decode(&pass)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userId := c.GetInt("userId")

	user, err := h.services.Login.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = h.services.Users.ChangePassword(pass.OldPassword, pass.NewPassword, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "пароль успешно изменен!")
}
