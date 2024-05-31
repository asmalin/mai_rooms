package service

import (
	"classrooms/internal/models"
	"classrooms/internal/repository"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type LoginService struct {
	repo repository.Login
}

func NewLoginService(repo repository.Login) *LoginService {
	return &LoginService{repo: repo}
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (s *LoginService) WebAuth(username, password string) (models.User, error) {
	user, err := s.repo.GetUser(username, password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *LoginService) GenerateAccessToken(user models.User) (string, error) {

	expirationTime := time.Now().Add(30 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(os.Getenv("Secret_key")))
}

func (s *LoginService) GenerateRefreshToken(userId int) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour * 60)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})

	return token.SignedString([]byte(os.Getenv("Secret_key")))
}

func (s *LoginService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("Secret_key")), nil
	})

	if err != nil {
		return 0, errors.New("wrong token")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil

}

func (s *LoginService) GetUserById(id int) (models.User, error) {
	user, err := s.repo.GetUserById(id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *LoginService) GetUserIdByTgUsername(tgUsername string) (int, error) {
	user, err := s.repo.GetUserByTgUsername(tgUsername)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}
