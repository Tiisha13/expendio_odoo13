package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"expensio-backend/internal/config"
	"expensio-backend/internal/routes"
	"expensio-backend/pkg/cache"
	"expensio-backend/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Println("✅ Configuration loaded successfully")

	// Initialize MongoDB
	if err := database.ConnectMongoDB(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer database.DisconnectMongoDB()
	log.Println("✅ Connected to MongoDB")

	// Initialize Redis
	if err := cache.InitRedis(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	defer cache.CloseRedis()
	log.Println("✅ Connected to Redis")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Expensio API v1.0.0",
		ServerHeader: "Expensio",
		ErrorHandler: customErrorHandler,
		// DisableStartupMessage: false,
	})

	// Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Expensio API is running",
		})
	})

	// Setup routes
	routes.SetupRoutes(app, cfg)

	// Create upload directories if they don't exist
	createDirectories(cfg)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("🛑 Shutting down gracefully...")
		app.Shutdown()
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

// createDirectories creates necessary directories for the application
func createDirectories(cfg *config.Config) {
	dirs := []string{
		cfg.FileUpload.UploadDir,
		cfg.OCR.TempDir,
		"./logs",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("⚠️  Warning: Failed to create directory %s: %v", dir, err)
		}
	}
}
