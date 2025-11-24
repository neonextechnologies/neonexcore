package main

import (
	"neonexcore/pkg/logger"
	"time"
)

func main() {
	// Test 1: Default text formatter (console)
	println("\n=== Test 1: Text Format (Console) ===")
	logger.Info("Application starting...")
	logger.Debug("Debug message - might not show (default level is INFO)")
	logger.Warn("This is a warning", logger.Fields{"user": "admin"})
	logger.Error("This is an error", logger.Fields{
		"error_code": 500,
		"timestamp":  time.Now().Unix(),
	})

	// Test 2: JSON formatter
	println("\n=== Test 2: JSON Format ===")
	logger.SetGlobalFormatter(logger.NewJSONFormatter())
	logger.Info("JSON formatted log", logger.Fields{
		"service": "test",
		"version": "1.0.0",
	})

	// Test 3: With fields
	println("\n=== Test 3: Logger with Fields ===")
	logger.SetGlobalFormatter(logger.NewTextFormatter())
	userLogger := logger.With(logger.Fields{
		"module": "user",
		"action": "login",
	})
	userLogger.Info("User logged in", logger.Fields{"user_id": 123})
	userLogger.Warn("Login attempt from suspicious IP")

	// Test 4: Different log levels
	println("\n=== Test 4: Different Log Levels ===")
	logger.SetGlobalLevel(logger.DebugLevel)
	logger.Debug("Now debug messages will show")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// Test 5: File output
	println("\n=== Test 5: File Output ===")
	fileWriter, err := logger.NewFileWriter(logger.FileWriterConfig{
		Filename:   "logs/test.log",
		MaxSize:    1, // 1MB
		MaxBackups: 3,
		MaxAge:     7,
	})
	if err != nil {
		logger.Error("Failed to create file writer", logger.Fields{"error": err.Error()})
	} else {
		logger.AddGlobalWriter(fileWriter)
		logger.Info("This message will be written to both console and file")
		logger.Info("Log file created at: logs/test.log")
	}

	// Test 6: Structured logging
	println("\n=== Test 6: Structured Logging ===")
	requestLogger := logger.With(logger.Fields{
		"request_id": "req-12345",
		"ip":         "192.168.1.1",
	})
	requestLogger.Info("Request received", logger.Fields{
		"method": "POST",
		"path":   "/api/users",
		"status": 201,
	})

	println("\n=== All Tests Completed ===")
}
