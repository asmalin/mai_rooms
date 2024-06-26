package repository

import (
	"classrooms/internal/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginPostgres struct {
	db *gorm.DB
}

func NewLoginPostgres(db *gorm.DB) *LoginPostgres {
	return &LoginPostgres{db: db}
}

func (l *LoginPostgres) GetUser(username, password string) (models.User, error) {

	var user models.User
	if err := l.db.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("wrong username or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, fmt.Errorf("wrong username or password")
	}

	return user, nil
}

func (l *LoginPostgres) GetUserById(id int) (models.User, error) {
	var user models.User
	result := l.db.Where("id = ?", id).Take(&user)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (l *LoginPostgres) GetUserByTgUsername(tgUsername string) (models.User, error) {
	var user models.User
	result := l.db.Where("tg_username = ?", tgUsername).Take(&user)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}
