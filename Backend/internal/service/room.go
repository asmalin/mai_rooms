package service

import (
	dto "classrooms/internal/DTO"
	"classrooms/internal/models"
	"classrooms/internal/repository"
)

type RoomService struct {
	repo repository.Room
}

func NewRoomService(repo repository.Room) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) GetAllBuildings() ([]models.Building, error) {
	return s.repo.GetAllBuildings()
}

func (s *RoomService) GetRoomsByBuildingId(buildingId int) ([]dto.RoomDto, error) {

	building, err := s.repo.GetBuildingById(buildingId)
	if err != nil {
		return nil, err
	}

	rooms, err := s.repo.GetRoomsByBuildingId(building.ID)
	if err != nil {
		return nil, err
	}

	roomsDTO := make([]dto.RoomDto, len(rooms))
	for i := 0; i < len(rooms); i++ {
		roomsDTO[i].Id = rooms[i].ID
		roomsDTO[i].Name = rooms[i].Name
	}

	return roomsDTO, nil
}

func (s *RoomService) GetRoomById(roomId int) (models.Room, error) {
	return s.repo.GetRoomById(roomId)
}

func (s *RoomService) GetRoomId(roomName string) (int, error) {
	return s.repo.GetRoomId(roomName)
}
