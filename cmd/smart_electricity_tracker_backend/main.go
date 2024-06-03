package main

import (
	"log"
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/database"
	"smart_electricity_tracker_backend/internal/routes"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrate(db, cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOriginsFunc: nil,
		AllowOrigins:     "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))
	// app.Use(middleware.ErrorMiddleware())
	routes.Setup(app, cfg, db)

	log.Fatal(app.Listen(":8080"))
}
