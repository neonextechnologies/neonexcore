package main

import (
	"context"
	"fmt"
	"math/rand"
	"neonexcore/pkg/metrics"
	"neonexcore/pkg/websocket"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	fmt.Println("=== NeonexCore Metrics Dashboard Example ===\n")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "NeonexCore Metrics Demo",
	})

	// Middleware
	app.Use(cors.New())
	app.Use(fiberlogger.New())

	// Create metrics collector
	collectorConfig := metrics.DefaultCollectorConfig()
	collectorConfig.CollectSystemMetrics = true
	collectorConfig.SystemMetricsInterval = 2 * time.Second
	collector := metrics.NewCollector(collectorConfig)
	defer collector.Close()

	// Add metrics middleware
	app.Use(metrics.Middleware(collector))
	app.Use(metrics.MethodMiddleware(collector))
	app.Use(metrics.ErrorMiddleware(collector))

	// Create WebSocket hub for real-time updates
	hubConfig := websocket.DefaultHubConfig()
	hub := websocket.NewHub(hubConfig)

	// Create dashboard
	dashConfig := metrics.DefaultDashboardConfig()
	dashConfig.BroadcastInterval = 1 * time.Second
	dashboard := metrics.NewDashboard(collector, hub, dashConfig)
	defer dashboard.Close()

	// Setup dashboard routes
	dashboard.SetupRoutes(app)

	// Setup WebSocket routes
	websocket.SetupRoutes(app, hub, nil)

	// Setup alerts
	setupAlerts(dashboard)

	// Create custom metrics
	setupCustomMetrics(collector)

	// Demo routes
	setupDemoRoutes(app, collector)

	// Start background simulation
	go simulateTraffic(app, collector)

	// Print info
	fmt.Println("✅ Metrics collector initialized")
	fmt.Println("✅ Real-time dashboard ready")
	fmt.Println("✅ WebSocket connection established")
	fmt.Println("✅ Alerts configured")
	fmt.Println("\n📊 Dashboard: http://localhost:3000/metrics/dashboard")
	fmt.Println("📡 WebSocket: ws://localhost:3000/ws")
	fmt.Println("📈 Metrics API: http://localhost:3000/metrics")
	fmt.Println("\nPress Ctrl+C to stop\n")

	// Start server
	if err := app.Listen(":3000"); err != nil {
		fmt.Printf("❌ Error starting server: %v\n", err)
	}
}

func setupAlerts(dashboard *metrics.Dashboard) {
	// High memory alert
	dashboard.AddAlert(metrics.Alert{
		Name:        "high_memory",
		Description: "Memory usage exceeded 100MB",
		Metric:      "system_memory_bytes",
		Condition:   metrics.ConditionGreaterThan,
		Threshold:   100 * 1024 * 1024, // 100MB
		Enabled:     true,
	})

	// High goroutine count
	dashboard.AddAlert(metrics.Alert{
		Name:        "high_goroutines",
		Description: "Too many goroutines running",
		Metric:      "system_goroutines",
		Condition:   metrics.ConditionGreaterThan,
		Threshold:   100,
		Enabled:     true,
	})

	// Slow response time
	dashboard.AddAlert(metrics.Alert{
		Name:        "slow_response",
		Description: "Average response time is too slow",
		Metric:      "http_request_duration_seconds",
		Condition:   metrics.ConditionGreaterThan,
		Threshold:   1.0, // 1 second
		Enabled:     true,
	})

	// High error rate
	dashboard.AddAlert(metrics.Alert{
		Name:        "high_errors",
		Description: "Too many HTTP errors",
		Metric:      "http_errors_total",
		Condition:   metrics.ConditionGreaterThan,
		Threshold:   50,
		Enabled:     true,
	})

	fmt.Println("✅ Configured 4 alerts")
}

func setupCustomMetrics(collector *metrics.Collector) {
	// Business metrics
	collector.NewCounter("user_signups_total", "Total user signups", nil)
	collector.NewCounter("user_logins_total", "Total user logins", nil)
	collector.NewCounter("orders_total", "Total orders", nil)
	collector.NewSummary("order_amount_usd", "Order amount in USD", nil)

	// Application metrics
	collector.NewGauge("active_users", "Number of active users", nil)
	collector.NewGauge("database_connections", "Active database connections", nil)
	collector.NewCounter("cache_hits_total", "Cache hits", nil)
	collector.NewCounter("cache_misses_total", "Cache misses", nil)

	// Performance metrics
	collector.NewHistogram(
		"db_query_duration_seconds",
		"Database query duration",
		nil,
		[]float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
	)

	fmt.Println("✅ Created 9 custom metrics")
}

func setupDemoRoutes(app *fiber.App, collector *metrics.Collector) {
	// Get custom metrics for demo
	activeUsers := collector.NewGauge("active_users", "", nil)
	signupsCounter := collector.NewCounter("user_signups_total", "", nil)
	loginsCounter := collector.NewCounter("user_logins_total", "", nil)
	ordersCounter := collector.NewCounter("orders_total", "", nil)
	orderAmount := collector.NewSummary("order_amount_usd", "", nil)
	cacheHits := collector.NewCounter("cache_hits_total", "", nil)
	cacheMisses := collector.NewCounter("cache_misses_total", "", nil)
	dbQuery := collector.NewHistogram("db_query_duration_seconds", "", nil, nil)
	dbConnections := collector.NewGauge("database_connections", "", nil)

	// Home page
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "NeonexCore Metrics Demo",
			"dashboard": "/metrics/dashboard",
			"metrics":   "/metrics",
			"websocket": "ws://localhost:3000/ws",
		})
	})

	// Simulate user signup
	app.Post("/api/signup", func(c *fiber.Ctx) error {
		signupsCounter.Inc()
		activeUsers.Inc()
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		return c.JSON(fiber.Map{"success": true, "message": "User created"})
	})

	// Simulate user login
	app.Post("/api/login", func(c *fiber.Ctx) error {
		loginsCounter.Inc()
		activeUsers.Inc()
		time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
		return c.JSON(fiber.Map{"success": true, "token": "abc123"})
	})

	// Simulate user logout
	app.Post("/api/logout", func(c *fiber.Ctx) error {
		activeUsers.Dec()
		return c.JSON(fiber.Map{"success": true})
	})

	// Simulate order creation
	app.Post("/api/orders", func(c *fiber.Ctx) error {
		ordersCounter.Inc()
		amount := float64(rand.Intn(500) + 10)
		orderAmount.Observe(amount)

		// Simulate DB query
		queryDuration := float64(rand.Intn(100)) / 1000.0
		dbQuery.Observe(queryDuration)

		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return c.JSON(fiber.Map{
			"success": true,
			"order_id": rand.Intn(10000),
			"amount":   amount,
		})
	})

	// Simulate cache lookup
	app.Get("/api/cache/:key", func(c *fiber.Ctx) error {
		if rand.Float32() < 0.8 { // 80% hit rate
			cacheHits.Inc()
			return c.JSON(fiber.Map{"cache": "hit", "value": "cached_data"})
		}
		cacheMisses.Inc()
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		return c.JSON(fiber.Map{"cache": "miss", "value": "fresh_data"})
	})

	// Simulate database query
	app.Get("/api/query", func(c *fiber.Ctx) error {
		dbConnections.Inc()
		defer dbConnections.Dec()

		queryDuration := float64(rand.Intn(200)) / 1000.0
		dbQuery.Observe(queryDuration)

		time.Sleep(time.Duration(queryDuration * 1000) * time.Millisecond)
		return c.JSON(fiber.Map{"success": true, "rows": rand.Intn(100)})
	})

	// Simulate error
	app.Get("/api/error", func(c *fiber.Ctx) error {
		return c.Status(500).JSON(fiber.Map{"error": "Something went wrong"})
	})

	// Simulate slow endpoint
	app.Get("/api/slow", func(c *fiber.Ctx) error {
		time.Sleep(2 * time.Second)
		return c.JSON(fiber.Map{"message": "This was slow"})
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"uptime": collector.GetUptime().Seconds(),
		})
	})

	fmt.Println("✅ Setup 9 demo routes")
}

func simulateTraffic(app *fiber.App, collector *metrics.Collector) {
	fmt.Println("🚀 Starting traffic simulation...\n")

	ctx := context.Background()
	activeUsers := collector.NewGauge("active_users", "", nil)

	// Set initial active users
	activeUsers.Set(int64(rand.Intn(20) + 10))

	endpoints := []string{
		"/api/signup",
		"/api/login",
		"/api/orders",
		"/api/cache/user:123",
		"/api/query",
		"/health",
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Random endpoint
			endpoint := endpoints[rand.Intn(len(endpoints))]

			// Simulate request (metrics are tracked by middleware)
			go func() {
				// We don't actually make HTTP requests here,
				// just update some metrics to simulate activity
				if endpoint == "/api/orders" {
					collector.NewCounter("orders_total", "", nil).Inc()
					amount := float64(rand.Intn(500) + 10)
					collector.NewSummary("order_amount_usd", "", nil).Observe(amount)
				}

				// Random cache hit/miss
				if rand.Float32() < 0.8 {
					collector.NewCounter("cache_hits_total", "", nil).Inc()
				} else {
					collector.NewCounter("cache_misses_total", "", nil).Inc()
				}

				// Random DB query
				duration := float64(rand.Intn(100)) / 1000.0
				collector.NewHistogram("db_query_duration_seconds", "", nil, nil).Observe(duration)

				// Random active users change
				if rand.Float32() < 0.1 {
					if rand.Float32() < 0.5 {
						activeUsers.Inc()
					} else if activeUsers.Get() > 5 {
						activeUsers.Dec()
					}
				}
			}()
		}
	}
}
