package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	ClassName  string `gorm:"column:class_name;type:varchar(32);not null"`
	TeacherID  uint   `gorm:"column:teacher_id;not null;index"`
	InviteCode string `gorm:"column:invite_code;type:varchar(100);not null;unique"`

	Teacher User `gorm:"foreignKey:TeacherID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (Class) TableName() string {
	return "classes"
}

// ClassStudent 学生班级关联表
type ClassStudent struct {
	gorm.Model
	ClassID   uint `gorm:"column:class_id;not null;index:idx_class_student,unique"`
	StudentID uint `gorm:"column:student_id;not null;index:idx_class_student,unique"`

	Class   Class `gorm:"foreignKey:ClassID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Student User  `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ClassStudent) TableName() string {
	return "class_students"
}
