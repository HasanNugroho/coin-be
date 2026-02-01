# Docker Setup Guide

## Overview

This guide provides instructions for containerizing and running the Coin Backend application using Docker and Docker Compose.

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Git

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/HasanNugroho/coin-be.git
cd coin-be
```

### 2. Configure Environment

Copy the appropriate environment file:

```bash
# For development
cp .env.development .env

# For staging
cp .env.staging .env

# For production
cp .env.production .env
```

### 3. Start Services

```bash
docker-compose up -d
```

This will start:
- MongoDB (port 27017)
- Redis (port 6379)
- Coin Backend API (port 8080)

### 4. Verify Services

```bash
# Check container status
docker-compose ps

# View logs
docker-compose logs -f app

# Test API
curl http://localhost:8080/api/health
```

### 5. Stop Services

```bash
docker-compose down
```

## Environment Configuration

### Environment Files

The project includes three pre-configured environment files:

#### `.env.development`
- **Purpose**: Local development
- **Log Level**: debug
- **Log Format**: text
- **Database**: Local MongoDB
- **Redis**: Local instance
- **JWT Secret**: Development key (change before production)

#### `.env.staging`
- **Purpose**: Staging/testing environment
- **Log Level**: info
- **Log Format**: JSON
- **Database**: Staging MongoDB
- **Redis**: Staging instance
- **JWT Secret**: Staging key (change before production)

#### `.env.production`
- **Purpose**: Production deployment
- **Log Level**: warn
- **Log Format**: JSON
- **Database**: Production MongoDB
- **Redis**: Production instance
- **JWT Secret**: Must be changed to strong random key

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | 8080 | Application port |
| `ENV` | development | Environment (development, staging, production) |
| `MONGO_URI` | mongodb://admin:password@localhost:27017 | MongoDB connection URI |
| `MONGO_DB` | coin_db | MongoDB database name |
| `MONGO_USER` | admin | MongoDB username |
| `MONGO_PASSWORD` | password | MongoDB password |
| `MONGO_PORT` | 27017 | MongoDB port |
| `REDIS_ADDR` | localhost:6379 | Redis address |
| `REDIS_PASSWORD` | (empty) | Redis password |
| `REDIS_DB` | 0 | Redis database number |
| `REDIS_PORT` | 6379 | Redis port |
| `JWT_SECRET` | dev-secret-key | JWT signing secret |
| `JWT_DURATION` | 24h | Access token duration |
| `JWT_REFRESH_DURATION` | 168h | Refresh token duration |
| `LOG_LEVEL` | info | Logging level (debug, info, warn, error) |
| `LOG_FORMAT` | json | Log format (json, text) |

## Docker Compose Services

### MongoDB

```yaml
mongodb:
  image: mongo:7.0
  ports:
    - "27017:27017"
  environment:
    MONGO_INITDB_ROOT_USERNAME: admin
    MONGO_INITDB_ROOT_PASSWORD: password
    MONGO_INITDB_DATABASE: coin_db
  volumes:
    - mongodb_data:/data/db
    - mongodb_config:/data/configdb
```

**Features**:
- Persistent data storage
- Health checks enabled
- Automatic initialization
- Network isolation

### Redis

```yaml
redis:
  image: redis:7-alpine
  ports:
    - "6379:6379"
  command: redis-server --appendonly yes
  volumes:
    - redis_data:/data
```

**Features**:
- AOF persistence enabled
- Health checks enabled
- Alpine image for smaller size
- Network isolation

### Application (Coin Backend)

```yaml
app:
  build:
    context: .
    dockerfile: Dockerfile
  ports:
    - "8080:8080"
  env_file:
    - .env
  depends_on:
    mongodb:
      condition: service_healthy
    redis:
      condition: service_healthy
```

**Features**:
- Multi-stage build for optimization
- Environment file support
- Service dependency management
- Health checks
- Automatic restart

## Building and Running

### Build Docker Image

```bash
# Build the image
docker build -t coin-be:latest .

# Build with specific tag
docker build -t coin-be:v1.0.0 .
```

### Run Container Manually

```bash
# Run with environment file
docker run --env-file .env -p 8080:8080 coin-be:latest

# Run with specific environment variables
docker run \
  -e MONGO_URI=mongodb://admin:password@mongodb:27017 \
  -e REDIS_ADDR=redis:6379 \
  -e JWT_SECRET=your-secret-key \
  -p 8080:8080 \
  coin-be:latest
```

### Docker Compose Commands

```bash
# Start services in background
docker-compose up -d

# Start services with logs
docker-compose up

# Stop services
docker-compose stop

# Stop and remove containers
docker-compose down

# Remove volumes (data loss)
docker-compose down -v

# View logs
docker-compose logs

# Follow logs for specific service
docker-compose logs -f app

# Rebuild images
docker-compose build

# Rebuild and start
docker-compose up -d --build
```

## Health Checks

### MongoDB Health Check

```bash
docker-compose exec mongodb mongosh localhost:27017/test --eval "db.runCommand('ping')"
```

### Redis Health Check

```bash
docker-compose exec redis redis-cli ping
```

### Application Health Check

```bash
curl http://localhost:8080/api/health
```

## Volumes

### Data Persistence

The following volumes are created for data persistence:

| Volume | Purpose | Path |
|--------|---------|------|
| `mongodb_data` | MongoDB data files | `/data/db` |
| `mongodb_config` | MongoDB configuration | `/data/configdb` |
| `redis_data` | Redis data files | `/data` |

### Backup Volumes

```bash
# Backup MongoDB data
docker run --rm -v coin-be_mongodb_data:/data -v $(pwd):/backup \
  alpine tar czf /backup/mongodb_backup.tar.gz -C /data .

# Backup Redis data
docker run --rm -v coin-be_redis_data:/data -v $(pwd):/backup \
  alpine tar czf /backup/redis_backup.tar.gz -C /data .
```

### Restore Volumes

```bash
# Restore MongoDB data
docker run --rm -v coin-be_mongodb_data:/data -v $(pwd):/backup \
  alpine tar xzf /backup/mongodb_backup.tar.gz -C /data

# Restore Redis data
docker run --rm -v coin-be_redis_data:/data -v $(pwd):/backup \
  alpine tar xzf /backup/redis_backup.tar.gz -C /data
```

## Networking

### Network Configuration

All services are connected via a custom bridge network `coin-network`:

```bash
# View network
docker network inspect coin-be_coin-network

# Services can communicate using service names:
# - mongodb:27017
# - redis:6379
# - app:8080
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker-compose logs app

# Inspect container
docker-compose ps

# Check resource usage
docker stats
```

### Connection Issues

```bash
# Test MongoDB connection
docker-compose exec app nc -zv mongodb 27017

# Test Redis connection
docker-compose exec app nc -zv redis 6379

# Check network
docker network inspect coin-be_coin-network
```

### Database Issues

```bash
# Access MongoDB shell
docker-compose exec mongodb mongosh

# Access Redis CLI
docker-compose exec redis redis-cli

# View MongoDB logs
docker-compose logs mongodb

# View Redis logs
docker-compose logs redis
```

### Port Conflicts

If ports are already in use, modify `.env`:

```bash
# Change ports in .env
APP_PORT=8081
MONGO_PORT=27018
REDIS_PORT=6380

# Or use docker-compose override
docker-compose -f docker-compose.yml -f docker-compose.override.yml up
```

## Production Deployment

### Pre-deployment Checklist

- [ ] Update `.env.production` with secure credentials
- [ ] Change `JWT_SECRET` to strong random value
- [ ] Set `MONGO_PASSWORD` to secure password
- [ ] Set `REDIS_PASSWORD` if needed
- [ ] Configure `LOG_LEVEL` to `warn`
- [ ] Set `LOG_FORMAT` to `json`
- [ ] Review resource limits
- [ ] Set up backup strategy
- [ ] Configure monitoring

### Production Deployment Steps

```bash
# 1. Copy production environment
cp .env.production .env

# 2. Build production image
docker build -t coin-be:v1.0.0 .

# 3. Push to registry (optional)
docker tag coin-be:v1.0.0 your-registry/coin-be:v1.0.0
docker push your-registry/coin-be:v1.0.0

# 4. Start services
docker-compose up -d

# 5. Verify deployment
docker-compose ps
curl http://localhost:8080/api/health

# 6. Check logs
docker-compose logs app
```

### Resource Limits

Add to `docker-compose.yml` for production:

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### Monitoring

```bash
# Monitor container stats
docker stats

# Monitor logs
docker-compose logs -f --tail=100

# Monitor specific service
docker-compose logs -f app
```

## Security Best Practices

1. **Change Default Credentials**
   - Update `MONGO_PASSWORD` and `REDIS_PASSWORD`
   - Use strong, random values

2. **JWT Secret**
   - Generate strong random secret: `openssl rand -base64 32`
   - Never commit secrets to version control

3. **Network Security**
   - Use custom bridge network (default)
   - Don't expose unnecessary ports
   - Use firewall rules

4. **Image Security**
   - Use specific base image versions
   - Scan images for vulnerabilities
   - Keep dependencies updated

5. **Secrets Management**
   - Use Docker secrets for production
   - Consider external secret management (Vault, etc.)
   - Rotate credentials regularly

## Performance Optimization

### Build Optimization

```bash
# Use BuildKit for faster builds
DOCKER_BUILDKIT=1 docker build -t coin-be:latest .
```

### Runtime Optimization

```bash
# Limit memory usage
docker-compose up -d --memory 512m

# Set CPU shares
docker update --cpus 1 coin-app
```

## Cleanup

```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Remove everything (careful!)
docker system prune -a
```

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [MongoDB Docker Image](https://hub.docker.com/_/mongo)
- [Redis Docker Image](https://hub.docker.com/_/redis)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)

## Support

For issues or questions:
1. Check logs: `docker-compose logs`
2. Review environment configuration
3. Verify service connectivity
4. Check Docker and Docker Compose versions
5. Consult project documentation
