package service

import (
	dto "classrooms/internal/DTO"
	"classrooms/internal/models"
	"classrooms/internal/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) GetAllUsers() ([]dto.UserDto, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var usersDto []dto.UserDto
	for _, user := range users {
		usersDto = append(usersDto, dto.UserDto{Id: user.Id, Username: user.Username, Fullname: user.Fullname, Email: user.Email, Role: user.Role})

	}
	return usersDto, nil
}

func (s *UsersService) CreateUser(user models.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(passwordHash)

	return s.repo.CreateUser(user)
}

func (s *UsersService) DeleteUser(userId int) error {
	return s.repo.DeleteUser(userId)
}

func (s *UsersService) UpdateUser(user models.User) error {
	return s.repo.UpdateUser(user)
}

func (s *UsersService) ChangePassword(oldPassword, newPassword string, user models.User) error {

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return fmt.Errorf("Неправильный старый пароль")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(passwordHash)

	return s.repo.UpdateUser(user)
}
