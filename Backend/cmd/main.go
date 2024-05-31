package main

import (
	"classrooms/internal/handler"
	"classrooms/internal/repository"
	"classrooms/internal/service"
	"classrooms/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func main() {

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

	if err != nil {
		log.Fatalf("Failed to init DB: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	err = utils.InitDBEntities(db)
	if err != nil {
		log.Fatalf("Failed to init DB entities: %s", err.Error())
	}

	setupCron(db)

	handlers.InitRoutes().Run("0.0.0.0:5001")

}

func setupCron(db *gorm.DB) {
	c := cron.New()

	c.AddFunc("0 3 * * 0", func() {
		utils.RefreshSchedule(db)
	})
	c.Start()
}
