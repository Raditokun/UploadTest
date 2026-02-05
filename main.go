package main

import (
	"log"

	"upload/config"
	"upload/repositories"
	"upload/routes"
	"upload/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.ConnectDB()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/uploads", "./uploads")

	fileRepo := repositories.NewFileRepository(config.DB)
	fileService := services.NewFileService(fileRepo, "./uploads", "http://localhost:3000")
	fileHandler := routes.NewFileHandler(fileService)

	fileHandler.RegisterRoutes(app)

	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
