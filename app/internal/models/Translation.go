package models

import "time"

type Translation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FAQID        uint      `gorm:"not null" json:"faq_id"`
	LanguageCode string    `gorm:"size:10;default:AR" json:"language_code"`
	Question     string    `gorm:"not null" json:"question"`
	Answer       string    `gorm:"not null" json:"answer"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	FAQ          FAQ       `json:"-"`
}

func (Translation) TableName() string {
	return "translations"
}
