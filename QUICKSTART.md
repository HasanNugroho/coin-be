# Quick Start Guide

## Prerequisites Setup

### 1. MongoDB

**Using Docker:**
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

**Or install locally:**
- Download from https://www.mongodb.com/try/download/community
- Follow installation guide for your OS

### 2. Redis

**Using Docker:**
```bash
docker run -d -p 6379:6379 --name redis redis:latest
```

**Or install locally:**
- Download from https://redis.io/download
- Follow installation guide for your OS

### 3. Environment Variables

Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

Update values if needed (defaults should work with Docker setup):
```env
APP_PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB=coin_db
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
JWT_SECRET=your-secret-key-here-change-in-production
```

## Development Setup

### Install Tools

```bash
make install-tools
```

This installs:
- **Air**: Hot reload development server
- **Swag**: Swagger documentation generator

### Start Development Server

```bash
make dev
```

The server will:
- Start on `http://localhost:8080`
- Auto-reload when you save files
- Show compilation errors in terminal

### Access Swagger UI

Open your browser and go to:
```
http://localhost:8080/swagger/index.html
```

You'll see interactive API documentation with:
- All available endpoints
- Request/response schemas
- Try-it-out functionality
- Authentication setup

## Common Workflows

### Testing an Endpoint

1. Open Swagger UI: `http://localhost:8080/swagger/index.html`
2. Click on an endpoint
3. Click "Try it out"
4. Fill in parameters
5. Click "Execute"
6. See response

### Making Code Changes

1. Edit your code
2. Save the file
3. Air automatically recompiles and restarts
4. Refresh your browser or re-run requests

### Updating Swagger Docs

After adding new endpoints or changing request/response structures:

```bash
make swagger-gen
```

Then refresh the Swagger UI in your browser.

## Building for Production

```bash
make build
```

This creates `bin/main` executable.

Run it:
```bash
./bin/main
```

## Useful Commands

```bash
# Development with hot reload
make dev

# Build for production
make build

# Run production build
make run

# Generate Swagger docs
make swagger-gen

# Format code
make fmt

# Run tests
make test

# Clean build artifacts
make clean

# Show all available commands
make help
```

## Troubleshooting

### Air not reloading changes

Check `build-errors.log`:
```bash
cat build-errors.log
```

Common issues:
- Syntax errors in code
- Missing imports
- Type mismatches

### Swagger docs not updating

Regenerate:
```bash
make swagger-gen
```

Clear browser cache (Ctrl+Shift+Delete or Cmd+Shift+Delete)

### MongoDB connection error

```
connection refused
```

Ensure MongoDB is running:
```bash
# Docker
docker ps | grep mongodb

# Or check if service is running locally
```

### Redis connection error

```
connection refused
```

Ensure Redis is running:
```bash
# Docker
docker ps | grep redis

# Or check if service is running locally
```

## API Examples

### Register User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe",
    "phone": "+62812345678"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response includes `access_token` and `refresh_token`.

### Get User Profile (Protected)

```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer <access_token>"
```

## Next Steps

1. **Explore Swagger UI** - See all available endpoints
2. **Try endpoints** - Use Swagger's "Try it out" feature
3. **Read code** - Check `internal/modules/` for implementation
4. **Add features** - Modify code and see hot reload in action
5. **Check logs** - Watch terminal for request logs

## Documentation

- **API Docs**: `http://localhost:8080/swagger/index.html`
- **Project Structure**: See `README.md`
- **Code Comments**: Check source files for implementation details

## Support

For issues:
1. Check `build-errors.log` for compilation errors
2. Check terminal output for runtime errors
3. Verify MongoDB and Redis are running
4. Check `.env` file for correct configuration
