package models

import "time"

type Store struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	UserID    uint   `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time

	User User `json:"-"`
}

func (Store) TableName() string {
	return "stores"
}
