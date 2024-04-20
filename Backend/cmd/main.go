package main

import (
	"classrooms/internal/handler"
	"classrooms/internal/repository"
	"classrooms/internal/service"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Group struct {
	Name   string `json:"name"`
	Fac    string `json:"fac"`
	Level  string `json:"level"`
	Course string `json:"course"`
}

type Pair struct {
	TimeStart string            `json:"time_start"`
	TimeEnd   string            `json:"time_end"`
	Name      string            `json:"name"`
	Groups    []string          `json:"groups"`
	Type      []string          `json:"types"`
	Rooms     map[string]string `json:"rooms"`
}

type DaySchedule struct {
	Date  string           `json:"date"`
	Day   string           `json:"day"`
	Pairs map[string]*Pair `json:"pairs"`
}

type Schedule struct {
	Name     string                  `json:"name"`
	Groups   map[string]int          `json:"groups"`
	Schedule map[string]*DaySchedule `json:"schedule"`
}

type Lesson struct {
	Lector    string
	Date      string
	TimeStart string
	TimeEnd   string
	Name      string
	Groups    []string
	Type      string
	Room      string
}

func main() {

	// groupNames, _ := utils.GetAllGroupNames()
	// var groups []string
	// for _, group := range groupNames {
	// 	groups = append(groups, group.Name)
	// }
	// utils.WriteToFile(groups, "allgroupNames.json")

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading env variables.")
	}

	db, err := repository.ConnectDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMode"),
	})

	// lectorSchedule, err := utils.GetLectorSchedule("0cb88141-658c-11e4-a65a-00155d79380a")
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// }

	// for key, _ := range lectorSchedule.Schedule {
	// 	fmt.Println("key: " + key)

	// }

	// os.Exit(1)

	if err != nil {
		log.Fatalf("Failed to init DB: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	// groupNames, _ := utils.GetAllGroupNames()

	// for _, group := range groupNames {
	// 	res, _ := utils.GetGroupSchedule(group)
	// 	utils.InsertGroupScheduleIntoDB(res, db)
	// }

	// os.Exit(1)
	/////////////////////
	///////////////
	/////////

	//services.Schedule.RefreshLectorSchedule()

	// filePathWithLectorsInfo := "../internal/localStorage/allLectors.json"

	// fileContent, err := os.ReadFile(filePathWithLectorsInfo)
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// var lectors map[string]string
	// if err := json.Unmarshal(fileContent, &lectors); err != nil {
	// 	fmt.Println("Ошибка при декодировании JSON:", err)
	// }

	// for _, lectorName := range lectors {
	// 	var lesson models.Lesson
	// 	err := db.Where("lector = ?", lectorName).First(&lesson).Error

	// 	if err != nil {
	// 		fmt.Println("OPA: ", lectorName, "\n", err)
	// 	}

	// }
	//err = services.Schedule.RefreshLectorSchedule()
	//fmt.Println(lectorSchedule)

	handlers.InitRoutes().Run("localhost:8080")

}
