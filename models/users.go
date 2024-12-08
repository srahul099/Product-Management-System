package models

type Users struct {
	UserID   int    `json:"user_id" gorm:"primary key;autoIncrement"`
	UserName string `json:"username"`
	
}
