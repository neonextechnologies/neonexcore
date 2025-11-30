# Monitoring

Comprehensive monitoring and observability setup for Neonex Core applications.

---

## Overview

Production monitoring provides:
- **Metrics** - Performance and health data
- **Logging** - Centralized log aggregation
- **Tracing** - Distributed request tracing
- **Alerting** - Proactive issue detection

---

## Metrics Collection

### Prometheus Integration

```go
// internal/middleware/prometheus.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/gofiber/adaptor/v2"
    "strconv"
    "time"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
    
    dbQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"operation", "table"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(dbQueryDuration)
}

func PrometheusMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        err := c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Response().StatusCode())
        
        httpRequestsTotal.WithLabelValues(
            c.Method(),
            c.Path(),
            status,
        ).Inc()
        
        httpRequestDuration.WithLabelValues(
            c.Method(),
            c.Path(),
        ).Observe(duration)
        
        return err
    }
}

func RegisterPrometheusEndpoint(app *fiber.App) {
    app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}
```

### Register Metrics Middleware

```go
// main.go
func main() {
    app := core.NewApp()
    
    // Register Prometheus middleware
    app.HTTP.Use(middleware.PrometheusMiddleware())
    
    // Expose metrics endpoint
    middleware.RegisterPrometheusEndpoint(app.HTTP)
    
    app.Run()
}
```

### Custom Business Metrics

```go
// metrics/business.go
package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
    UserRegistrations = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "user_registrations_total",
            Help: "Total number of user registrations",
        },
    )
    
    OrdersCreated = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orders_created_total",
            Help: "Total number of orders created",
        },
        []string{"status"},
    )
    
    ActiveSessions = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_sessions",
            Help: "Number of active user sessions",
        },
    )
)

func init() {
    prometheus.MustRegister(UserRegistrations)
    prometheus.MustRegister(OrdersCreated)
    prometheus.MustRegister(ActiveSessions)
}

// Usage in service
func (s *UserService) Register(ctx context.Context, data RegisterData) error {
    // ... registration logic ...
    
    metrics.UserRegistrations.Inc()
    
    return nil
}
```

---

## Prometheus Setup

### Docker Compose

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
      - '--web.console.templates=/usr/share/prometheus/consoles'
    ports:
      - "9090:9090"
    restart: unless-stopped

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    ports:
      - "3000:3000"
    restart: unless-stopped
    depends_on:
      - prometheus

volumes:
  prometheus_data:
  grafana_data:
```

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'production'
    environment: 'prod'

scrape_configs:
  - job_name: 'neonex-app'
    static_configs:
      - targets: ['app:8080']
        labels:
          service: 'api'
    
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']
  
  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - 'alerts.yml'
```

### Alert Rules

```yaml
# alerts.yml
groups:
  - name: application
    interval: 30s
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} requests/sec"
      
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "95th percentile latency is {{ $value }}s"
      
      - alert: DatabaseConnectionPoolExhausted
        expr: db_connections_in_use / db_max_connections > 0.9
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Database connection pool almost exhausted"
          description: "{{ $value | humanizePercentage }} of connections in use"
      
      - alert: ServiceDown
        expr: up{job="neonex-app"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service is down"
          description: "{{ $labels.instance }} is not responding"
```

---

## Grafana Dashboards

### Dashboard Provisioning

```yaml
# grafana/provisioning/dashboards/dashboard.yml
apiVersion: 1

providers:
  - name: 'Neonex Dashboards'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /etc/grafana/provisioning/dashboards
```

### Data Source Provisioning

```yaml
# grafana/provisioning/datasources/prometheus.yml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
```

### Application Dashboard JSON

```json
{
  "dashboard": {
    "title": "Neonex Application Metrics",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{path}}"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "5xx errors"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Response Time (p95)",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "{{path}}"
          }
        ],
        "type": "graph"
      }
    ]
  }
}
```

---

## Log Aggregation

### ELK Stack (Elasticsearch, Logstash, Kibana)

```yaml
# docker-compose.elk.yml
version: '3.8'

services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.10.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
      - xpack.security.enabled=false
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data
    ports:
      - "9200:9200"
    restart: unless-stopped

  logstash:
    image: docker.elastic.co/logstash/logstash:8.10.0
    volumes:
      - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf
    depends_on:
      - elasticsearch
    ports:
      - "5000:5000"
    restart: unless-stopped

  kibana:
    image: docker.elastic.co/kibana/kibana:8.10.0
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
    restart: unless-stopped

volumes:
  elasticsearch_data:
```

### Logstash Configuration

```ruby
# logstash.conf
input {
  tcp {
    port => 5000
    codec => json
  }
  
  file {
    path => "/var/log/neonex/*.log"
    start_position => "beginning"
    codec => json
  }
}

filter {
  json {
    source => "message"
  }
  
  date {
    match => ["timestamp", "ISO8601"]
  }
  
  if [level] == "error" or [level] == "fatal" {
    mutate {
      add_tag => ["alert"]
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "neonex-%{+YYYY.MM.dd}"
  }
  
  if "alert" in [tags] {
    http {
      url => "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
      http_method => "post"
      content_type => "application/json"
      format => "message"
      message => '{"text": "Error: %{message}"}'
    }
  }
}
```

### Structured Logging for ELK

```go
// Configure Zap to output JSON
config := zap.NewProductionConfig()
config.OutputPaths = []string{
    "stdout",
    "/var/log/neonex/app.log",
}

logger, _ := config.Build()

// Add context fields
logger.Info("User registered",
    zap.String("user_id", userID),
    zap.String("email", email),
    zap.String("ip", clientIP),
    zap.Duration("elapsed", elapsed),
)
```

---

## Distributed Tracing

### OpenTelemetry Setup

```go
// internal/tracing/tracer.go
package tracing

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer(serviceName string) (*trace.TracerProvider, error) {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    if err != nil {
        return nil, err
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
            semconv.DeploymentEnvironmentKey.String("production"),
        )),
    )
    
    otel.SetTracerProvider(tp)
    
    return tp, nil
}
```

### Tracing Middleware

```go
// internal/middleware/tracing.go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

func TracingMiddleware() fiber.Handler {
    tracer := otel.Tracer("neonex-http")
    
    return func(c *fiber.Ctx) error {
        ctx, span := tracer.Start(c.UserContext(), c.Path())
        defer span.End()
        
        span.SetAttributes(
            attribute.String("http.method", c.Method()),
            attribute.String("http.path", c.Path()),
            attribute.String("http.client_ip", c.IP()),
        )
        
        c.SetUserContext(ctx)
        
        err := c.Next()
        
        span.SetAttributes(
            attribute.Int("http.status_code", c.Response().StatusCode()),
        )
        
        if err != nil {
            span.RecordError(err)
        }
        
        return err
    }
}
```

### Service Tracing

```go
func (s *UserService) Create(ctx context.Context, data CreateUserData) (*User, error) {
    ctx, span := otel.Tracer("user-service").Start(ctx, "CreateUser")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("user.email", data.Email),
    )
    
    // Validate
    ctx, validateSpan := otel.Tracer("user-service").Start(ctx, "ValidateUser")
    if err := s.validator.Validate(data); err != nil {
        validateSpan.RecordError(err)
        validateSpan.End()
        return nil, err
    }
    validateSpan.End()
    
    // Create in DB
    ctx, dbSpan := otel.Tracer("user-service").Start(ctx, "SaveUserToDB")
    user, err := s.repo.Create(ctx, data)
    if err != nil {
        dbSpan.RecordError(err)
        dbSpan.End()
        return nil, err
    }
    dbSpan.End()
    
    return user, nil
}
```

### Jaeger Docker Setup

```yaml
services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    environment:
      - COLLECTOR_ZIPKIN_HOST_PORT=:9411
    ports:
      - "5775:5775/udp"
      - "6831:6831/udp"
      - "6832:6832/udp"
      - "5778:5778"
      - "16686:16686"
      - "14268:14268"
      - "14250:14250"
      - "9411:9411"
    restart: unless-stopped
```

---

## Health Checks

### Health Endpoint

```go
// internal/health/handler.go
package health

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    "time"
)

type HealthHandler struct {
    db *gorm.DB
}

type HealthResponse struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Services  map[string]string `json:"services"`
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
    resp := HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Services:  make(map[string]string),
    }
    
    // Check database
    sqlDB, _ := h.db.DB()
    if err := sqlDB.Ping(); err != nil {
        resp.Services["database"] = "unhealthy"
        resp.Status = "unhealthy"
    } else {
        resp.Services["database"] = "healthy"
    }
    
    // Check Redis (if using)
    // if err := redisClient.Ping(c.Context()).Err(); err != nil {
    //     resp.Services["redis"] = "unhealthy"
    //     resp.Status = "unhealthy"
    // } else {
    //     resp.Services["redis"] = "healthy"
    // }
    
    status := fiber.StatusOK
    if resp.Status == "unhealthy" {
        status = fiber.StatusServiceUnavailable
    }
    
    return c.Status(status).JSON(resp)
}

func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
    // More strict check for Kubernetes readiness
    sqlDB, _ := h.db.DB()
    if err := sqlDB.Ping(); err != nil {
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
            "ready": false,
            "error": "database not ready",
        })
    }
    
    return c.JSON(fiber.Map{"ready": true})
}

func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
    // Simple liveness check
    return c.JSON(fiber.Map{"alive": true})
}
```

---

## Alerting

### AlertManager Configuration

```yaml
# alertmanager.yml
global:
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: 'alerts@example.com'
  smtp_auth_username: 'alerts@example.com'
  smtp_auth_password: 'password'

route:
  receiver: 'team-notifications'
  group_by: ['alertname', 'cluster']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  
  routes:
    - match:
        severity: critical
      receiver: 'pagerduty'
      continue: true
    
    - match:
        severity: warning
      receiver: 'slack'

receivers:
  - name: 'team-notifications'
    email_configs:
      - to: 'team@example.com'
  
  - name: 'slack'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'
        channel: '#alerts'
        title: '{{ .GroupLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
  
  - name: 'pagerduty'
    pagerduty_configs:
      - service_key: 'YOUR_PAGERDUTY_KEY'
```

### Slack Notifications

```go
// internal/alerting/slack.go
package alerting

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type SlackMessage struct {
    Text        string       `json:"text"`
    Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
    Color  string `json:"color"`
    Title  string `json:"title"`
    Text   string `json:"text"`
    Fields []Field `json:"fields"`
}

type Field struct {
    Title string `json:"title"`
    Value string `json:"value"`
    Short bool   `json:"short"`
}

func SendSlackAlert(webhookURL, message string, level string) error {
    color := "good"
    if level == "error" || level == "critical" {
        color = "danger"
    } else if level == "warning" {
        color = "warning"
    }
    
    msg := SlackMessage{
        Text: "⚠️ Alert from Neonex",
        Attachments: []Attachment{
            {
                Color: color,
                Title: "Application Alert",
                Text:  message,
                Fields: []Field{
                    {Title: "Environment", Value: "production", Short: true},
                    {Title: "Level", Value: level, Short: true},
                },
            },
        },
    }
    
    payload, _ := json.Marshal(msg)
    
    _, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
    return err
}
```

---

## Application Performance Monitoring (APM)

### New Relic Integration

```go
import (
    "github.com/newrelic/go-agent/v3/newrelic"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := core.NewApp()
    
    // Initialize New Relic
    nrApp, err := newrelic.NewApplication(
        newrelic.ConfigAppName("Neonex Core"),
        newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Middleware
    app.HTTP.Use(func(c *fiber.Ctx) error {
        txn := nrApp.StartTransaction(c.Path())
        defer txn.End()
        
        c.Locals("nrTxn", txn)
        return c.Next()
    })
    
    app.Run()
}
```

### DataDog APM

```go
import (
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
    "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
)

func main() {
    // Start DataDog tracer
    tracer.Start(
        tracer.WithEnv("production"),
        tracer.WithService("neonex-core"),
        tracer.WithServiceVersion("1.0.0"),
    )
    defer tracer.Stop()
    
    // Instrument GORM
    db, err := gormtracer.Open("postgres", dsn, &gorm.Config{})
    
    app := core.NewApp()
    app.Run()
}
```

---

## Performance Metrics Dashboard

### Key Metrics to Monitor

```go
// Recommended metrics structure
type Metrics struct {
    // Application
    RequestRate         float64
    ErrorRate           float64
    ResponseTimeP50     float64
    ResponseTimeP95     float64
    ResponseTimeP99     float64
    
    // Database
    DBConnections       int
    DBQueryTime         float64
    DBSlowQueries       int
    
    // System
    CPUUsage           float64
    MemoryUsage        float64
    DiskUsage          float64
    GoroutineCount     int
    
    // Business
    ActiveUsers        int
    RequestsPerMinute  int
    ErrorsPerMinute    int
}
```

---

## Best Practices

### ✅ DO

- Monitor golden signals (latency, traffic, errors, saturation)
- Set up alerts for critical thresholds
- Use structured logging
- Track business metrics
- Set SLOs (Service Level Objectives)
- Review metrics regularly
- Test monitoring system
- Document runbooks

### ❌ DON'T

- Alert on every minor issue
- Ignore monitoring costs
- Over-complicate dashboards
- Log sensitive data
- Forget to rotate logs
- Skip health checks
- Neglect capacity planning

---

## Monitoring Checklist

- [ ] Prometheus collecting metrics
- [ ] Grafana dashboards configured
- [ ] Log aggregation setup (ELK/Loki)
- [ ] Health/readiness endpoints
- [ ] Alert rules defined
- [ ] Notification channels configured
- [ ] Distributed tracing enabled
- [ ] APM tool integrated
- [ ] SLOs defined
- [ ] Runbooks documented
- [ ] On-call rotation setup

---

## Next Steps

- [**Production Setup**](production-setup.md) - Production deployment
- [**Environment Variables**](environment-variables.md) - Configuration
- [**Performance**](../advanced/performance.md) - Optimization
- [**Docker**](docker.md) - Container monitoring

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
