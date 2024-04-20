package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"tg_bot/models"

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
	}

	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		chatID := tu.ID(update.Message.Chat.ID)

		buildings := Buildings()
		buttonsData := make([]ButtonData, len(buildings))
		for index, building := range buildings {
			buttonsData[index] = ButtonData{ButtonName: building.Name, CallBackData: buttonsPrefixes.Building + "-" + fmt.Sprint(building.Id)}
		}

		buttonsInRow := 3

		msg := tu.Message(
			chatID,
			"Для просмотра аудиторий выберете нужный корпус:",
		).WithReplyMarkup(CreateKeyboard(buttonsData, buttonsInRow))

		bot.SendMessage(msg)
	}, th.CommandEqual("start"))

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		queryData := strings.Split(query.Data, "-")
		btnPrefix := queryData[0]
		callBackData := queryData[1]

		if btnPrefix == buttonsPrefixes.Building {
			rooms := Rooms(callBackData)

			buttonsData := make([]ButtonData, len(rooms))
			for index, room := range rooms {
				buttonsData[index] = ButtonData{ButtonName: room.Name, CallBackData: buttonsPrefixes.Room + "-" + fmt.Sprint(room.ID)}
			}

			buttonsInRow := 4

			_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID))

			chatId := telego.ChatID{ID: query.Message.GetChat().ID}
			_, err := bot.EditMessageText(&telego.EditMessageTextParams{
				ChatID:      chatId,
				MessageID:   query.Message.GetMessageID(),
				Text:        "Выберете аудиторию:",
				ReplyMarkup: CreateKeyboard(buttonsData, buttonsInRow),
			})
			if err != nil {
				fmt.Println("Ошибка при редактировании сообщения:", err)
			}

		} else if btnPrefix == buttonsPrefixes.Room {
			roomId := callBackData
			date := "20.04.2024"
			schedule := Schedule(callBackData, date)

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
					ButtonName:   fmt.Sprint(lessons_timeStart[lesson.StartTime]) + ". " + lesson.StartTime + " - " + lesson.EndTime + " По расписанию",
					CallBackData: buttonsPrefixes.Lesson + "-" + roomId + "-" + lesson.StartTime + "-" + date}
			}

			for timeStart, lessonNumber := range lessons_timeStart {

				if buttonsData[lessonNumber-1].ButtonName == "" {
					buttonsData[lessonNumber-1] = ButtonData{
						ButtonName:   fmt.Sprint(lessonNumber) + ". " + timeStart + " - " + lessons_timeEnd[lessonNumber] + " Свободно",
						CallBackData: buttonsPrefixes.Lesson + "-" + roomId + "-" + timeStart + "-" + date}
				}
			}

			buttonsInRow := 1

			_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID))

			chatId := telego.ChatID{ID: query.Message.GetChat().ID}
			_, err := bot.EditMessageText(&telego.EditMessageTextParams{
				ChatID:      chatId,
				MessageID:   query.Message.GetMessageID(),
				Text:        "Выберете занятие:",
				ReplyMarkup: CreateKeyboard(buttonsData, buttonsInRow),
			})
			if err != nil {
				fmt.Println("Ошибка при редактировании сообщения:", err)
			}
		} else if btnPrefix == buttonsPrefixes.Lesson {

			roomId := callBackData
			timeStart := queryData[2]
			date := queryData[3]

			lesson, lessonType := LessonInfo(roomId, date, timeStart)

			msg := roomId + " (" + date + ")\n\n" + lesson.GetStartTime() + " - " + lesson.GetEndTime() + "\n\n"

			if lessonType == "schedule" {
				msg += lesson.GetLector() + "\n\n" +
					lesson.GetSubject() + " (" + lesson.GetType() + ")\n\n" +
					lesson.GetGroups() + "\n\n"
			}
			chatId := telego.ChatID{ID: query.Message.GetChat().ID}
			_, err := bot.EditMessageText(&telego.EditMessageTextParams{
				ChatID:    chatId,
				MessageID: query.Message.GetMessageID(),
				Text:      msg,
			})
			if err != nil {
				fmt.Println("Ошибка при редактировании сообщения:", err)
			}
		}

	}, th.AnyCallbackQueryWithMessage())

	bh.Start()

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
