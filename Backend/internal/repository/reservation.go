package repository

import (
	"classrooms/internal/models"
	"errors"

	"gorm.io/gorm"
)

type ReservationPostgres struct {
	db *gorm.DB
}

func NewReservationPostgres(db *gorm.DB) *ReservationPostgres {
	return &ReservationPostgres{db: db}
}

func (r *ReservationPostgres) ReserveRoom(reservedRoom models.ReservedLesson) error {

	result := r.db.Create(&reservedRoom)
	if result.Error != nil {
		return errors.New("insert error")
	}

	return nil
}

func (r *ReservationPostgres) CancelReservation(reservedLesson_id int) error {

	result := r.db.Delete(models.ReservedLesson{}, "id = ?", reservedLesson_id)
	if result.Error != nil {
		return errors.New("невозможно удалить эту запись")
	}
	return nil
}

func (r *ReservationPostgres) GetAllReservedLessons() ([]models.ReservedLesson, error) {

	var reservedLessons []models.ReservedLesson
	result := r.db.Preload("User").Preload("Room").Find(&reservedLessons)
	if result.Error != nil {
		return nil, errors.New("ошибка при получении данных")
	}

	return reservedLessons, nil
}

func (r *ReservationPostgres) GetReservedLessons(roomId int, date string) ([]models.ReservedLesson, error) {

	var reservedLessons []models.ReservedLesson
	result := r.db.Where("room_id = ? AND date = ?", roomId, date).Find(&reservedLessons)
	if result.Error != nil {
		return nil, errors.New("ошибка при получении данных")
	}

	return reservedLessons, nil
}

func (r *ReservationPostgres) GetReservedLesson(roomId int, date string, startTime string) (models.ReservedLesson, error) {
	var reservedLessons models.ReservedLesson
	result := r.db.Where("room_id = ? AND date = ? AND time_start = ?", roomId, date, startTime).Find(&reservedLessons)
	if result.Error != nil {
		return models.ReservedLesson{}, errors.New("ошибка при получении данных")
	}

	return reservedLessons, nil
}
