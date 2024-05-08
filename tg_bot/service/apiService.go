package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"tg_bot/models"

	"github.com/dgrijalva/jwt-go"
)

type LoginInput struct {
	Username string
	Password string
	ChatId   int64
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId   int   `json:"user_id"`
	TgChatId int64 `json:"tgChat_id"`
}

func Buildings() []models.Building {
	resp, err := http.Get("/api/buildings")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var Buildings []models.Building
	if err := json.NewDecoder(resp.Body).Decode(&Buildings); err != nil {
		log.Fatal(err)
	}

	return Buildings
}

func Rooms(building_id int) []models.Room {
	resp, err := http.Get("/api/rooms/" + fmt.Sprint(building_id))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var rooms []models.Room
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		log.Fatal(err)
	}

	return rooms
}

func Schedule(room_id string, date string) []models.ScheduleLesson {
	resp, err := http.Get(fmt.Sprintf("/api/schedule?room=%s&date=%s", room_id, date))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var scheduleLessons []models.ScheduleLesson
	if err := json.NewDecoder(resp.Body).Decode(&scheduleLessons); err != nil {
		log.Fatal(err)
	}

	return scheduleLessons
}

func ReservedLessons(room_id string, date string) []models.ReservedLesson {
	resp, err := http.Get(fmt.Sprintf("/api/reserved_lesssons?room=%s&date=%s", room_id, date))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var reservedLessons []models.ReservedLesson
	if err := json.NewDecoder(resp.Body).Decode(&reservedLessons); err != nil {
		log.Fatal(err)
	}

	return reservedLessons
}

func Auth(username string, password string, chatId int64) (err error) {

	jsonData, _ := json.Marshal(LoginInput{Username: username, Password: password, ChatId: chatId})
	requestBody := bytes.NewBuffer(jsonData)

	resp, err := http.Post("/api/auth/tg_login", "application/json", requestBody)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("неправильный логин или пароль")
	}

	return nil
}

func Reserve(chatId int64, lesson models.LessonForReservationJSON) error {

	jsonData, _ := json.Marshal(lesson)
	requestBody := bytes.NewBuffer(jsonData)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "/api/reserve", requestBody)
	if err != nil {
		log.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
		0,
		chatId,
	})

	tokenStr, _ := token.SignedString([]byte(os.Getenv("Secret_key")))

	req.Header.Set("Authorization", "Bearer "+tokenStr)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return errors.New(resp.Status)
	}

	return nil
}
