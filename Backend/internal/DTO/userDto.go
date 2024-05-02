package dto

type UserDto struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}
