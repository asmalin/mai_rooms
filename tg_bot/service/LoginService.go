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

	"github.com/dgrijalva/jwt-go"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type UserData struct {
	ID int `json:"userId"`
}

func GetUserIdByTgUsername(tgUsername string) (int, error) {

	if len(tgUsername) == 0 {
		return 0, errors.New("empty tg username")
	}

	jsonData, _ := json.Marshal(map[string]string{
		"tgUsername": tgUsername,
	})
	requestBody := bytes.NewBuffer(jsonData)

	resp, err := http.Post(os.Getenv("BASE_URL")+"/api/auth/tg/getUserId", "application/json", requestBody)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var userData UserData

		if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
			log.Fatal(err)
		}

		return userData.ID, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	return 0, errors.New(string(bodyBytes))

}
