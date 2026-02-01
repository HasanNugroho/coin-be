# Module Integration Guide

## Overview

This guide explains how to integrate the new `CategoryTemplate` and `UserCategory` modules into the main application.

## Step 1: Register Modules in DI Container

In `cmd/api/main.go`, add the following imports at the top of the file:

```go
import (
	// ... existing imports ...
	"github.com/HasanNugroho/coin-be/internal/modules/category_template"
	"github.com/HasanNugroho/coin-be/internal/modules/user_category"
)
```

## Step 2: Register Modules with DI Builder

In the `main()` function, after building the container with `container.BuildContainer()`, register the new modules:

```go
func main() {
	// ... existing setup code ...

	// Build DI container
	ctn, err := container.BuildContainer()
	if err != nil {
		log.Fatalf("Failed to build container: %v", err)
	}
	defer ctn.Close()

	// Register modules
	builder, _ := di.NewBuilder()
	
	// Register existing modules
	auth.Register(builder)
	user.Register(builder)
	category.Register(builder)
	pocket.Register(builder)
	pocket_template.Register(builder)
	transaction.Register(builder)
	platform.Register(builder)
	
	// Register new modules
	category_template.Register(builder)
	user_category.Register(builder)

	// ... rest of the setup code ...
}
```

## Step 3: Setup Routes

In the router setup section of `main()`, add the new routes:

```go
func main() {
	// ... existing setup code ...

	// Setup Gin engine
	engine := gin.Default()

	// Apply middleware
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())

	// Setup routes
	v1 := engine.Group("/v1")

	// Auth routes (no auth middleware)
	authController := ctn.Get("authController").(*auth.Controller)
	jwtManager := ctn.Get("jwtManager").(*utils.JWTManager)
	db := ctn.Get("mongo").(*mongo.Client).Database(cfg.MongoDB)
	auth.RegisterRoutes(v1.Group("/auth"), authController, jwtManager, db)

	// Protected routes (with auth middleware)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager, db))

	// User routes
	userController := ctn.Get("userController").(*user.Controller)
	user.RegisterRoutes(protected.Group("/users"), userController)

	// Category routes (existing)
	categoryController := ctn.Get("categoryController").(*category.Controller)
	category.RegisterRoutes(protected.Group("/categories"), categoryController)

	// Category Template routes (new)
	categoryTemplateController := ctn.Get("categoryTemplateController").(*category_template.Controller)
	category_template.RegisterRoutes(protected.Group("/category-templates"), categoryTemplateController)

	// User Category routes (new)
	userCategoryController := ctn.Get("userCategoryController").(*user_category.Controller)
	user_category.RegisterRoutes(protected.Group("/user-categories"), userCategoryController)

	// Pocket routes
	pocketController := ctn.Get("pocketController").(*pocket.Controller)
	pocket.RegisterRoutes(protected.Group("/pockets"), pocketController)

	// Pocket Template routes
	pocketTemplateController := ctn.Get("pocketTemplateController").(*pocket_template.Controller)
	pocket_template.RegisterRoutes(protected.Group("/pocket-templates"), pocketTemplateController)

	// Transaction routes
	transactionController := ctn.Get("transactionController").(*transaction.Controller)
	transaction.RegisterRoutes(protected.Group("/transactions"), transactionController)

	// Platform routes
	platformController := ctn.Get("platformController").(*platform.Controller)
	platform.RegisterRoutes(protected.Group("/platforms"), platformController)

	// ... rest of the setup code ...
}
```

## Step 4: Verify Integration

After making the changes, verify that the application compiles:

```bash
go build ./cmd/api
```

And run the application:

```bash
go run ./cmd/api/main.go
```

## API Endpoints

After integration, the following endpoints will be available:

### Category Template Endpoints
- `GET /v1/category-templates` - Get all category templates
- `GET /v1/category-templates/{id}` - Get a specific category template
- `POST /v1/category-templates` - Create a new category template (admin only)
- `PUT /v1/category-templates/{id}` - Update a category template (admin only)
- `DELETE /v1/category-templates/{id}` - Delete a category template (admin only)

### User Category Endpoints
- `GET /v1/user-categories` - Get all user categories
- `GET /v1/user-categories/{id}` - Get a specific user category
- `POST /v1/user-categories` - Create a new user category
- `PUT /v1/user-categories/{id}` - Update a user category
- `DELETE /v1/user-categories/{id}` - Delete a user category

## Documentation

For detailed API documentation, see:
- `CATEGORY_TEMPLATE_API.md` - Category Template API documentation
- `USER_CATEGORY_API.md` - User Category API documentation

## Notes

- Both modules require authentication via JWT bearer token
- CategoryTemplate endpoints require admin privileges for create/update/delete operations
- UserCategory endpoints are user-scoped; users can only access their own categories
- All routes are protected by the `AuthMiddleware` which validates JWT tokens
- The modules use the same DI container pattern as existing modules
