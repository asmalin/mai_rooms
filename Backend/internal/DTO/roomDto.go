package dto

import "time"

type RoomDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ScheduleLessonDto struct {
	Lector    string `json:"lector"`
	StartTime string `json:"time_start"`
	EndTime   string `json:"time_end"`
	Subject   string `json:"subject"`
	Groups    string `json:"groups"`
	Type      string `json:"type"`
}

type ReservedLessonDto struct {
	ReserverName string `json:"reserver"`
	ReserverId   int    `json:"reserver_id"`
	Date         string `json:"date"`
	StartTime    string `json:"time_start"`
	EndTime      string `json:"time_end"`
	Comment      string `json:"comment"`
}

type LessonForCancelReservationDto struct {
	ReserverId int       `json:"reserver_id"`
	Room_id    int       `json:"room_id"`
	Date       time.Time `json:"date"`
	StartTime  string    `json:"time_start"`
}
