# Docker Setup Guide

## Overview

This guide explains how to run the Coin Backend application using Docker. Docker provides a containerized environment that ensures consistency across development, testing, and production.

## Prerequisites

- **Docker** (version 20.10+)
- **Docker Compose** (version 1.29+)

### Installation

#### macOS
```bash
# Using Homebrew
brew install docker docker-compose

# Or download Docker Desktop
# https://www.docker.com/products/docker-desktop
```

#### Linux (Ubuntu/Debian)
```bash
# Install Docker
sudo apt-get update
sudo apt-get install docker.io docker-compose

# Add user to docker group (optional)
sudo usermod -aG docker $USER
```

#### Windows
- Download Docker Desktop: https://www.docker.com/products/docker-desktop
- Follow installation wizard

## Quick Start

### Option 1: Using Make Commands (Recommended)

```bash
# Build Docker image
make docker-build

# Start all services
make docker-up

# View logs
make docker-logs

# Seed database
make docker-seed

# Stop services
make docker-down
```

### Option 2: Using Docker Compose Directly

```bash
# Build and start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Services

The `docker-compose.yml` includes three services:

### 1. MongoDB
- **Image**: mongo:7.0
- **Port**: 27017
- **Container**: coin-mongodb
- **Volumes**: mongodb_data, mongodb_config
- **Credentials**: 
  - Username: admin (default)
  - Password: password (default)

### 2. Redis
- **Image**: redis:7-alpine
- **Port**: 6379
- **Container**: coin-redis
- **Volumes**: redis_data

### 3. Application (Go API)
- **Image**: Built from Dockerfile
- **Port**: 8080
- **Container**: coin-app
- **Depends on**: MongoDB, Redis

## Environment Variables

### Docker Compose Environment

Create a `.env` file in the project root:

```env
# MongoDB
MONGO_USER=admin
MONGO_PASSWORD=password
MONGODB_NAME=coin_db

# Redis (optional)
REDIS_PASSWORD=

# Application
JWT_SECRET=your-secret-key-change-in-production
ENV=development
```

### Available Variables

| Variable | Default | Description |
|----------|---------|-------------|
| MONGO_USER | admin | MongoDB username |
| MONGO_PASSWORD | password | MongoDB password |
| MONGODB_NAME | coin_db | Database name |
| JWT_SECRET | your-secret-key-change-in-production | JWT signing key |
| ENV | development | Environment (development/production) |
| PORT | 8080 | Application port |

## Docker Commands

### Build Commands

```bash
# Build Docker image
make docker-build

# Build without cache
docker-compose build --no-cache
```

### Lifecycle Commands

```bash
# Start services in background
make docker-up

# Start services with logs
docker-compose up

# Stop services
make docker-down

# Stop and remove volumes
make docker-clean

# Restart services
docker-compose restart
```

### Logging and Debugging

```bash
# View all logs
make docker-logs

# View logs for specific service
docker-compose logs app
docker-compose logs mongodb
docker-compose logs redis

# Follow logs in real-time
docker-compose logs -f

# View last 100 lines
docker-compose logs --tail=100
```

### Database Operations

```bash
# Seed database
make docker-seed

# Connect to MongoDB
docker-compose exec mongodb mongosh -u admin -p password

# Connect to Redis
docker-compose exec redis redis-cli
```

### Container Management

```bash
# List running containers
docker-compose ps

# Execute command in container
docker-compose exec app /bin/sh

# View container stats
docker stats

# Remove containers
docker-compose down

# Remove containers and volumes
docker-compose down -v
```

## Access Points

Once services are running:

| Service | URL/Address |
|---------|------------|
| API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/swagger/index.html |
| Health Check | http://localhost:8080/api/health |
| MongoDB | localhost:27017 |
| Redis | localhost:6379 |

## File Structure

```
coin-be/
├── Dockerfile              # Multi-stage build for Go app
├── docker-compose.yml      # Docker Compose configuration
├── .dockerignore           # Files to exclude from Docker build
├── .env.example            # Example environment variables
├── cmd/
│   ├── api/main.go        # API application
│   └── seeder/main.go     # Database seeder
└── internal/
    └── seeder/            # Seeder package
```

## Dockerfile Explanation

The Dockerfile uses a **multi-stage build** for optimization:

### Build Stage
```dockerfile
FROM golang:1.21-alpine AS builder
# - Compiles Go application
# - Builds both API and seeder binaries
# - Reduces final image size
```

### Final Stage
```dockerfile
FROM alpine:latest
# - Minimal base image
# - Only includes compiled binaries
# - Includes health check
# - Exposes port 8080
```

### Benefits
- **Small image size** (~50MB vs 300MB+)
- **Security** - No build tools in final image
- **Fast startup** - Pre-compiled binaries

## Health Checks

The Docker setup includes health checks:

### Application Health Check
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health
```

### MongoDB Health Check
```yaml
healthcheck:
  test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
  interval: 10s
  timeout: 5s
  retries: 5
```

### Redis Health Check
```yaml
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 10s
  timeout: 5s
  retries: 5
```

## Volumes

### MongoDB Volumes
- `mongodb_data` - Database files
- `mongodb_config` - Configuration files

### Redis Volumes
- `redis_data` - Redis persistence

### Benefits
- **Data persistence** - Survives container restarts
- **Performance** - Faster than writing to host filesystem
- **Isolation** - Separate from host system

## Networking

All services connect via `coin-network` bridge network:

```yaml
networks:
  coin-network:
    driver: bridge
```

### Benefits
- **Service discovery** - Services can reference each other by name
- **Isolation** - Separate from other Docker networks
- **Security** - Only exposed ports are accessible

## Common Tasks

### Development Workflow

```bash
# 1. Start services
make docker-up

# 2. Seed database
make docker-seed

# 3. View logs
make docker-logs

# 4. Make code changes (files are not synced in Docker)
# Note: Rebuild image for changes to take effect

# 5. Rebuild and restart
make docker-build
make docker-up

# 6. Stop services
make docker-down
```

### Testing

```bash
# Run tests in container
docker-compose exec app go test ./...

# Run specific test
docker-compose exec app go test ./internal/modules/user
```

### Database Inspection

```bash
# Connect to MongoDB
docker-compose exec mongodb mongosh -u admin -p password

# List databases
show dbs

# Use database
use coin_db

# List collections
show collections

# Query data
db.users.find()
db.categories.find()
```

### Redis Inspection

```bash
# Connect to Redis
docker-compose exec redis redis-cli

# Get all keys
KEYS *

# Get specific key
GET key_name

# Monitor commands
MONITOR
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker-compose logs app

# Verify services are healthy
docker-compose ps

# Check port conflicts
lsof -i :8080
lsof -i :27017
lsof -i :6379
```

### Connection Refused

```bash
# Ensure all services are running
docker-compose ps

# Check service health
docker-compose exec app wget -O- http://localhost:8080/api/health

# Verify network connectivity
docker-compose exec app ping mongodb
docker-compose exec app ping redis
```

### Database Issues

```bash
# Clear database and restart
docker-compose down -v
docker-compose up -d

# Reseed database
make docker-seed
```

### Permission Denied

```bash
# On Linux, add user to docker group
sudo usermod -aG docker $USER

# Or use sudo
sudo docker-compose up
```

### Out of Disk Space

```bash
# Clean up Docker resources
docker system prune -a

# Remove unused volumes
docker volume prune
```

## Production Deployment

### Security Considerations

1. **Change Default Credentials**
   ```env
   MONGO_USER=your_secure_user
   MONGO_PASSWORD=your_secure_password
   JWT_SECRET=your_secure_jwt_secret
   ```

2. **Use Environment Variables**
   - Never commit `.env` file
   - Use `.env.example` as template
   - Manage secrets securely

3. **Network Security**
   - Don't expose MongoDB/Redis to public
   - Use firewall rules
   - Enable authentication

4. **Image Security**
   - Use specific version tags
   - Scan images for vulnerabilities
   - Keep base images updated

### Production Commands

```bash
# Build production image
docker build -t coin-api:1.0.0 .

# Push to registry
docker tag coin-api:1.0.0 your-registry/coin-api:1.0.0
docker push your-registry/coin-api:1.0.0

# Deploy with docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

## Performance Optimization

### Memory Limits

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

### CPU Limits

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
        reservations:
          cpus: '0.5'
```

## Monitoring

### View Resource Usage

```bash
# Real-time stats
docker stats

# Specific container
docker stats coin-app
```

### View Logs with Timestamps

```bash
docker-compose logs --timestamps
```

## Advanced Usage

### Custom Docker Network

```bash
# Create custom network
docker network create coin-network

# Run services on custom network
docker-compose --network coin-network up
```

### Multi-Stage Builds

The Dockerfile uses multi-stage builds:
1. **Builder stage** - Compiles application
2. **Final stage** - Runs application

This reduces image size from ~300MB to ~50MB.

### Docker Compose Override

Create `docker-compose.override.yml` for local development:

```yaml
version: '3.8'
services:
  app:
    environment:
      ENV: development
      DEBUG: "true"
```

## Support

For issues or questions:
1. Check troubleshooting section
2. Review Docker logs
3. Verify environment variables
4. Check Docker and Docker Compose versions
5. Consult Docker documentation: https://docs.docker.com
