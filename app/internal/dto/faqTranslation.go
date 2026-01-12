package dto

type CreateFAQRequest struct {
	CategoryID   uint                     `json:"category_id" binding:"required"`
	StoreID      *uint                    `json:"store_id"`
	IsGlobal     bool                     `json:"is_global"`
	Translations []TranslationItemRequest `json:"translations" binding:"required,min=1,dive"`
}

type TranslationItemRequest struct {
	ID           *uint  `json:"id"`
	LanguageCode string `json:"language_code" binding:"required,oneof=EN AR"`
	Question     string `json:"question" binding:"required,min=1"`
	Answer       string `json:"answer" binding:"required,min=1"`
}

type UpdateFAQRequest struct {
	IsGlobal bool `json:"is_global"`
}

type UpdateTranslationsRequest struct {
	Translations []TranslationUpdateItem `json:"translations" binding:"required,min=1,dive"`
}

type TranslationUpdateItem struct {
	ID           *uint  `json:"id"`
	LanguageCode string `json:"language_code" binding:"required,oneof=EN AR"`
	Question     string `json:"question" binding:"required,min=1"`
	Answer       string `json:"answer" binding:"required,min=1"`
}

type FAQResponse struct {
	ID           uint                  `json:"id"`
	CategoryID   uint                  `json:"category_id"`
	CategoryName string                `json:"category_name"`
	StoreID      *uint                 `json:"store_id,omitempty"`
	StoreName    string                `json:"store_name,omitempty"`
	IsGlobal     bool                  `json:"is_global"`
	CreatedAt    string                `json:"created_at"`
	Translations []TranslationResponse `json:"translations"`
}

type TranslationResponse struct {
	ID           uint   `json:"id"`
	LanguageCode string `json:"language_code"`
	Question     string `json:"question"`
	Answer       string `json:"answer"`
}

type GetFAQsRequest struct {
	LanguageCode string `form:"lang"`
	CategoryID   *uint  `form:"category_id"`
}

type GroupedFAQResponse struct {
	CategoryName string        `json:"category_name"`
	FAQs         []FAQResponse `json:"faqs"`
}
