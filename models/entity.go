package models

type UserCredential struct {
	ID       string `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

type UserModel struct {
	ID       string `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Phone    string `json:"phone" gorm:"column:phone"`
}