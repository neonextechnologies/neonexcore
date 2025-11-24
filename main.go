package main

import (
	"context"
	"fmt"
	"log"

	"neonexcore/internal/config"
	"neonexcore/internal/core"
	"neonexcore/modules/user"
	"neonexcore/pkg/database"
	"neonexcore/pkg/logger"
)

func main() {
	fmt.Println("Neonex Core v0.1 starting...")

	// Register module factories
	core.ModuleMap["user"] = func() core.Module { return user.New() }

	app := core.NewApp()

	// Initialize Logger
	loggerConfig := logger.LoadConfig()
	if err := app.InitLogger(loggerConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize Database
	if err := app.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Register models for auto-migration
	app.RegisterModels(
		&user.User{},
	)

	// Run auto-migration
	if err := app.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed database (optional)
	seeder := database.NewSeederManager(config.DB.GetDB())
	seeder.Register(&user.UserSeeder{})
	if err := seeder.Run(context.Background()); err != nil {
		log.Printf("Warning: Seeding failed: %v", err)
	}

	// Load modules
	app.Registry.AutoDiscover()
	app.Boot()
	app.Registry.Load()

	// Start HTTP server
	app.StartHTTP()
}
