package controllers

import (
	"github.com/ClothShop/auth-service/internal/dtos"
	"github.com/ClothShop/auth-service/internal/repositories"
	"github.com/ClothShop/auth-service/internal/services"
	"github.com/ClothShop/auth-service/internal/services/oauth"
	"github.com/ClothShop/auth-service/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Register(ctx *gin.Context) {
	userRegisterReqAny, exists := ctx.Get("validatedUserReq")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "❌ Invalid register req",
		})
		return
	}
	userRegisterReq := userRegisterReqAny.(dtos.AuthReq)
	user, err := services.Register(userRegisterReq.Email, userRegisterReq.Password)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

func Login(ctx *gin.Context) {
	userLoginReqAny, exists := ctx.Get("validatedUserReq")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "❌ Invalid login req",
		})
		return
	}
	userLoginReq := userLoginReqAny.(dtos.AuthReq)
	access, refresh, err := services.Login(userLoginReq.Email, userLoginReq.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setAuthCookies(ctx, access, refresh)
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func GoogleLogin(ctx *gin.Context) {
	authURL, err := oauth.GetGoogleAuthURL("state-123")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auth url"})
		return
	}

	ctx.Redirect(http.StatusFound, authURL)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is invalid"})
		return
	}
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state is invalid"})
	}

	googleUser, err := oauth.ExchangeCodeForUser(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code"})
		return
	}

	user, err := services.FindOrCreateFromGoogle(googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find or create user"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	if err := repositories.Save(user.ID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store refresh token"})
		return
	}

	setAuthCookies(c, accessToken, refreshToken)
	c.Redirect(http.StatusFound, "http://localhost:3000/login")
}

func HandleRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenRecord, err := repositories.FindByToken(refreshToken)
	if err != nil || tokenRecord.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "expired or invalid refresh token"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	setAuthCookies(c, accessToken, refreshToken)
	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

func setAuthCookies(c *gin.Context, access, refresh string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "access_token",
		Value: access,
		Path:  "/",
		//Domain:   os.Getenv("DOMAIN"),
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "refresh_token",
		Value: refresh,
		Path:  "/",
		//Domain:   os.Getenv("DOMAIN"),
		MaxAge:   7 * 24 * 3600,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func Logout(c *gin.Context) {
	token, _ := c.Cookie("access_token")
	go services.Logout(token)
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
