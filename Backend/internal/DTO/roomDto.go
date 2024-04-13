package dto

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
	StartTime    string `json:"time_start"`
	EndTime      string `json:"time_end"`
	Comment      string `json:"comment"`
	Reserved     bool   `json:"reserved"`
}
