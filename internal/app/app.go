package app

import (
	"fmt"
	"log"

	"github.com/Go11Group/url_shortner/config"
	v1 "github.com/Go11Group/url_shortner/internal/controller/http/v1"
	"github.com/Go11Group/url_shortner/internal/repo/storage"
	"github.com/Go11Group/url_shortner/internal/usecase"
	"github.com/Go11Group/url_shortner/pkg/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/gofiber/swagger"

	_ "github.com/Go11Group/url_shortner/api/openapi"
)

// Run initializes and starts the application.
// @title URL Shortener API
// @version 1.0
// @description This is a simple URL shortener API.
// @host localhost:8080
func Run(cfg *config.Config) {
	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	storageLayer := storage.NewStorage(pg)

	// UseCase
	urlUseCase := usecase.NewUrlUseCase(storageLayer)

	// HTTP Server
	app := fiber.New()
	app.Use(cors.New())

	// Controller
	urlController := v1.NewController(urlUseCase)

	// Routes
	api := app.Group("/api")
	api.Post("/shorten", urlController.Shorten)
	app.Get("/:code", urlController.Redirect)
	app.Get("/swagger/*", fiberSwagger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("URL Shortener API is running!")
	})

	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.HTTP.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.HTTP.Port)))
}
