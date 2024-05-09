package repository

import (
	"classrooms/internal/models"
	"time"

	"gorm.io/gorm"
)

type Login interface {
	GetUser(username, password string) (models.User, error)
	GetUserById(id int) (models.User, error)

	AddUserTgChatRelation(userId int, tgChatId int64) error
	GetUserChatRelation(tgChatId int64) (models.UserTgChatRelation, error)
}

type Reservation interface {
	ReserveRoom(models.ReservedLesson) error
	CancelReservation(reservedLesson_id int) error
	GetReservedLesson(roomId int, date time.Time, startTime string) (models.ReservedLesson, error)
	GetReservedLessons(roomId int, date time.Time) ([]models.ReservedLesson, error)
	GetAllReservedLessons() ([]models.ReservedLesson, error)
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

type Users interface {
	GetAllUsers() ([]models.User, error)
	CreateUser(user models.User) error
	DeleteUser(userId int) error
	UpdateUser(user models.User) error
}

type Repository struct {
	Login
	Reservation
	Room
	Lesson
	Users
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Login:       NewLoginPostgres(db),
		Reservation: NewReservationPostgres(db),
		Room:        NewRoomPostgres(db),
		Lesson:      NewLessonPostgres(db),
		Users:       NewUsersPostgres(db),
	}
}
