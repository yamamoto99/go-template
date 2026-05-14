package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/yamamoto99/go-template/app/infrastructure/db"
	"github.com/yamamoto99/go-template/app/internal/handler"
	"github.com/yamamoto99/go-template/app/internal/repository"
	"github.com/yamamoto99/go-template/app/internal/router"
	"github.com/yamamoto99/go-template/app/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	dbConnection := db.NewDB()
	welcomeRepository := repository.NewWelcomeRepository(dbConnection)
	welcomeUsecase := usecase.NewWelcomeUsecase(welcomeRepository)
	welcomeHandler := handler.NewWelcomeHandler(welcomeUsecase)
	e := router.NewRouter(welcomeHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
