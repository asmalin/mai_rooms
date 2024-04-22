package models

// const (
// 	userRole  = "User"
// 	adminRole = "Admin"
// )

type User struct {
	Id       int    `json:"id"`
	Username string `json:"name" binding:"required" gorm:"size:255"`
	Password string `json:"password" binding:"required" gorm:"size:255"`
	Fullname string `json:"fullname" binding:"required" gorm:"size:255"`
	Role     string `json:"role" binding:"required" gorm:"size:50"`
}

type UserTgChatRelation struct {
	UserId   int
	TgChatId int64
	User     User `gorm:"foreignKey:UserId"`
}
