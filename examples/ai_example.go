package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"neonexcore/pkg/ai"
)

func main() {
	fmt.Println("=== AI/ML Integration Examples ===\n")

	// Example 1: Basic Model Management
	fmt.Println("1. Model Management Example...")
	runModelManagement()

	// Example 2: OpenAI Integration
	fmt.Println("\n2. OpenAI Integration Example...")
	runOpenAIExample()

	// Example 3: Feature Store
	fmt.Println("\n3. Feature Store Example...")
	runFeatureStoreExample()

	// Example 4: ML Pipeline
	fmt.Println("\n4. ML Pipeline Example...")
	runPipelineExample()

	// Example 5: Inference Caching
	fmt.Println("\n5. Inference Caching Example...")
	runCachingExample()

	// Example 6: Batch Processing
	fmt.Println("\n6. Batch Processing Example...")
	runBatchProcessing()
}

// Example 1: Model Management
func runModelManagement() {
	manager := ai.NewModelManager()

	// Register mock provider
	mockProvider := &MockProvider{}
	manager.RegisterProvider("mock", mockProvider)

	// Load models
	models := []ai.ModelConfig{
		{
			ID:       "classifier-v1",
			Name:     "Text Classifier",
			Version:  "1.0.0",
			Type:     ai.ModelTypeTextClassification,
			Provider: "mock",
		},
		{
			ID:       "sentiment-v1",
			Name:     "Sentiment Analyzer",
			Version:  "1.0.0",
			Type:     ai.ModelTypeSentiment,
			Provider: "mock",
		},
		{
			ID:       "embedder-v1",
			Name:     "Text Embedder",
			Version:  "1.0.0",
			Type:     ai.ModelTypeEmbedding,
			Provider: "mock",
		},
	}

	for _, config := range models {
		model, err := manager.LoadModel(&config)
		if err != nil {
			log.Printf("Failed to load model: %v", err)
			continue
		}
		fmt.Printf("✓ Loaded model: %s (type: %s, status: %s)\n", 
			model.Name, model.Type, model.Status)
	}

	// List all models
	fmt.Println("\nLoaded models:")
	for _, model := range manager.ListModels() {
		fmt.Printf("  - %s (%s) - %s\n", model.Name, model.ID, model.Type)
	}

	// Perform inference
	ctx := context.Background()
	output, err := manager.Predict(ctx, &ai.InferenceInput{
		ModelID: "classifier-v1",
		Data:    "This is a test sentence",
	})

	if err != nil {
		log.Printf("Prediction failed: %v", err)
		return
	}

	fmt.Printf("\n✓ Prediction result: %v\n", output.Result)
	fmt.Printf("  Latency: %v\n", output.Latency)

	// Get metrics
	metrics := manager.GetAllMetrics()
	fmt.Println("\nModel metrics:")
	for modelID, m := range metrics {
		fmt.Printf("  %s: %d requests, avg latency: %v\n", 
			modelID, m.RequestCount, m.AvgLatency)
	}
}

// Example 2: OpenAI Integration (demo without actual API key)
func runOpenAIExample() {
	manager := ai.NewModelManager()

	// Note: This example shows the API structure
	// Replace with actual API key to use
	fmt.Println("OpenAI API Structure Demo:")
	
	config := &ai.ModelConfig{
		ID:       "gpt-4",
		Name:     "GPT-4",
		Version:  "latest",
		Type:     ai.ModelTypeTextGeneration,
		Provider: "openai",
		APIKey:   "your-api-key-here",
	}

	fmt.Printf("✓ Model config: %s\n", config.Name)
	fmt.Printf("  Provider: %s\n", config.Provider)
	fmt.Printf("  Type: %s\n", config.Type)

	// Example input structure
	input := &ai.InferenceInput{
		ModelID: "gpt-4",
		Data:    "Explain machine learning in simple terms",
		Parameters: map[string]interface{}{
			"type":        "chat",
			"temperature": 0.7,
			"max_tokens":  200,
			"system":      "You are a helpful AI assistant",
		},
	}

	fmt.Println("\n✓ Input parameters:")
	fmt.Printf("  Temperature: %v\n", input.Parameters["temperature"])
	fmt.Printf("  Max tokens: %v\n", input.Parameters["max_tokens"])
	fmt.Printf("  Type: %v\n", input.Parameters["type"])

	fmt.Println("\n📝 To use OpenAI:")
	fmt.Println("  1. Set OPENAI_API_KEY environment variable")
	fmt.Println("  2. Create provider: ai.NewOpenAIProvider(&ai.OpenAIConfig{...})")
	fmt.Println("  3. Register: manager.RegisterProvider(\"openai\", provider)")
	fmt.Println("  4. Call: manager.Predict(ctx, input)")
}

// Example 3: Feature Store (in-memory demo)
func runFeatureStoreExample() {
	// Note: This is a simplified demo without database
	fmt.Println("Feature Store Demo (in-memory):")

	features := map[string]map[string]interface{}{
		"user-123": {
			"age":             25,
			"location":        "US",
			"activity_score":  85.5,
			"purchase_count":  12,
			"last_login_days": 2,
		},
		"user-456": {
			"age":             32,
			"location":        "UK",
			"activity_score":  92.3,
			"purchase_count":  45,
			"last_login_days": 1,
		},
	}

	fmt.Println("\n✓ User features:")
	for userID, featureVector := range features {
		fmt.Printf("\n  %s:\n", userID)
		for k, v := range featureVector {
			fmt.Printf("    %s: %v\n", k, v)
		}
	}

	// Feature groups
	groups := map[string][]string{
		"user-profile": {
			"age",
			"location",
			"activity_score",
		},
		"user-behavior": {
			"purchase_count",
			"last_login_days",
		},
	}

	fmt.Println("\n✓ Feature groups:")
	for groupName, featureNames := range groups {
		fmt.Printf("  %s: %v\n", groupName, featureNames)
	}

	fmt.Println("\n📝 With database:")
	fmt.Println("  featureStore := ai.NewFeatureStore(db)")
	fmt.Println("  featureStore.SetFeature(ctx, feature)")
	fmt.Println("  vector, _ := featureStore.GetFeatureVector(ctx, type, id, names)")
}

// Example 4: ML Pipeline
func runPipelineExample() {
	manager := ai.NewModelManager()
	mockProvider := &MockProvider{}
	manager.RegisterProvider("mock", mockProvider)

	// Load sentiment model
	manager.LoadModel(&ai.ModelConfig{
		ID:       "sentiment-model",
		Name:     "Sentiment Analyzer",
		Type:     ai.ModelTypeSentiment,
		Provider: "mock",
	})

	// Create pipeline manager
	pipelineManager := ai.NewPipelineManager(manager)

	// Define sentiment analysis pipeline
	pipeline := &ai.Pipeline{
		ID:          "sentiment-pipeline",
		Name:        "Sentiment Analysis Pipeline",
		Description: "Preprocessing → Model → Extract score",
		Steps: []ai.PipelineStep{
			{
				Name: "preprocess",
				Type: ai.StepTypePreprocess,
				Transform: func(ctx context.Context, input interface{}) (interface{}, error) {
					text := fmt.Sprintf("%v", input)
					fmt.Printf("  [Preprocess] Input: %s\n", text)
					return text, nil
				},
			},
			{
				Name:    "sentiment-model",
				Type:    ai.StepTypeModel,
				ModelID: "sentiment-model",
			},
			{
				Name: "extract-result",
				Type: ai.StepTypePostprocess,
				Transform: func(ctx context.Context, input interface{}) (interface{}, error) {
					result := input.(map[string]interface{})
					fmt.Printf("  [Postprocess] Extracted: %v\n", result["sentiment"])
					return result["sentiment"], nil
				},
			},
		},
	}

	pipelineManager.CreatePipeline(pipeline)
	fmt.Printf("✓ Created pipeline: %s\n", pipeline.Name)

	// Execute pipeline
	ctx := context.Background()
	testInputs := []string{
		"This product is amazing!",
		"Terrible experience, very disappointed",
		"It's okay, nothing special",
	}

	fmt.Println("\nExecuting pipeline:")
	for _, input := range testInputs {
		fmt.Printf("\n📝 Input: %s\n", input)
		result, err := pipelineManager.Execute(ctx, "sentiment-pipeline", input)
		if err != nil {
			log.Printf("  ❌ Error: %v", err)
			continue
		}

		fmt.Printf("  ✓ Output: %v\n", result.Output)
		fmt.Printf("  ⏱️  Latency: %v\n", result.Latency)
		fmt.Printf("  📊 Steps: %d\n", len(result.StepResults))
	}
}

// Example 5: Inference Caching
func runCachingExample() {
	cache := ai.NewInferenceCache(100, 1*time.Hour)

	// Simulate predictions
	predictions := []struct {
		modelID string
		data    string
		result  string
	}{
		{"model-1", "hello", "positive"},
		{"model-1", "hello", "positive"}, // Cache hit
		{"model-1", "world", "neutral"},
		{"model-1", "hello", "positive"}, // Cache hit
		{"model-2", "test", "negative"},
	}

	fmt.Println("Simulating inference with caching:")
	for i, pred := range predictions {
		input := &ai.InferenceInput{
			ModelID: pred.modelID,
			Data:    pred.data,
		}

		// Check cache
		cached := cache.Get(input)
		if cached != nil {
			fmt.Printf("  %d. ⚡ Cache HIT - Model: %s, Input: %s\n", 
				i+1, pred.modelID, pred.data)
		} else {
			// Simulate inference
			output := &ai.InferenceOutput{
				ModelID: pred.modelID,
				Result:  pred.result,
				Latency: 100 * time.Millisecond,
			}
			cache.Set(input, output)
			fmt.Printf("  %d. 🔄 Cache MISS - Model: %s, Input: %s\n", 
				i+1, pred.modelID, pred.data)
		}
	}

	// Show cache stats
	stats := cache.GetStats()
	fmt.Println("\n✓ Cache statistics:")
	fmt.Printf("  Size: %v / %v\n", stats["size"], stats["max_size"])
	fmt.Printf("  Hits: %v\n", stats["hits"])
	fmt.Printf("  Misses: %v\n", stats["misses"])
	fmt.Printf("  Hit rate: %.1f%%\n", stats["hit_rate"].(float64)*100)
}

// Example 6: Batch Processing
func runBatchProcessing() {
	manager := ai.NewModelManager()
	mockProvider := &MockProvider{}
	manager.RegisterProvider("mock", mockProvider)

	manager.LoadModel(&ai.ModelConfig{
		ID:       "batch-classifier",
		Type:     ai.ModelTypeTextClassification,
		Provider: "mock",
	})

	pipelineManager := ai.NewPipelineManager(manager)

	// Create batch processing pipeline
	pipeline := &ai.Pipeline{
		ID:   "batch-pipeline",
		Name: "Batch Text Classification",
		Steps: []ai.PipelineStep{
			{
				Name: "batch-classify",
				Type: ai.StepTypeTransform,
				Transform: ai.BatchProcessor(2, func(ctx context.Context, item interface{}) (interface{}, error) {
					// Classify each item
					return map[string]interface{}{
						"text":  item,
						"label": "positive", // Mock result
					}, nil
				}),
			},
		},
	}

	pipelineManager.CreatePipeline(pipeline)

	// Batch data
	texts := []interface{}{
		"Great product!",
		"Not bad",
		"Excellent service",
		"Could be better",
		"Absolutely love it",
	}

	fmt.Printf("✓ Processing %d items in batches of 2:\n", len(texts))
	
	ctx := context.Background()
	result, err := pipelineManager.Execute(ctx, "batch-pipeline", texts)
	if err != nil {
		log.Printf("Batch processing failed: %v", err)
		return
	}

	fmt.Printf("\n✓ Batch processing complete\n")
	fmt.Printf("  Total latency: %v\n", result.Latency)
	fmt.Printf("  Items processed: %d\n", len(texts))
	
	if results, ok := result.Output.([]interface{}); ok {
		fmt.Println("\n  Results:")
		for i, r := range results {
			if m, ok := r.(map[string]interface{}); ok {
				fmt.Printf("    %d. %s → %s\n", i+1, m["text"], m["label"])
			}
		}
	}
}

// MockProvider for demonstration
type MockProvider struct{}

func (m *MockProvider) LoadModel(config *ai.ModelConfig) (*ai.Model, error) {
	return &ai.Model{
		ID:       config.ID,
		Name:     config.Name,
		Version:  config.Version,
		Type:     config.Type,
		Status:   ai.ModelStatusReady,
		Provider: "mock",
		LoadedAt: time.Now(),
	}, nil
}

func (m *MockProvider) UnloadModel(modelID string) error {
	return nil
}

func (m *MockProvider) Predict(ctx context.Context, modelID string, input *ai.InferenceInput) (*ai.InferenceOutput, error) {
	// Mock predictions based on model type
	var result interface{}

	switch {
	case modelID == "classifier-v1":
		result = map[string]interface{}{
			"label":      "positive",
			"confidence": 0.95,
		}
	case modelID == "sentiment-model":
		result = map[string]interface{}{
			"sentiment": "positive",
			"score":     0.85,
		}
	default:
		result = map[string]interface{}{
			"prediction": "mock_result",
		}
	}

	return &ai.InferenceOutput{
		ModelID:   modelID,
		Result:    result,
		Latency:   50 * time.Millisecond,
		Timestamp: time.Now(),
	}, nil
}

func (m *MockProvider) GetMetrics(modelID string) *ai.ModelMetrics {
	return &ai.ModelMetrics{
		ModelID:      modelID,
		RequestCount: 10,
		AvgLatency:   50 * time.Millisecond,
	}
}
