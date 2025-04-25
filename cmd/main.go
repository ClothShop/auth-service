package main

import (
	"github.com/ClothShop/auth-service/internal/config"
	"github.com/ClothShop/auth-service/internal/config/db"
	"github.com/ClothShop/auth-service/internal/routes"
	"log"
	"os"
)

func main() {
	config.LoadEnv()
	db.InitDB()
	r := routes.SetupAuthRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 Server started on: ", port)
	log.Fatal(r.Run(":" + port))
}
