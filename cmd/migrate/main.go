package main

import (
	"fmt"
	"os"

	"github.com/alireza-akbarzadeh/ginflow/internal/config"
	"github.com/alireza-akbarzadeh/ginflow/internal/console"
	"github.com/alireza-akbarzadeh/ginflow/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	c := console.New()

	c.Line()
	c.Info("🔄", "GinFlow Database Migration")
	c.Line()

	// Load database URL from environment
	dbUrl := config.GetEnvString("DATABASE_URL", "")
	if dbUrl == "" {
		c.Error("❌", "DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	// Connect to database
	c.Info("🗄️ ", "Connecting to database...")
	db, err := database.Connect(dbUrl)
	if err != nil {
		c.Error("❌", fmt.Sprintf("Database connection failed: %v", err))
		os.Exit(1)
	}
	c.Success("✅", "Database connected")

	// Close connection when done
	sqlDB, err := db.DB()
	if err != nil {
		c.Error("❌", fmt.Sprintf("Failed to get SQL DB: %v", err))
		os.Exit(1)
	}
	defer sqlDB.Close()

	// Run migrations
	c.Info("📦", "Running database migrations...")
	if err := database.Migrate(db); err != nil {
		c.Error("❌", fmt.Sprintf("Migration failed: %v", err))
		os.Exit(1)
	}

	c.Line()
	c.Success("✅", "Migrations completed successfully!")
	c.Line()
}
