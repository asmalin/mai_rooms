package handler

import (
	dto "classrooms/internal/DTO"
	"classrooms/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type lessonForReservationJSON_WEB struct {
	RoomId    string `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Comment   string `json:"comment"`
}

type lessonForReservationJSON_TG struct {
	ChatId    string `json:"userId"`
	RoomId    string `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type lessonForCancelReservationJSON struct {
	RoomId    string `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
}

func (h *Handler) Reserve(c *gin.Context) {

	var lessonForReservation lessonForReservationJSON_WEB
	if err := json.NewDecoder(c.Request.Body).Decode(&lessonForReservation); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")
	roomId, _ := strconv.Atoi(lessonForReservation.RoomId)
	date := lessonForReservation.Date
	startTime := lessonForReservation.StartTime
	endTime := lessonForReservation.EndTime
	comment := lessonForReservation.Comment

	parsedDate, dateErr := time.Parse("02.01.2006", date)

	if userId == 0 || roomId == 0 || dateErr != nil {
		c.JSON(http.StatusBadRequest, "Invalid request")
	}

	reservedRoom := models.ReservedLesson{User_id: userId, Room_id: roomId, Date: parsedDate, TimeStart: startTime, TimeEnd: endTime, Comment: comment}

	h.services.Reservation.ReserveRoom(reservedRoom)

	c.JSON(http.StatusOK, "")
}

func (h *Handler) TgReserve(c *gin.Context) {

	var lessonForReservation lessonForReservationJSON_TG
	if err := json.NewDecoder(c.Request.Body).Decode(&lessonForReservation); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	chatId, _ := strconv.ParseInt(lessonForReservation.ChatId, 10, 64)
	userId, _ := h.services.Login.UserIdByChatId(chatId)
	roomId, _ := strconv.Atoi(lessonForReservation.RoomId)
	date := lessonForReservation.Date
	parsedDate, dateErr := time.Parse("02.01.2006", date)
	startTime := lessonForReservation.StartTime
	endTime := lessonForReservation.EndTime
	comment := ""

	if userId == 0 || roomId == 0 || dateErr != nil {
		c.JSON(http.StatusBadRequest, "Invalid request")
	}

	reservedRoom := models.ReservedLesson{User_id: userId, Room_id: roomId, Date: parsedDate, TimeStart: startTime, TimeEnd: endTime, Comment: comment}

	err := h.services.Reservation.ReserveRoom(reservedRoom)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, "service problem")
		return
	}

	c.JSON(http.StatusOK, "")
}

func (h *Handler) CancelReservation(c *gin.Context) {
	var lessonForCancelReservation lessonForCancelReservationJSON
	if err := json.NewDecoder(c.Request.Body).Decode(&lessonForCancelReservation); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")
	roomId, _ := strconv.Atoi(lessonForCancelReservation.RoomId)
	date := lessonForCancelReservation.Date
	startTime := lessonForCancelReservation.StartTime

	parsedDate, dateErr := time.Parse("02.01.2006", date)

	if userId == 0 || roomId == 0 || dateErr != nil {
		c.JSON(http.StatusBadRequest, "Invalid request")
	}

	lesson := dto.LessonForCancelReservationDto{ReserverId: userId, Room_id: roomId, Date: parsedDate, StartTime: startTime}

	h.services.Reservation.CancelReservation(lesson)

	c.JSON(http.StatusOK, "canceled!")
}
