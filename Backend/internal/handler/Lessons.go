package handler

import (
	dto "classrooms/internal/DTO"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ILesson interface {
}

type LessonsTime struct {
	Start string
	End   string
}

type FreeLesson struct {
	RoomId    int    `json:"roomId"`
	StartTime string `json:"time_start"`
	EndTime   string `json:"time_end"`
	Free      bool   `json:"free"`
}

func (h *Handler) GetScheduleByRoomAndDate(c *gin.Context) {
	roomId, err := ParseRoomId(c.Query("room"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	date, err := ParseDate(c.Query("date"))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	}

	scheduleLessons, err := h.services.Schedule.GetScheduleByRoomIdAndDate(roomId, date)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scheduleLessons)

}

func (h *Handler) GetReservedLessonsByRoomAndDate(c *gin.Context) {
	roomId, err := ParseRoomId(c.Query("room"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	date, err := ParseDate(c.Query("date"))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	}

	reservedLessons, err := h.services.Reservation.GetReservationRoom(roomId, date)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	var reservedLessonsDTO []dto.ReservedLessonDto
	for _, lesson := range reservedLessons {
		reserver, _ := h.services.Login.GetUserById(lesson.User_id)
		date := lesson.Date.Format("02.01.2006")
		startTime := lesson.TimeStart[:5]
		endTime := lesson.TimeEnd[:5]
		comment := lesson.Comment
		reservedLessonsDTO = append(reservedLessonsDTO, dto.ReservedLessonDto{ReserverName: reserver.Fullname, ReserverId: reserver.Id, Date: date, StartTime: startTime,
			EndTime: endTime, Comment: comment})
	}
	c.JSON(http.StatusOK, reservedLessonsDTO)
}

func ParseRoomId(roomIdStr string) (int, error) {
	if roomIdStr == "" {
		return 0, errors.New("аудитория не выбрана")
	}

	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		return 0, errors.New("неправильный ID аудитории")
	}
	return roomId, nil

}

func ParseDate(date string) (string, error) {
	if date == "" || date == "null" {
		currentTime := time.Now()
		location, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return "", errors.New("ошибка при автоматической загрузке текущей даты")
		}

		date = currentTime.In(location).Format("02.01.2006")
	}
	return date, nil
}

func (h *Handler) GetAllReservedLessons(c *gin.Context) {
	lessons, err := h.services.Reservation.GetAllReservedLessons()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, err)
		return
	}

	var reservedLessonsDTO []dto.ReservedLessonDto
	for _, lesson := range lessons {
		dateStr := lesson.Date.Format("02.01.2006")
		startTime := lesson.TimeStart[:5]
		endTime := lesson.TimeEnd[:5]
		comment := lesson.Comment

		reservedLessonsDTO = append(reservedLessonsDTO, dto.ReservedLessonDto{
			ReserverName: lesson.User.Fullname,
			ReserverId:   lesson.User_id,
			RoomName:     lesson.Room.Name,
			RoomId:       lesson.Room.ID,
			Date:         dateStr,
			StartTime:    startTime,
			EndTime:      endTime,
			Comment:      comment})
	}

	c.JSON(http.StatusOK, reservedLessonsDTO)
}
