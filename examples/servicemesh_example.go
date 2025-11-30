package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"neonexcore/pkg/servicemesh"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("=== Service Mesh Example ===\n")

	// Example 1: Basic Sidecar Proxy
	fmt.Println("1. Starting Sidecar Proxy...")
	runBasicSidecar()

	// Example 2: Service Discovery
	fmt.Println("\n2. Service Discovery Example...")
	runServiceDiscovery()

	// Example 3: Traffic Management - Canary Deployment
	fmt.Println("\n3. Canary Deployment Example...")
	runCanaryDeployment()

	// Example 4: A/B Testing
	fmt.Println("\n4. A/B Testing Example...")
	runABTesting()

	// Example 5: Circuit Breaker
	fmt.Println("\n5. Circuit Breaker Example...")
	runCircuitBreaker()

	// Example 6: Complete Service Mesh Setup
	fmt.Println("\n6. Complete Service Mesh Setup...")
	runCompleteSetup()
}

// Example 1: Basic Sidecar Proxy
func runBasicSidecar() {
	// Create a simple HTTP service
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/api/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"users": []string{"Alice", "Bob", "Charlie"},
		})
	})

	// Start service on port 8080
	go app.Listen(":8080")
	time.Sleep(500 * time.Millisecond)

	// Create sidecar proxy
	config := &servicemesh.SidecarConfig{
		ServiceName:  "user-service",
		ServicePort:  8080,
		ProxyPort:    8081,
		EnableMetrics: true,
		EnableTracing: true,
		CircuitBreakerCfg: &servicemesh.CircuitBreakerConfig{
			FailureThreshold: 5,
			SuccessThreshold: 2,
			Timeout:          30 * time.Second,
		},
	}

	proxy, err := servicemesh.NewSidecarProxy(config)
	if err != nil {
		log.Printf("Failed to create sidecar: %v", err)
		return
	}

	// Start proxy in background
	go proxy.Start()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("✓ Sidecar proxy running on :8081")
	fmt.Println("✓ Application running on :8080")
	fmt.Println("✓ Health check: http://localhost:8081/health")
	fmt.Println("✓ Metrics: http://localhost:8081/metrics")

	// Stop after demo
	time.Sleep(2 * time.Second)
	proxy.Stop(context.Background())
}

// Example 2: Service Discovery
func runServiceDiscovery() {
	registry := servicemesh.NewServiceRegistry("")

	// Register multiple service instances
	services := []struct {
		name string
		port int
	}{
		{"user-service", 8080},
		{"user-service", 8081},
		{"order-service", 8082},
		{"product-service", 8083},
	}

	for _, svc := range services {
		instance := &servicemesh.ServiceInstance{
			ServiceName: svc.name,
			Host:        "localhost",
			Port:        svc.port,
			Protocol:    "http",
			Metadata: map[string]string{
				"version": "1.0",
				"region":  "us-west",
			},
		}
		if err := registry.Register(instance); err != nil {
			log.Printf("Failed to register: %v", err)
			continue
		}
		fmt.Printf("✓ Registered %s on port %d\n", svc.name, svc.port)
	}

	// Discover services
	fmt.Println("\nDiscovering services:")
	allServices := registry.ListServices()
	for _, name := range allServices {
		instances := registry.GetServiceInstances(name)
		fmt.Printf("  %s: %d instances\n", name, len(instances))
		for _, inst := range instances {
			fmt.Printf("    - %s:%d (health: %s)\n", inst.Host, inst.Port, inst.Health)
		}
	}

	// Discover specific service
	instance, err := registry.Discover("user-service")
	if err != nil {
		log.Printf("Discovery failed: %v", err)
		return
	}
	fmt.Printf("\n✓ Discovered user-service at %s:%d\n", instance.Host, instance.Port)
}

// Example 3: Canary Deployment
func runCanaryDeployment() {
	tm := servicemesh.NewTrafficManager()

	// Configure canary deployment
	policy := &servicemesh.TrafficPolicy{
		ServiceName: "user-service",
		Canary: &servicemesh.CanaryConfig{
			Enabled:        true,
			NewVersion:     "v2",
			StableVersion:  "v1",
			InitialWeight:  10,  // Start with 10% traffic
			IncrementStep:  10,  // Increase by 10% each step
			IncrementDelay: 60,  // Wait 60s between steps
			MaxWeight:      100,
			SuccessRate:    0.99, // Require 99% success
		},
	}

	if err := tm.SetPolicy(policy); err != nil {
		log.Printf("Failed to set policy: %v", err)
		return
	}

	fmt.Println("✓ Canary deployment configured")
	fmt.Println("  Initial: 10% v2, 90% v1")

	// Simulate traffic routing
	fmt.Println("\nSimulating 100 requests:")
	v1Count := 0
	v2Count := 0

	for i := 0; i < 100; i++ {
		version := tm.SelectVersion("user-service", nil, "")
		if version == "v2" {
			v2Count++
		} else {
			v1Count++
		}
	}

	fmt.Printf("  v1: %d requests (%.0f%%)\n", v1Count, float64(v1Count))
	fmt.Printf("  v2: %d requests (%.0f%%)\n", v2Count, float64(v2Count))

	// Increment canary
	fmt.Println("\n✓ Incrementing canary weight...")
	tm.IncrementCanary("user-service")
	currentPolicy := tm.GetPolicy("user-service")
	fmt.Printf("  New weight: %d%% v2\n", currentPolicy.Canary.InitialWeight)

	// Promote canary
	fmt.Println("\n✓ Promoting canary to stable...")
	tm.PromoteCanary("user-service")
	fmt.Println("  v2 is now stable version")
}

// Example 4: A/B Testing
func runABTesting() {
	tm := servicemesh.NewTrafficManager()

	// Configure A/B test
	policy := &servicemesh.TrafficPolicy{
		ServiceName: "user-service",
		ABTest: &servicemesh.ABTestConfig{
			Enabled:  true,
			VersionA: "v1",
			VersionB: "v2",
			SplitKey: "X-User-Cohort",
			WeightA:  50, // 50% each
			WeightB:  50,
		},
	}

	tm.SetPolicy(policy)

	fmt.Println("✓ A/B test configured (50/50 split)")

	// Test with sticky sessions
	fmt.Println("\nTesting sticky sessions:")
	
	cohorts := []string{"A", "B"}
	for _, cohort := range cohorts {
		headers := map[string]string{
			"X-User-Cohort": cohort,
		}
		version := tm.SelectVersion("user-service", headers, "")
		fmt.Printf("  Cohort %s → Version %s (sticky)\n", cohort, version)
	}

	// Test random assignment
	fmt.Println("\nSimulating 100 new users (random assignment):")
	v1Count := 0
	v2Count := 0

	for i := 0; i < 100; i++ {
		version := tm.SelectVersion("user-service", map[string]string{}, "")
		if version == "v1" {
			v1Count++
		} else {
			v2Count++
		}
	}

	fmt.Printf("  v1: %d users (%.0f%%)\n", v1Count, float64(v1Count))
	fmt.Printf("  v2: %d users (%.0f%%)\n", v2Count, float64(v2Count))
}

// Example 5: Circuit Breaker
func runCircuitBreaker() {
	config := &servicemesh.CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          10 * time.Second,
		HalfOpenRequests: 3,
	}

	cb := servicemesh.NewCircuitBreaker(config)

	fmt.Println("✓ Circuit breaker created")
	fmt.Printf("  State: %s\n", cb.GetState())

	// Simulate failures
	fmt.Println("\nSimulating 5 failures...")
	for i := 0; i < 5; i++ {
		cb.RecordFailure()
		fmt.Printf("  Failure %d recorded (state: %s)\n", i+1, cb.GetState())
	}

	// Circuit should be open
	fmt.Printf("\n✓ Circuit breaker is now: %s\n", cb.GetState())
	
	if cb.IsOpen() {
		fmt.Println("  ⚠️  Requests are being blocked")
	}

	// Wait for half-open
	fmt.Println("\nWaiting 3 seconds for half-open state...")
	time.Sleep(3 * time.Second)

	// Check state after timeout (in real scenario, would be longer)
	fmt.Printf("  State after timeout: %s\n", cb.GetState())

	// Simulate recovery
	fmt.Println("\nSimulating 2 successes in half-open...")
	cb.RecordSuccess()
	fmt.Printf("  Success 1 (state: %s)\n", cb.GetState())
	cb.RecordSuccess()
	fmt.Printf("  Success 2 (state: %s)\n", cb.GetState())

	fmt.Printf("\n✓ Circuit breaker recovered: %s\n", cb.GetState())

	// Show metrics
	metrics := cb.GetMetrics()
	fmt.Println("\nCircuit Breaker Metrics:")
	fmt.Printf("  State: %v\n", metrics["state"])
	fmt.Printf("  Failure Count: %v\n", metrics["failure_count"])
	fmt.Printf("  Success Count: %v\n", metrics["success_count"])
}

// Example 6: Complete Service Mesh Setup
func runCompleteSetup() {
	fmt.Println("Setting up complete service mesh...")

	// 1. Create service registry
	registry := servicemesh.NewServiceRegistry("")

	// 2. Register services
	services := []string{"user-service", "order-service", "product-service"}
	for i, name := range services {
		instance := &servicemesh.ServiceInstance{
			ServiceName: name,
			Host:        "localhost",
			Port:        8080 + i,
			Protocol:    "http",
			Metadata: map[string]string{
				"version": "1.0",
			},
		}
		registry.Register(instance)
		fmt.Printf("✓ Registered %s\n", name)
	}

	// 3. Create traffic manager
	tm := servicemesh.NewTrafficManager()

	// 4. Configure traffic policies
	policies := []servicemesh.TrafficPolicy{
		{
			ServiceName: "user-service",
			Splits: []servicemesh.TrafficSplit{
				{Version: "v1", Weight: 80},
				{Version: "v2", Weight: 20},
			},
		},
		{
			ServiceName: "order-service",
			Canary: &servicemesh.CanaryConfig{
				Enabled:       true,
				NewVersion:    "v2",
				StableVersion: "v1",
				InitialWeight: 10,
			},
		},
	}

	for _, policy := range policies {
		p := policy // Create new variable for pointer
		tm.SetPolicy(&p)
		fmt.Printf("✓ Traffic policy set for %s\n", policy.ServiceName)
	}

	// 5. Create circuit breakers
	cb := servicemesh.NewCircuitBreaker(&servicemesh.CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          30 * time.Second,
	})
	fmt.Println("✓ Circuit breaker configured")

	// 6. Create sidecar proxies
	config := &servicemesh.SidecarConfig{
		ServiceName:   "user-service",
		ServicePort:   8080,
		ProxyPort:     8090,
		EnableMetrics: true,
		EnableTracing: true,
		CircuitBreakerCfg: &servicemesh.CircuitBreakerConfig{
			FailureThreshold: 5,
			SuccessThreshold: 2,
			Timeout:          30 * time.Second,
		},
	}

	proxy, err := servicemesh.NewSidecarProxy(config)
	if err != nil {
		log.Printf("Failed to create proxy: %v", err)
		return
	}
	fmt.Println("✓ Sidecar proxy created")

	fmt.Println("\n🎉 Service Mesh Setup Complete!")
	fmt.Println("\nComponents:")
	fmt.Printf("  - Services Registered: %d\n", len(registry.ListServices()))
	fmt.Printf("  - Traffic Policies: %d\n", len(tm.ListPolicies()))
	fmt.Printf("  - Circuit Breaker: %s\n", cb.GetState())
	fmt.Printf("  - Sidecar Proxy: Ready\n")

	fmt.Println("\nFeatures Enabled:")
	fmt.Println("  ✓ Service Discovery")
	fmt.Println("  ✓ Load Balancing")
	fmt.Println("  ✓ Traffic Splitting")
	fmt.Println("  ✓ Canary Deployment")
	fmt.Println("  ✓ Circuit Breaking")
	fmt.Println("  ✓ Health Checking")
	fmt.Println("  ✓ Metrics Collection")
	fmt.Println("  ✓ Distributed Tracing")

	// Cleanup
	time.Sleep(2 * time.Second)
	proxy.Stop(context.Background())
}
