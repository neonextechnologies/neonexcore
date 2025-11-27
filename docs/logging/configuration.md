# Logger Configuration

Complete guide to configuring the logging system in Neonex Core.

---

## Overview

Neonex Core's logger is highly configurable through **environment variables** and **programmatic configuration**, supporting various formats, outputs, and rotation strategies.

---

## Environment Variables

### Basic Configuration

```bash
# .env or config/.env

# Log Level (debug, info, warn, error, fatal)
LOG_LEVEL=info

# Log Format (json, text, console)
LOG_FORMAT=json

# Log Output (stdout, stderr, file)
LOG_OUTPUT=stdout
```

### File Logging Configuration

```bash
# Enable file logging
LOG_OUTPUT=file

# File path
LOG_FILE_PATH=logs/neonex.log

# Rotation settings
LOG_FILE_MAX_SIZE=100        # Maximum size in MB before rotation
LOG_FILE_MAX_BACKUPS=5       # Number of old log files to keep
LOG_FILE_MAX_AGE=30          # Maximum days to retain old logs
LOG_FILE_COMPRESS=true       # Compress rotated files (gzip)

# File permissions (Unix)
LOG_FILE_PERMISSIONS=0644
```

### Advanced Settings

```bash
# Caller information (file:line)
LOG_ENABLE_CALLER=false

# Stack trace for errors
LOG_ENABLE_STACKTRACE=true

# Development mode (prettier output)
LOG_DEVELOPMENT=false

# Sampling (log every Nth message at high volume)
LOG_SAMPLING_ENABLED=false
LOG_SAMPLING_INITIAL=100
LOG_SAMPLING_THEREAFTER=100
```

---

## Configuration Struct

### Logger Config

```go
// internal/config/logger.go
package config

import "github.com/caarlos0/env/v11"

type LoggerConfig struct {
    // Core settings
    Level  string `env:"LOG_LEVEL" envDefault:"info"`
    Format string `env:"LOG_FORMAT" envDefault:"json"`
    Output string `env:"LOG_OUTPUT" envDefault:"stdout"`
    
    // File settings
    File FileLogConfig
    
    // Advanced settings
    EnableCaller     bool `env:"LOG_ENABLE_CALLER" envDefault:"false"`
    EnableStacktrace bool `env:"LOG_ENABLE_STACKTRACE" envDefault:"true"`
    Development      bool `env:"LOG_DEVELOPMENT" envDefault:"false"`
    
    // Sampling
    Sampling SamplingConfig
}

type FileLogConfig struct {
    Path        string `env:"LOG_FILE_PATH" envDefault:"logs/neonex.log"`
    MaxSize     int    `env:"LOG_FILE_MAX_SIZE" envDefault:"100"`
    MaxBackups  int    `env:"LOG_FILE_MAX_BACKUPS" envDefault:"5"`
    MaxAge      int    `env:"LOG_FILE_MAX_AGE" envDefault:"30"`
    Compress    bool   `env:"LOG_FILE_COMPRESS" envDefault:"true"`
    Permissions int    `env:"LOG_FILE_PERMISSIONS" envDefault:"0644"`
}

type SamplingConfig struct {
    Enabled    bool `env:"LOG_SAMPLING_ENABLED" envDefault:"false"`
    Initial    int  `env:"LOG_SAMPLING_INITIAL" envDefault:"100"`
    Thereafter int  `env:"LOG_SAMPLING_THEREAFTER" envDefault:"100"`
}

// Load configuration
func LoadLoggerConfig() (*LoggerConfig, error) {
    cfg := &LoggerConfig{}
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}
```

---

## Log Levels

### Level Priority

```go
const (
    DebugLevel = "debug"  // -1: Verbose debugging
    InfoLevel  = "info"   //  0: General info (default)
    WarnLevel  = "warn"   //  1: Warning messages
    ErrorLevel = "error"  //  2: Error conditions
    FatalLevel = "fatal"  //  3: Fatal errors (exits)
)
```

### Level Hierarchy

```
┌─────────────────────────────────────────────────┐
│                                                 │
│  DEBUG → INFO → WARN → ERROR → FATAL           │
│    ↓       ↓      ↓       ↓       ↓            │
│  Most          Medium           Least           │
│  Verbose                       Verbose          │
│                                                 │
└─────────────────────────────────────────────────┘
```

**When `LOG_LEVEL=info`:**
- ✅ Info, Warn, Error, Fatal logs shown
- ❌ Debug logs hidden

**When `LOG_LEVEL=debug`:**
- ✅ All logs shown

### Environment-Specific Levels

```bash
# Development
LOG_LEVEL=debug

# Staging
LOG_LEVEL=info

# Production
LOG_LEVEL=warn
```

---

## Log Formats

### JSON Format

**Best for:** Production, log aggregation, parsing

```bash
LOG_FORMAT=json
```

**Output:**
```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:45.123Z",
  "caller": "user/service.go:42",
  "message": "User created",
  "user_id": 12345,
  "email": "user@example.com"
}
```

**Configuration:**
```go
// Zap JSON encoder
encoderConfig := zapcore.EncoderConfig{
    TimeKey:        "timestamp",
    LevelKey:       "level",
    NameKey:        "logger",
    CallerKey:      "caller",
    FunctionKey:    zapcore.OmitKey,
    MessageKey:     "message",
    StacktraceKey:  "stacktrace",
    LineEnding:     zapcore.DefaultLineEnding,
    EncodeLevel:    zapcore.LowercaseLevelEncoder,
    EncodeTime:     zapcore.ISO8601TimeEncoder,
    EncodeDuration: zapcore.SecondsDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,
}
```

### Text/Console Format

**Best for:** Development, human readability

```bash
LOG_FORMAT=text
# or
LOG_FORMAT=console
```

**Output:**
```
2024-01-15T10:30:45.123Z  INFO  user/service.go:42  User created  user_id=12345 email=user@example.com
```

**Configuration:**
```go
// Zap console encoder
encoderConfig := zapcore.EncoderConfig{
    TimeKey:        "T",
    LevelKey:       "L",
    NameKey:        "N",
    CallerKey:      "C",
    FunctionKey:    zapcore.OmitKey,
    MessageKey:     "M",
    StacktraceKey:  "S",
    LineEnding:     zapcore.DefaultLineEnding,
    EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Colored
    EncodeTime:     zapcore.ISO8601TimeEncoder,
    EncodeDuration: zapcore.StringDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,
}
```

---

## Output Destinations

### Console Output

#### Stdout (Default)

```bash
LOG_OUTPUT=stdout
```

**Use cases:**
- Cloud platforms (Kubernetes, Docker)
- CI/CD pipelines
- Development

#### Stderr

```bash
LOG_OUTPUT=stderr
```

**Use cases:**
- Separate errors from normal output
- Unix pipe operations

### File Output

#### Basic File Logging

```bash
LOG_OUTPUT=file
LOG_FILE_PATH=logs/neonex.log
```

#### Rotating File Logs

```bash
LOG_OUTPUT=file
LOG_FILE_PATH=logs/neonex.log
LOG_FILE_MAX_SIZE=100         # Rotate at 100MB
LOG_FILE_MAX_BACKUPS=5        # Keep 5 old files
LOG_FILE_MAX_AGE=30           # Delete files older than 30 days
LOG_FILE_COMPRESS=true        # Compress old files
```

**Directory structure:**
```
logs/
├── neonex.log                    # Current log
├── neonex-2024-01-14.log        # Yesterday
├── neonex-2024-01-13.log.gz     # Compressed
├── neonex-2024-01-12.log.gz
└── neonex-2024-01-11.log.gz
```

#### Implementation

```go
// pkg/logger/file.go
package logger

import (
    "gopkg.in/natefinch/lumberjack.v2"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewFileLogger(config *FileLogConfig) (*zap.Logger, error) {
    // File writer with rotation
    writer := &lumberjack.Logger{
        Filename:   config.Path,
        MaxSize:    config.MaxSize,    // MB
        MaxBackups: config.MaxBackups,
        MaxAge:     config.MaxAge,     // days
        Compress:   config.Compress,
    }
    
    // Encoder
    encoderConfig := zap.NewProductionEncoderConfig()
    encoder := zapcore.NewJSONEncoder(encoderConfig)
    
    // Core
    core := zapcore.NewCore(
        encoder,
        zapcore.AddSync(writer),
        zap.InfoLevel,
    )
    
    return zap.New(core), nil
}
```

### Multiple Outputs

```go
// Write to both console and file
func NewMultiLogger(config *LoggerConfig) (*zap.Logger, error) {
    // Console writer
    consoleWriter := zapcore.Lock(os.Stdout)
    consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
    
    // File writer
    fileWriter := &lumberjack.Logger{
        Filename:   config.File.Path,
        MaxSize:    config.File.MaxSize,
        MaxBackups: config.File.MaxBackups,
        MaxAge:     config.File.MaxAge,
        Compress:   config.File.Compress,
    }
    fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
    
    // Combine cores
    core := zapcore.NewTee(
        zapcore.NewCore(consoleEncoder, consoleWriter, zap.DebugLevel),
        zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), zap.InfoLevel),
    )
    
    return zap.New(core, zap.AddCaller()), nil
}
```

---

## Initialization

### Bootstrap Logger

```go
// cmd/neonex/main.go
package main

import (
    "neonexcore/internal/config"
    "neonexcore/pkg/logger"
)

func main() {
    // Load logger config
    loggerCfg, err := config.LoadLoggerConfig()
    if err != nil {
        panic("Failed to load logger config: " + err.Error())
    }
    
    // Initialize logger
    log, err := logger.New(loggerCfg)
    if err != nil {
        panic("Failed to initialize logger: " + err.Error())
    }
    defer log.Sync()
    
    // Use logger
    log.Info("Application starting", logger.Fields{
        "version": "0.1.0",
        "env":     "production",
    })
    
    // ... rest of application
}
```

### Logger Factory

```go
// pkg/logger/factory.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func New(config *Config) (Logger, error) {
    // Parse log level
    level, err := zapcore.ParseLevel(config.Level)
    if err != nil {
        return nil, err
    }
    
    // Build encoder config
    encoderConfig := buildEncoderConfig(config)
    
    // Build encoder
    var encoder zapcore.Encoder
    if config.Format == "json" {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }
    
    // Build writer
    writer := buildWriter(config)
    
    // Build core
    core := zapcore.NewCore(encoder, writer, level)
    
    // Apply sampling if enabled
    if config.Sampling.Enabled {
        core = zapcore.NewSamplerWithOptions(
            core,
            time.Second,
            config.Sampling.Initial,
            config.Sampling.Thereafter,
        )
    }
    
    // Build logger
    zapLogger := zap.New(core)
    
    if config.EnableCaller {
        zapLogger = zapLogger.WithOptions(zap.AddCaller())
    }
    
    if config.EnableStacktrace {
        zapLogger = zapLogger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
    }
    
    return &zapLoggerImpl{logger: zapLogger}, nil
}

func buildWriter(config *Config) zapcore.WriteSyncer {
    switch config.Output {
    case "stdout":
        return zapcore.Lock(os.Stdout)
    case "stderr":
        return zapcore.Lock(os.Stderr)
    case "file":
        writer := &lumberjack.Logger{
            Filename:   config.File.Path,
            MaxSize:    config.File.MaxSize,
            MaxBackups: config.File.MaxBackups,
            MaxAge:     config.File.MaxAge,
            Compress:   config.File.Compress,
        }
        return zapcore.AddSync(writer)
    default:
        return zapcore.Lock(os.Stdout)
    }
}
```

---

## Development vs Production

### Development Configuration

```bash
# .env.development
LOG_LEVEL=debug
LOG_FORMAT=console
LOG_OUTPUT=stdout
LOG_ENABLE_CALLER=true
LOG_DEVELOPMENT=true
```

**Features:**
- Colored output
- Human-readable format
- Verbose logging
- Caller information

### Production Configuration

```bash
# .env.production
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=file
LOG_FILE_PATH=/var/log/neonex/app.log
LOG_FILE_MAX_SIZE=100
LOG_FILE_MAX_BACKUPS=10
LOG_FILE_MAX_AGE=30
LOG_FILE_COMPRESS=true
LOG_ENABLE_STACKTRACE=true
LOG_SAMPLING_ENABLED=true
```

**Features:**
- JSON format (parsing)
- File rotation
- Error stack traces
- Log sampling (high volume)

---

## Dynamic Configuration

### Runtime Level Changes

```go
// Create atomic level
atomicLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)

// Build logger with atomic level
core := zapcore.NewCore(encoder, writer, atomicLevel)
logger := zap.New(core)

// Change level at runtime
atomicLevel.SetLevel(zapcore.DebugLevel)

// HTTP endpoint to change level
http.HandleFunc("/log/level", func(w http.ResponseWriter, r *http.Request) {
    level := r.URL.Query().Get("level")
    
    var zapLevel zapcore.Level
    if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
        http.Error(w, "Invalid level", http.StatusBadRequest)
        return
    }
    
    atomicLevel.SetLevel(zapLevel)
    w.Write([]byte("Log level changed to " + level))
})
```

---

## Testing Configuration

### Test Logger

```go
// Test with memory buffer
func setupTestLogger(t *testing.T) (*zap.Logger, *observer.ObservedLogs) {
    core, logs := observer.New(zap.DebugLevel)
    logger := zap.New(core)
    return logger, logs
}

func TestLogging(t *testing.T) {
    logger, logs := setupTestLogger(t)
    
    logger.Info("test message", zap.String("key", "value"))
    
    assert.Equal(t, 1, logs.Len())
    assert.Equal(t, "test message", logs.All()[0].Message)
}
```

---

## Best Practices

### ✅ DO:

**1. Use Environment-Specific Configs**
```bash
# Development
LOG_LEVEL=debug LOG_FORMAT=console

# Production
LOG_LEVEL=info LOG_FORMAT=json
```

**2. Enable File Rotation**
```bash
LOG_FILE_MAX_SIZE=100
LOG_FILE_MAX_BACKUPS=5
LOG_FILE_COMPRESS=true
```

**3. Configure Sampling for High Volume**
```bash
LOG_SAMPLING_ENABLED=true
LOG_SAMPLING_INITIAL=100
LOG_SAMPLING_THEREAFTER=100
```

### ❌ DON'T:

**1. Debug in Production**
```bash
# Bad
LOG_LEVEL=debug  # ❌ Too verbose

# Good
LOG_LEVEL=info
```

**2. Unlimited Log Files**
```bash
# Bad
LOG_FILE_MAX_BACKUPS=0  # ❌ Fills disk

# Good
LOG_FILE_MAX_BACKUPS=10
LOG_FILE_MAX_AGE=30
```

---

## Next Steps

- [**Usage**](usage.md) - Logging in your code
- [**Middleware**](middleware.md) - HTTP request logging
- [**Best Practices**](best-practices.md) - Logging guidelines
- [**Overview**](overview.md) - Logger architecture

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
