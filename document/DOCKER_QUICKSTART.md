# Docker Quick Start Guide

## 30-Second Setup

```bash
# 1. Copy environment file
cp .env.development .env

# 2. Start all services
docker-compose up -d

# 3. Verify it's running
curl http://localhost:8080/api/health
```

## Common Commands

### Start Services
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

### View Logs
```bash
docker-compose logs -f app
```

### Check Status
```bash
docker-compose ps
```

## Environment Files

| File | Purpose | When to Use |
|------|---------|------------|
| `.env.development` | Local development | `cp .env.development .env` |
| `.env.staging` | Staging/testing | `cp .env.staging .env` |
| `.env.production` | Production | `cp .env.production .env` |

## Access Services

| Service | URL | Credentials |
|---------|-----|-------------|
| API | http://localhost:8080 | - |
| MongoDB | localhost:27017 | admin / password |
| Redis | localhost:6379 | - |
| Swagger Docs | http://localhost:8080/swagger/index.html | - |

## Database Access

### MongoDB Shell
```bash
docker-compose exec mongodb mongosh
```

### Redis CLI
```bash
docker-compose exec redis redis-cli
```

## Troubleshooting

### Services won't start
```bash
# Check logs
docker-compose logs

# Rebuild images
docker-compose build --no-cache
docker-compose up -d
```

### Port already in use
Edit `.env` and change ports:
```
APP_PORT=8081
MONGO_PORT=27018
REDIS_PORT=6380
```

### Clean everything
```bash
docker-compose down -v
docker system prune -a
```

## Using Makefile

```bash
# View all docker commands
make docker-help

# Start with development environment
make docker-dev

# View logs
make docker-logs

# Check health
make docker-health

# Open app shell
make docker-shell
```

## Next Steps

1. Read `DOCKER_SETUP.md` for detailed documentation
2. Check `.env.example` for all available variables
3. Review `docker-compose.yml` for service configuration
4. See `Dockerfile` for build configuration

## Support

- Logs: `docker-compose logs`
- Status: `docker-compose ps`
- Health: `docker-compose exec app wget -O- http://localhost:8080/api/health`
