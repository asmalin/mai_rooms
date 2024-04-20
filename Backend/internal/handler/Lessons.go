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

	// lessons := []ILesson{
	// 	FreeLesson{RoomId: roomId, StartTime: "09:00", EndTime: "10:30", Free: true},
	// 	FreeLesson{RoomId: roomId, StartTime: "10:45", EndTime: "12:15", Free: true},
	// 	FreeLesson{RoomId: roomId, StartTime: "13:00", EndTime: "14:30", Free: true},
	// 	FreeLesson{RoomId: roomId, StartTime: "14:45", EndTime: "16:15", Free: true},
	// 	FreeLesson{RoomId: roomId, StartTime: "16:30", EndTime: "18:00", Free: true},
	// 	FreeLesson{RoomId: roomId, StartTime: "18:15", EndTime: "19:45", Free: true},
	// 	FreeLesson{RoomId: roomId, StartTime: "20:00", EndTime: "21:30", Free: true},
	// }

	// lessonsStartTime := map[string]int{
	// 	"09:00": 0,
	// 	"10:45": 1,
	// 	"13:00": 2,
	// 	"14:45": 3,
	// 	"16:30": 4,
	// 	"18:15": 5,
	// 	"20:00": 6,
	// }

	// for _, lesson := range reservedLessonsDTO {
	// 	lessons[lessonsStartTime[lesson.StartTime]] = lesson
	// }

	// for _, lesson := range scheduleLessons {
	// 	lessons[lessonsStartTime[lesson.StartTime]] = lesson
	// }

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
