package models

import (
	"time"
)

type Building struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Room struct {
	ID          int    `json:"id"`
	Building_id int    `json:"building_id" gorm:"foreignKey:Building.ID"`
	Name        string `json:"name"`
}

type Lesson struct {
	ID        int       `json:"-"`
	Lector    string    `json:"lector"`
	Date      time.Time `json:"date" gorm:"type: date"`
	TimeStart string    `json:"time_start" gorm:"type: timestamp without time zone"`
	TimeEnd   string    `json:"time_end" gorm:"type: timestamp without time zone"`
	Name      string    `json:"name"`
	Groups    string    `json:"groups"`
	Type      string    `json:"type"`
	RoomId    int       `json:"room"`
	Room      Room      `json:"-" gorm:"foreignKey:RoomId;"`
}

type ReservedLesson struct {
	ID        int       `json:"-"`
	User_id   int       `json:"user_id"`
	Room_id   int       `json:"room_id"`
	Date      time.Time `json:"date" gorm:"type: date"`
	TimeStart string    `json:"time_start" gorm:"type: time without time zone"`
	TimeEnd   string    `json:"time_end" gorm:"type: time without time zone"`
	Comment   string    `json:"comment"`
	User      User      `json:"-" gorm:"foreignKey:User_id;"`
	Room      Room      `json:"-" gorm:"foreignKey:Room_id;"`
}
