package repository

import (
	"yamm-project/app/internal/models"

	"gorm.io/gorm"
)

type StoreRepo interface {
	Create(store *models.Store) error
	GetByUserid(id uint) (*models.Store, error)
}

type storeRepo struct {
	db *gorm.DB
}

func NewStoreRepo(db *gorm.DB) StoreRepo {
	return &storeRepo{db: db}
}

func (s *storeRepo) Create(store *models.Store) error {

	err := s.db.Create(store).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *storeRepo) GetByUserid(id uint) (*models.Store, error) {

	var store models.Store
	err := s.db.Where("user_id = ?", id).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil

}

func (s *storeRepo) Update(store *models.Store) error {
	return s.db.Save(store).Error
}

func (s *storeRepo) Delete(id uint) error {
	return s.db.Delete(&models.Store{}, id).Error
}
