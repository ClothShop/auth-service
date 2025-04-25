package routes

import (
	"fmt"
	"github.com/ClothShop/auth-service/internal/controllers"
	"github.com/ClothShop/auth-service/internal/middlewares"
	"github.com/gin-gonic/gin"
	"os"
)

func SetupAuthRoutes() *gin.Engine {
	apiVersion := os.Getenv("API_VERSION")
	baseURL := fmt.Sprintf("/api/%s/auth", apiVersion)

	r := gin.Default()

	authRoute := r.Group(baseURL)
	{
		authRoute.POST("/register", middlewares.AuthReqMiddleware(), controllers.Register)
		authRoute.POST("/login", middlewares.AuthReqMiddleware(), controllers.Login)
		authRoute.GET("/refresh-token", controllers.HandleRefreshToken)
		authRoute.POST("/logout", controllers.Logout)
	}
	oauthGoogleRoute := r.Group(baseURL + "/oauth/google")
	{
		oauthGoogleRoute.GET("/", controllers.GoogleLogin)
		oauthGoogleRoute.GET("/callback", controllers.GoogleCallback)
	}
	userRoutes := r.Group(baseURL + "/users")
	{
		userRoutes.GET("/me", controllers.GetMe)
	}

	return r
}
