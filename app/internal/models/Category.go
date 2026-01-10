package models

import "time"

type Category struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time

	FAQs []FAQ
}

func (Category) TableName() string {
	return "categories"
}
