# Category Template API Documentation

## Overview

The Category Template API manages system-level category templates that serve as blueprints for user categories. These are admin-managed templates that users can reference when creating their own categories.

## Base URL

```
/v1/category-templates
```

## Authentication

All endpoints require JWT authentication via Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Endpoints

### 1. Get All Category Templates

Retrieve all system templates and templates owned by the authenticated user.

**Endpoint:** `GET /v1/category-templates`

**Authentication:** Required (any authenticated user)

**Query Parameters:** None

**Response:**

```json
{
  "success": true,
  "statusCode": 200,
  "message": "Category templates retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Salary",
      "transaction_type": "income",
      "is_default": true,
      "color": "#4CAF50",
      "icon": "💰",
      "description": "Monthly salary income",
      "parent_id": null,
      "user_id": null,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "deleted_at": null
    }
  ]
}
```

**Status Codes:**
- `200 OK` - Templates retrieved successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication token

---

### 2. Get Category Template by ID

Retrieve a specific category template by its ID.

**Endpoint:** `GET /v1/category-templates/{id}`

**Authentication:** Required (any authenticated user)

**Path Parameters:**
- `id` (string, required) - The category template ID (24-character hex string)

**Response:**

```json
{
  "success": true,
  "statusCode": 200,
  "message": "Category template retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Salary",
    "transaction_type": "income",
    "is_default": true,
    "color": "#4CAF50",
    "icon": "💰",
    "description": "Monthly salary income",
    "parent_id": null,
    "user_id": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "deleted_at": null
  }
}
```

**Status Codes:**
- `200 OK` - Template retrieved successfully
- `400 Bad Request` - Invalid template ID format
- `401 Unauthorized` - Missing or invalid authentication token
- `404 Not Found` - Template not found

---

### 3. Create Category Template

Create a new system-level category template. **Admin only**.

**Endpoint:** `POST /v1/category-templates`

**Authentication:** Required (admin only)

**Request Body:**

```json
{
  "name": "Freelance Income",
  "transaction_type": "income",
  "is_default": false,
  "parent_id": "",
  "description": "Income from freelance projects",
  "icon": "💼",
  "color": "#2196F3"
}
```

**Request Fields:**
- `name` (string, required) - Template name (1-255 characters)
- `transaction_type` (string, optional) - Either "income" or "expense"
- `is_default` (boolean, optional) - Whether this is a default template
- `parent_id` (string, optional) - Parent template ID (24-character hex string)
- `description` (string, optional) - Template description (max 500 characters)
- `icon` (string, optional) - Icon/emoji representation (max 100 characters)
- `color` (string, optional) - Color code (max 50 characters)

**Response:**

```json
{
  "success": true,
  "statusCode": 201,
  "message": "Category template created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439012",
    "name": "Freelance Income",
    "transaction_type": "income",
    "is_default": false,
    "color": "#2196F3",
    "icon": "💼",
    "description": "Income from freelance projects",
    "parent_id": null,
    "user_id": null,
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z",
    "deleted_at": null
  }
}
```

**Status Codes:**
- `201 Created` - Template created successfully
- `400 Bad Request` - Invalid request data or validation error
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Admin access required

---

### 4. Update Category Template

Update an existing category template. **Admin only**.

**Endpoint:** `PUT /v1/category-templates/{id}`

**Authentication:** Required (admin only)

**Path Parameters:**
- `id` (string, required) - The category template ID (24-character hex string)

**Request Body:**

```json
{
  "name": "Updated Freelance Income",
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
  "message": "Category template updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439012",
    "name": "Updated Freelance Income",
    "transaction_type": "income",
    "is_default": true,
    "color": "#FF9800",
    "icon": "💼",
    "description": "Updated description",
    "parent_id": null,
    "user_id": null,
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:30:00Z",
    "deleted_at": null
  }
}
```

**Status Codes:**
- `200 OK` - Template updated successfully
- `400 Bad Request` - Invalid request data or validation error
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Admin access required
- `404 Not Found` - Template not found

---

### 5. Delete Category Template

Soft delete a category template. **Admin only**.

**Endpoint:** `DELETE /v1/category-templates/{id}`

**Authentication:** Required (admin only)

**Path Parameters:**
- `id` (string, required) - The category template ID (24-character hex string)

**Response:**

```json
{
  "success": true,
  "statusCode": 200,
  "message": "Category template deleted successfully",
  "data": null
}
```

**Status Codes:**
- `200 OK` - Template deleted successfully
- `400 Bad Request` - Invalid template ID format
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Admin access required
- `404 Not Found` - Template not found

---

## Data Model

### CategoryTemplate

```go
type CategoryTemplate struct {
    ID              string     `json:"id"`                    // MongoDB ObjectID
    Name            string     `json:"name"`                  // Template name
    TransactionType *string    `json:"transaction_type"`      // "income" or "expense"
    ParentID        *string    `json:"parent_id"`             // Parent template ID
    UserID          *string    `json:"user_id"`               // null for system templates
    Description     *string    `json:"description"`           // Optional description
    Icon            *string    `json:"icon"`                  // Optional icon/emoji
    Color           *string    `json:"color"`                 // Optional color code
    IsDefault       bool       `json:"is_default"`            // Default template flag
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

- `"category template name is required"` - Name field is missing
- `"invalid parent id"` - Parent ID format is invalid
- `"parent category template not found"` - Referenced parent template doesn't exist
- `"invalid category template id"` - Template ID format is invalid
- `"category template not found"` - Template doesn't exist or is deleted
- `"admin access required"` - User is not an admin

---

## Examples

### Create a system template for salary income

```bash
curl -X POST http://localhost:8080/v1/category-templates \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Salary",
    "transaction_type": "income",
    "is_default": true,
    "description": "Monthly salary income",
    "icon": "💰",
    "color": "#4CAF50"
  }'
```

### Get all category templates

```bash
curl -X GET http://localhost:8080/v1/category-templates \
  -H "Authorization: Bearer <token>"
```

### Update a category template

```bash
curl -X PUT http://localhost:8080/v1/category-templates/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Salary",
    "color": "#FF9800"
  }'
```

### Delete a category template

```bash
curl -X DELETE http://localhost:8080/v1/category-templates/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>"
```

---

## Notes

- All timestamps are in UTC format (ISO 8601)
- Soft deletes are used; deleted templates are not returned in list queries
- Parent-child relationships are supported for template hierarchies
- System templates have `user_id` set to null
- All IDs are MongoDB ObjectIDs represented as 24-character hex strings
