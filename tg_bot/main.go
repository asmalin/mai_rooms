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
		Building:      "building",
		Room:          "room",
		Lesson:        "lesson",
		Reserve:       "reserve",
		CancelReserve: "cancel",
	}

	var allBuildings = make(map[string]models.Building)
	var allRooms = make(map[string]models.Room)

	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		chatID := tu.ID(update.Message.Chat.ID)

		bot.SendMessage(tu.Message(
			chatID,
			"Здесь можно посмотреть расписание занятий на аудиторию и занять свободную.",
		))

		ShowBuildingsPage(chatID, bot, buttonsPrefixes.Building)

		buildings := service.Buildings()
		for _, building := range buildings {
			allBuildings[fmt.Sprint(building.Id)] = building
			rooms := service.Rooms(building.Id)
			for _, room := range rooms {
				room.Building_id = building.Id
				allRooms[fmt.Sprint(room.ID)] = room
			}
		}

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

			userId, _ := service.GetUserIdByTgUsername(query.Message.GetChat().Username)

			if userId == lesson.GetReserverId() {
				keyRow = append(keyRow, tu.InlineKeyboardButton("Отменить бронь").WithCallbackData(buttonsPrefixes.CancelReserve+"|"+roomId+"|"+date+"|"+timeStart))
			}

		} else {

			_, err := service.GetUserIdByTgUsername(query.Message.GetChat().Username)

			if err != nil {
				msg += "Аудитория свободна.\nДля ее бронирования необходимо быть авторизованным (Подробнее /info)"

			} else {
				keyRow = append(keyRow, tu.InlineKeyboardButton("Забронировать").WithCallbackData(buttonsPrefixes.Reserve+"|"+roomId+"|"+date+"|"+timeStart+"|"+timeEnd))
			}
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

		roomIdStr := queryData[1]

		roomId, _ := strconv.Atoi(roomIdStr)
		date := queryData[2]
		timeStart := queryData[3]
		timeEnd := queryData[4]

		err = service.Reserve(query.Message.GetChat().Username, models.LessonForReservationJSON{
			RoomId:    roomId,
			Date:      date,
			StartTime: timeStart,
			EndTime:   timeEnd,
		})

		if err != nil {
			bot.SendMessage(tu.Message(
				chatId,
				err.Error(),
			))
			return
		}

		_, err := bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    chatId,
			MessageID: query.Message.GetMessageID(),
			Text:      "Аудитория успешно забронирована!",
			ReplyMarkup: tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Назад").WithCallbackData(buttonsPrefixes.Room + "|" + roomIdStr),
				),
			),
		})

		if err != nil {
			fmt.Println(err)
		}

	}, th.CallbackDataContains(buttonsPrefixes.Reserve))

	//cancel reserve lesson
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		chatId := telego.ChatID{ID: query.Message.GetChat().ID}
		queryData := strings.Split(query.Data, "|")

		roomIdStr := queryData[1]

		roomId, _ := strconv.Atoi(roomIdStr)
		date := queryData[2]
		timeStart := queryData[3]

		err = service.CancelReserve(query.Message.GetChat().Username, models.LessonForCancelReservationJSON{
			RoomId:    roomId,
			Date:      date,
			StartTime: timeStart,
		})

		if err != nil {
			bot.SendMessage(tu.Message(
				chatId,
				err.Error(),
			))
			return
		}

		_, err := bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    chatId,
			MessageID: query.Message.GetMessageID(),
			Text:      "Отмена брони произведена успешно!",
			ReplyMarkup: tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Назад").WithCallbackData(buttonsPrefixes.Room + "|" + roomIdStr),
				),
			),
		})

		if err != nil {
			fmt.Println(err)
		}

	}, th.CallbackDataContains(buttonsPrefixes.CancelReserve))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		ShowBuildingsPage(update.Message.Chat.ChatID(), bot, buttonsPrefixes.Building)
	}, th.CommandEqual("home"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		msgText := "Здесь можно посмотреть расписание занятий на аудиторию, а также забронировать свободную.\n\n" +
			"Для того чтобы забронировать аудиторию необходимо чтобы имя пользователя в телеграмме было привязано к аккаунту на веб-сайте." +
			"Это можно сделать на <a href=\"" + os.Getenv("BASE_URL") + "/profile\">странице своего профиля на веб-сайте</a> " +
			"и после авторизации ввести в поле \"tg username\" свое имя пользователя из телеграмма."

		bot.SendMessage(&telego.SendMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			Text:      msgText,
			ParseMode: "HTML",
		})
	}, th.CommandEqual("info"))

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
