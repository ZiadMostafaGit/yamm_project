package service

import (
	"errors"
	"os"
	"time"
	"yamm-project/app/internal/models"
	"yamm-project/app/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	UserID  uint   `json:"user_id"`
	Role    string `json:"role"`
	StoreID uint   `json:"store_id,omitempty"`
	jwt.RegisteredClaims
}

func (a *authService) generateJWT(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "P1cKbmLREFCTs2uioP4o8jjkNqs4h4sqgwlL/YsMA4s="
	}

	var storeID uint
	if user.Store != nil {
		storeID = user.Store.ID
	}

	claims := &Claims{
		UserID:  user.ID,
		Role:    user.Role,
		StoreID: storeID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type AuthService interface {
	Register(user *models.User, storeName string) (string, error)
	LogIn(email string, password string) (*string, error)
	RegisterAdmin(user *models.User) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	storeRepo repository.StoreRepo
	db        *gorm.DB
}

func NewAuthService(userRepo repository.UserRepository, storeRepo repository.StoreRepo, db *gorm.DB) AuthService {
	return &authService{
		userRepo:  userRepo,
		storeRepo: storeRepo,
		db:        db,
	}
}

func (a *authService) Register(user *models.User, storeName string) (string, error) {
	_, err := a.userRepo.GetByEmail(user.Email)
	if err == nil {
		return "", errors.New("user already registered")
	}
	if user.Role != "customer" && user.Role != "merchant" && user.Role != "admin" {
		return "", errors.New("invalid role")
	}
	hashedPassword, newErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if newErr != nil {
		return "", newErr
	}
	user.Password = string(hashedPassword)

	transactionErr := a.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if user.Role == "merchant" {
			store := &models.Store{
				UserID: user.ID,
				Name:   storeName,
			}
			if err := tx.Create(store).Error; err != nil {
				return err
			}
			user.Store = store
		}

		return nil
	})

	if transactionErr != nil {
		return "", transactionErr
	}

	token, tokenErr := a.generateJWT(user)
	if tokenErr != nil {
		return "", tokenErr
	}

	return token, nil
}
func (a *authService) RegisterAdmin(user *models.User) (string, error) {
	_, err := a.userRepo.GetByEmail(user.Email)
	if err == nil {
		return "", errors.New("user already registered")
	}

	hashedPassword, newErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if newErr != nil {
		return "", newErr
	}
	user.Password = string(hashedPassword)
	user.Role = "admin"

	otherErr := a.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(user).Error
		if err != nil {
			return err
		}
		return nil
	})

	if otherErr != nil {
		return "", otherErr
	}

	token, tokenErr := a.generateJWT(user)
	if tokenErr != nil {
		return "", tokenErr
	}

	return token, nil
}

func (a authService) LogIn(email, password string) (*string, error) {
	user, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	isValidUser := CheckPasswordHash(password, user.Password)
	if !isValidUser {
		return nil, errors.New("invalid user name or password")
	}

	jwt, err := a.generateJWT(user)
	if err != nil {
		return nil, err
	}
	return &jwt, nil

}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
