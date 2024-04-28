package utils

import (
	"fmt"
	"math"
	"tg_bot/models"
	"tg_bot/service"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CreateButtonsForRoom(roomId, date, buildingIdStr string) (keyboard *telego.InlineKeyboardMarkup) {

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

	schedule := service.Schedule(roomId, date)
	reserved := service.ReservedLessons(roomId, date)

	buttonsData := make([]ButtonData, len(lessons_timeStart))

	for _, lesson := range schedule {
		buttonsData[lessons_timeStart[lesson.StartTime]-1] = ButtonData{
			ButtonName:   fmt.Sprint(lessons_timeStart[lesson.StartTime]) + ". " + lesson.StartTime + " - " + lesson.EndTime + " По расписанию 📅",
			CallBackData: "lesson" + "|" + roomId + "|" + lesson.StartTime + "|" + lesson.EndTime + "|" + date}
	}

	for _, lesson := range reserved {

		buttonsData[lessons_timeStart[lesson.StartTime]-1] = ButtonData{
			ButtonName:   fmt.Sprint(lessons_timeStart[lesson.StartTime]) + ". " + lesson.StartTime + " - " + lesson.EndTime + " Забронировано 🔶",
			CallBackData: "lesson" + "|" + roomId + "|" + lesson.StartTime + "|" + lesson.EndTime + "|" + date}
	}

	for timeStart, lessonNumber := range lessons_timeStart {

		if buttonsData[lessonNumber-1].ButtonName == "" {
			buttonsData[lessonNumber-1] = ButtonData{
				ButtonName:   fmt.Sprint(lessonNumber) + ". " + timeStart + " - " + lessons_timeEnd[lessonNumber] + " Свободно ✅",
				CallBackData: "lesson" + "|" + roomId + "|" + timeStart + "|" + lessons_timeEnd[lessonNumber] + "|" + date}
		}
	}

	buttonsInRow := 1

	keyBoard := CreateKeyboard(buttonsData, buttonsInRow)

	prevBtn := []telego.InlineKeyboardButton{tu.InlineKeyboardButton("Назад").WithCallbackData("building" + "|" + buildingIdStr)}
	keyBoard.InlineKeyboard = append(keyBoard.InlineKeyboard, prevBtn)

	return keyBoard

}

func LessonInfo(room_id string, date string, startTime string) (ILesson models.Lesson, lessonType string) {

	for _, lesson := range service.Schedule(room_id, date) {
		if lesson.StartTime == startTime {
			return lesson, "schedule"
		}
	}

	for _, lesson := range service.ReservedLessons(room_id, date) {
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
