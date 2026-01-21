# Coin Backend API

A comprehensive user management system with authentication, authorization, financial profiles, and role-based access control built with Go and Gin.

## Features

- **User Management**: Registration, login, profile management
- **Authentication**: JWT-based authentication with refresh tokens
- **Authorization**: Role-based access control (RBAC)
- **Financial Profiles**: User salary and payment information management
- **Hot Reload**: Development with Air for instant feedback
- **API Documentation**: Swagger/OpenAPI documentation
- **Middleware**: CORS, logging, panic recovery

## Prerequisites

- Go 1.25.4 or higher
- MongoDB
- Redis
- Make (optional, for using Makefile commands)

## Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd coin-be
```

### 2. Install dependencies

```bash
go mod download
go mod tidy
```

### 3. Install development tools

```bash
make install-tools
```

Or manually:

```bash
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### 4. Environment setup

Create a `.env` file in the project root:

```env
APP_PORT=8080

MONGO_URI=mongodb://localhost:27017
MONGO_DB=coin_db

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=your-secret-key-here-change-in-production
```

## Running the Application

### Development with Hot Reload

```bash
make dev
```

Or directly with Air:

```bash
air -c .air.toml
```

The application will restart automatically when you save changes.

### Production Build

```bash
make build
make run
```

Or:

```bash
go build -o bin/main ./cmd/api
./bin/main
```

## API Documentation

### Swagger UI

Access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

### Generate/Update Swagger Docs

```bash
make swagger-gen
```

Or manually:

```bash
swag init -g cmd/api/main.go --output docs
```

## Project Structure

```
coin-be/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── core/
│   │   ├── config/
│   │   │   └── config.go          # Configuration management
│   │   ├── database/
│   │   │   ├── mongodb.go         # MongoDB connection
│   │   │   └── redis.go           # Redis connection
│   │   ├── middleware/
│   │   │   ├── auth.go            # JWT authentication
│   │   │   ├── cors.go            # CORS configuration
│   │   │   ├── logger.go          # Request logging
│   │   │   └── recovery.go        # Panic recovery
│   │   ├── container/
│   │   │   └── container.go       # Dependency injection
│   │   └── utils/
│   │       ├── jwt.go             # JWT utilities
│   │       └── password.go        # Password hashing
│   │
│   └── modules/
│       ├── auth/
│       │   ├── dto/
│       │   │   ├── request.go
│       │   │   └── response.go
│       │   ├── service.go
│       │   ├── controller.go
│       │   ├── routes.go
│       │   └── module.go
│       │
│       ├── user/
│       │   ├── dto/
│       │   │   ├── request.go
│       │   │   └── response.go
│       │   ├── model.go
│       │   ├── repository.go
│       │   ├── service.go
│       │   ├── controller.go
│       │   ├── routes.go
│       │   └── module.go
│       │
│       └── health/
│           ├── controller.go
│           ├── routes.go
│           └── module.go
│
├── pkg/
│   ├── errors/
│   │   └── errors.go              # Custom error types
│   └── constants/
│       └── constants.go           # Application constants
│
├── docs/
│   └── swagger.go                 # Swagger documentation
│
├── .air.toml                      # Air configuration
├── .env.example                   # Environment variables template
├── .gitignore                     # Git ignore rules
├── Makefile                       # Build commands
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
└── README.md                      # This file
```

## API Endpoints

### Authentication (Public)

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - User logout

### User Management (Protected)

- `GET /api/users/profile` - Get current user profile
- `PUT /api/users/profile` - Update current user profile
- `GET /api/users/financial-profile` - Get financial profile
- `POST /api/users/financial-profile` - Create financial profile
- `PUT /api/users/financial-profile` - Update financial profile
- `DELETE /api/users/financial-profile` - Delete financial profile

### Admin Operations (Protected + Admin Role)

- `GET /api/users` - List all users
- `GET /api/users/:id` - Get specific user
- `DELETE /api/users/:id` - Delete user
- `POST /api/users/roles` - Create role
- `GET /api/users/roles` - List roles
- `GET /api/users/roles/:id` - Get specific role
- `POST /api/users/:id/roles` - Assign role to user
- `GET /api/users/:id/roles` - Get user's roles
- `DELETE /api/users/:id/roles/:role_id` - Remove role from user

### Health Check (Public)

- `GET /api/health` - Health check

## Authentication

### Bearer Token

Include the access token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### Token Expiration

- **Access Token**: 15 minutes
- **Refresh Token**: 7 days

## Makefile Commands

```bash
make help              # Show all available commands
make install-tools    # Install Air and Swagger
make build            # Build the application
make run              # Run the application
make dev              # Run with hot reload
make swagger-gen      # Generate Swagger docs
make swagger          # Generate and open Swagger docs
make clean            # Clean build artifacts
make test             # Run tests
make fmt              # Format code
make lint             # Run linter
make deps             # Download dependencies
```

## Development

### Code Formatting

```bash
make fmt
```

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

## Troubleshooting

### Air not reloading

1. Check `.air.toml` configuration
2. Ensure file changes are saved
3. Check `build-errors.log` for compilation errors

### Swagger docs not showing

1. Run `make swagger-gen` to regenerate docs
2. Clear browser cache
3. Check that `docs/swagger.go` exists

### MongoDB connection issues

1. Ensure MongoDB is running
2. Check `MONGO_URI` in `.env`
3. Verify MongoDB is accessible on the specified host/port

### Redis connection issues

1. Ensure Redis is running
2. Check `REDIS_ADDR` in `.env`
3. Verify Redis is accessible on the specified host/port

## License

MIT

## Support

For issues and questions, please create an issue in the repository.
