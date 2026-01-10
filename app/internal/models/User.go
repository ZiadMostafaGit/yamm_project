package models

import "time"

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Email     string    `gorm:"size:255;uniqueindex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"size:20;not null" json:"role"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Age       uint8     `gorm:"not null" json:"age"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Store     *Store    `gorm:"constraint:onDelete:CASCADE" json:"store,omitempty"`
}

func (User) TableName() string {
	return "users"
}
