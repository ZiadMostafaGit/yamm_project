package service

import (
	"errors"
	"fmt"
	"time"
	"yamm-project/app/internal/dto"
	"yamm-project/app/internal/models"
	"yamm-project/app/internal/repository"

	"gorm.io/gorm"
)

type FAQService interface {
	CreateFAQ(req *dto.CreateFAQRequest, userRole string, userStoreID *uint) (*uint, error)
	GetFAQByID(id uint) (*dto.FAQResponse, error)
	UpdateFAQVisibility(id uint, req *dto.UpdateFAQRequest, userRole string) error
	UpdateTranslations(faqID uint, req *dto.UpdateTranslationsRequest, userRole string, userStoreID *uint) error
	DeleteFAQ(id uint, userRole string, userStoreID *uint) error
	GetFAQsForCustomer(storeID *uint, languageCode string, categoryID *uint) ([]dto.FAQResponse, error)
	GetAllFAQs(languageCode string) ([]dto.FAQResponse, error)
	GetGroupedFAQs(storeID *uint, lang string) ([]dto.GroupedFAQResponse, error)
}

type faqService struct {
	faqRepo         repository.FAQRepository
	translationRepo repository.TranslationRepository
	categoryRepo    repository.CategoryRepository
	storeRepo       repository.StoreRepo
	db              *gorm.DB
}

func NewFAQService(faqRepo repository.FAQRepository, transRepo repository.TranslationRepository, categoryRepo repository.CategoryRepository, storeRepo repository.StoreRepo, db *gorm.DB) FAQService {

	return &faqService{
		faqRepo:         faqRepo,
		translationRepo: transRepo,
		categoryRepo:    categoryRepo,
		storeRepo:       storeRepo,
		db:              db,
	}
}

func (s *faqService) validateFAQAccess(faq *models.FAQ, userRole string, userStoreID *uint) error {
	if userRole == "admin" {
		return nil
	}

	if userRole == "merchant" {
		if faq.IsGlobal {
			return errors.New("merchants cannot modify global FAQs")
		}

		if faq.StoreID == nil || userStoreID == nil || *faq.StoreID != *userStoreID {
			return errors.New("unauthorized to access this FAQ")
		}

		return nil
	}

	return errors.New("insufficient permissions")
}

func (s *faqService) validateFAQCreation(req *dto.CreateFAQRequest, userRole string, userStoreID *uint) error {
	if userRole == "admin" {
		if req.IsGlobal && req.StoreID != nil {
			return errors.New("global FAQs cannot be associated with a specific store")
		}
		return nil
	}

	if userRole == "merchant" {
		if req.IsGlobal {
			return errors.New("merchants cannot create global FAQs")
		}

		if req.StoreID == nil {
			return errors.New("store_id is required for merchant FAQs")
		}

		if userStoreID == nil || *req.StoreID != *userStoreID {
			return errors.New("merchants can only create FAQs for their own store")
		}

		return nil
	}

	return errors.New("insufficient permissions to create FAQs")
}

func (s *faqService) CreateFAQ(req *dto.CreateFAQRequest, userRole string, userStoreID *uint) (*uint, error) {

	err := s.validateFAQCreation(req, userRole, userStoreID)

	if err != nil {
		return nil, err

	}

	errChan := make(chan error, 2)

	go func() {

		_, err2 := s.categoryRepo.GetById(req.CategoryID)

		if err2 != nil {
			errChan <- err2
		} else {
			errChan <- nil
		}
	}()

	go func() {
		if req.StoreID != nil {
			_, err := s.storeRepo.GetByUserid(*req.StoreID)
			if err != nil {
				errChan <- err
			}
		} else {
			errChan <- nil
		}

	}()

	for i := 0; i < 2; i++ {
		checkError := <-errChan
		if checkError != nil {
			return nil, checkError
		}
	}

	var faqID *uint
	err = s.db.Transaction(func(tx *gorm.DB) error {
		faq := &models.FAQ{
			CategoryID: req.CategoryID,
			StoreID:    req.StoreID,
			IsGlobal:   req.IsGlobal,
		}
		id, err := s.faqRepo.Create(faq)
		if err != nil {
			return fmt.Errorf("failed to create FAQ: %w", err)
		}
		faqID = id

		translations := make([]models.Translation, len(req.Translations))
		for i, transDTO := range req.Translations {
			translations[i] = models.Translation{
				FAQID:        *faqID,
				LanguageCode: transDTO.LanguageCode,
				Question:     transDTO.Question,
				Answer:       transDTO.Answer,
			}
		}

		if err := s.translationRepo.Create(translations); err != nil {
			return fmt.Errorf("failed to create translations: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return faqID, nil
}

func (s *faqService) GetFAQByID(id uint) (*dto.FAQResponse, error) {
	faq, err := s.faqRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("FAQ not found")
		}
		return nil, err
	}

	return s.mapFAQToResponse(faq), nil
}
func (s *faqService) UpdateFAQVisibility(id uint, req *dto.UpdateFAQRequest, role string) error {
	if role != "admin" {
		return errors.New("only admins can change visibility")
	}

	newErr := s.faqRepo.UpdateVisibility(id, req.IsGlobal)
	if newErr != nil {
		return newErr
	}
	return nil
}

func (s *faqService) UpdateTranslations(faqID uint, req *dto.UpdateTranslationsRequest, userRole string, userStoreID *uint) error {
	faq, err := s.faqRepo.GetByID(faqID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("FAQ not found")
		}
		return err
	}

	if err := s.validateFAQAccess(faq, userRole, userStoreID); err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, transDTO := range req.Translations {
			if transDTO.ID != nil {
				err := s.translationRepo.Update(*transDTO.ID, transDTO.LanguageCode, transDTO.Question, transDTO.Answer)
				if err != nil {
					return errors.New("failed to update translation")
				}
			} else {
				newTrans := []models.Translation{{
					FAQID:        faqID,
					LanguageCode: transDTO.LanguageCode,
					Question:     transDTO.Question,
					Answer:       transDTO.Answer,
					CreatedAt:    time.Now(),
				}}
				if err := s.translationRepo.Create(newTrans); err != nil {
					return fmt.Errorf("failed to create translation: %w", err)
				}
			}
		}
		return nil
	})
}

func (s *faqService) DeleteFAQ(id uint, userRole string, userStoreID *uint) error {
	faq, err := s.faqRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("FAQ not found")
		}
		return err
	}
	if err := s.validateFAQAccess(faq, userRole, userStoreID); err != nil {
		return err
	}
	return s.faqRepo.Delete(id)
}

func (s *faqService) GetFAQsForCustomer(storeID *uint, languageCode string, categoryID *uint) ([]dto.FAQResponse, error) {
	faqs, err := s.faqRepo.GetForCustomer(storeID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.FAQResponse, 0)
	for _, faq := range faqs {
		if categoryID != nil && faq.CategoryID != *categoryID {
			continue
		}

		response := s.mapFAQToResponse(&faq)

		if languageCode != "" {
			response.Translations = s.filterTranslationsByLanguage(response.Translations, languageCode)
		}

		responses = append(responses, *response)
	}

	return responses, nil
}

func (s *faqService) GetGroupedFAQs(storeID *uint, lang string) ([]dto.GroupedFAQResponse, error) {
	faqs, err := s.GetFAQsForCustomer(storeID, lang, nil)
	if err != nil {
		return nil, err
	}

	groups := make(map[string][]dto.FAQResponse)
	for _, f := range faqs {
		groups[f.CategoryName] = append(groups[f.CategoryName], f)
	}

	var result []dto.GroupedFAQResponse
	for name, list := range groups {
		result = append(result, dto.GroupedFAQResponse{CategoryName: name, FAQs: list})
	}
	return result, nil
}

func (s *faqService) GetAllFAQs(languageCode string) ([]dto.FAQResponse, error) {
	faqs, err := s.faqRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.FAQResponse, len(faqs))
	for i, faq := range faqs {
		response := s.mapFAQToResponse(&faq)

		if languageCode != "" {
			response.Translations = s.filterTranslationsByLanguage(response.Translations, languageCode)
		}

		responses[i] = *response
	}

	return responses, nil
}

func (s *faqService) mapFAQToResponse(faq *models.FAQ) *dto.FAQResponse {
	response := &dto.FAQResponse{
		ID:         faq.ID,
		CategoryID: faq.CategoryID,
		StoreID:    faq.StoreID,
		IsGlobal:   faq.IsGlobal,
		CreatedAt:  faq.CreatedAt.Format(time.RFC3339),
	}

	if faq.Category.ID != 0 {
		response.CategoryName = faq.Category.Name
	}

	if faq.Store != nil && faq.Store.ID != 0 {
		response.StoreName = faq.Store.Name
	}

	response.Translations = make([]dto.TranslationResponse, len(faq.Translations))
	for i, trans := range faq.Translations {
		response.Translations[i] = dto.TranslationResponse{
			ID:           trans.ID,
			LanguageCode: trans.LanguageCode,
			Question:     trans.Question,
			Answer:       trans.Answer,
		}
	}

	return response
}

func (s *faqService) filterTranslationsByLanguage(translations []dto.TranslationResponse, languageCode string) []dto.TranslationResponse {
	filtered := make([]dto.TranslationResponse, 0)
	for _, trans := range translations {
		if trans.LanguageCode == languageCode {
			filtered = append(filtered, trans)
		}
	}
	return filtered
}
