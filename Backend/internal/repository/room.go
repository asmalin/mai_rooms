package repository

import (
	"classrooms/internal/models"
	"errors"

	"gorm.io/gorm"
)

type RoomPostgres struct {
	db *gorm.DB
}

func NewRoomPostgres(db *gorm.DB) *RoomPostgres {
	return &RoomPostgres{db: db}
}

func (r *RoomPostgres) GetAllBuildings() ([]models.Building, error) {
	var buildings []models.Building
	if result := r.db.Order("name ASC").Find(&buildings); result.Error != nil {
		return nil, errors.New("ошибка при получении данных")
	}

	return buildings, nil
}

func (r *RoomPostgres) GetBuildingById(buildingId int) (models.Building, error) {
	var building models.Building
	result := r.db.Where("id = ?", buildingId).First(&building)
	if result.Error != nil {
		return models.Building{}, errors.New("ошибка при получении данных")
	}

	return building, nil
}

func (r *RoomPostgres) GetRoomsByBuildingId(buildingId int) ([]models.Room, error) {
	var building models.Building
	result := r.db.Where("id = ?", buildingId).First(&building)
	if result.Error != nil {
		return nil, errors.New("корпус с таким id не найден")
	}

	var rooms []models.Room
	result = r.db.Order("name ASC").Where("building_id = ?", buildingId).Find(&rooms)
	if result.Error != nil {
		return nil, errors.New("корпус с таким id не найден")
	}

	return rooms, nil
}

func (r *RoomPostgres) GetRoomById(roomId int) (models.Room, error) {
	var room models.Room
	result := r.db.Take(&room, "id = ?", roomId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Room{}, errors.New("аудитория не найдена")
	} else if result.Error != nil {
		return models.Room{}, errors.New("произошла ошибка при выполнении запроса")
	} else {
		return room, nil
	}
}

func (r *RoomPostgres) GetRoomId(roomName string) (int, error) {
	var roomId int
	result := r.db.Model(&models.Room{}).Select("id").Where("name = ?", roomName).First(&roomId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, errors.New("аудитория не найдена")
	} else if result.Error != nil {
		return 0, result.Error
	} else {
		return roomId, nil
	}
}
