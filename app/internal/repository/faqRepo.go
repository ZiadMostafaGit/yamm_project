/*
this is the hardest part so far becase the relation ebtwean the translation and faq and all the operation performed on it so

	execuse my long comment its for me to remeber why i did this
*/
package repository

import (
	"errors"
	"yamm-project/app/internal/models"

	"gorm.io/gorm"
)

type FAQRepository interface {
	Create(faq *models.FAQ) (*uint, error)
	GetByID(id uint) (*models.FAQ, error)
	Update(id uint, state bool) error
	UpdateVisibility(id uint, isGlobal bool) error
	GetForCustomer(storeID *uint) ([]models.FAQ, error) //since the user get the global and store specific then i nead to query all the faq for global

	// and the one with the store id so i nead to get both all global and the one with store id in the task asking for groupby iam not sure if it means
	// that user can filter by the category or just that the order of them become based on the category
	// i think i may just order by category id and thats it
	Delete(id uint) error
	GetAll() ([]models.FAQ, error)
}

type faqRepository struct {
	db *gorm.DB
}

func NewFAQRepository(db *gorm.DB) FAQRepository {
	return &faqRepository{db: db}
}

func (r *faqRepository) Create(faq *models.FAQ) (*uint, error) {
	err := r.db.Create(&faq).Error
	if err != nil {
		return nil, err
	}
	return &faq.ID, nil
}

func (r *faqRepository) GetByID(id uint) (*models.FAQ, error) {
	var faq models.FAQ
	err := r.db.Preload("Translations").Preload("Category").First(&faq, id).Error

	/*
		the preload is not so important and dont do performance inp becase its one id no n+1 prolem but its cleaner than
		one query fpr faq and one to pobulate the tranlation slice

	*/
	return &faq, err
}

func (r *faqRepository) GetForCustomer(storeID *uint) ([]models.FAQ, error) {
	var faqs []models.FAQ
	query := r.db.Preload("Translations").Preload("Category")

	if storeID != nil {
		query = query.Where("is_global = ? OR store_id = ?", true, *storeID)
	} else {
		query = query.Where("is_global = ?", true)
	}

	err := query.Find(&faqs).Error
	return faqs, err
}

func (r *faqRepository) Delete(id uint) error {
	res := r.db.Delete(&models.FAQ{}, id)
	if res.Error == nil {
		return res.Error

	}
	if res.RowsAffected == 0 {
		return errors.New("Faq not found")
	}
	return nil
}
func (r *faqRepository) UpdateVisibility(id uint, isGlobal bool) error {
	if isGlobal == true {
		err := r.db.Model(&models.FAQ{}).Where("id = ?", id).Updates(map[string]interface{}{
			"is_global": true,
			"store_id":  nil,
		}).Error
		if err != nil {
			return err
		}
		return nil
	}

	err := r.db.Model(&models.FAQ{}).Where("id = ?", id).Update("is_global", false).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *faqRepository) Update(id uint, state bool) error {

	/*
	   at this part in the task "Create, update, delete FAQs under these categories"
	   i think this means that the faq category should not be chened like change category itself in the table but cant change categoryid in the faq and think this is the
	   correct way



	   so what can be chenged the isGlobal only becase at this db desgin changing the faq quetion or answer or both is translation matter
	*/
	err := r.db.Model(&models.FAQ{}).Where("id = ? ", id).Update("is_global", state).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *faqRepository) GetAll() ([]models.FAQ, error) {
	var Faqs []models.FAQ
	err := r.db.Preload("Translations").Preload("Category").Find(&Faqs).Error
	if err != nil {
		return nil, err
	}
	return Faqs, nil

}
