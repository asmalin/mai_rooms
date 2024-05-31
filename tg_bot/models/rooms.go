package models

type Building struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Room struct {
	ID          int    `json:"id"`
	Building_id int    `json:"building_id"`
	Name        string `json:"name"`
}

type LessonForReservationJSON struct {
	RoomId    int    `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Comment   string `json:"comment"`
}

type LessonForCancelReservationJSON struct {
	RoomId    int    `json:"roomId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
}

type Lesson interface {
	GetLector() string
	GetStartTime() string
	GetEndTime() string
	GetSubject() string
	GetGroups() string
	GetType() string
	GetReserverName() string
	GetReserverId() int
	GetComment() string
}

type ScheduleLesson struct {
	Lector    string `json:"lector"`
	StartTime string `json:"time_start"`
	EndTime   string `json:"time_end"`
	Subject   string `json:"subject"`
	Groups    string `json:"groups"`
	Type      string `json:"type"`
}

func (s ScheduleLesson) GetLector() string {
	return s.Lector
}

func (s ScheduleLesson) GetStartTime() string {
	return s.StartTime
}

func (s ScheduleLesson) GetEndTime() string {
	return s.EndTime
}

func (s ScheduleLesson) GetSubject() string {
	return s.Subject
}

func (s ScheduleLesson) GetGroups() string {
	return s.Groups
}

func (s ScheduleLesson) GetType() string {
	return s.Type
}

func (s ScheduleLesson) GetReserverName() string {
	return "" // Возвращаем пустую строку, так как для расписания урока резервирующийся человек не определен
}

func (s ScheduleLesson) GetReserverId() int {
	return 0 // Возвращаем 0, так как для расписания урока ID резервирующего человека не определен
}

func (s ScheduleLesson) GetComment() string {
	return "" // Возвращаем пустую строку, так как для расписания урока комментарий не определен
}

type ReservedLesson struct {
	ReserverName string `json:"reserver"`
	ReserverId   int    `json:"reserver_id"`
	Date         string `json:"date"`
	StartTime    string `json:"time_start"`
	EndTime      string `json:"time_end"`
	Comment      string `json:"comment"`
}

func (r ReservedLesson) GetLector() string {
	return "" // Возвращаем пустую строку, так как для зарезервированного урока лектор не определен
}

func (r ReservedLesson) GetStartTime() string {
	return r.StartTime
}

func (r ReservedLesson) GetEndTime() string {
	return r.EndTime
}

func (r ReservedLesson) GetSubject() string {
	return "" // Возвращаем пустую строку, так как для зарезервированного урока предмет не определен
}

func (r ReservedLesson) GetGroups() string {
	return "" // Возвращаем пустую строку, так как для зарезервированного урока группы не определены
}

func (r ReservedLesson) GetType() string {
	return "" // Возвращаем пустую строку, так как для зарезервированного урока тип не определен
}

func (r ReservedLesson) GetReserverName() string {
	return r.ReserverName
}

func (r ReservedLesson) GetReserverId() int {
	return r.ReserverId
}

func (r ReservedLesson) GetComment() string {
	return r.Comment
}
