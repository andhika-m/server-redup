package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"redup/internal/api/v1/rest"
	"redup/internal/database"
	"redup/internal/middleware"

	bRepo "redup/internal/repository/bookmark"
	hRepo "redup/internal/repository/history"
	uRepo "redup/internal/repository/user"
	vRepo "redup/internal/repository/video"

	rUsecase "redup/internal/usecase/redup"
)

func main() {
	e := echo.New()

	e.Static("/public", "public")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	secret := "AES256Key-32Characters1234567890"
	signKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	storagePath := "./public/upload/videos"

	videoRepo := vRepo.GetRepository(db, storagePath)
	userRepo, err := uRepo.GetRepository(db, secret, 1, 64*1024, 4, 32, signKey, 60*time.Hour)
	if err != nil {
		panic(err)
	}
	
	bookmarkRepo := bRepo.GetRepository(db)
	historyRepo := hRepo.GetRepository(db)

	redupUsecase := rUsecase.RedupUsecase(videoRepo, userRepo, bookmarkRepo, historyRepo)
	redupHandler := rest.RedupHandler(redupUsecase)

	middleware.LoadMiddlewares(e)
	rest.InitRoutes(e, redupHandler)

	e.Start(":8080")
}
