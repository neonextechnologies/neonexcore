package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/neonexcore/pkg/workflow"
)

func main() {
	fmt.Println("=== Workflow Engine Examples ===\n")

	// Example 1: Basic Workflow
	basicWorkflowExample()

	// Example 2: YAML Workflow
	yamlWorkflowExample()

	// Example 3: Conditional Execution
	conditionalExecutionExample()

	// Example 4: Loop Execution
	loopExecutionExample()

	// Example 5: Parallel Execution
	parallelExecutionExample()

	// Example 6: State Persistence (Demo structure)
	statePersistenceExample()
}

// Example 1: Basic Workflow
func basicWorkflowExample() {
	fmt.Println("1. Basic Workflow Example")
	fmt.Println("--------------------------")

	engine := workflow.NewWorkflowEngine()

	// Build workflow using fluent API
	wf := workflow.NewWorkflowBuilder("order-processing").
		Description("Process customer orders").
		Version("1.0.0").
		AddStep("validate", "Validate Order").
		Action(func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("  → Validating order...")
			orderID, _ := execCtx.Get("order_id")
			fmt.Printf("  → Order ID: %v\n", orderID)
			return map[string]interface{}{"valid": true}, nil
		}).
		Retry(3, 1*time.Second, 2.0).
		Timeout(30 * time.Second).
		Then("process", "Process Payment").
		Action(func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("  → Processing payment...")
			amount, _ := execCtx.Get("amount")
			fmt.Printf("  → Amount: $%.2f\n", amount.(float64))
			return map[string]interface{}{"payment_id": "PAY123"}, nil
		}).
		Then("notify", "Send Notification").
		Action(func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("  → Sending notification...")
			return map[string]interface{}{"sent": true}, nil
		}).
		End().
		Build()

	engine.RegisterWorkflow(wf)

	execution, err := engine.StartExecution(context.Background(), wf.ID, map[string]interface{}{
		"order_id": "ORD123",
		"amount":   99.99,
	})

	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	// Wait for completion
	time.Sleep(2 * time.Second)

	fmt.Printf("  ✓ Execution Status: %s\n", execution.Status)
	fmt.Printf("  ✓ Steps Completed: %d\n\n", len(execution.StepResults))
}

// Example 2: YAML Workflow
func yamlWorkflowExample() {
	fmt.Println("2. YAML Workflow Example")
	fmt.Println("------------------------")

	yamlContent := `
name: data-processing
description: Process data pipeline
version: 1.0.0
config:
  timeout: 300s

steps:
  - id: extract
    name: Extract Data
    type: task
    action_type: extract_data
    timeout: 30s
    on_success:
      - transform

  - id: transform
    name: Transform Data
    type: task
    action_type: transform_data
    timeout: 60s
    retry:
      max_attempts: 3
      delay: 2s
      backoff_rate: 2.0
    on_success:
      - load

  - id: load
    name: Load Data
    type: task
    action_type: load_data
    timeout: 30s
`

	actionRegistry := map[string]workflow.ActionFunc{
		"extract_data": func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("  → Extracting data...")
			return map[string]interface{}{"records": 1000}, nil
		},
		"transform_data": func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("  → Transforming data...")
			return map[string]interface{}{"transformed": 1000}, nil
		},
		"load_data": func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("  → Loading data...")
			return map[string]interface{}{"loaded": 1000}, nil
		},
	}

	wf, err := workflow.FromYAML([]byte(yamlContent), actionRegistry)
	if err != nil {
		log.Printf("Error parsing YAML: %v\n", err)
		return
	}

	engine := workflow.NewWorkflowEngine()
	engine.RegisterWorkflow(wf)

	execution, err := engine.StartExecution(context.Background(), wf.ID, map[string]interface{}{
		"source": "database",
	})

	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	time.Sleep(2 * time.Second)
	fmt.Printf("  ✓ Pipeline Status: %s\n\n", execution.Status)
}

// Example 3: Conditional Execution
func conditionalExecutionExample() {
	fmt.Println("3. Conditional Execution Example")
	fmt.Println("---------------------------------")

	ctx := context.Background()
	execCtx := &workflow.ExecutionContext{
		Variables:   map[string]interface{}{"amount": 150.0},
		StepResults: make(map[string]interface{}),
		Metadata:    make(map[string]string),
	}

	condExecutor := workflow.NewConditionalExecutor()

	// If-Then-Else
	fmt.Println("  → If-Then-Else:")
	condition := func(execCtx *workflow.ExecutionContext) (bool, error) {
		amount, _ := execCtx.Get("amount")
		return amount.(float64) > 100, nil
	}

	thenStep := workflow.Step{
		ID:   "premium",
		Name: "Premium Processing",
		Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("    → Executing premium service")
			return "Premium service", nil
		},
	}

	elseStep := workflow.Step{
		ID:   "standard",
		Name: "Standard Processing",
		Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("    → Executing standard service")
			return "Standard service", nil
		},
	}

	result := condExecutor.IfThenElse(ctx, condition, thenStep, &elseStep, execCtx)
	fmt.Printf("    ✓ Result: %v\n", result.Output)

	// Switch Statement
	fmt.Println("  → Switch Statement:")
	memberType := "gold"
	
	goldStep := workflow.Step{
		ID:   "gold",
		Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("    → Gold member benefits applied")
			return "20% discount", nil
		},
	}

	silverStep := workflow.Step{
		ID:   "silver",
		Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			fmt.Println("    → Silver member benefits applied")
			return "10% discount", nil
		},
	}

	cases := map[interface{}]workflow.Step{
		"gold":   goldStep,
		"silver": silverStep,
	}

	result = condExecutor.Switch(ctx, memberType, cases, nil, execCtx)
	fmt.Printf("    ✓ Discount: %v\n\n", result.Output)
}

// Example 4: Loop Execution
func loopExecutionExample() {
	fmt.Println("4. Loop Execution Example")
	fmt.Println("-------------------------")

	ctx := context.Background()
	execCtx := &workflow.ExecutionContext{
		Variables:   make(map[string]interface{}),
		StepResults: make(map[string]interface{}),
		Metadata:    make(map[string]string),
	}

	loopExecutor := workflow.NewLoopExecutor()

	// ForEach Loop
	fmt.Println("  → ForEach Loop:")
	items := []interface{}{"task1", "task2", "task3"}
	
	step := workflow.Step{
		ID:   "process_item",
		Name: "Process Item",
		Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			item, _ := execCtx.Get("current_item")
			index, _ := execCtx.Get("current_index")
			fmt.Printf("    → Processing item %d: %v\n", index, item)
			return item, nil
		},
	}

	results := loopExecutor.ForEach(ctx, step, items, execCtx)
	fmt.Printf("    ✓ Processed %d items\n", len(results))

	// While Loop
	fmt.Println("  → While Loop:")
	execCtx.Set("counter", 0)
	
	whileCondition := func(execCtx *workflow.ExecutionContext) (bool, error) {
		counter, _ := execCtx.Get("counter")
		return counter.(int) < 3, nil
	}

	whileStep := workflow.Step{
		ID:   "increment",
		Name: "Increment Counter",
		Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
			counter, _ := execCtx.Get("counter")
			newCounter := counter.(int) + 1
			execCtx.Set("counter", newCounter)
			fmt.Printf("    → Counter: %d\n", newCounter)
			return newCounter, nil
		},
	}

	results = loopExecutor.While(ctx, whileStep, whileCondition, execCtx, 10)
	fmt.Printf("    ✓ Loop iterations: %d\n\n", len(results))
}

// Example 5: Parallel Execution
func parallelExecutionExample() {
	fmt.Println("5. Parallel Execution Example")
	fmt.Println("-----------------------------")

	ctx := context.Background()
	execCtx := &workflow.ExecutionContext{
		Variables:   make(map[string]interface{}),
		StepResults: make(map[string]interface{}),
		Metadata:    make(map[string]string),
	}

	parallelExecutor := workflow.NewParallelExecutor(3)

	steps := []workflow.Step{
		{
			ID:   "task1",
			Name: "Fetch User Data",
			Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
				fmt.Println("  → Fetching user data...")
				time.Sleep(500 * time.Millisecond)
				return map[string]interface{}{"user": "john"}, nil
			},
		},
		{
			ID:   "task2",
			Name: "Fetch Orders",
			Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
				fmt.Println("  → Fetching orders...")
				time.Sleep(300 * time.Millisecond)
				return map[string]interface{}{"orders": 5}, nil
			},
		},
		{
			ID:   "task3",
			Name: "Fetch Analytics",
			Action: func(ctx context.Context, execCtx *workflow.ExecutionContext) (interface{}, error) {
				fmt.Println("  → Fetching analytics...")
				time.Sleep(400 * time.Millisecond)
				return map[string]interface{}{"views": 1000}, nil
			},
		},
	}

	startTime := time.Now()
	results := parallelExecutor.Execute(ctx, steps, execCtx)
	duration := time.Since(startTime)

	fmt.Printf("  ✓ Completed %d tasks in parallel\n", len(results))
	fmt.Printf("  ✓ Total duration: %v\n\n", duration)
}

// Example 6: State Persistence (Demo structure)
func statePersistenceExample() {
	fmt.Println("6. State Persistence Example")
	fmt.Println("----------------------------")

	fmt.Println("  → State persistence enables:")
	fmt.Println("    • Save workflow state to database")
	fmt.Println("    • Resume paused/failed executions")
	fmt.Println("    • Query execution history")
	fmt.Println("    • Event logging and tracking")
	fmt.Println("    • Cleanup old completed states")
	
	fmt.Println("\n  Example usage:")
	fmt.Println("    // Create state store")
	fmt.Println("    stateStore, _ := workflow.NewStateStore(db)")
	fmt.Println("    ")
	fmt.Println("    // Create stateful engine")
	fmt.Println("    engine := workflow.NewStatefulWorkflowEngine(stateStore)")
	fmt.Println("    ")
	fmt.Println("    // Start execution (auto-saved)")
	fmt.Println("    execution, _ := engine.StartExecution(ctx, workflowID, input)")
	fmt.Println("    ")
	fmt.Println("    // Resume execution")
	fmt.Println("    engine.ResumeExecution(ctx, executionID)")
	fmt.Println("    ")
	fmt.Println("    // Query states")
	fmt.Println("    states, _ := stateStore.ListStates(workflowID, StatusRunning, 10)")
	fmt.Println("    ")
	fmt.Println("    // Get event logs")
	fmt.Println("    events, _ := stateStore.GetEvents(executionID, 100)")
	
	fmt.Println("\n  ✓ Requires database connection (PostgreSQL/MySQL)")
	fmt.Println()
}
