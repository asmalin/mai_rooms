package repository

import (
	"classrooms/internal/models"
	"time"

	"gorm.io/gorm"
)

type Login interface {
	GetUser(username, password string) (models.User, error)
	GetUserById(id int) (models.User, error)
}

type Reservation interface {
	ReserveRoom(models.ReservedLesson) error
	GetReservedLessons(roomId int, date string) ([]models.ReservedLesson, error)
}

type QRCode interface {
}

type Room interface {
	GetAllBuildings() ([]models.Building, error)
	GetBuildingById(buildingId int) (models.Building, error)
	GetRoomsByBuildingId(buildingId int) ([]models.Room, error)
	GetRoomById(roomId int) (models.Room, error)
	GetRoomId(roomName string) (int, error)
}

type Lesson interface {
	GetScheduleLessons(roomId int, date time.Time) ([]models.Lesson, error)
	InsertLessonToDB(lesson models.Lesson) error
	DeleteAll() error
}

type Repository struct {
	Login
	Reservation
	QRCode
	Room
	Lesson
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Login:       NewLoginPostgres(db),
		Reservation: NewReservationPostgres(db),
		QRCode:      nil,
		Room:        NewRoomPostgres(db),
		Lesson:      NewLessonPostgres(db),
	}
}
