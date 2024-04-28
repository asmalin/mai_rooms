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
	UserId   int   `json:"user_id"`
	TgChatId int64 `json:"tgChat_id"`
}

func (s *LoginService) TgLogin(username, password string, tgChatId int64) (err error) {
	user, err := s.repo.GetUser(username, password)
	if err != nil {
		return err
	}

	err = s.repo.AddUserTgChatRelation(user.Id, tgChatId)

	if err != nil {
		return err
	}

	return nil
}

func (s *LoginService) UserIdByChatId(tgChatId int64) (userId int, err error) {

	userChatRelation, err := s.repo.GetUserChatRelation(tgChatId)

	if err != nil {
		return 0, err
	}

	return userChatRelation.UserId, nil
}

func (s *LoginService) GenerateToken(username, password string) (string, error) {

	user, err := s.repo.GetUser(username, password)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(30 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		0,
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

	if claims.UserId != 0 {
		return claims.UserId, nil
	}

	if claims.TgChatId != 0 {
		userTgChatRelation, err := s.repo.GetUserChatRelation(claims.TgChatId)
		if err != nil {
			return 0, err
		}
		return userTgChatRelation.UserId, nil
	}

	return 0, errors.New("неизвестная ошибка")
}

func (s *LoginService) GetUserById(id int) (models.User, error) {
	user, err := s.repo.GetUserById(id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
