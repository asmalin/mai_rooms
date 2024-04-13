package utils

import (
	"classrooms/internal/service"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func uniqueStrings(input []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueSlice []string

	for _, str := range input {
		if !uniqueMap[str] {
			uniqueMap[str] = true
			uniqueSlice = append(uniqueSlice, str)
		}
	}

	return uniqueSlice
}

func GetAllRooms() []string {

	folderPath := "../internal/localStorage/LectorsDataJSON/"

	// Чтение всех файлов в указанной директории
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	var fullNames []string
	var buildings []string
	var rooms []string

	for _, file := range files {
		//fmt.Println(file.Name())
		filePath := folderPath + file.Name()
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}
		var schedule service.LectorSchedule
		if err := json.Unmarshal(fileContent, &schedule); err != nil {
			fmt.Println("Ошибка при декодировании JSON:", err)
			continue
		}
		for _, val := range schedule.Schedule {
			for _, val2 := range val.Pairs {
				for _, val3 := range val2.Rooms {
					if !strings.Contains(val3, "-") || val3[0] == '-' {
						continue
					}

					//fmt.Println(val3)
					fullnameRoom := strings.SplitN(val3, "-", 2)
					fullNames = append(fullNames, val3)
					buildings = append(buildings, fullnameRoom[0])
					rooms = append(rooms, fullnameRoom[1])

				}
			}
		}
	}
	return uniqueStrings(fullNames)

}

func insertRooms() {
	// db.AutoMigrate(&models.Room{}, &models.Building{})

	// allRooms := utils.GetAllRooms()

	// for i := 0; i < len(allRooms); i++ {

	// 	fullnameRoom := strings.SplitN(allRooms[i], "-", 2)
	// 	fmt.Println(fullnameRoom[0], "\t", fullnameRoom[1])
	// 	var building models.Building
	// 	db.Take(&building, "name = ?", fullnameRoom[0])
	// 	fmt.Println(building)
	// 	room := models.Room{Name: fullnameRoom[1], Building_id: building.ID}
	// 	db.Create(&room)
	// }
}
