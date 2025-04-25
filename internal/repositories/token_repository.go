package repositories

import (
	"github.com/ClothShop/auth-service/internal/config/db"
	"github.com/ClothShop/auth-service/internal/models"
	"github.com/google/uuid"
	"time"
)

func SaveRefreshToken(userID uuid.UUID, token string) error {
	rt := models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	return db.DB.Create(&rt).Error
}

func DeleteRefreshToken(token string) error {
	return db.DB.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func Save(userID uuid.UUID, token string) error {
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return db.DB.Create(&refreshToken).Error
}

func FindByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := db.DB.Where("token = ?", token).First(&rt).Error
	return &rt, err
}
