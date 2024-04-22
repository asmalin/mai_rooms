package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"tg_bot/models"
	"time"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type ButtonData struct {
	ButtonName   string
	CallBackData string
}

type ButtonsPrefixes struct {
	Building string
	Room     string
	Lesson   string
	Reserve  string
}

type LoginInput struct {
	Username string
	Password string
	ChatId   int64
}

type lessonForReservationJSON struct {
	ChatId    string `json:"userId"`
	RoomId    string `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := telego.NewBot(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	defer bot.StopLongPolling()
	fmt.Println("Started!")

	var buttonsPrefixes = ButtonsPrefixes{
		Building: "buildingId",
		Room:     "roomId",
		Lesson:   "lesson",
		Reserve:  "reserve",
	}

	if err != nil {
		log.Fatalf("Failed to init DB: %s", err.Error())
	}

	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		chatID := tu.ID(update.Message.Chat.ID)

		bot.SendMessage(tu.Message(
			chatID,
			"Здесь можно посмотреть расписание занятий на аудиторию и занять свободную.",
		).WithReplyMarkup(
			tu.Keyboard(
				tu.KeyboardRow(
					tu.KeyboardButton("home"),
				),
			).WithResizeKeyboard(),
		))

		ShowBuildingsPage(chatID, bot, buttonsPrefixes.Building)

	}, th.CommandEqual("start"))

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		queryData := strings.Split(query.Data, "|")
		btnPrefix := queryData[0]

		if btnPrefix == buttonsPrefixes.Building {
			buildingName := queryData[1]
			buidldingId := queryData[2]
			rooms := Rooms(buidldingId)

			buttonsData := make([]ButtonData, len(rooms))
			for index, room := range rooms {
				buttonsData[index] = ButtonData{
					ButtonName:   room.Name,
					CallBackData: buttonsPrefixes.Room + "|" + room.Name + "|" + fmt.Sprint(room.ID) + "|" + buildingName + "|" + buidldingId}
			}

			buttonsInRow := 3

			_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID))

			chatId := telego.ChatID{ID: query.Message.GetChat().ID}

			_, err := bot.EditMessageText(&telego.EditMessageTextParams{
				ChatID:      chatId,
				MessageID:   query.Message.GetMessageID(),
				Text:        "Корпус: " + buildingName + "\nВыберете аудиторию:",
				ReplyMarkup: CreateKeyboard(buttonsData, buttonsInRow),
			})
			if err != nil {
				fmt.Println("Ошибка при редактировании сообщения:", err)
			}

		} else if btnPrefix == buttonsPrefixes.Room {
			roomName := queryData[1]
			roomId := queryData[2]

			buildingName := queryData[3]
			buidldingId := queryData[4]

			currentTime := time.Now()
			location, _ := time.LoadLocation("Europe/Moscow")
			date := currentTime.In(location).Format("02.01.2006")

			schedule := Schedule(roomId, date)
			reserved := ReservedLessons(roomId, date)

			lessons_timeStart := map[string]int{
				"09:00": 1,
				"10:45": 2,
				"13:00": 3,
				"14:45": 4,
				"16:30": 5,
				"18:15": 6,
				"20:00": 7,
			}

			lessons_timeEnd := map[int]string{
				1: "10:30",
				2: "12:15",
				3: "14:30",
				4: "16:15",
				5: "18:00",
				6: "19:45",
				7: "21:30",
			}

			buttonsData := make([]ButtonData, len(lessons_timeStart))

			for _, lesson := range schedule {
				buttonsData[lessons_timeStart[lesson.StartTime]-1] = ButtonData{
					ButtonName:   fmt.Sprint(lessons_timeStart[lesson.StartTime]) + ". " + lesson.StartTime + " - " + lesson.EndTime + " По расписанию 📅",
					CallBackData: buttonsPrefixes.Lesson + "|" + roomName + "|" + roomId + "|" + lesson.StartTime + "|" + lesson.EndTime + "|" + date + "|" + buildingName + "|" + buidldingId}
			}

			for _, lesson := range reserved {
				buttonsData[lessons_timeStart[lesson.StartTime]-1] = ButtonData{
					ButtonName:   fmt.Sprint(lessons_timeStart[lesson.StartTime]) + ". " + lesson.StartTime + " - " + lesson.EndTime + " Забронировано 🔶",
					CallBackData: buttonsPrefixes.Lesson + "|" + roomName + "|" + roomId + "|" + lesson.StartTime + "|" + lesson.EndTime + "|" + date + "|" + buildingName + "|" + buidldingId}
			}

			for timeStart, lessonNumber := range lessons_timeStart {

				if buttonsData[lessonNumber-1].ButtonName == "" {
					buttonsData[lessonNumber-1] = ButtonData{
						ButtonName:   fmt.Sprint(lessonNumber) + ". " + timeStart + " - " + lessons_timeEnd[lessonNumber] + " Свободно ✅",
						CallBackData: buttonsPrefixes.Lesson + "|" + roomName + "|" + roomId + "|" + timeStart + "|" + lessons_timeEnd[lessonNumber] + "|" + date + "|" + buildingName + "|" + buidldingId}
				}
			}

			buttonsInRow := 1

			_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID))

			keyboard := CreateKeyboard(buttonsData, buttonsInRow)
			prevBtn := []telego.InlineKeyboardButton{tu.InlineKeyboardButton("Назад").WithCallbackData(buttonsPrefixes.Building + "|" + buildingName + "|" + buidldingId)}
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, prevBtn)

			chatId := telego.ChatID{ID: query.Message.GetChat().ID}
			_, err := bot.EditMessageText(&telego.EditMessageTextParams{
				ChatID:      chatId,
				MessageID:   query.Message.GetMessageID(),
				Text:        "Аудитория: " + roomName + "\nВыберете занятие:",
				ReplyMarkup: keyboard,
			})
			if err != nil {
				fmt.Println("Ошибка при редактировании сообщения:", err)
			}
		} else if btnPrefix == buttonsPrefixes.Lesson {

			chatId := telego.ChatID{ID: query.Message.GetChat().ID}

			roomName := queryData[1]
			roomId := queryData[2]
			timeStart := queryData[3]
			timeEnd := queryData[4]
			date := queryData[5]
			buildingName := queryData[6]
			buildingId := queryData[7]

			lesson, lessonType := LessonInfo(roomId, date, timeStart)

			msg := roomName + " (" + date + ")\n\n" + timeStart + " - " + timeEnd + "\n\n"

			var buttons [][]telego.InlineKeyboardButton
			var keyRow []telego.InlineKeyboardButton

			keyRow = append(keyRow, tu.InlineKeyboardButton("Назад").WithCallbackData(buttonsPrefixes.Room+"|"+roomName+"|"+fmt.Sprint(roomId)+"|"+buildingName+"|"+buildingId))

			if lessonType == "schedule" {
				msg += lesson.GetLector() + "\n\n" +
					lesson.GetSubject() + " (" + lesson.GetType() + ")\n\n" +
					lesson.GetGroups()

			} else if lessonType == "reserved" {
				msg += lesson.GetReserverName() + "\n\n" + lesson.GetComment()

			} else {
				msg += "Свободно\n\nБронирование только для авторизованных пользователей:"

				keyRow = append(keyRow, tu.InlineKeyboardButton("Забронировать").WithCallbackData(buttonsPrefixes.Reserve+"|"+roomId+"|"+date+"|"+timeStart+"|"+timeEnd))

			}

			buttons = append(buttons, keyRow)

			_, err := bot.EditMessageText(&telego.EditMessageTextParams{
				ChatID:      chatId,
				MessageID:   query.Message.GetMessageID(),
				Text:        msg,
				ReplyMarkup: &telego.InlineKeyboardMarkup{InlineKeyboard: buttons},
			})
			if err != nil {
				fmt.Println(err)
			}

		} else if btnPrefix == buttonsPrefixes.Reserve {

			roomId := queryData[1]
			date := queryData[2]
			timeStart := queryData[3]
			timeEnd := queryData[4]

			chatId := telego.ChatID{ID: query.Message.GetChat().ID}

			err = Reserve(lessonForReservationJSON{
				ChatId:    fmt.Sprint(chatId.ID),
				RoomId:    roomId,
				Date:      date,
				StartTime: timeStart,
				EndTime:   timeEnd,
			})

			if err != nil {
				fmt.Println(err)
				bot.SendMessage(tu.Message(
					chatId,
					"Для аутентификации требуется ввести логин и пароль в следующем формате:\n/auth <логин> <пароль>",
				))
				return
			}

			bot.SendMessage(tu.Message(
				chatId,
				"Аудитория успешно забронирована!",
			))

		}

	}, th.AnyCallbackQueryWithMessage())

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {

		query := strings.Split(message.Text, " ")

		if len(query) != 3 {
			bot.SendMessage(tu.Message(
				message.Chat.ChatID(),
				"Неверный формат.\nДля аутентификации требуется логин и пароль в следующем формате:\n/auth <логин> <пароль>",
			))
		} else {
			err := Auth(query[1], query[2], message.Chat.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}, th.CommandEqual("auth"))

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		if message.Text == "home" {
			ShowBuildingsPage(message.Chat.ChatID(), bot, buttonsPrefixes.Building)
		}
	}, th.AnyMessage())

	bh.Start()

}

func ShowBuildingsPage(chatID telego.ChatID, bot *telego.Bot, btnsPrefix string) {
	buildings := Buildings()
	buttonsData := make([]ButtonData, len(buildings))
	for index, building := range buildings {
		buttonsData[index] = ButtonData{ButtonName: building.Name, CallBackData: btnsPrefix + "|" + building.Name + "|" + fmt.Sprint(building.Id)}
	}

	buttonsInRow := 3

	msg := tu.Message(
		chatID,
		"Для просмотра аудиторий выберете нужный корпус:",
	).WithReplyMarkup(CreateKeyboard(buttonsData, buttonsInRow))

	bot.SendMessage(msg)
}

func ShowRoomPage() {

}

func Buildings() []models.Building {
	resp, err := http.Get("http://localhost:8080/api/buildings")
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

func Rooms(building_id string) []models.Room {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/rooms/%s", building_id))
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
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/schedule?room=%s&date=%s", room_id, date))
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
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/reserved_lesssons?room=%s&date=%s", room_id, date))
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

func LessonInfo(room_id string, date string, startTime string) (ILesson models.Lesson, lessonType string) {

	for _, lesson := range Schedule(room_id, date) {
		if lesson.StartTime == startTime {
			return lesson, "schedule"
		}
	}

	for _, lesson := range ReservedLessons(room_id, date) {
		if lesson.StartTime == startTime {
			return lesson, "reserved"
		}
	}

	return nil, ""
}

func Auth(username string, password string, chatId int64) (err error) {

	jsonData, _ := json.Marshal(LoginInput{Username: username, Password: password, ChatId: chatId})
	requestBody := bytes.NewBuffer(jsonData)

	resp, err := http.Post("http://localhost:8080/tg_login", "application/json", requestBody)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
	fmt.Println(http.StatusOK)
	if resp.StatusCode != http.StatusOK {
		return errors.New("неправильный логин или пароль")
	}

	return nil
}

func Reserve(lesson lessonForReservationJSON) error {

	jsonData, _ := json.Marshal(lesson)
	requestBody := bytes.NewBuffer(jsonData)

	resp, err := http.Post("http://localhost:8080/tg/reserve", "application/json", requestBody)
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

func CreateKeyboard(buttonsData []ButtonData, buttonsInRow int) *telego.InlineKeyboardMarkup {

	buttonsCount := len(buttonsData)

	if buttonsCount == 0 {
		return nil
	}

	var buttons [][]telego.InlineKeyboardButton
	rowCount := int(math.Ceil(float64(buttonsCount) / float64(buttonsInRow)))

	for i := 0; i < rowCount; i++ {
		var row []telego.InlineKeyboardButton
		for j := 0; j < buttonsInRow; j++ {
			if i*buttonsInRow+j < buttonsCount {
				btnData := buttonsData[i*buttonsInRow+j]
				row = append(row, tu.InlineKeyboardButton(btnData.ButtonName).WithCallbackData(btnData.CallBackData))
			}

		}
		buttons = append(buttons, row)
	}

	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	return keyboard
}
