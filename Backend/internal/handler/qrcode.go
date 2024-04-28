package handler

import (
	"archive/zip"
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func (h *Handler) GetQRCodes(c *gin.Context) {

	roomIds := c.QueryArray("roomId")

	// Создание буфера для хранения данных архива
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Базовый URL для создания QR кодов
	baseUrl := "http://localhost:3000/rooms?room="

	for _, roomId := range roomIds {
		// Создание URL для текущего roomId
		url := baseUrl + roomId

		// Создание QR кода
		qrCode, err := qrcode.New(url, qrcode.Medium)
		if err != nil {
			log.Fatal(err)
		}

		// Имя файла PNG в архиве
		fileName := "qrcode_room_" + roomId + ".png"

		// Добавление PNG файла в архив
		file, err := zipWriter.Create(fileName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка создания файла в архиве")
			return
		}

		// Запись данных QR кода в файл
		err = qrCode.Write(512, file)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Закрытие архива
	err := zipWriter.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка закрытия архива")
		return
	}

	err = os.WriteFile("qrcodes.zip", buf.Bytes(), 0666)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка сохранения архива на диск")
		return
	}

	// Отправка архива клиенту
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=qrcodes.zip")
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}
