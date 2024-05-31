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

func Buildings() []models.Building {
	resp, err := http.Get(os.Getenv("BASE_URL") + "/api/buildings")
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
	resp, err := http.Get(os.Getenv("BASE_URL") + "/api/rooms/" + fmt.Sprint(building_id))
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
	resp, err := http.Get(fmt.Sprintf(os.Getenv("BASE_URL")+"/api/schedule?room=%s&date=%s", room_id, date))
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
	resp, err := http.Get(fmt.Sprintf(os.Getenv("BASE_URL")+"/api/reserved_lesssons?room=%s&date=%s", room_id, date))
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

func Reserve(tgUsername string, lesson models.LessonForReservationJSON) error {

	userId, err := GetUserIdByTgUsername(tgUsername)

	if err != nil {
		return err
	}

	jsonData, _ := json.Marshal(lesson)
	requestBody := bytes.NewBuffer(jsonData)

	client := &http.Client{}

	req, err := http.NewRequest("POST", os.Getenv("BASE_URL")+"/api/reserve", requestBody)
	if err != nil {
		log.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
		userId,
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

func CancelReserve(tgUsername string, lesson models.LessonForCancelReservationJSON) error {

	userId, err := GetUserIdByTgUsername(tgUsername)

	if err != nil {
		return err
	}

	jsonData, _ := json.Marshal(lesson)
	requestBody := bytes.NewBuffer(jsonData)

	client := &http.Client{}

	req, err := http.NewRequest("POST", os.Getenv("BASE_URL")+"/api/cancelReservation", requestBody)
	if err != nil {
		log.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
		userId,
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
