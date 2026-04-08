package model

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Class{},
		&ClassStudent{},
	)

}
