# Category Module Documentation

## Overview
The Category module provides CRUD operations for managing transaction and pocket categories in the Coin financial management system. It supports hierarchical categories with parent-child relationships and implements soft delete functionality.

## Features
- **Create Categories**: Admin-only endpoint to create new categories
- **Read Categories**: Retrieve categories by ID, type, or as subcategories
- **Update Categories**: Admin-only endpoint to modify category details
- **Delete Categories**: Soft delete implementation (marks as deleted without removing from database)
- **Hierarchical Support**: Support for parent-child category relationships
- **Type Filtering**: Filter categories by type (transaction or pocket)
- **Default Categories**: Mark categories as default for quick access

## Data Model

### Category Schema
```
{
  _id: ObjectID (Primary Key)
  name: String (Unique, Required)
  type: String (Enum: "transaction", "pocket", Required)
  is_default: Boolean (Default: false)
  is_deleted: Boolean (Default: false, for soft delete)
  parent_id: ObjectID (Optional, reference to parent category)
  created_at: DateTime (Auto-set on creation)
  updated_at: DateTime (Auto-updated on modification)
  deleted_at: DateTime (Set when soft deleted)
}
```

### Field Descriptions
- **_id**: MongoDB ObjectID, automatically generated
- **name**: Category name, must be unique across all categories
- **type**: Category type - either "transaction" (for expense/income tracking) or "pocket" (for savings pockets)
- **is_default**: Boolean flag to mark as default category
- **is_deleted**: Boolean flag for soft delete (true = deleted, false = active)
- **parent_id**: Reference to parent category for subcategories (null for top-level categories)
- **created_at**: Timestamp when category was created
- **updated_at**: Timestamp of last modification
- **deleted_at**: Timestamp when category was soft deleted (null if not deleted)

## API Endpoints

### 1. Create Category
**POST** `/v1/categories`

**Authentication**: Required (Admin only)

**Request Body**:
```json
{
  "name": "Food & Dining",
  "type": "transaction",
  "is_default": true,
  "parent_id": ""
}
```

**Response**: 201 Created
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Category created successfully",
  "data": {
    "id": "697ce57ea135e8c451bb2b46",
    "name": "Food & Dining",
    "type": "transaction",
    "is_default": true,
    "parent_id": null,
    "created_at": "2026-01-31T10:30:00Z",
    "updated_at": "2026-01-31T10:30:00Z",
    "deleted_at": null
  }
}
```

### 2. Get Category by ID
**GET** `/v1/categories/{id}`

**Authentication**: Required

**Response**: 200 OK
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Category retrieved successfully",
  "data": {
    "id": "697ce57ea135e8c451bb2b46",
    "name": "Food & Dining",
    "type": "transaction",
    "is_default": true,
    "parent_id": null,
    "created_at": "2026-01-31T10:30:00Z",
    "updated_at": "2026-01-31T10:30:00Z",
    "deleted_at": null
  }
}
```

### 3. Update Category
**PUT** `/v1/categories/{id}`

**Authentication**: Required (Admin only)

**Request Body**:
```json
{
  "name": "Food & Dining Updated",
  "type": "transaction",
  "is_default": false,
  "parent_id": ""
}
```

**Response**: 200 OK

### 4. Delete Category (Soft Delete)
**DELETE** `/v1/categories/{id}`

**Authentication**: Required (Admin only)

**Response**: 200 OK
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Category deleted successfully",
  "data": null
}
```

### 5. List All Categories
**GET** `/v1/categories?limit=10&skip=0`

**Authentication**: Required

**Query Parameters**:
- `limit`: Number of results per page (default: 10)
- `skip`: Number of results to skip (default: 0)

**Response**: 200 OK
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": "697ce57ea135e8c451bb2b46",
      "name": "Food & Dining",
      "type": "transaction",
      "is_default": true,
      "parent_id": null,
      "created_at": "2026-01-31T10:30:00Z",
      "updated_at": "2026-01-31T10:30:00Z",
      "deleted_at": null
    }
  ]
}
```

### 6. List Categories by Type
**GET** `/v1/categories/type/{type}?limit=10&skip=0`

**Authentication**: Required

**Path Parameters**:
- `type`: "transaction" or "pocket"

**Query Parameters**:
- `limit`: Number of results per page (default: 10)
- `skip`: Number of results to skip (default: 0)

**Response**: 200 OK

### 7. List Subcategories
**GET** `/v1/categories/{parent_id}/subcategories`

**Authentication**: Required

**Path Parameters**:
- `parent_id`: Parent category ID

**Response**: 200 OK
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Subcategories retrieved successfully",
  "data": [
    {
      "id": "697ce57ea135e8c451bb2b48",
      "name": "Restaurants",
      "type": "transaction",
      "is_default": false,
      "parent_id": "697ce57ea135e8c451bb2b46",
      "created_at": "2026-01-31T10:32:00Z",
      "updated_at": "2026-01-31T10:32:00Z",
      "deleted_at": null
    }
  ]
}
```

## Module Structure

### Files
```
internal/modules/category/
├── models.go           # Data models and constants
├── repository.go       # Database operations
├── service.go          # Business logic
├── controller.go       # HTTP handlers
├── routes.go           # Route definitions
├── module.go           # Dependency injection
└── dto/
    ├── request.go      # Request DTOs
    └── response.go     # Response DTOs
```

### Key Components

#### Models (models.go)
- `Category`: Main category struct
- Constants: `TypeTransaction`, `TypePocket`

#### Repository (repository.go)
- `CreateCategory()`: Insert new category
- `GetCategoryByID()`: Fetch by ID (excludes deleted)
- `GetCategoryByName()`: Fetch by name (excludes deleted)
- `UpdateCategory()`: Update category fields
- `DeleteCategory()`: Soft delete (sets is_deleted=true)
- `ListCategories()`: Paginated list of active categories
- `ListCategoriesByType()`: Filter by type with pagination
- `ListSubcategories()`: Get all subcategories of a parent
- `CountCategories()`: Count active categories

#### Service (service.go)
- `CreateCategory()`: Validates and creates category
- `GetCategoryByID()`: Retrieves category with validation
- `UpdateCategory()`: Updates with validation
- `DeleteCategory()`: Soft deletes category
- `ListCategories()`: Lists with pagination
- `ListCategoriesByType()`: Lists filtered by type
- `ListSubcategories()`: Lists subcategories

#### Controller (controller.go)
- HTTP request handlers for all endpoints
- Request validation and response formatting
- Swagger documentation comments

#### Routes (routes.go)
- Public routes (authenticated users): GET endpoints
- Admin routes: POST, PUT, DELETE endpoints

#### DTOs
- **CreateCategoryRequest**: Fields for creating category
- **UpdateCategoryRequest**: Optional fields for updating
- **CategoryResponse**: Response format with all fields

## Usage Examples

### Create a Category
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Food & Dining",
    "type": "transaction",
    "is_default": true
  }'
```

### Create a Subcategory
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Restaurants",
    "type": "transaction",
    "parent_id": "697ce57ea135e8c451bb2b46"
  }'
```

### Get Category by ID
```bash
curl -X GET http://localhost:8080/api/v1/categories/697ce57ea135e8c451bb2b46 \
  -H "Authorization: Bearer <token>"
```

### List All Categories
```bash
curl -X GET "http://localhost:8080/api/v1/categories?limit=20&skip=0" \
  -H "Authorization: Bearer <token>"
```

### List Transaction Categories
```bash
curl -X GET "http://localhost:8080/api/v1/categories/type/transaction?limit=20" \
  -H "Authorization: Bearer <token>"
```

### Update Category
```bash
curl -X PUT http://localhost:8080/api/v1/categories/697ce57ea135e8c451bb2b46 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Food & Dining Updated",
    "is_default": false
  }'
```

### Delete Category (Soft Delete)
```bash
curl -X DELETE http://localhost:8080/api/v1/categories/697ce57ea135e8c451bb2b46 \
  -H "Authorization: Bearer <token>"
```

### List Subcategories
```bash
curl -X GET http://localhost:8080/api/v1/categories/697ce57ea135e8c451bb2b46/subcategories \
  -H "Authorization: Bearer <token>"
```

## Soft Delete Implementation

The category module uses soft delete to maintain data integrity and allow for recovery:

1. **Deletion Process**: When a category is deleted, the `is_deleted` field is set to `true` and `deleted_at` timestamp is recorded
2. **Query Filtering**: All queries automatically exclude deleted categories (filter by `is_deleted: false`)
3. **Data Preservation**: Original data remains in the database for audit trails and recovery
4. **Cascade Behavior**: Parent categories can be deleted independently of their subcategories

## Validation Rules

### Create Category
- `name`: Required, must be unique
- `type`: Required, must be "transaction" or "pocket"
- `parent_id`: Optional, must reference an existing category if provided

### Update Category
- `name`: Optional, must be unique if provided
- `type`: Optional, must be "transaction" or "pocket" if provided
- `parent_id`: Optional, must reference an existing category if provided

## Error Handling

### Common Errors
- `400 Bad Request`: Invalid input data or validation failure
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Admin access required for write operations
- `404 Not Found`: Category not found or has been deleted

## Integration with Main API

The category module is registered in `cmd/api/main.go`:

```go
// Register module
category.Register(builder)

// Setup routes
categoryController := appContainer.Get("categoryController").(*category.Controller)
categoryRoutes := api.Group("/v1/categories")
categoryRoutes.Use(middleware.AuthMiddleware(jwtManager, db))
category.RegisterRoutes(categoryRoutes, categoryController)
```

## Database Indexes

Recommended MongoDB indexes for optimal performance:

```javascript
// Index on name for uniqueness
db.categories.createIndex({ "name": 1 }, { unique: true })

// Index on type for filtering
db.categories.createIndex({ "type": 1, "is_deleted": 1 })

// Index on parent_id for subcategory queries
db.categories.createIndex({ "parent_id": 1, "is_deleted": 1 })

// Index on is_deleted for soft delete queries
db.categories.createIndex({ "is_deleted": 1 })
```

## Future Enhancements

- Category icons/colors for UI display
- Category ordering/sorting preferences
- Bulk operations (create/update/delete multiple)
- Category templates for quick setup
- Category usage statistics
- Category archival instead of deletion
- Multi-language category names
