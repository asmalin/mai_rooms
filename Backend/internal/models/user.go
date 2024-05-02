package models

type User struct {
	Id              int              `json:"id"`
	Username        string           `json:"username" binding:"required" gorm:"size:255"`
	Password        string           `json:"password" binding:"required" gorm:"size:255"`
	Fullname        string           `json:"fullname" binding:"required" gorm:"size:255"`
	Role            string           `json:"role" binding:"required" gorm:"size:50"`
	Email           string           `json:"email" gorm:"size:255"`
	ReservedLessons []ReservedLesson `json:"-"`
}

type UserTgChatRelation struct {
	UserId   int
	TgChatId int64
	User     User `gorm:"foreignKey:UserId"`
}
