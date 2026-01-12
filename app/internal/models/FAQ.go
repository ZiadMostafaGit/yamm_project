package models

import "time"

type FAQ struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `gorm:"not null" json:"category_id"`
	StoreID    *uint     `json:"store_id"`
	IsGlobal   bool      `json:"is_global"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	Category     Category      `json:"-"`
	Store        *Store        `json:"store,omitempty"`
	Translations []Translation `gorm:"foreignKey:FAQID;constraint:OnDelete:CASCADE" json:"translations,omitempty"`
}

func (FAQ) TableName() string {
	return "faqs"
}
