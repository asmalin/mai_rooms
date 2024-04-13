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

	r.db.AutoMigrate(models.ReservedLesson{})
	result := r.db.Create(&reservedRoom)
	if result.Error != nil {
		return errors.New("insert error")
	}

	return nil
}

func (r *ReservationPostgres) GetReservedLessons(roomId int, date string) ([]models.ReservedLesson, error) {
	var reservedLessons []models.ReservedLesson
	result := r.db.Where("room_id = ? AND DATE(date) = ?", roomId, date).Find(&reservedLessons)
	if result.Error != nil {
		return nil, errors.New("ошибка при получении данных")
	}

	return reservedLessons, nil
}
