# User Category API Documentation

## Overview

The User Category API manages per-user categories that allow users to create and customize their own expense and income categories. Users can optionally reference a CategoryTemplate when creating their categories.

## Base URL

```
/v1/user-categories
```

## Authentication

All endpoints require JWT authentication via Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Endpoints

### 1. Get All User Categories

Retrieve all categories belonging to the authenticated user.

**Endpoint:** `GET /v1/user-categories`

**Authentication:** Required (authenticated user)

**Query Parameters:** None

**Response:**

```json
{
  "success": true,
  "statusCode": 200,
  "message": "User categories retrieved successfully",
  "data": [
    {
      "id": "607f1f77bcf86cd799439011",
      "user_id": "507f1f77bcf86cd799439010",
      "template_id": "507f1f77bcf86cd799439011",
      "name": "My Salary",
      "transaction_type": "income",
      "is_default": true,
      "color": "#4CAF50",
      "icon": "💰",
      "description": "My monthly salary",
      "parent_id": null,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "deleted_at": null
    }
  ]
}
```

**Status Codes:**
- `200 OK` - Categories retrieved successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication token

---

### 2. Get User Category by ID

Retrieve a specific category belonging to the authenticated user.

**Endpoint:** `GET /v1/user-categories/{id}`

**Authentication:** Required (authenticated user)

**Path Parameters:**
- `id` (string, required) - The user category ID (24-character hex string)

**Response:**

```json
{
  "success": true,
  "statusCode": 200,
  "message": "User category retrieved successfully",
  "data": {
    "id": "607f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439010",
    "template_id": "507f1f77bcf86cd799439011",
    "name": "My Salary",
    "transaction_type": "income",
    "is_default": true,
    "color": "#4CAF50",
    "icon": "💰",
    "description": "My monthly salary",
    "parent_id": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "deleted_at": null
  }
}
```

**Status Codes:**
- `200 OK` - Category retrieved successfully
- `400 Bad Request` - Invalid category ID format
- `401 Unauthorized` - Missing or invalid authentication token
- `404 Not Found` - Category not found or belongs to another user

---

### 3. Create User Category

Create a new category for the authenticated user.

**Endpoint:** `POST /v1/user-categories`

**Authentication:** Required (authenticated user)

**Request Body:**

```json
{
  "name": "Freelance Income",
  "template_id": "507f1f77bcf86cd799439012",
  "transaction_type": "income",
  "is_default": false,
  "parent_id": "",
  "description": "Income from freelance projects",
  "icon": "💼",
  "color": "#2196F3"
}
```

**Request Fields:**
- `name` (string, required) - Category name (1-255 characters)
- `template_id` (string, optional) - Reference to a CategoryTemplate ID (24-character hex string)
- `transaction_type` (string, optional) - Either "income" or "expense"
- `is_default` (boolean, optional) - Whether this is a default category for the user
- `parent_id` (string, optional) - Parent category ID (24-character hex string, must belong to same user)
- `description` (string, optional) - Category description (max 500 characters)
- `icon` (string, optional) - Icon/emoji representation (max 100 characters)
- `color` (string, optional) - Color code (max 50 characters)

**Response:**

```json
{
  "success": true,
  "statusCode": 201,
  "message": "User category created successfully",
  "data": {
    "id": "607f1f77bcf86cd799439012",
    "user_id": "507f1f77bcf86cd799439010",
    "template_id": "507f1f77bcf86cd799439012",
    "name": "Freelance Income",
    "transaction_type": "income",
    "is_default": false,
    "color": "#2196F3",
    "icon": "💼",
    "description": "Income from freelance projects",
    "parent_id": null,
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z",
    "deleted_at": null
  }
}
```

**Status Codes:**
- `201 Created` - Category created successfully
- `400 Bad Request` - Invalid request data or validation error
- `401 Unauthorized` - Missing or invalid authentication token

---

### 4. Update User Category

Update an existing category belonging to the authenticated user.

**Endpoint:** `PUT /v1/user-categories/{id}`

**Authentication:** Required (authenticated user)

**Path Parameters:**
- `id` (string, required) - The user category ID (24-character hex string)

**Request Body:**

```json
{
  "name": "Updated Freelance Income",
  "template_id": "507f1f77bcf86cd799439012",
  "transaction_type": "income",
  "is_default": true,
  "parent_id": "",
  "description": "Updated description",
  "icon": "💼",
  "color": "#FF9800"
}
```

**Request Fields:** Same as Create endpoint, all optional

**Response:**

```json
{
  "success": true,
  "statusCode": 200,
  "message": "User category updated successfully",
  "data": {
    "id": "607f1f77bcf86cd799439012",
    "user_id": "507f1f77bcf86cd799439010",
    "template_id": "507f1f77bcf86cd799439012",
    "name": "Updated Freelance Income",
    "transaction_type": "income",
    "is_default": true,
    "color": "#FF9800",
    "icon": "💼",
    "description": "Updated description",
    "parent_id": null,
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:30:00Z",
    "deleted_at": null
  }
}
```

**Status Codes:**
- `200 OK` - Category updated successfully
- `400 Bad Request` - Invalid request data or validation error
- `401 Unauthorized` - Missing or invalid authentication token
- `404 Not Found` - Category not found or belongs to another user

---

### 5. Delete User Category

Soft delete a category belonging to the authenticated user.

**Endpoint:** `DELETE /v1/user-categories/{id}`

**Authentication:** Required (authenticated user)

**Path Parameters:**
- `id` (string, required) - The user category ID (24-character hex string)

**Response:**

```json
{
  "success": true,
  "statusCode": 200,
  "message": "User category deleted successfully",
  "data": null
}
```

**Status Codes:**
- `200 OK` - Category deleted successfully
- `400 Bad Request` - Invalid category ID format
- `401 Unauthorized` - Missing or invalid authentication token
- `404 Not Found` - Category not found or belongs to another user

---

## Data Model

### UserCategory

```go
type UserCategory struct {
    ID              string     `json:"id"`                    // MongoDB ObjectID
    UserID          string     `json:"user_id"`               // Owner user ID (mandatory)
    TemplateID      *string    `json:"template_id"`           // Reference to CategoryTemplate (optional)
    Name            string     `json:"name"`                  // Category name
    TransactionType *string    `json:"transaction_type"`      // "income" or "expense"
    ParentID        *string    `json:"parent_id"`             // Parent category ID (same user)
    Description     *string    `json:"description"`           // Optional description
    Icon            *string    `json:"icon"`                  // Optional icon/emoji
    Color           *string    `json:"color"`                 // Optional color code
    IsDefault       bool       `json:"is_default"`            // Default category flag
    IsDeleted       bool       `json:"is_deleted"`            // Soft delete flag
    CreatedAt       time.Time  `json:"created_at"`            // Creation timestamp
    UpdatedAt       time.Time  `json:"updated_at"`            // Last update timestamp
    DeletedAt       *time.Time `json:"deleted_at"`            // Soft delete timestamp
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "statusCode": 400,
  "message": "Error description",
  "data": null
}
```

### Common Error Messages

- `"user category name is required"` - Name field is missing
- `"invalid template id"` - Template ID format is invalid
- `"category template not found"` - Referenced template doesn't exist
- `"invalid parent id"` - Parent ID format is invalid
- `"parent user category not found"` - Referenced parent category doesn't exist or belongs to another user
- `"invalid user category id"` - Category ID format is invalid
- `"user category not found"` - Category doesn't exist, is deleted, or belongs to another user
- `"user id not found in context"` - Authentication context issue

---

## Examples

### Create a user category from a template

```bash
curl -X POST http://localhost:8080/v1/user-categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Salary",
    "template_id": "507f1f77bcf86cd799439011",
    "transaction_type": "income",
    "is_default": true,
    "description": "My monthly salary",
    "icon": "💰",
    "color": "#4CAF50"
  }'
```

### Get all user categories

```bash
curl -X GET http://localhost:8080/v1/user-categories \
  -H "Authorization: Bearer <token>"
```

### Create a custom category without a template

```bash
curl -X POST http://localhost:8080/v1/user-categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Custom Expense",
    "transaction_type": "expense",
    "is_default": false,
    "description": "My custom expense category",
    "icon": "💸",
    "color": "#F44336"
  }'
```

### Update a user category

```bash
curl -X PUT http://localhost:8080/v1/user-categories/607f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Category Name",
    "color": "#FF9800"
  }'
```

### Delete a user category

```bash
curl -X DELETE http://localhost:8080/v1/user-categories/607f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>"
```

---

## Notes

- All timestamps are in UTC format (ISO 8601)
- Soft deletes are used; deleted categories are not returned in list queries
- Parent-child relationships are supported for category hierarchies within a user's categories
- Categories are scoped to the authenticated user; users can only access their own categories
- TemplateID is optional; users can create categories without referencing a template
- All IDs are MongoDB ObjectIDs represented as 24-character hex strings
- When updating, only provide fields that need to be changed
