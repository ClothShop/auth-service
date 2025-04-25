package controllers

import (
	"github.com/ClothShop/auth-service/internal/services"
	"github.com/ClothShop/auth-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetMe(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access token not found"})
		return
	}

	claims, err := utils.ValidateAccessToken(accessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := services.GetUser(uuid.MustParse(claims.UserID))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": "User info retrieved successfully",
	})
}
