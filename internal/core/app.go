package core

import (
	"fmt"

	"neonexcore/internal/config"
	"neonexcore/pkg/database"

	"github.com/gofiber/fiber/v2"
)

// -----------------------------------------------------------
// 1) App Struct
// -----------------------------------------------------------
type App struct {
	Registry  *ModuleRegistry
	Container *Container
	Migrator  *database.Migrator
}

// -----------------------------------------------------------
// 2) NewApp() - สร้าง App + โหลด ModuleRegistry
// -----------------------------------------------------------
func NewApp() *App {
	return &App{
		Registry:  NewModuleRegistry(),
		Container: NewContainer(),
	}
}

// -----------------------------------------------------------
// 3) InitDatabase() - เริ่ม Database + Migrator
// -----------------------------------------------------------
func (a *App) InitDatabase() error {
	dbConfig := config.LoadDatabaseConfig()
	_, err := config.InitDatabase(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize migrator
	a.Migrator = database.NewMigrator(config.DB.GetDB())

	return nil
}

// -----------------------------------------------------------
// 4) RegisterModels() - Register models for auto-migration
// -----------------------------------------------------------
func (a *App) RegisterModels(models ...interface{}) {
	if a.Migrator != nil {
		a.Migrator.RegisterModels(models...)
	}
}

// -----------------------------------------------------------
// 5) AutoMigrate() - Run auto-migration
// -----------------------------------------------------------
func (a *App) AutoMigrate() error {
	if a.Migrator != nil {
		return a.Migrator.AutoMigrate()
	}
	return nil
}

// -----------------------------------------------------------
// 6) Boot() - เริ่มระบบพื้นฐาน
// -----------------------------------------------------------
func (a *App) Boot() {
	fmt.Println("⚙️  Booting Neonex Core...")
}

// -----------------------------------------------------------
// 7) StartHTTP() - HTTP Server Engine
// -----------------------------------------------------------
func (a *App) StartHTTP() {
	// Configure Fiber with custom branding
	app := fiber.New(fiber.Config{
		AppName:               "Neonex Core v0.1-alpha",
		DisableStartupMessage: true, // Disable default Fiber banner
	})

	// โหลด routes จากทุก module
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

	app.Listen(":8080")
}
