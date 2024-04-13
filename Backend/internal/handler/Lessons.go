package handler

import (
	dto "classrooms/internal/DTO"
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

func (h *Handler) GetLessonsByRoomAndDate(c *gin.Context) {
	roomIdStr := c.Query("room")
	date := c.Query("date")

	if roomIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не выбрана аудитория"})
		return
	}

	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неправильный ID аудитории"})
		return
	}

	if date == "" || date == "null" {
		currentTime := time.Now()
		location, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ошибка при загрузке текущей даты"})
			return
		}

		date = currentTime.In(location).Format("02.01.2006")
	}

	scheduleLessons, err := h.services.Schedule.GetScheduleByRoomIdAndDate(roomId, date)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	reservedLessons, err := h.services.Reservation.GetReservationRoom(roomId, date)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	var reservedLessonsDTO []dto.ReservedLessonDto
	for _, lesson := range reservedLessons {
		reserver, _ := h.services.Login.GetUserById(lesson.User_id)
		comment := lesson.Comment
		startTime := lesson.Date.Format("15:04")
		endTime := lesson.Date.Add(90 * time.Minute).Format("15:04")
		reservedLessonsDTO = append(reservedLessonsDTO, dto.ReservedLessonDto{ReserverName: reserver.Fullname, StartTime: startTime,
			EndTime: endTime, Comment: comment, Reserved: true})
	}

	lessons := []ILesson{
		FreeLesson{RoomId: roomId, StartTime: "09:00", EndTime: "10:30", Free: true},
		FreeLesson{RoomId: roomId, StartTime: "10:45", EndTime: "12:15", Free: true},
		FreeLesson{RoomId: roomId, StartTime: "13:00", EndTime: "14:30", Free: true},
		FreeLesson{RoomId: roomId, StartTime: "14:45", EndTime: "16:15", Free: true},
		FreeLesson{RoomId: roomId, StartTime: "16:30", EndTime: "18:00", Free: true},
		FreeLesson{RoomId: roomId, StartTime: "18:15", EndTime: "19:45", Free: true},
		FreeLesson{RoomId: roomId, StartTime: "20:00", EndTime: "21:30", Free: true},
	}

	lessonsStartTime := map[string]int{
		"09:00": 0,
		"10:45": 1,
		"13:00": 2,
		"14:45": 3,
		"16:30": 4,
		"18:15": 5,
		"20:00": 6,
	}

	for _, lesson := range reservedLessonsDTO {
		lessons[lessonsStartTime[lesson.StartTime]] = lesson
	}

	for _, lesson := range scheduleLessons {
		lessons[lessonsStartTime[lesson.StartTime]] = lesson
	}

	c.JSON(http.StatusOK, lessons)

}
