package repository

import (
	"yamm-project/app/internal/models"

	"gorm.io/gorm"
)

type TranslationRepository interface {
	Create(translations []models.Translation) error
	GetById(id uint) (*models.Translation, error)
	GetByFAQId(faqId uint) ([]models.Translation, error)
	Update(faqId uint, lang string, question string, answer string) error
	Delete(id uint) error
}

type translationRepository struct {
	db *gorm.DB
}

func NewTranslationRepository(db *gorm.DB) TranslationRepository {
	return &translationRepository{db: db}
}

func (r *translationRepository) Create(translations []models.Translation) error {
	err := r.db.Create(&translations).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *translationRepository) GetById(id uint) (*models.Translation, error) {
	var translation models.Translation
	err := r.db.First(&translation, id).Error
	if err != nil {
		return nil, err
	}
	return &translation, nil
}

func (r *translationRepository) GetByFAQId(faqId uint) ([]models.Translation, error) {
	var translations []models.Translation
	err := r.db.Where("faq_id = ?", faqId).Find(&translations).Error
	if err != nil {
		return nil, err
	}
	return translations, nil
}

func (r *translationRepository) Update(id uint, lang string, question string, answer string) error {
	err := r.db.Model(&models.Translation{}).Where("id = ? AND language_code = ?", id, lang).Updates(map[string]interface{}{
		"question": question,
		"answer":   answer,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *translationRepository) Delete(id uint) error {
	err := r.db.Delete(&models.Translation{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
