package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"tg_bot/models"
	"time"

	"tg_bot/service"
	"tg_bot/utils"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

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

	var buttonsPrefixes = utils.ButtonsPrefixes{
		Building: "building",
		Room:     "room",
		Lesson:   "lesson",
		Reserve:  "reserve",
	}

	var allBuildings = make(map[string]models.Building)
	var allRooms = make(map[string]models.Room)

	buildings := service.Buildings()
	for _, building := range buildings {
		allBuildings[fmt.Sprint(building.Id)] = building
		rooms := service.Rooms(building.Id)
		for _, room := range rooms {
			room.Building_id = building.Id
			allRooms[fmt.Sprint(room.ID)] = room
		}
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

	//rooms by buildling id
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		queryData := strings.Split(query.Data, "|")
		buildingIdStr := queryData[1]

		buildingId, err := strconv.Atoi(buildingIdStr)
		if err != nil {
			fmt.Println("Converting string to int error: ", err)
			return
		}
		rooms := service.Rooms(buildingId)

		buttonsData := make([]utils.ButtonData, len(rooms))
		for index, room := range rooms {
			buttonsData[index] = utils.ButtonData{
				ButtonName:   room.Name,
				CallBackData: buttonsPrefixes.Room + "|" + fmt.Sprint(room.ID)}
		}

		buttonsInRow := 3

		chatId := telego.ChatID{ID: query.Message.GetChat().ID}

		_, err = bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   query.Message.GetMessageID(),
			Text:        "Корпус: " + allBuildings[buildingIdStr].Name + "\nВыберете аудиторию:",
			ReplyMarkup: utils.CreateKeyboard(buttonsData, buttonsInRow),
		})

		if err != nil {
			fmt.Println("Ошибка при редактировании сообщения:", err)
		}

	}, th.CallbackDataContains(buttonsPrefixes.Building))

	//lessons by room id
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		chatId := telego.ChatID{ID: query.Message.GetChat().ID}

		queryData := strings.Split(query.Data, "|")
		roomId := queryData[1]

		currentTime := time.Now()
		location, _ := time.LoadLocation("Europe/Moscow")
		date := currentTime.In(location).Format("02.01.2006")

		keyboard := utils.CreateButtonsForRoom(roomId, date, fmt.Sprint(allRooms[roomId].Building_id))

		_, err := bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   query.Message.GetMessageID(),
			Text:        "Аудитория: " + allRooms[roomId].Name + "\nВыберете занятие:",
			ReplyMarkup: keyboard,
		})
		if err != nil {
			fmt.Println("Ошибка при редактировании сообщения:", err)
		}

	}, th.CallbackDataContains(buttonsPrefixes.Room))

	// lesson by roomId,startTime,date
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		chatId := telego.ChatID{ID: query.Message.GetChat().ID}
		queryData := strings.Split(query.Data, "|")

		roomId := queryData[1]
		timeStart := queryData[2]
		timeEnd := queryData[3]
		date := queryData[4]

		lesson, lessonType := utils.LessonInfo(roomId, date, timeStart)

		msg := "<b>" + allRooms[roomId].Name + " (" + date + ")</b>\n\n" + timeStart + " - " + timeEnd + "\n\n"

		var buttons [][]telego.InlineKeyboardButton
		var keyRow []telego.InlineKeyboardButton

		keyRow = append(keyRow, tu.InlineKeyboardButton("Назад").WithCallbackData(buttonsPrefixes.Room+"|"+fmt.Sprint(roomId)))

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
			ParseMode:   "HTML",
		})
		if err != nil {
			fmt.Println(err)
		}
	}, th.CallbackDataContains(buttonsPrefixes.Lesson))

	//reserve free lesson
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		chatId := telego.ChatID{ID: query.Message.GetChat().ID}
		queryData := strings.Split(query.Data, "|")

		roomId := queryData[1]
		date := queryData[2]
		timeStart := queryData[3]
		timeEnd := queryData[4]

		err = service.Reserve(models.LessonForReservationJSON{
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

		_, err := bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    chatId,
			MessageID: query.Message.GetMessageID(),
			Text:      "Аудитория успешно забронирована!",
			ReplyMarkup: tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Назад").WithCallbackData(buttonsPrefixes.Room + "|" + roomId),
				),
			),
		})

		if err != nil {
			fmt.Println(err)
		}

	}, th.CallbackDataContains(buttonsPrefixes.Reserve))

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {

		query := strings.Split(message.Text, " ")

		if len(query) != 3 {
			bot.SendMessage(tu.Message(
				message.Chat.ChatID(),
				"Неверный формат.\nДля аутентификации требуется логин и пароль в следующем формате:\n/auth <логин> <пароль>",
			))
		} else {
			err := service.Auth(query[1], query[2], message.Chat.ID)
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
	buildings := service.Buildings()
	buttonsData := make([]utils.ButtonData, len(buildings))
	for index, building := range buildings {
		buttonsData[index] = utils.ButtonData{ButtonName: building.Name, CallBackData: btnsPrefix + "|" + fmt.Sprint(building.Id)}
	}

	buttonsInRow := 3

	msg := tu.Message(
		chatID,
		"Для просмотра аудиторий выберете нужный корпус:",
	).WithReplyMarkup(utils.CreateKeyboard(buttonsData, buttonsInRow))

	bot.SendMessage(msg)
}
