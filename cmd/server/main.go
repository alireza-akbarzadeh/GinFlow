package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/alireza-akbarzadeh/ginflow/docs"
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/alireza-akbarzadeh/ginflow/internal/api/routers"
	"github.com/alireza-akbarzadeh/ginflow/internal/config"
	"github.com/alireza-akbarzadeh/ginflow/internal/console"
	"github.com/alireza-akbarzadeh/ginflow/internal/database"
	"github.com/alireza-akbarzadeh/ginflow/internal/repository"
	_ "github.com/joho/godotenv/autoload"
)

// @title Go Gin REST API
// @version 1.0
// @description This is a REST API server for managing events and attendees.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize console for colored output
	console := console.New()

	// 0. Setup Structured Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	console.Line()
	console.Info("🔧", "Starting GinFlow API Server...")
	console.Line()

	// 1. Load Configuration
	console.Info("⚙️ ", "Loading configuration...")
	port := config.GetEnvInt("PORT", 8080)
	jwtSecret := config.GetEnvString("JWT_SECRET", "some-secret-123456")
	dbUrl := config.GetEnvString("DATABASE_URL", "")
	console.Success("✅", "Configuration loaded")

	// 2. Initialize Database
	console.Info("🗄️ ", "Connecting to database...")
	db, err := database.Connect(dbUrl)
	if err != nil {
		console.Error("❌", fmt.Sprintf("Database connection failed: %v", err))
		os.Exit(1)
	}
	console.Success("✅", "Database connected")

	// Get underlying SQL DB for graceful shutdown
	sqlDB, err := db.DB()
	if err != nil {
		console.Error("❌", fmt.Sprintf("Failed to get SQL DB: %v", err))
		os.Exit(1)
	}
	defer func() {
		console.Info("🔌", "Closing database connection...")
		if err := sqlDB.Close(); err != nil {
			console.Error("❌", fmt.Sprintf("Error closing database: %v", err))
		}
		console.Success("✅", "Database connection closed")
	}()

	// 3. Initialize Dependencies
	console.Info("🔗", "Initializing dependencies...")
	repos := repository.NewModels(db)
	handler := handlers.NewHandler(repos, jwtSecret)
	router := routers.SetupRouter(handler, jwtSecret, repos.Users)
	console.Success("✅", "Dependencies initialized")

	// 4. Configure Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 5. Start Server in a Goroutine
	go func() {
		console.Line()
		console.Divider()
		console.Success("🚀", fmt.Sprintf("Server running on port %d", port))
		console.Divider()
		console.Line()
		console.URL("📚", "Swagger UI", fmt.Sprintf("http://localhost:%d/swagger/index.html", port))
		console.URL("🏥", "Health Check", fmt.Sprintf("http://localhost:%d/health", port))
		console.URL("🌐", "API Base", fmt.Sprintf("http://localhost:%d/api/v1", port))
		console.Line()
		console.Info("💡", "Press Ctrl+C to stop the server")
		console.Line()

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			console.Error("❌", fmt.Sprintf("Failed to start server: %v", err))
			os.Exit(1)
		}
	}()

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	console.Line()
	console.Warning("⚠️ ", "Shutdown signal received...")
	console.Info("🛑", "Gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		console.Error("❌", fmt.Sprintf("Server forced to shutdown: %v", err))
		os.Exit(1)
	}

	console.Line()
	console.Success("👋", "Server stopped gracefully. Goodbye!")
	console.Line()
}
