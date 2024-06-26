package service

import (
	dto "classrooms/internal/DTO"
	"classrooms/internal/models"
	"classrooms/internal/repository"
	"errors"
	"time"
)

type ReservationService struct {
	repo repository.Reservation
}

func NewReservationService(repo repository.Reservation) *ReservationService {
	return &ReservationService{repo: repo}
}

func (r *ReservationService) ReserveRoom(reservedRoom models.ReservedLesson) error {

	r.repo.ReserveRoom(reservedRoom)
	return nil
}

func (r *ReservationService) CancelReservation(lessonForCancelReservation dto.LessonForCancelReservationDto) error {
	user := lessonForCancelReservation.User
	roomId := lessonForCancelReservation.Room_id

	startTime := lessonForCancelReservation.StartTime

	lesson, err := r.repo.GetReservedLesson(roomId, lessonForCancelReservation.Date, startTime)
	if err != nil {
		return err
	}

	if user.Id == lesson.User_id || user.Role == "admin" {

		err = r.repo.CancelReservation(lesson.ID)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("недостаточно прав")

}

func (r *ReservationService) GetReservationRoom(roomId int, date string) ([]models.ReservedLesson, error) {

	parsedDate, err := time.Parse("02.01.2006", date)

	if err != nil {
		return nil, err
	}

	result, err := r.repo.GetReservedLessons(roomId, parsedDate)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ReservationService) GetAllReservedLessons() ([]models.ReservedLesson, error) {
	result, err := r.repo.GetAllReservedLessons()
	if err != nil {
		return nil, err
	}
	return result, nil
}
