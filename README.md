# Gopher Foody Gateway Service

A high-performance API Gateway built with Go and Gin. It acts as the centralized entry point for the Gopher Foody microservices, handling JWT authentication and request proxying.

## 📂 Project Structure

```text
.
├── cmd/server/main.go       # Application entry point
├── internal/
│   ├── config/              # Configuration loading
│   └── presentation/        # HTTP Handlers & Middleware
└── pkg/
    ├── jwt/                 # Token validation
    └── logger/              # Structured logging
```

## 🚀 Getting Started

1. **Configure Environment**:
   Create a `.env` file from the template:
   ```env
   APP_HTTP_PORT=8000
   JWT_ACCESS_SECRET=your-secret-key
   IDENTITY_SERVICE_URL=http://localhost:8080
   RESTAURANT_SERVICE_URL=http://localhost:8081
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the Service**:
   ```bash
   go run cmd/server/main.go
   ```

The gateway will be accessible at `http://localhost:8000`.
