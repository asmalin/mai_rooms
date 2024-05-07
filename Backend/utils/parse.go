package utils

import (
	"classrooms/internal/models"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Group struct {
	Name string `json:"name"`
}

type Schedule struct {
	Group string `json:"group"`
	Dates map[string]DaySchedule
}

type DaySchedule struct {
	Day   string                     `json:"day"`
	Pairs map[string]map[string]Pair `json:"pairs"`
}

type Pair struct {
	TimeStart string            `json:"time_start"`
	TimeEnd   string            `json:"time_end"`
	Lector    map[string]string `json:"lector"`
	Type      map[string]int    `json:"type"`
	Room      map[string]string `json:"room"`
}

func RefreshSchedule(db *gorm.DB) {

	groupNames, err := GetAllGroupNames()
	if err != nil {
		return
	}

	if err := db.Exec("DELETE FROM lessons; DELETE FROM rooms; DELETE FROM buildings;").Error; err != nil {
		return
	}

	fmt.Println("DELETED!")
	for _, group := range groupNames {
		res, _ := GetGroupSchedule(group)
		InsertGroupScheduleIntoDB(res, db)
	}

}

func GetAllGroupNames() ([]string, error) {
	response, err := http.Get("https://public.mai.ru/schedule/data/groups.json")
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %s\n", err)
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var groups []Group
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}

	var groupNames []string
	for _, group := range groups {
		groupNames = append(groupNames, group.Name)
	}

	return groupNames, nil
}

func GetGroupSchedule(groupName string) (Schedule, error) {

	groupHash := fmt.Sprintf("%x", md5.Sum([]byte(groupName)))
	response, err := http.Get(fmt.Sprintf("https://public.mai.ru/schedule/data/%s.json", groupHash))
	if err != nil {
		return Schedule{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Schedule{}, err
	}

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(body, &temp); err != nil {
		return Schedule{}, err
	}

	var schedule Schedule

	json.Unmarshal(temp["group"], &schedule.Group)

	schedule.Dates = make(map[string]DaySchedule)
	delete(temp, "group")

	for date, value := range temp {
		var daySchedule DaySchedule
		if err := json.Unmarshal(value, &daySchedule); err != nil {
			log.Fatal(err)
		}
		schedule.Dates[date] = daySchedule
	}

	return schedule, nil
}

func InsertGroupScheduleIntoDB(schedule Schedule, db *gorm.DB) error {

	db.AutoMigrate(&models.Lesson{}, &models.Building{}, &models.Room{})
	fmt.Println(schedule.Group)
	for date, pairs := range schedule.Dates {
		for _, pairTime := range pairs.Pairs {
			for subjectName, subject := range pairTime {
				var lectors string
				for _, lectorName := range subject.Lector {
					if len(lectorName) > 0 {
						lectors += lectorName + " / "
					}

				}
				if len(lectors) > 3 {
					lectors = lectors[:len(lectors)-3]
				}

				for lessonType, _ := range subject.Type {
					for _, roomName := range subject.Room {
						if strings.HasPrefix(roomName, "--") || roomName == "" || roomName == "дистанционно" || strings.HasPrefix(roomName, "<Объект не найден>") {
							continue
						}

						parsedDate, err := time.Parse("02.01.2006", date)
						if err != nil {
							fmt.Println("Parse date error: ", err)
							continue
						}

						var room models.Room
						result := db.Take(&room, "name = ?", roomName)
						fmt.Println(room)
						fmt.Println(result.Error)
						if result.Error != nil {
							fmt.Println(roomName)
							fullnameRoom := strings.SplitN(roomName, "-", 2)
							var building models.Building
							result = db.Take(&building, "name = ?", fullnameRoom[0])
							if result.Error != nil {
								building.Name = fullnameRoom[0]
								db.Create(&building)
							}
							room = models.Room{Building_id: building.ID, Name: roomName}
							db.Create(&room)
						}

						lesson := models.Lesson{Lector: lectors, Date: parsedDate, TimeStart: subject.TimeStart, TimeEnd: subject.TimeEnd,
							Name: subjectName, Groups: schedule.Group, Type: lessonType, RoomId: room.ID, Room: room}
						db.Create(&lesson)
					}
				}

			}

		}
	}
	return nil
}
