package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"column:username;not null;type:varchar(32)"`
	Email    string `gorm:"column:email;not null;type:varchar(32)"`
	Password string `gorm:"column:password;not null;type:varchar(100)"`
	Role     string `gorm:"column:role;not null;type:enum('teacher','student');default:'student'"`
}

func (User) TableName() string {
	return "users"
}
