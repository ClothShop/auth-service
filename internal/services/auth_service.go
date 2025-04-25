package services

import (
	"errors"
	"github.com/ClothShop/auth-service/internal/dtos"
	"github.com/ClothShop/auth-service/internal/mapper"
	"github.com/ClothShop/auth-service/internal/models"
	"github.com/ClothShop/auth-service/internal/repositories"
	"github.com/ClothShop/auth-service/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

func Register(email, password string) (*models.User, error) {
	_, err := repositories.FindByEmail(email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	hashed, _ := utils.HashPassword(password)
	return repositories.CreateUser(email, hashed)
}

func Login(email, password string) (access string, refresh string, err error) {
	user, err := repositories.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, _ := utils.GenerateAccessToken(user.ID.String(), string(user.Role))
	refreshToken, _ := utils.GenerateRefreshToken(user.ID.String())

	err = repositories.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}
	log.Println("access token:", accessToken, "refresh token:", refreshToken)
	return accessToken, refreshToken, nil
}

func Logout(refreshToken string) error {
	return repositories.DeleteRefreshToken(refreshToken)
}

func FindOrCreateFromGoogle(googleUser *models.User) (*models.User, error) {
	user, err := repositories.FindByEmail(googleUser.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &models.User{
			ID:              uuid.New(),
			Email:           googleUser.Email,
			Name:            googleUser.Name,
			Role:            "User",
			IsEmailVerified: true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		err = repositories.SaveUser(user)
	}
	return user, err
}

func GetUser(id uuid.UUID) (*dtos.MeResponse, error) {
	user, err := repositories.FindByID(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToMeResponse(user), nil
}
