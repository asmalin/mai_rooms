package service

import (
	dto "classrooms/internal/DTO"
	"classrooms/internal/models"
	"classrooms/internal/repository"
)

type Login interface {
	WebAuth(username, password string) (models.User, error)
	TgAuth(username, password string, tgChatId int64) (err error)
	UserIdByChatId(tgChatId int64) (userId int, err error)

	GenerateAccessToken(models.User) (string, error)
	GenerateRefreshToken(userId int) (string, error)

	ParseToken(token string) (int, error)
	GetUserById(id int) (models.User, error)
}

type Reservation interface {
	ReserveRoom(reservedRoom models.ReservedLesson) error
	CancelReservation(lessonForCancelReservation dto.LessonForCancelReservationDto) error
	GetReservationRoom(roomId int, date string) ([]models.ReservedLesson, error)
	GetAllReservedLessons() ([]models.ReservedLesson, error)
}

type QRCode interface {
}

type Room interface {
	GetAllBuildings() ([]models.Building, error)
	GetRoomsByBuildingId(buildingId int) ([]dto.RoomDto, error)
	GetRoomById(roomId int) (models.Room, error)
	GetRoomId(roomName string) (int, error)
}

type Schedule interface {
	GetScheduleByRoomIdAndDate(roomId int, Date string) ([]dto.ScheduleLessonDto, error)
	WriteScheduleToDB(lectorSchedule LectorSchedule) error
	RefreshLectorSchedule() error
}

type Users interface {
	GetAllUsers() ([]dto.UserDto, error)
	CreateUser(models.User) error
	DeleteUser(userId int) error
	UpdateUser(user models.User) error
	ChangePassword(oldPassword, newPassword string, user models.User) error
}

type Service struct {
	Login
	Reservation
	QRCode
	Room
	Schedule
	Users
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Login:       NewLoginService(repo.Login),
		Reservation: NewReservationService(repo.Reservation),
		QRCode:      nil,
		Room:        NewRoomService(repo.Room),
		Schedule:    NewScheduleService(repo.Room, repo.Lesson),
		Users:       NewUsersService(repo.Users),
	}
}
