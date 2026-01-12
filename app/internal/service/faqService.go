
import (
	"errors"
	"yamm-project/app/internal/models"
	"yamm-project/app/internal/repository"
	"gorm.io/gorm"
)

type FAQService interface {
	CreateFAQ(categoryID uint, isGlobal bool, storeID *uint, role string, translations []*models.Translation) error
	UpdateFAQVisibility(id uint, isGlobal bool, role string) error
	DeleteFAQ(id uint) error
	GetCustomerFAQs(storeID *uint) ([]models.FAQ, error)
}

type faqService struct {
	faqRepo         repository.FAQRepository
	translationRepo repository.TranslationRepository
	db              *gorm.DB
}


func NewFAQService(FaqRepo repository.FAQRepository,transRepo repository.TranslationRepository,db *gorm.DB)FAQService{
	return &faqService{
		faqRepo:FaqRepo,
		translationRepo,transRepo,
		db:db
	}
}



func (f *faqService)CreateFAQ(categoryID uint, isGlobal bool, storeID *uint, role string, translations []*models.Translation){

faq := &models.FAQ{
			CategoryID: categoryID,
		}
	if role=="admin"{
		faq.isGlobal=isGlobal
		feq.storeID=nil

		
	}else{
		feq.isGlobal=false
		faq.storeID=storeID
	}
	ID,err:=f.faqRepo.Create(faq)

}

