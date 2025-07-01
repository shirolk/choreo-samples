# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a simple Go-based REST API greeting service designed for deployment on WSO2 Choreo platform. The service provides a single HTTP endpoint that greets users by name.

## Common Development Commands

### Local Development
```bash
# Run the service locally
go run main.go

# Build the executable
go build -o go-greeter .

# Test the service locally
curl "http://localhost:9090/greeter/greet?name=World"
```

### Docker Operations
```bash
# Build Docker image
docker build -t greeting-service-go .

# Run containerized service
docker run -p 9090:9090 greeting-service-go

# Test containerized service
curl "http://localhost:9090/greeter/greet?name=Docker"
```

## Architecture

### Service Structure
- **Single-file application**: All logic contained in `main.go`
- **HTTP server**: Uses Go's standard `net/http` package
- **Graceful shutdown**: Implements proper signal handling for clean termination
- **Security**: Runs as non-root user (UID 10014) in container

### Key Components
- **HTTP Handler**: `greet()` function handles GET requests to `/greeter/greet`
- **Query Parameter**: Accepts optional `name` parameter, defaults to "Stranger"
- **Server Configuration**: Listens on port 9090 with 10-second shutdown timeout

### Choreo Integration
- **Component Configuration**: `.choreo/component.yaml` defines endpoint exposure
- **API Contract**: `openapi.yaml` provides REST API specification
- **Container Security**: Dockerfile creates minimal Alpine-based image with non-root user

## File Structure

| File | Purpose |
|------|---------|
| `main.go` | Complete service implementation |
| `go.mod` | Go module definition (Go 1.21) |
| `Dockerfile` | Multi-stage container build |
| `openapi.yaml` | REST API specification |
| `.choreo/component.yaml` | Choreo platform configuration |

## Development Notes

### Port Configuration
- Service runs on port 9090 (hardcoded)
- Choreo endpoint configuration matches this port

### API Endpoint
- **Path**: `/greeter/greet`
- **Method**: GET
- **Parameter**: `name` (optional query parameter)
- **Response**: Plain text greeting message

### Container Requirements
- Non-root user execution (Choreo requirement)
- Alpine Linux base for minimal footprint
- Static binary compilation for portability