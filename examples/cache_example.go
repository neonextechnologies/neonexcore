package main

import (
	"context"
	"fmt"
	"neonexcore/pkg/cache"
	"time"
)

// User represents a sample user model
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func main() {
	ctx := context.Background()

	fmt.Println("=== NeonexCore Cache Examples ===\n")

	// Example 1: Memory Cache
	runMemoryCacheExample(ctx)

	// Example 2: Redis Cache
	runRedisCacheExample(ctx)

	// Example 3: Multi-Tier Cache
	runMultiTierCacheExample(ctx)

	// Example 4: Practical Use Cases
	runUseCaseExamples(ctx)
}

// Example 1: Memory Cache (L1)
func runMemoryCacheExample(ctx context.Context) {
	fmt.Println("📦 Example 1: Memory Cache (In-Memory LRU)")
	fmt.Println("-------------------------------------------")

	// Create memory cache with custom config
	config := cache.DefaultMemoryCacheConfig()
	config.MaxSize = 1000
	config.CleanupInterval = 30 * time.Second

	memCache := cache.NewMemoryCache(config)
	defer memCache.Close()

	// Basic operations
	user := User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
		Role:     "admin",
	}

	// Set with TTL
	err := memCache.Set(ctx, "user:1", user, 5*time.Minute)
	if err != nil {
		fmt.Printf("❌ Error setting cache: %v\n", err)
		return
	}
	fmt.Println("✅ Set user:1 with 5-minute TTL")

	// Get
	value, err := memCache.Get(ctx, "user:1")
	if err != nil {
		fmt.Printf("❌ Error getting cache: %v\n", err)
		return
	}
	cachedUser := value.(User)
	fmt.Printf("✅ Get user:1: %s (%s)\n", cachedUser.Username, cachedUser.Email)

	// Check TTL
	ttl, _ := memCache.TTL(ctx, "user:1")
	fmt.Printf("⏱️  TTL remaining: %v\n", ttl.Round(time.Second))

	// Increment counter
	memCache.Set(ctx, "views:post:123", int64(0), 1*time.Hour)
	newCount, _ := memCache.Increment(ctx, "views:post:123", 1)
	fmt.Printf("📊 Incremented views: %d\n", newCount)

	// Batch operations
	items := map[string]interface{}{
		"config:timeout": 30,
		"config:retries": 3,
		"config:enabled": true,
	}
	memCache.SetMulti(ctx, items, 10*time.Minute)
	fmt.Println("✅ Set 3 config items in batch")

	values, _ := memCache.GetMulti(ctx, []string{"config:timeout", "config:retries"})
	fmt.Printf("✅ GetMulti: timeout=%v, retries=%v\n", values["config:timeout"], values["config:retries"])

	// Statistics
	if sp, ok := memCache.(cache.StatsProvider); ok {
		stats, _ := sp.Stats(ctx)
		fmt.Printf("📈 Stats: Hits=%d, Misses=%d, Keys=%d\n\n",
			stats.Hits, stats.Misses, stats.Keys)
	}
}

// Example 2: Redis Cache (L2)
func runRedisCacheExample(ctx context.Context) {
	fmt.Println("🔴 Example 2: Redis Cache (Distributed)")
	fmt.Println("----------------------------------------")

	// Create Redis cache
	config := cache.DefaultRedisCacheConfig()
	config.Addr = "localhost:6379"
	config.DB = 0

	redisCache, err := cache.NewRedisCache(config)
	if err != nil {
		fmt.Printf("⚠️  Warning: Redis not available: %v\n", err)
		fmt.Println("   Skipping Redis examples...\n")
		return
	}
	defer redisCache.Close()

	// Basic operations (same interface as memory cache)
	session := map[string]interface{}{
		"user_id":   123,
		"token":     "abc123xyz",
		"expires":   time.Now().Add(24 * time.Hour).Unix(),
		"ip":        "192.168.1.100",
		"user_agent": "Mozilla/5.0",
	}

	// Set session
	err = redisCache.Set(ctx, "session:abc123", session, 24*time.Hour)
	if err != nil {
		fmt.Printf("❌ Error setting session: %v\n", err)
		return
	}
	fmt.Println("✅ Set session:abc123 with 24-hour TTL")

	// Get session
	value, err := redisCache.Get(ctx, "session:abc123")
	if err != nil {
		fmt.Printf("❌ Error getting session: %v\n", err)
		return
	}
	cachedSession := value.(map[string]interface{})
	fmt.Printf("✅ Get session: user_id=%v, token=%v\n",
		cachedSession["user_id"], cachedSession["token"])

	// Pattern matching
	redisCache.Set(ctx, "user:1", "John", 10*time.Minute)
	redisCache.Set(ctx, "user:2", "Jane", 10*time.Minute)
	redisCache.Set(ctx, "user:3", "Bob", 10*time.Minute)

	userKeys, _ := redisCache.Keys(ctx, "user:*")
	fmt.Printf("🔍 Found %d user keys: %v\n", len(userKeys), userKeys)

	// Atomic counters
	redisCache.Set(ctx, "api:calls:today", int64(0), 24*time.Hour)
	count1, _ := redisCache.Increment(ctx, "api:calls:today", 1)
	count2, _ := redisCache.Increment(ctx, "api:calls:today", 5)
	fmt.Printf("📊 API calls: %d → %d\n", count1, count2)

	// Statistics
	if sp, ok := redisCache.(cache.StatsProvider); ok {
		stats, _ := sp.Stats(ctx)
		fmt.Printf("📈 Redis Stats: Keys=%d, Memory=%d bytes\n\n",
			stats.Keys, stats.Memory)
	}
}

// Example 3: Multi-Tier Cache (L1 + L2)
func runMultiTierCacheExample(ctx context.Context) {
	fmt.Println("🏗️  Example 3: Multi-Tier Cache (L1 Memory + L2 Redis)")
	fmt.Println("-------------------------------------------------------")

	// Create L1 (Memory)
	memConfig := cache.DefaultMemoryCacheConfig()
	memConfig.MaxSize = 100
	memCache := cache.NewMemoryCache(memConfig)
	defer memCache.Close()

	// Create L2 (Redis)
	redisConfig := cache.DefaultRedisCacheConfig()
	redisConfig.Addr = "localhost:6379"
	redisCache, err := cache.NewRedisCache(redisConfig)
	if err != nil {
		fmt.Printf("⚠️  Warning: Redis not available, using memory cache only\n\n")
		return
	}
	defer redisCache.Close()

	// Create multi-tier cache
	config := cache.DefaultMultiTierConfig()
	config.PromoteL1 = true  // Promote to L1 on hit
	config.WriteThru = true   // Write to all tiers

	multiCache := cache.NewMultiTierCache(config)
	multiCache.AddTier(memCache, cache.TierL1)
	multiCache.AddTier(redisCache, cache.TierL2)

	fmt.Println("✅ Created 2-tier cache: L1 (Memory) + L2 (Redis)")

	// Write data (goes to both L1 and L2)
	product := map[string]interface{}{
		"id":    456,
		"name":  "Laptop Pro",
		"price": 1299.99,
		"stock": 10,
	}

	err = multiCache.Set(ctx, "product:456", product, 10*time.Minute)
	if err != nil {
		fmt.Printf("❌ Error setting cache: %v\n", err)
		return
	}
	fmt.Println("✅ Set product:456 (written to L1 and L2)")

	// Clear L1 to test L2 retrieval
	memCache.Delete(ctx, "product:456")
	fmt.Println("🗑️  Cleared product:456 from L1")

	// First access - retrieves from L2, promotes to L1
	value, err := multiCache.Get(ctx, "product:456")
	if err != nil {
		fmt.Printf("❌ Error getting cache: %v\n", err)
		return
	}
	fmt.Println("✅ Get product:456 from L2 (cache hit)")

	cachedProduct := value.(map[string]interface{})
	fmt.Printf("   Product: %s - $%.2f\n", cachedProduct["name"], cachedProduct["price"])

	// Second access - now in L1 (faster)
	value2, _ := multiCache.Get(ctx, "product:456")
	fmt.Println("✅ Get product:456 from L1 (promoted, faster!)")

	// Check if promoted to L1
	_, err = memCache.Get(ctx, "product:456")
	if err == nil {
		fmt.Println("✅ Confirmed: product:456 is now in L1")
	}

	// Statistics from both tiers
	if sp, ok := multiCache.(cache.StatsProvider); ok {
		stats, _ := sp.Stats(ctx)
		fmt.Printf("📈 Multi-tier Stats: Hits=%d, Misses=%d, Keys=%d\n\n",
			stats.Hits, stats.Misses, stats.Keys)
	}
}

// Example 4: Practical Use Cases
func runUseCaseExamples(ctx context.Context) {
	fmt.Println("💡 Example 4: Practical Use Cases")
	fmt.Println("----------------------------------")

	memCache := cache.NewMemoryCache(cache.DefaultMemoryCacheConfig())
	defer memCache.Close()

	// Use Case 1: HTTP Response Caching
	fmt.Println("\n1️⃣  HTTP Response Caching")
	cacheKey := "api:users:list:page:1"
	users := []User{
		{ID: 1, Username: "john", Email: "john@example.com", Role: "admin"},
		{ID: 2, Username: "jane", Email: "jane@example.com", Role: "user"},
	}

	memCache.Set(ctx, cacheKey, users, 5*time.Minute)
	fmt.Printf("✅ Cached API response for 5 minutes\n")

	cached, _ := memCache.Get(ctx, cacheKey)
	cachedUsers := cached.([]User)
	fmt.Printf("✅ Retrieved %d users from cache\n", len(cachedUsers))

	// Use Case 2: Session Management
	fmt.Println("\n2️⃣  Session Management")
	sessionID := "sess_xyz789"
	sessionData := map[string]interface{}{
		"user_id":    123,
		"username":   "john_doe",
		"role":       "admin",
		"created_at": time.Now().Unix(),
	}

	memCache.Set(ctx, "session:"+sessionID, sessionData, 30*time.Minute)
	fmt.Printf("✅ Stored session with 30-minute TTL\n")

	session, _ := memCache.Get(ctx, "session:"+sessionID)
	s := session.(map[string]interface{})
	fmt.Printf("✅ Retrieved session for user: %s (role: %s)\n",
		s["username"], s["role"])

	// Use Case 3: Rate Limiting
	fmt.Println("\n3️⃣  Rate Limiting")
	ip := "192.168.1.100"
	rateLimitKey := "ratelimit:" + ip

	// Simulate 5 requests
	for i := 1; i <= 5; i++ {
		count, _ := memCache.Increment(ctx, rateLimitKey, 1)
		if i == 1 {
			memCache.Expire(ctx, rateLimitKey, 1*time.Minute)
		}

		if count > 10 {
			fmt.Printf("❌ Request %d: Rate limit exceeded (%d/10)\n", i, count)
		} else {
			fmt.Printf("✅ Request %d: Allowed (%d/10)\n", i, count)
		}
	}

	// Use Case 4: Cache Invalidation
	fmt.Println("\n4️⃣  Cache Invalidation")
	memCache.Set(ctx, "user:1", "John", 10*time.Minute)
	memCache.Set(ctx, "user:1:profile", "Profile data", 10*time.Minute)
	memCache.Set(ctx, "user:1:permissions", "Permissions data", 10*time.Minute)

	fmt.Println("✅ Cached user data and related resources")

	// Invalidate all related caches
	memCache.DeleteMulti(ctx, []string{
		"user:1",
		"user:1:profile",
		"user:1:permissions",
	})
	fmt.Println("✅ Invalidated user and all related caches")

	// Use Case 5: Query Result Caching
	fmt.Println("\n5️⃣  Query Result Caching")
	queryKey := "query:products:category:electronics:page:1"
	results := []map[string]interface{}{
		{"id": 1, "name": "Laptop", "price": 999.99},
		{"id": 2, "name": "Mouse", "price": 29.99},
	}

	memCache.Set(ctx, queryKey, results, 2*time.Minute)
	fmt.Println("✅ Cached database query results for 2 minutes")

	cached2, err := memCache.Get(ctx, queryKey)
	if err == nil {
		products := cached2.([]map[string]interface{})
		fmt.Printf("✅ Retrieved %d products from cache (avoided DB query)\n", len(products))
	}

	// Use Case 6: Feature Flags
	fmt.Println("\n6️⃣  Feature Flags")
	flags := map[string]interface{}{
		"feature:new_ui":       true,
		"feature:beta_api":     false,
		"feature:dark_mode":    true,
		"feature:ai_assistant": true,
	}

	memCache.SetMulti(ctx, flags, 1*time.Hour)
	fmt.Println("✅ Cached 4 feature flags for 1 hour")

	newUIEnabled, _ := memCache.Get(ctx, "feature:new_ui")
	if newUIEnabled.(bool) {
		fmt.Println("✅ Feature 'new_ui' is enabled")
	}

	fmt.Println("\n✨ All examples completed successfully!")
}
