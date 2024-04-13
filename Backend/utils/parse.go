package utils

import (
	"classrooms/internal/service"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Group struct {
	Name   string `json:"name"`
	Fac    string `json:"fac"`
	Level  string `json:"level"`
	Course string `json:"course"`
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

// type ScheduleItem struct {
// 	Time_start string
// 	Time_end   string
// 	Name       string
// 	Groups     []string
// 	Types      string
// 	Rooms      map[string]string
// }

// type ScheduleDate struct {
// 	Date  string
// 	Day   string
// 	Pairs map[string]ScheduleItem
// }

// type LectorSchedule struct {
// 	Name     string
// 	Groups   map[string]int
// 	Schedule map[string]ScheduleDate
// }

func GetHashAndNameLectors() (map[string]string, error) {
	lecturers := make(map[string]string)

	groups, _ := GetAllGroups()

	for _, group := range groups {
		groupHash := fmt.Sprintf("%x", md5.Sum([]byte(group.Name)))
		response, err := http.Get(fmt.Sprintf("https://public.mai.ru/schedule/data/%s.json", groupHash))
		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var temp map[string]json.RawMessage
		if err := json.Unmarshal(body, &temp); err != nil {
			return nil, err
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

		for _, value := range schedule.Dates {
			for _, value2 := range value.Pairs {
				for _, value3 := range value2 {
					for keyLector, lector := range value3.Lector {
						lecturers[keyLector] = lector
					}
				}

			}
		}
	}

	return lecturers, nil
}

func GetAllGroups() ([]Group, error) {
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

	return groups, nil
}

func GetLectorSchedule(LectorHash string) (service.LectorSchedule, error) {
	response, err := http.Get(fmt.Sprintf("https://public.mai.ru/schedule/data/%s.json", LectorHash))
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %s\n", err)
		return service.LectorSchedule{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return service.LectorSchedule{}, err
	}

	var lectorSchedule service.LectorSchedule
	if err := json.Unmarshal(body, &lectorSchedule); err != nil {
		return service.LectorSchedule{}, err
	}

	return lectorSchedule, nil
}

func WriteToFile(data any, filename string) error {

	folderPath := "../internal/localStorage/"

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	file, err := os.Create(folderPath + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
