package service

import (
	dto "classrooms/internal/DTO"
	"classrooms/internal/models"
	"classrooms/internal/repository"
	"strings"
	"time"
)

type ScheduleService struct {
	RoomRepo   repository.Room
	LessonRepo repository.Lesson
}

func NewScheduleService(roomRepo repository.Room, lessonRepo repository.Lesson) *ScheduleService {
	return &ScheduleService{RoomRepo: roomRepo, LessonRepo: lessonRepo}
}

type ScheduleItem struct {
	Time_start string
	Time_end   string
	Name       string
	Groups     []string
	Types      []string
	Rooms      map[string]string
}

type ScheduleDate struct {
	Date  string
	Day   string
	Pairs map[string]ScheduleItem
}

type LectorSchedule struct {
	Name     string
	Groups   map[string]int
	Schedule map[string]ScheduleDate
}

func (s *ScheduleService) GetScheduleByRoomIdAndDate(roomId int, date string) ([]dto.ScheduleLessonDto, error) {

	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return nil, err
	}

	lessons, err := s.LessonRepo.GetScheduleLessons(roomId, parsedDate)

	if err != nil {

		return nil, err
	}

	lessonsDTO := make([]dto.ScheduleLessonDto, len(lessons))

	for index, lesson := range lessons {
		lessonsDTO[index].Lector = lesson.Lector
		lessonsDTO[index].Groups = lesson.Groups
		lessonsDTO[index].Subject = lesson.Name
		lessonsDTO[index].Type = lesson.Type
		lessonsDTO[index].StartTime = lesson.TimeStart[:5]
		lessonsDTO[index].EndTime = lesson.TimeEnd[:5]
	}
	return lessonsDTO, nil
}

func (s *ScheduleService) WriteScheduleToDB(lectorSchedule LectorSchedule) error {

	lectorName := lectorSchedule.Name

	for _, schedule := range lectorSchedule.Schedule {
		for _, pair := range schedule.Pairs {
			subject := pair.Name
			groups := strings.Join(pair.Groups, ", ")
			types := pair.Types
			for _, roomName := range pair.Rooms {

				if !strings.Contains(roomName, "-") || roomName[0] == '-' {
					continue
				}

				roomId, _ := s.RoomRepo.GetRoomId(roomName)

				lesson := models.Lesson{Lector: lectorName, TimeStart: "startTime",
					TimeEnd: "endTime", Name: subject, Groups: groups, Type: types[0], RoomId: roomId}

				err := s.LessonRepo.InsertLessonToDB(lesson)
				if err != nil {
					return err
				}

			}

		}
	}

	return nil
}
