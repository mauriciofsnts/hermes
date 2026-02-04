<div align="center">
  <h1 align="center">⚡ HERMES</h1>

  <h3>Production-ready email notification service with queue-based processing, circuit breakers, and automatic retry logic</h3>

  <p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat-square&logo=Go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Redis-DC382D.svg?style=flat-square&logo=Redis&logoColor=white" alt="Redis" />
  <img src="https://img.shields.io/badge/Docker-2496ED.svg?style=flat-square&logo=Docker&logoColor=white" alt="Docker" />
  <img src="https://img.shields.io/badge/Prometheus-E6522C.svg?style=flat-square&logo=Prometheus&logoColor=white" alt="Prometheus" />
  </p>
  <img src="https://img.shields.io/github/license/mauriciofsnts/hermes?style=flat-square&color=5D6D7E" alt="GitHub license" />
  <img src="https://img.shields.io/github/last-commit/mauriciofsnts/hermes?style=flat-square&color=5D6D7E" alt="git-last-commit" />
  <img src="https://img.shields.io/github/commit-activity/m/mauriciofsnts/hermes?style=flat-square&color=5D6D7E" alt="GitHub commit activity" />
  <img src="https://img.shields.io/github/languages/top/mauriciofsnts/hermes?style=flat-square&color=5D6D7E" alt="GitHub top language" />
</div>

---

## 📖 Table of Contents

- [Why Hermes?](#-why-hermes)
- [Quick Start](#-quick-start)
- [Architecture](#-architecture)
- [Features](#-features)
- [Configuration](#-configuration)
- [API Reference](#-api-reference)
- [Development](#-development)
- [Deployment](#-deployment)
- [Observability](#-observability)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🎯 Why Hermes?

Hermes transforms email sending from a fragile, blocking operation into a resilient, observable microservice. Unlike simple SMTP wrappers, Hermes provides:

- **🔄 Guaranteed Delivery**: Dead Letter Queue with automatic retry (up to 5 attempts)
- **🛡️ Production Resilience**: Circuit breakers prevent cascading SMTP failures
- **📊 Full Observability**: Prometheus metrics for email success rates, queue depth, and latency
- **⚖️ Horizontal Scaling**: Redis-backed distributed queue and rate limiting
- **🎨 Template Management**: Dynamic HTML templates with caching
- **🔐 Multi-App Support**: Isolated API keys and rate limits per application

### When to Use Hermes

✅ **Perfect for:**
- Microservices needing reliable transactional emails
- Multi-tenant applications requiring isolated email sending
- High-volume notification systems (marketing, alerts, reports)
- Teams wanting email observability without vendor lock-in

❌ **Not ideal for:**
- Simple scripts needing one-off emails (use `net/smtp` directly)
- Real-time chat applications (consider WebSockets/SSE instead)

---

## � Quick Start

### Prerequisites
- Go 1.25+
- (Optional) Redis for distributed features
- SMTP server credentials (Gmail, SendGrid, Mailgun, etc.)

### 1. Install and Configure

```bash
# Clone the repository
git clone https://github.com/mauriciofsnts/hermes
cd hermes

# Install dependencies
go mod download

# Create config from example
make start  # Auto-creates config.yaml
```

### 2. Configure Your SMTP & App

Edit `config.yaml`:

```yaml
smtp:
  host: "smtp.gmail.com"
  port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  sender: "noreply@yourapp.com"

apps:
  my-app:
    enabled: true
    apiKey: "7a28c3e0-83e4-426f-89a4-d932cdcadac4"  # Change this!
    limitPerIPPerHour: 1000
    enabledFeatures: [email]
```

### 3. Create a Template

```bash
# Create templates/welcome.html
cat > templates/welcome.html << 'EOF'
<!DOCTYPE html>
<html>
<body>
  <h1>Welcome, {{.Name}}!</h1>
  <p>{{.Message}}</p>
</body>
</html>
EOF
```

### 4. Send Your First Email

```bash
# Start the server
make dev

# Send email via API
curl -X POST http://localhost:3000/api/v1/app/notify/notification \
  -H "x-api-key: 7a28c3e0-83e4-426f-89a4-d932cdcadac4" \
  -H "Content-Type: application/json" \
  -d '{
    "templateId": "welcome",
    "subject": "Welcome to Our Service!",
    "recipients": [{
      "type": "mail",
      "data": {
        "to": "user@example.com",
        "Name": "Alice",
        "Message": "Thanks for joining us!"
      }
    }]
  }'
```

✅ **Response:** `{"message": "Email sent successfully"}`

---

## 🏗️ Architecture

Hermes follows a clean architecture with dependency injection and interface-based providers:

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ POST /api/v1/app/notify/notification
       ▼
┌─────────────────────────────────────────┐
│         HTTP Server (Chi)               │
│  ┌──────────────────────────────┐       │
│  │ Middleware Chain:            │       │
│  │  → Auth (API Key)            │       │
│  │  → Rate Limiter              │       │
│  │  → Metrics (Prometheus)      │       │
│  └──────────────────────────────┘       │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│      Template Service                   │
│  • Parse HTML with dynamic data         │
│  • In-memory cache (sync.RWMutex)       │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│      Queue (Redis/Memory)               │
│  • Async processing                     │
│  • Worker reads from queue              │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│      SMTP Provider                      │
│  • Circuit breaker (3 failures → open) │
│  • Automatic retry logic                │
└──────┬──────────────────────────────────┘
       │
       ├─ Success ✓
       │
       └─ Failure ✗
          │
          ▼
    ┌─────────────────────────────┐
    │  Dead Letter Queue (SQLite) │
    │  • Max 5 retry attempts     │
    │  • Background worker        │
    │  • Admin API for monitoring │
    └─────────────────────────────┘
```

### Key Design Patterns

- **Provider Interface Pattern**: All external services (SMTP, Queue, Templates) implement interfaces for easy testing/mocking
- **Circuit Breaker**: Prevents cascading SMTP failures; opens after 3 failures, half-opens after 30s
- **Template Caching**: Parsed templates cached in-memory with thread-safe access
- **Queue Abstraction**: Swap between Redis (distributed) and Memory (development) seamlessly
- **WrappedHandler**: Custom router pattern that returns `Response` objects instead of writing directly to `http.ResponseWriter`

---

## �📦 Features

## 📦 Features

### 🔄 Dead Letter Queue (DLQ)

Automatic failure handling with persistent retry logic:

```
Email Send Failed → DLQ (SQLite)
                     ↓
            Background Worker (5min interval)
                     ↓
        Retry Attempt (max 5 times)
                     ↓
         Success ✓ or Permanent Failure ✗
```

**Admin API:**
- `GET /api/v1/admin/dlq/stats` - View retry statistics
- `GET /api/v1/admin/dlq/pending` - List pending retries
- `GET /api/v1/admin/dlq/failed` - View permanently failed emails

### 🛡️ Circuit Breaker

Protects against cascading SMTP failures:

- **Closed** (normal): Requests pass through
- **Open** (failing): Fast-fail for 30s after 3 failures
- **Half-Open** (testing): Allow 1 request to test recovery

```go
// Distributed Redis version shares state across instances
type CircuitBreaker interface {
    CanExecute() bool
    RecordSuccess()
    RecordFailure()
    GetState() string  // "closed", "open", "half-open"
}
```

### 📊 Prometheus Metrics

Production-grade observability out of the box:

```prometheus
# Email metrics
hermes_emails_sent_total{status="success|failed"}
hermes_email_send_duration_seconds

# Queue metrics
hermes_queue_depth
hermes_queue_processing_duration_seconds

# Circuit breaker
hermes_circuit_breaker_state{state="closed|open|half-open"}

# Rate limiting
hermes_rate_limit_events_total{action="allowed|blocked"}
```

Access at: `http://localhost:3000/metrics`

### ⚖️ Distributed Features

Run multiple Hermes instances with shared state:

| Feature | Single Instance | Multi-Instance (Redis) |
|---------|----------------|------------------------|
| Queue Processing | ✅ Memory | ✅ Redis (shared jobs) |
| Circuit Breaker | ✅ Local state | ✅ Redis (cluster-wide) |
| Rate Limiting | ✅ In-memory | ✅ Redis (global limits) |
| DLQ | ✅ SQLite | ✅ SQLite (per-instance) |

**Enable Redis:**
```yaml
redis:
  address: "localhost:6379"
  password: "your-password"
  topic: hermes
```

### 🎨 Dynamic Templates

Go template engine with caching:

```html
<!-- templates/invoice.html -->
<!DOCTYPE html>
<html>
<body>
  <h1>Invoice #{{.InvoiceID}}</h1>
  <p>Dear {{.CustomerName}},</p>
  <p>Amount due: ${{.Amount}}</p>
  {{if .IsPastDue}}
    <p style="color: red;">⚠️ Payment overdue!</p>
  {{end}}
</body>
</html>
```

**Template API:**
- `POST /api/v1/app/templates` - Upload template
- `GET /api/v1/app/templates/{id}` - Retrieve template
- `DELETE /api/v1/app/templates/{id}` - Delete template

### 🔐 Multi-App Support

Isolate email sending per application:

```yaml
apps:
  app-production:
    enabled: true
    apiKey: "prod-key-xxx"
    limitPerIPPerHour: 5000
    allowedOrigins: ["https://app.example.com"]

  app-staging:
    enabled: true
    apiKey: "staging-key-yyy"
    limitPerIPPerHour: 100
    allowedOrigins: ["https://staging.example.com"]
```

Each app gets:
- ✅ Unique API key for authentication
- ✅ Independent rate limits
- ✅ Custom CORS origins
- ✅ Feature flags (email, discord)

---


## 📡 API Reference

### Send Notification

**Endpoint:** `POST /api/v1/app/notify/notification`

**Headers:**
```
x-api-key: your-api-key
Content-Type: application/json
```

**Request Body:**
```json
{
  "templateId": "welcome",
  "subject": "Welcome to Our Service",
  "recipients": [
    {
      "type": "mail",
      "data": {
        "to": "user@example.com",
        "Name": "John Doe",
        "CustomField": "Any value you need in template"
      }
    }
  ]
}
```

**Success Response (200):**
```json
{
  "message": "Email sent successfully"
}
```

**Error Response (4xx/5xx):**
```json
{
  "error": "Failed to send email: template not found"
}
```

### Health Check

**Endpoint:** `GET /api/v1/health`

**Response:**
```json
{
  "status": "healthy",
  "queue": "redis connected"
}
```

### Template Management

**Upload Template:**
```bash
POST /api/v1/app/templates
x-api-key: your-api-key
Content-Type: application/json

{
  "name": "welcome",
  "content": "<html>...</html>"
}
```

**Get Template:**
```bash
GET /api/v1/app/templates/welcome
x-api-key: your-api-key
```

### DLQ Management

**View Statistics:**
```bash
GET /api/v1/admin/dlq/stats
```

**Response:**
```json
{
  "pending": 5,
  "processing": 2,
  "failed": 1,
  "succeeded": 234
}
```

### Swagger Documentation

Interactive API docs available at: `http://localhost:3000/swagger/index.html`

---

## 🛠️ Development

### Local Development

```bash
# Hot reload with Air
make dev

# Run tests
make test

# Integration tests (requires Docker)
make test-integration

# Code quality checks
make inspect  # Runs revive + gosec + staticcheck

# Generate Swagger docs
make swagger
```

### Project Structure

```
hermes/
├── cmd/hermes/          # Application entry point
├── internal/
│   ├── bootstrap/       # Initialization logic
│   ├── config/          # Configuration loading & validation
│   ├── metrics/         # Prometheus metrics
│   ├── providers/       # External service interfaces
│   │   ├── database/    # DLQ persistence (SQLite)
│   │   ├── discord/     # Discord webhook integration
│   │   ├── queue/       # Queue abstraction (Redis/Memory)
│   │   ├── smtp/        # Email sending with circuit breaker
│   │   └── template/    # Template parsing & caching
│   ├── server/          # HTTP server & middleware
│   │   ├── api/         # Controllers & routing
│   │   ├── middleware/  # Auth, rate limiting, logging
│   │   └── router/      # Route definitions
│   └── types/           # Shared data structures
├── templates/           # Email HTML templates
├── config.yaml          # Runtime configuration
└── Makefile             # Build & dev commands
```

### Adding a New Endpoint

1. **Create controller** in `internal/server/api/your-feature/`
2. **Implement handler** returning `api.Response`
3. **Register route** in `internal/server/router/main.go`
4. **Add Swagger comments** and run `make swagger`

Example:
```go
// internal/server/api/myfeature/controller.go
type MyController struct {
    provider providers.SomeProvider
}

func (c *MyController) Route(r api.Router) {
    r.Post("/my-endpoint", c.HandleRequest)
}

func (c *MyController) HandleRequest(r *http.Request) api.Response {
    // Your logic here
    return api.SuccessResponse("Done!")
}
```


## 📚 Advanced Examples

### Multi-Recipient Email

```bash
curl -X POST http://localhost:3000/api/v1/app/notify/notification \
  -H "x-api-key: your-key" \
  -H "Content-Type: application/json" \
  -d '{
    "templateId": "newsletter",
    "subject": "Monthly Update",
    "recipients": [
      {
        "type": "mail",
        "data": {
          "to": "alice@example.com",
          "Name": "Alice",
          "Content": "Custom content for Alice"
        }
      },
      {
        "type": "mail",
        "data": {
          "to": "bob@example.com",
          "Name": "Bob",
          "Content": "Custom content for Bob"
        }
      }
    ]
  }'
```

### Conditional Template Logic

```html
<!-- templates/order-confirmation.html -->
<!DOCTYPE html>
<html>
<body>
  <h1>Order #{{.OrderID}} Confirmed</h1>

  {{if .IsExpressShipping}}
    <p style="color: green;">⚡ Express shipping - arrives tomorrow!</p>
  {{else}}
    <p>Standard shipping - arrives in 3-5 days</p>
  {{end}}

  <h2>Items ({{len .Items}}):</h2>
  <ul>
    {{range .Items}}
      <li>{{.Name}} - ${{.Price}}</li>
    {{end}}
  </ul>

  <p><strong>Total: ${{.Total}}</strong></p>
</body>
</html>
```

---

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create feature branch: `git checkout -b feature/my-feature`
3. Make changes and add tests
4. Run quality checks: `make inspect`
5. Commit: `git commit -m 'Add feature X'`
6. Push: `git push origin feature/my-feature`
7. Open Pull Request

### Areas for Contribution

- 🔌 **New Providers**: SMS, Slack, Teams integrations
- 📊 **Enhanced Metrics**: Custom business metrics
- 🧪 **Test Coverage**: Integration tests, benchmarks
- 📚 **Documentation**: Tutorials, architecture diagrams
- 🐛 **Bug Fixes**: Check [Issues](https://github.com/mauriciofsnts/hermes/issues)

---

## 📄 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
