package handler

import (
	dto "classrooms/internal/DTO"
	"classrooms/internal/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type lessonForReservationJSON struct {
	RoomId    int    `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Comment   string `json:"comment"`
}

type lessonForCancelReservationJSON struct {
	RoomId    int    `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
}

func (h *Handler) Reserve(c *gin.Context) {

	var lessonForReservation lessonForReservationJSON
	if err := json.NewDecoder(c.Request.Body).Decode(&lessonForReservation); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")
	roomId := lessonForReservation.RoomId
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

func (h *Handler) CancelReservation(c *gin.Context) {

	var lessonForCancelReservation lessonForCancelReservationJSON
	if err := json.NewDecoder(c.Request.Body).Decode(&lessonForCancelReservation); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")
	roomId := lessonForCancelReservation.RoomId
	date := lessonForCancelReservation.Date
	startTime := lessonForCancelReservation.StartTime

	parsedDate, dateErr := time.Parse("02.01.2006", date)

	if userId == 0 || roomId == 0 || dateErr != nil {
		c.JSON(http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := h.services.Login.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	lesson := dto.LessonForCancelReservationDto{User: user, Room_id: roomId, Date: parsedDate, StartTime: startTime}

	err = h.services.Reservation.CancelReservation(lesson)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, "")
}
