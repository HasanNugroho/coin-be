# Pocket API Documentation

## Overview
The Pocket API provides endpoints to manage user pockets (kantong). Pockets are containers for storing user balance with different types and purposes. Each user must have exactly one MAIN pocket.

## Base URL
```
/v1/pockets (User endpoints)
/v1/admin/pockets (Admin endpoints)
```

## Authentication
All endpoints require Bearer token authentication via the `Authorization` header:
```
Authorization: Bearer <access_token>
```

## Data Models

### Pocket
```json
{
  "id": "507f1f77bcf86cd799439011",
  "user_id": "507f1f77bcf86cd799439010",
  "name": "Main Pocket",
  "type": "main",
  "category_id": "507f1f77bcf86cd799439012",
  "balance": 1000000.00,
  "is_default": true,
  "is_active": true,
  "is_locked": false,
  "icon": "wallet",
  "icon_color": "#FF6B6B",
  "background_color": "#FFE5E5",
  "created_at": "2024-01-31T12:00:00Z",
  "updated_at": "2024-01-31T12:00:00Z",
  "deleted_at": null
}
```

### Field Descriptions
- **id**: MongoDB ObjectID (24-character hex string)
- **user_id**: Reference to User document (24-character hex string)
- **name**: Pocket name (2-255 characters)
- **type**: Pocket type - one of: `main`, `allocation`, `saving`, `debt`, `system`
- **category_id**: Reference to Category document (optional, 24-character hex string)
- **balance**: Current balance (decimal, default: 0, read-only)
- **is_default**: Whether this is the default pocket (boolean)
- **is_active**: Whether this pocket is active (boolean)
- **is_locked**: Whether this pocket is locked (boolean, prevents updates/deletes)
- **icon**: Icon name/identifier (max 100 characters, optional)
- **icon_color**: Icon color in hex format (max 50 characters, optional)
- **background_color**: Background color in hex format (max 50 characters, optional)
- **created_at**: Creation timestamp (ISO 8601)
- **updated_at**: Last update timestamp (ISO 8601)
- **deleted_at**: Soft delete timestamp (ISO 8601, null if not deleted)

### Pocket Types
- **main**: Default pocket (exactly ONE per user, cannot be deleted)
- **allocation**: Allocation pocket for budget allocation
- **saving**: Saving pocket for savings goals
- **debt**: Debt pocket for tracking debts (cicilan)
- **system**: System pocket (admin-only, locked, cannot be deleted by users)

## User Endpoints

### 1. Create Pocket
**POST** `/v1/pockets`

Creates a new pocket for the authenticated user.

#### Request Body
```json
{
  "name": "Emergency Fund",
  "type": "saving",
  "category_id": "507f1f77bcf86cd799439012",
  "icon": "piggy-bank",
  "icon_color": "#FF6B6B",
  "background_color": "#FFE5E5"
}
```

#### Request Validation
- `name`: required, 2-255 characters
- `type`: required, enum (main, allocation, saving, debt)
- `category_id`: optional, 24-character hexadecimal string
- `icon`: optional, max 100 characters
- `icon_color`: optional, max 50 characters
- `background_color`: optional, max 50 characters

#### Response
**Status: 201 Created**
```json
{
  "success": true,
  "message": "Pocket created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439010",
    "name": "Emergency Fund",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "balance": 0,
    "is_default": false,
    "is_active": true,
    "is_locked": false,
    "icon": "piggy-bank",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5",
    "created_at": "2024-01-31T12:00:00Z",
    "updated_at": "2024-01-31T12:00:00Z",
    "deleted_at": null
  }
}
```

#### Error Responses
- **400 Bad Request**: Invalid input or validation error
- **401 Unauthorized**: User not authenticated
- **409 Conflict**: User already has a MAIN pocket (if type is main)

#### Example cURL
```bash
curl -X POST http://localhost:8080/v1/pockets \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Emergency Fund",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "icon": "piggy-bank",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5"
  }'
```

---

### 2. Get Pocket by ID
**GET** `/v1/pockets/{id}`

Retrieves a specific pocket by ID. User can only access their own pockets.

#### Path Parameters
- `id`: Pocket ID (24-character hex string)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pocket retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439010",
    "name": "Emergency Fund",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "balance": 500000,
    "is_default": false,
    "is_active": true,
    "is_locked": false,
    "icon": "piggy-bank",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5",
    "created_at": "2024-01-31T12:00:00Z",
    "updated_at": "2024-01-31T12:00:00Z",
    "deleted_at": null
  }
}
```

#### Error Responses
- **400 Bad Request**: Invalid ID format
- **401 Unauthorized**: User not authenticated
- **404 Not Found**: Pocket not found or unauthorized access

#### Example cURL
```bash
curl -X GET http://localhost:8080/v1/pockets/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>"
```

---

### 3. Update Pocket
**PUT** `/v1/pockets/{id}`

Updates an existing pocket. Cannot update MAIN pocket or locked pockets.

#### Path Parameters
- `id`: Pocket ID (24-character hex string)

#### Request Body
```json
{
  "name": "Emergency Fund Updated",
  "type": "saving",
  "category_id": "507f1f77bcf86cd799439012",
  "icon": "piggy-bank",
  "icon_color": "#FF6B6B",
  "background_color": "#FFE5E5",
  "is_active": true
}
```

#### Request Validation
- `name`: optional, 2-255 characters
- `type`: optional, enum (main, allocation, saving, debt)
- `category_id`: optional, 24-character hexadecimal string
- `icon`: optional, max 100 characters
- `icon_color`: optional, max 50 characters
- `background_color`: optional, max 50 characters
- `is_active`: optional, boolean

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pocket updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439010",
    "name": "Emergency Fund Updated",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "balance": 500000,
    "is_default": false,
    "is_active": true,
    "is_locked": false,
    "icon": "piggy-bank",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5",
    "created_at": "2024-01-31T12:00:00Z",
    "updated_at": "2024-01-31T12:05:00Z",
    "deleted_at": null
  }
}
```

#### Error Responses
- **400 Bad Request**: Invalid input or validation error
- **401 Unauthorized**: User not authenticated
- **404 Not Found**: Pocket not found or unauthorized access
- **409 Conflict**: Cannot update MAIN pocket or locked pocket

#### Example cURL
```bash
curl -X PUT http://localhost:8080/v1/pockets/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Emergency Fund Updated",
    "is_active": true
  }'
```

---

### 4. Delete Pocket
**DELETE** `/v1/pockets/{id}`

Soft deletes a pocket. Cannot delete MAIN pocket or locked pockets.

#### Path Parameters
- `id`: Pocket ID (24-character hex string)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pocket deleted successfully",
  "data": null
}
```

#### Error Responses
- **400 Bad Request**: Invalid ID format or cannot delete MAIN/locked pocket
- **401 Unauthorized**: User not authenticated
- **404 Not Found**: Pocket not found or unauthorized access

#### Example cURL
```bash
curl -X DELETE http://localhost:8080/v1/pockets/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>"
```

---

### 5. List User Pockets
**GET** `/v1/pockets`

Retrieves a paginated list of user's pockets.

#### Query Parameters
- `limit`: Number of results per page (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pockets retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "user_id": "507f1f77bcf86cd799439010",
      "name": "Main Pocket",
      "type": "main",
      "category_id": null,
      "balance": 1000000,
      "is_default": true,
      "is_active": true,
      "is_locked": false,
      "icon": "wallet",
      "icon_color": "#FF6B6B",
      "background_color": "#FFE5E5",
      "created_at": "2024-01-31T12:00:00Z",
      "updated_at": "2024-01-31T12:00:00Z",
      "deleted_at": null
    }
  ]
}
```

#### Error Responses
- **400 Bad Request**: Invalid pagination parameters
- **401 Unauthorized**: User not authenticated

#### Example cURL
```bash
curl -X GET "http://localhost:8080/v1/pockets?limit=20&skip=0" \
  -H "Authorization: Bearer <token>"
```

---

### 6. List Active Pockets
**GET** `/v1/pockets/active`

Retrieves a paginated list of user's active pockets only.

#### Query Parameters
- `limit`: Number of results per page (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Active pockets retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "user_id": "507f1f77bcf86cd799439010",
      "name": "Main Pocket",
      "type": "main",
      "category_id": null,
      "balance": 1000000,
      "is_default": true,
      "is_active": true,
      "is_locked": false,
      "icon": "wallet",
      "icon_color": "#FF6B6B",
      "background_color": "#FFE5E5",
      "created_at": "2024-01-31T12:00:00Z",
      "updated_at": "2024-01-31T12:00:00Z",
      "deleted_at": null
    }
  ]
}
```

#### Error Responses
- **400 Bad Request**: Invalid pagination parameters
- **401 Unauthorized**: User not authenticated

#### Example cURL
```bash
curl -X GET "http://localhost:8080/v1/pockets/active?limit=10" \
  -H "Authorization: Bearer <token>"
```

---

### 7. Get Main Pocket
**GET** `/v1/pockets/main`

Retrieves the user's main pocket.

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Main pocket retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439010",
    "name": "Main Pocket",
    "type": "main",
    "category_id": null,
    "balance": 1000000,
    "is_default": true,
    "is_active": true,
    "is_locked": false,
    "icon": "wallet",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5",
    "created_at": "2024-01-31T12:00:00Z",
    "updated_at": "2024-01-31T12:00:00Z",
    "deleted_at": null
  }
}
```

#### Error Responses
- **401 Unauthorized**: User not authenticated
- **404 Not Found**: Main pocket not found

#### Example cURL
```bash
curl -X GET http://localhost:8080/v1/pockets/main \
  -H "Authorization: Bearer <token>"
```

---

## Admin Endpoints

### 1. Create System Pocket (Admin Only)
**POST** `/v1/admin/pockets/{user_id}`

Creates a system pocket for a user. System pockets are locked and cannot be deleted by users.

#### Path Parameters
- `user_id`: User ID (24-character hex string)

#### Request Body
```json
{
  "name": "System Pocket",
  "category_id": "507f1f77bcf86cd799439012",
  "icon": "lock",
  "icon_color": "#000000",
  "background_color": "#CCCCCC"
}
```

#### Request Validation
- `name`: required, 2-255 characters
- `category_id`: optional, 24-character hexadecimal string
- `icon`: optional, max 100 characters
- `icon_color`: optional, max 50 characters
- `background_color`: optional, max 50 characters

#### Response
**Status: 201 Created**
```json
{
  "success": true,
  "message": "System pocket created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "user_id": "507f1f77bcf86cd799439010",
    "name": "System Pocket",
    "type": "system",
    "category_id": "507f1f77bcf86cd799439012",
    "balance": 0,
    "is_default": false,
    "is_active": true,
    "is_locked": true,
    "icon": "lock",
    "icon_color": "#000000",
    "background_color": "#CCCCCC",
    "created_at": "2024-01-31T12:00:00Z",
    "updated_at": "2024-01-31T12:00:00Z",
    "deleted_at": null
  }
}
```

#### Error Responses
- **400 Bad Request**: Invalid input or validation error
- **403 Forbidden**: Admin access required
- **404 Not Found**: User not found

#### Example cURL
```bash
curl -X POST http://localhost:8080/v1/admin/pockets/507f1f77bcf86cd799439010 \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "System Pocket",
    "icon": "lock",
    "icon_color": "#000000",
    "background_color": "#CCCCCC"
  }'
```

---

### 2. List All Pockets (Admin Only)
**GET** `/v1/admin/pockets`

Retrieves a paginated list of all pockets (admin only).

#### Query Parameters
- `limit`: Number of results per page (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

#### Response
**Status: 200 OK**
```json
{
  "success": true,
  "message": "Pockets retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "user_id": "507f1f77bcf86cd799439010",
      "name": "Main Pocket",
      "type": "main",
      "category_id": null,
      "balance": 1000000,
      "is_default": true,
      "is_active": true,
      "is_locked": false,
      "icon": "wallet",
      "icon_color": "#FF6B6B",
      "background_color": "#FFE5E5",
      "created_at": "2024-01-31T12:00:00Z",
      "updated_at": "2024-01-31T12:00:00Z",
      "deleted_at": null
    }
  ]
}
```

#### Error Responses
- **400 Bad Request**: Invalid pagination parameters
- **403 Forbidden**: Admin access required

#### Example cURL
```bash
curl -X GET "http://localhost:8080/v1/admin/pockets?limit=20&skip=0" \
  -H "Authorization: Bearer <admin_token>"
```

---

## Business Rules & Constraints

### MAIN Pocket
- Each user MUST have exactly ONE MAIN pocket
- MAIN pocket cannot be deleted
- MAIN pocket cannot be updated
- MAIN pocket is created automatically during user registration or via migration

### SYSTEM Pocket
- Can only be created by admin
- Cannot be deleted by users
- Always locked (is_locked = true)
- Cannot be updated by users

### Locked Pockets
- Cannot be updated
- Cannot be deleted
- Balance can only be modified via transaction logic

### Balance
- Read-only field
- Cannot be updated directly via pocket CRUD
- Balance changes only via transaction logic
- Default value: 0

### Soft Deletes
- Deleted pockets are marked with `deleted_at` timestamp
- Deleted pockets are excluded from all queries
- Can be permanently deleted from database by administrators

## Validation Rules

### Name
- Required for create
- Optional for update
- Length: 2-255 characters

### Type
- Required for create
- Optional for update
- Allowed values: `main`, `allocation`, `saving`, `debt`
- System pockets can only be created by admin

### Category ID
- Optional
- Format: 24-character hexadecimal string (MongoDB ObjectID)
- Must reference a valid Category

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
- **401 Unauthorized**: User not authenticated
- **403 Forbidden**: Insufficient permissions (admin required)
- **404 Not Found**: Resource not found
- **409 Conflict**: Business rule violation (e.g., duplicate MAIN pocket)
- **500 Internal Server Error**: Server error

## Authorization

- Users can only access their own pockets
- Admin can access all pockets
- System pocket endpoints restricted to admin role
- MAIN pocket cannot be modified by users

## Data Initialization

When a new user is created, a MAIN pocket is automatically created with:
- name: "Main Pocket"
- type: "main"
- is_default: true
- balance: 0
- is_active: true
- is_locked: false

A migration script is available to create MAIN pockets for existing users:
```bash
go run cmd/seeder/main.go
```

## Examples

### Create a Saving Pocket
```bash
curl -X POST http://localhost:8080/v1/pockets \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Emergency Fund",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439012",
    "icon": "piggy-bank",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5"
  }'
```

### Get User's Main Pocket
```bash
curl -X GET http://localhost:8080/v1/pockets/main \
  -H "Authorization: Bearer eyJhbGc..."
```

### List All User Pockets
```bash
curl -X GET "http://localhost:8080/v1/pockets?limit=20" \
  -H "Authorization: Bearer eyJhbGc..."
```

### Admin: Create System Pocket
```bash
curl -X POST http://localhost:8080/v1/admin/pockets/507f1f77bcf86cd799439010 \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "System Pocket",
    "icon": "lock"
  }'
```

## Support
For issues or questions, contact the API support team.
