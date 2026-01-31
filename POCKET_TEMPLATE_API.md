# Pocket Template API Documentation

## Overview
The Pocket Template API provides endpoints to manage pocket templates that define the structure and configuration for user pockets. Pocket templates include information about pocket types, categories, and ordering.

## Base URL
```
/v1/pocket-templates
```

## Authentication
All endpoints require Bearer token authentication via the `Authorization` header:
```
Authorization: Bearer <access_token>
```

## Data Models

### PocketTemplate
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Emergency Fund",
  "type": "saving",
  "category_id": "507f1f77bcf86cd799439012",
  "is_default": true,
  "is_active": true,
  "order": 1,
  "icon": "piggy-bank",
  "icon_color": "#FF6B6B",
  "background_color": "#FFE5E5",
  "created_at": "2024-01-31T12:00:00Z",
  "updated_at": "2024-01-31T12:00:00Z",
  "deleted_at": null
}
```

### Field Descriptions
- **id**: MongoDB ObjectID (24-character hex string)
- **name**: Template name (1-255 characters)
- **type**: Pocket type - one of: `main`, `saving`, `allocation`
- **category_id**: Reference to Category document with type "pocket" (24-character hex string)
- **is_default**: Whether this is a default template (boolean)
- **is_active**: Whether this template is active (boolean)
- **order**: Display order (0-10000)
- **icon**: Icon name/identifier (max 100 characters, optional)
- **icon_color**: Icon color in hex format (max 50 characters, optional)
- **background_color**: Background color in hex format (max 50 characters, optional)
- **created_at**: Creation timestamp (ISO 8601)
- **updated_at**: Last update timestamp (ISO 8601)
- **deleted_at**: Soft delete timestamp (ISO 8601, null if not deleted)

## Endpoints

### 1. Create Pocket Template
**POST** `/v1/pocket-templates`

Creates a new pocket template. Requires admin access.

#### Request Body
```json
{
  "name": "Emergency Fund",
  "type": "saving",
  "category_id": "507f1f77bcf86cd799439012",
  "is_default": true,
  "is_active": true,
  "order": 1,
  "icon": "piggy-bank",
  "icon_color": "#FF6B6B",
  "background_color": "#FFE5E5"
}
```

#### Request Validation
- `name`: required, 1-255 characters
- `type`: required, enum (main, saving, allocation)
- `category_id`: required, 24-character hexadecimal string
- `is_default`: optional, boolean
- `is_active`: optional, boolean
- `order`: optional, 0-10000
- `icon`: optional, max 100 characters
- `icon_color`: optional, max 50 characters (hex color format recommended)
- `background_color`: optional, max 50 characters (hex color format recommended)

#### Response
**Status: 201 Created**
```json
{
  "success": true,
  "message": "Pocket template created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Emergency Fund",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "is_default": true,
    "is_active": true,
    "order": 1,
    "created_at": "2024-01-31T12:00:00Z",
    "updated_at": "2024-01-31T12:00:00Z",
    "deleted_at": null
  }
}
```

#### Error Responses
- **400 Bad Request**: Invalid input or validation error
  ```json
  {
    "success": false,
    "message": "validation failed: name is required; type must be one of: main, saving, allocation"
  }
  ```
- **403 Forbidden**: Admin access required
- **409 Conflict**: Pocket template name already exists

#### Example cURL
```bash
curl -X POST http://localhost:8080/v1/pocket-templates \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Emergency Fund",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "is_default": true,
    "is_active": true,
    "order": 1
  }'
```

---

### 2. Get Pocket Template by ID
**GET** `/v1/pocket-templates/{id}`

Retrieves a specific pocket template by ID.

#### Path Parameters
- `id`: Pocket template ID (24-character hex string)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pocket template retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Emergency Fund",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "is_default": true,
    "is_active": true,
    "order": 1,
    "created_at": "2024-01-31T12:00:00Z",
    "updated_at": "2024-01-31T12:00:00Z",
    "deleted_at": null
  }
}
```

#### Error Responses
- **400 Bad Request**: Invalid ID format
- **404 Not Found**: Pocket template not found

#### Example cURL
```bash
curl -X GET http://localhost:8080/v1/pocket-templates/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>"
```

---

### 3. Update Pocket Template
**PUT** `/v1/pocket-templates/{id}`

Updates an existing pocket template. Requires admin access.

#### Path Parameters
- `id`: Pocket template ID (24-character hex string)

#### Request Body
```json
{
  "name": "Emergency Fund Updated",
  "type": "saving",
  "category_id": "507f1f77bcf86cd799439012",
  "is_default": false,
  "is_active": true,
  "order": 2
}
```

#### Request Validation
- `name`: optional, 1-255 characters
- `type`: optional, enum (main, saving, allocation)
- `category_id`: optional, 24-character hexadecimal string
- `is_default`: optional, boolean
- `is_active`: optional, boolean
- `order`: optional, 0-10000

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pocket template updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Emergency Fund Updated",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "is_default": false,
    "is_active": true,
    "order": 2,
    "created_at": "2024-01-31T12:00:00Z",
    "updated_at": "2024-01-31T12:05:00Z",
    "deleted_at": null
  }
}
```

#### Error Responses
- **400 Bad Request**: Invalid input or validation error
- **403 Forbidden**: Admin access required
- **404 Not Found**: Pocket template not found
- **409 Conflict**: Pocket template name already exists

#### Example cURL
```bash
curl -X PUT http://localhost:8080/v1/pocket-templates/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Emergency Fund Updated",
    "order": 2
  }'
```

---

### 4. Delete Pocket Template
**DELETE** `/v1/pocket-templates/{id}`

Soft deletes a pocket template. Requires admin access.

#### Path Parameters
- `id`: Pocket template ID (24-character hex string)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pocket template deleted successfully",
  "data": null
}
```

#### Error Responses
- **400 Bad Request**: Invalid ID format
- **403 Forbidden**: Admin access required
- **404 Not Found**: Pocket template not found

#### Example cURL
```bash
curl -X DELETE http://localhost:8080/v1/pocket-templates/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>"
```

---

### 5. List All Pocket Templates
**GET** `/v1/pocket-templates`

Retrieves a paginated list of all pocket templates.

#### Query Parameters
- `limit`: Number of results per page (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pocket templates retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Emergency Fund",
      "type": "saving",
      "category_id": "507f1f77bcf86cd799439012",
      "is_default": true,
      "is_active": true,
      "order": 1,
      "created_at": "2024-01-31T12:00:00Z",
      "updated_at": "2024-01-31T12:00:00Z",
      "deleted_at": null
    },
    {
      "id": "507f1f77bcf86cd799439013",
      "name": "Vacation Fund",
      "type": "saving",
      "category_id": "507f1f77bcf86cd799439014",
      "is_default": false,
      "is_active": true,
      "order": 2,
      "created_at": "2024-01-31T12:01:00Z",
      "updated_at": "2024-01-31T12:01:00Z",
      "deleted_at": null
    }
  ]
}
```

#### Error Responses
- **400 Bad Request**: Invalid pagination parameters

#### Example cURL
```bash
curl -X GET "http://localhost:8080/v1/pocket-templates?limit=20&skip=0" \
  -H "Authorization: Bearer <token>"
```

---

### 6. List Active Pocket Templates
**GET** `/v1/pocket-templates/active`

Retrieves a paginated list of active pocket templates only.

#### Query Parameters
- `limit`: Number of results per page (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Active pocket templates retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Emergency Fund",
      "type": "saving",
      "category_id": "507f1f77bcf86cd799439012",
      "is_default": true,
      "is_active": true,
      "order": 1,
      "created_at": "2024-01-31T12:00:00Z",
      "updated_at": "2024-01-31T12:00:00Z",
      "deleted_at": null
    }
  ]
}
```

#### Error Responses
- **400 Bad Request**: Invalid pagination parameters

#### Example cURL
```bash
curl -X GET "http://localhost:8080/v1/pocket-templates/active?limit=10" \
  -H "Authorization: Bearer <token>"
```

---

### 7. List Pocket Templates by Type
**GET** `/v1/pocket-templates/type/{type}`

Retrieves a paginated list of pocket templates filtered by type.

#### Path Parameters
- `type`: Pocket type filter - one of: `main`, `saving`, `allocation`

#### Query Parameters
- `limit`: Number of results per page (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pocket templates retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Emergency Fund",
      "type": "saving",
      "category_id": "507f1f77bcf86cd799439012",
      "is_default": true,
      "is_active": true,
      "order": 1,
      "created_at": "2024-01-31T12:00:00Z",
      "updated_at": "2024-01-31T12:00:00Z",
      "deleted_at": null
    }
  ]
}
```

#### Error Responses
- **400 Bad Request**: Invalid type or pagination parameters

#### Example cURL
```bash
curl -X GET "http://localhost:8080/v1/pocket-templates/type/saving?limit=10" \
  -H "Authorization: Bearer <token>"
```

---

## Validation Rules

### Name
- Required for create
- Optional for update
- Length: 1-255 characters
- Must be unique (cannot duplicate existing names)

### Type
- Required for create
- Optional for update
- Allowed values: `main`, `saving`, `allocation`

### Category ID
- Required for create
- Optional for update
- Format: 24-character hexadecimal string (MongoDB ObjectID)
- Must reference a valid Category with type "pocket"

### Order
- Optional
- Range: 0-10000
- Used for sorting templates in UI

### Pagination
- **limit**: 1-1000 (default: 10)
- **skip**: ≥0 (default: 0)

## Error Handling

All error responses follow this format:
```json
{
  "success": false,
  "message": "Error description"
}
```

### Common Error Codes
- **400 Bad Request**: Invalid input, validation failed, or malformed request
- **403 Forbidden**: Insufficient permissions (admin required)
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists (e.g., duplicate name)
- **500 Internal Server Error**: Server error

## Rate Limiting
No rate limiting is currently implemented. Contact support for high-volume API usage.

## Sorting
Results are sorted by the `order` field in ascending order.

## Soft Deletes
Deleted pocket templates are soft-deleted (marked with `deleted_at` timestamp) and excluded from all queries. They can be permanently deleted from the database by administrators.

## Examples

### Create a Saving Pocket Template
```bash
curl -X POST http://localhost:8080/v1/pocket-templates \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Emergency Fund",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "is_default": true,
    "is_active": true,
    "order": 1
  }'
```

### Get All Active Saving Templates
```bash
curl -X GET "http://localhost:8080/v1/pocket-templates/type/saving?limit=20" \
  -H "Authorization: Bearer eyJhbGc..."
```

### Update Template Order
```bash
curl -X PUT http://localhost:8080/v1/pocket-templates/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "order": 5
  }'
```

## Support
For issues or questions, contact the API support team.
