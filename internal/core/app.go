package core

import (
	"fmt"

	"neonexcore/internal/config"
	"neonexcore/pkg/database"
	"neonexcore/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// -----------------------------------------------------------
// 1) App Struct
// -----------------------------------------------------------
type App struct {
	Registry  *ModuleRegistry
	Container *Container
	Migrator  *database.Migrator
	Logger    logger.Logger
}

// -----------------------------------------------------------
// 2) NewApp() - สร้าง App + โหลด ModuleRegistry
// -----------------------------------------------------------
func NewApp() *App {
	return &App{
		Registry:  NewModuleRegistry(),
		Container: NewContainer(),
		Logger:    logger.NewLogger(),
	}
}

// -----------------------------------------------------------
// 3) InitLogger() - Initialize Logger
// -----------------------------------------------------------
func (a *App) InitLogger(cfg logger.Config) error {
	if err := logger.Setup(cfg); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	a.Logger = logger.NewLogger()
	a.Logger.Info("Logger initialized", logger.Fields{
		"level":  cfg.Level,
		"format": cfg.Format,
		"output": cfg.Output,
	})
	return nil
}

// -----------------------------------------------------------
// 4) InitDatabase() - เริ่ม Database + Migrator
// -----------------------------------------------------------
func (a *App) InitDatabase() error {
	dbConfig := config.LoadDatabaseConfig()
	_, err := config.InitDatabase(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize migrator
	a.Migrator = database.NewMigrator(config.DB.GetDB())
	a.Logger.Info("Database initialized", logger.Fields{"driver": dbConfig.Driver})

	return nil
}

// -----------------------------------------------------------
// 5) RegisterModels() - Register models for auto-migration
// -----------------------------------------------------------
func (a *App) RegisterModels(models ...interface{}) {
	if a.Migrator != nil {
		a.Migrator.RegisterModels(models...)
		a.Logger.Info("Models registered for migration", logger.Fields{"count": len(models)})
	}
}

// -----------------------------------------------------------
// 6) AutoMigrate() - Run auto-migration
// -----------------------------------------------------------
func (a *App) AutoMigrate() error {
	if a.Migrator != nil {
		a.Logger.Info("Running auto-migration...")
		if err := a.Migrator.AutoMigrate(); err != nil {
			a.Logger.Error("Auto-migration failed", logger.Fields{"error": err.Error()})
			return err
		}
		a.Logger.Info("Auto-migration completed")
	}
	return nil
}

// -----------------------------------------------------------
// 7) Boot() - เริ่มระบบพื้นฐาน
// -----------------------------------------------------------
func (a *App) Boot() {
	fmt.Println("⚙️  Booting Neonex Core...")
	a.Logger.Info("Neonex Core booting...")
}

// -----------------------------------------------------------
// 8) StartHTTP() - HTTP Server Engine
// -----------------------------------------------------------
func (a *App) StartHTTP() {
	// Configure Fiber with custom branding
	app := fiber.New(fiber.Config{
		AppName:               "Neonex Core v0.1-alpha",
		DisableStartupMessage: true, // Disable default Fiber banner
	})

	// Add logger middleware
	app.Use(logger.RequestIDMiddleware(a.Logger))
	app.Use(logger.HTTPMiddleware(a.Logger))

	// โหลด routes จากทุก module
	a.Logger.Info("Registering modules...")
	a.Registry.RegisterModuleServices(a.Container)
	a.Registry.LoadRoutes(app, a.Container)

	// default homepage
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"framework": "Neonex Core",
			"version":   "0.1-alpha",
			"status":    "running",
			"engine":    "Fiber (fasthttp)",
		})
	})

	// Custom Neonex startup banner
	fmt.Println()
	fmt.Println("┌───────────────────────────────────────────────────┐")
	fmt.Println("│              Neonex Core v0.1-alpha               │")
	fmt.Println("│               http://127.0.0.1:8080               │")
	fmt.Println("│       (bound on host 0.0.0.0 and port 8080)       │")
	fmt.Println("│                                                   │")
	fmt.Println("│ Framework .... Neonex  Engine ..... Fiber/fasthttp│")
	fmt.Println("└───────────────────────────────────────────────────┘")
	fmt.Println()

	a.Logger.Info("HTTP server starting", logger.Fields{"port": 8080})
	if err := app.Listen(":8080"); err != nil {
		a.Logger.Fatal("Failed to start server", logger.Fields{"error": err.Error()})
	}
}
