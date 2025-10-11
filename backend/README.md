# Flow Backend

A Go-based backend API for the Flow project.

## Prerequisites

- Go 1.24+ (or compatible version)
- PostgreSQL database

## Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Configure environment:**
   ```bash
   cp env.example .env
   # Edit .env file with your database credentials
   ```

3. **Start PostgreSQL database:**
   Make sure PostgreSQL is running on the configured host and port (default: localhost:5432)

## Running the Application

### Option 1: Using the run script
```bash
./run.sh
```

### Option 2: Direct Go commands
```bash
# Build the application
go build -o bin/api ./cmd/api

# Run the application
./bin/api
```

### Option 3: Run directly with Go
```bash
go run ./cmd/api
```

## Configuration

The application uses environment variables for configuration. See `env.example` for available options:

- `SERVER_HOST`: Server host (default: localhost)
- `SERVER_PORT`: Server port (default: 8080)
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password (default: password)
- `DB_NAME`: Database name (default: flow)
- `DB_SSL_MODE`: Database SSL mode (default: disable)
- `JWT_SECRET`: JWT secret key (default: your-secret-key-change-this-in-production)

## Project Structure

```
backend/
├── cmd/api/           # Main application entry point
├── internal/          # Internal application code
│   ├── app/          # Application services
│   ├── domain/       # Domain models and interfaces
│   └── infrastructure/ # External concerns (DB, HTTP, etc.)
├── pkg/              # Shared packages
│   ├── config/       # Configuration management
│   ├── middleware/   # HTTP middleware
│   └── security/     # Security utilities
├── go.mod           # Go module file
├── env.example      # Environment variables example
└── run.sh           # Run script
```