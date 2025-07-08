package models

type Student struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"not null"`
	RollNumber string `json:"roll_number" gorm:"unique;not null"`
	ClassID    uint   `json:"class_id" gorm:"not null"`
	Class      Class  `json:"class" gorm:"foreignKey:ClassID"`
}