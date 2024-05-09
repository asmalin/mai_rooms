package repository

import (
	"classrooms/internal/models"
	"errors"

	"gorm.io/gorm"
)

type UsersPostgres struct {
	db *gorm.DB
}

func NewUsersPostgres(db *gorm.DB) *RoomPostgres {
	return &RoomPostgres{db: db}
}

func (r *RoomPostgres) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if result := r.db.Find(&users); result.Error != nil {
		return nil, errors.New("ошибка при получении данных")
	}

	return users, nil
}

func (r *RoomPostgres) CreateUser(user models.User) (err error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *RoomPostgres) DeleteUser(userId int) error {

	result := r.db.Delete(models.User{}, "id = ?", userId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RoomPostgres) UpdateUser(user models.User) error {
	err := r.db.Save(&user).Error

	if err != nil {
		return err
	}
	return nil
}
