package handler

import (
	"classrooms/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type lessonForReservationJSON struct {
	RoomId    string `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	Comment   string `json:"comment"`
}

func (h *Handler) Reserve(c *gin.Context) {

	var lessonForReservation lessonForReservationJSON
	if err := json.NewDecoder(c.Request.Body).Decode(&lessonForReservation); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)

		return
	}
	fmt.Println(lessonForReservation.RoomId)
	fmt.Println(lessonForReservation.Date)
	fmt.Println(lessonForReservation.StartTime)
	fmt.Println(lessonForReservation.Comment)

	userId := c.GetInt("userId")
	roomId, _ := strconv.Atoi(lessonForReservation.RoomId)
	date := lessonForReservation.Date
	startTime := lessonForReservation.StartTime
	comment := lessonForReservation.Comment

	parsedDate, dateErr := time.Parse("02.01.2006 15:04", date+" "+startTime)

	if userId == 0 || roomId == 0 || dateErr != nil {
		c.JSON(http.StatusBadRequest, "Invalid request")
	}

	//parsedTime, _ := time.Parse("15:04", startTime)

	reservedRoom := models.ReservedLesson{User_id: userId, Room_id: roomId, Date: parsedDate, Comment: comment}
	h.services.Reservation.ReserveRoom(reservedRoom)

	c.JSON(http.StatusOK, "reserved!")
}

func (h *Handler) CancelReservation(c *gin.Context) {
	// Ваш код обработки запроса для /api/cancelReservation
}

func (h *Handler) GetReservedRooms(c *gin.Context) {
	// Ваш код обработки запроса для /api/reservedRooms
}
