package models

type UserModel struct {
	ID       string `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Phone    string `json:"phone" gorm:"column:phone"`
}
