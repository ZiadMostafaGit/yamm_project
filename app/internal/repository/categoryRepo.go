package repository

import (
	"errors"
	"yamm-project/app/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	GetById(id uint) (*models.Category, error)
	Update(id uint, category string) error
	Delete(id uint) error
	GetAll() ([]models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (c *categoryRepository) Create(category *models.Category) error {
	err := c.db.Create(category).Error
	if err != nil {
		return err
	}
	return nil

}

func (c *categoryRepository) GetById(id uint) (*models.Category, error) {
	var category models.Category
	err := c.db.Preload("FAQs").Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *categoryRepository) Update(id uint, category string) error {

	res := c.db.Model(&models.Category{}).Where("id = ? ", id).Update("name", category)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("invalid operation")
	}
	return nil

}

func (c *categoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := c.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *categoryRepository) Delete(id uint) error {
	res := c.db.Delete(&models.Category{}, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}
