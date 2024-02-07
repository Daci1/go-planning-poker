package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go-planning-poker/handler"
	"log"
	"os"
)

func main() {
	loadDotEnv()

	port := fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT"))
	app := echo.New()
	app.GET("/ws", handler.GameWSHandler)
	app.Logger.Fatal(app.Start(port))
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}
