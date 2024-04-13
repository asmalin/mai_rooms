package service

import (
	"classrooms/internal/models"
	"classrooms/internal/repository"
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

func (r *ReservationService) GetReservationRoom(roomId int, date string) ([]models.ReservedLesson, error) {
	result, err := r.repo.GetReservedLessons(roomId, date)
	if err != nil {
		return nil, err
	}
	return result, nil
}
