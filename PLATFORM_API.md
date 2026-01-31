# Platform API Documentation

## Overview

The Platform module represents the source or destination of financial transactions. It answers the question: "Uang datang / pergi lewat mana?" (Where does money come from / go to?).

Platforms are system-defined entities that can be referenced by transactions to track which channel (bank, e-wallet, cash, ATM) was used for financial activities.

## Platform Types

- **BANK**: Bank transfers (BCA, BRI, Mandiri, etc.)
- **E_WALLET**: Digital wallets (GoPay, OVO, Dana, ShopeePay, etc.)
- **CASH**: Physical cash transactions
- **ATM**: ATM withdrawals

## Platform Fields

```json
{
  "id": "ObjectID",
  "name": "string",
  "type": "BANK|E_WALLET|CASH|ATM",
  "is_active": "boolean",
  "created_at": "ISO8601 datetime",
  "updated_at": "ISO8601 datetime",
  "deleted_at": "ISO8601 datetime (nullable, soft delete)"
}
```

## API Endpoints

### Create Platform (Admin Only)

**POST** `/api/v1/platforms/admin`

Creates a new platform.

**Request Body:**
```json
{
  "name": "BCA",
  "type": "BANK",
  "is_active": true
}
```

**Response (201):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Platform created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "BCA",
    "type": "BANK",
    "is_active": true,
    "created_at": "2024-01-31T10:35:00Z",
    "updated_at": "2024-01-31T10:35:00Z",
    "deleted_at": null
  }
}
```

**Validation Rules:**
- `name`: Required, 1-255 characters
- `type`: Required, must be one of: BANK, E_WALLET, CASH, ATM
- `is_active`: Defaults to true
- Platform name must be unique

**Error Responses:**
- 400: Invalid request, validation failed, duplicate name
- 403: Admin access required
- 401: Unauthorized

---

### Get Platform

**GET** `/api/v1/platforms/{id}`

Retrieves a specific platform by ID.

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Platform retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "BCA",
    "type": "BANK",
    "is_active": true,
    "created_at": "2024-01-31T10:35:00Z",
    "updated_at": "2024-01-31T10:35:00Z",
    "deleted_at": null
  }
}
```

**Error Responses:**
- 401: Unauthorized
- 404: Platform not found

---

### Update Platform (Admin Only)

**PUT** `/api/v1/platforms/admin/{id}`

Updates an existing platform.

**Request Body:**
```json
{
  "name": "BCA Transfer",
  "type": "BANK",
  "is_active": false
}
```

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Platform updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "BCA Transfer",
    "type": "BANK",
    "is_active": false,
    "created_at": "2024-01-31T10:35:00Z",
    "updated_at": "2024-01-31T10:40:00Z",
    "deleted_at": null
  }
}
```

**Validation Rules:**
- `name`: Optional, 1-255 characters if provided
- `type`: Optional, must be valid type if provided
- `is_active`: Optional boolean

**Error Responses:**
- 400: Invalid request, validation failed, duplicate name
- 403: Admin access required
- 401: Unauthorized
- 404: Platform not found

---

### Delete Platform (Admin Only)

**DELETE** `/api/v1/platforms/admin/{id}`

Soft deletes a platform (sets `deleted_at` timestamp).

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Platform deleted successfully",
  "data": null
}
```

**Error Responses:**
- 403: Admin access required
- 401: Unauthorized
- 404: Platform not found

---

### List All Platforms

**GET** `/api/v1/platforms?limit=10&skip=0`

Retrieves all platforms (including inactive ones) with pagination.

**Query Parameters:**
- `limit`: Number of results (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Platforms retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "BCA",
      "type": "BANK",
      "is_active": true,
      "created_at": "2024-01-31T10:35:00Z",
      "updated_at": "2024-01-31T10:35:00Z",
      "deleted_at": null
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "name": "GoPay",
      "type": "E_WALLET",
      "is_active": true,
      "created_at": "2024-01-31T10:36:00Z",
      "updated_at": "2024-01-31T10:36:00Z",
      "deleted_at": null
    }
  ]
}
```

**Error Responses:**
- 401: Unauthorized

---

### List Active Platforms

**GET** `/api/v1/platforms/active?limit=10&skip=0`

Retrieves only active platforms with pagination.

**Query Parameters:**
- `limit`: Number of results (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Active platforms retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "BCA",
      "type": "BANK",
      "is_active": true,
      "created_at": "2024-01-31T10:35:00Z",
      "updated_at": "2024-01-31T10:35:00Z",
      "deleted_at": null
    }
  ]
}
```

**Error Responses:**
- 401: Unauthorized

---

### List Platforms by Type

**GET** `/api/v1/platforms/type/{type}?limit=10&skip=0`

Retrieves platforms filtered by type with pagination.

**Path Parameters:**
- `type`: Platform type (BANK, E_WALLET, CASH, ATM)

**Query Parameters:**
- `limit`: Number of results (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Platforms retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "BCA",
      "type": "BANK",
      "is_active": true,
      "created_at": "2024-01-31T10:35:00Z",
      "updated_at": "2024-01-31T10:35:00Z",
      "deleted_at": null
    },
    {
      "id": "507f1f77bcf86cd799439013",
      "name": "BRI",
      "type": "BANK",
      "is_active": true,
      "created_at": "2024-01-31T10:37:00Z",
      "updated_at": "2024-01-31T10:37:00Z",
      "deleted_at": null
    }
  ]
}
```

**Error Responses:**
- 400: Invalid platform type
- 401: Unauthorized

---

## Validation Rules

### Create Platform
- `name` is required and must be 1-255 characters
- `type` is required and must be one of: BANK, E_WALLET, CASH, ATM
- `is_active` defaults to true if not provided
- Platform name must be unique (case-sensitive)

### Update Platform
- `name` is optional, 1-255 characters if provided
- `type` is optional, must be valid if provided
- `is_active` is optional boolean
- If updating name, must not duplicate existing platform names

### General Rules
- Platforms are system-defined (no user_id field)
- Soft delete behavior: `deleted_at` is set, platform becomes invisible in queries
- Deleted platforms remain in database for audit/historical transaction purposes
- Inactive platforms (`is_active: false`) are excluded from active platform queries
- Platforms can be referenced by transactions via `platform_id`

---

## Database Schema

### Platforms Collection

```javascript
db.createCollection("platforms", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["name", "type", "is_active"],
      properties: {
        _id: { bsonType: "objectId" },
        name: { bsonType: "string", minLength: 1, maxLength: 255 },
        type: { enum: ["BANK", "E_WALLET", "CASH", "ATM"] },
        is_active: { bsonType: "bool" },
        created_at: { bsonType: "date" },
        updated_at: { bsonType: "date" },
        deleted_at: { bsonType: ["date", "null"] }
      }
    }
  }
});

// Indexes for performance
db.platforms.createIndex({ name: 1 }, { unique: true });
db.platforms.createIndex({ type: 1, deleted_at: 1 });
db.platforms.createIndex({ is_active: 1, deleted_at: 1 });
db.platforms.createIndex({ deleted_at: 1 });
```

---

## Error Handling

### Common Error Codes

| Error | HTTP Status | Description |
|-------|------------|-------------|
| Platform name is required | 400 | Name field is empty |
| Invalid platform type | 400 | Type not in allowed list |
| Platform name already exists | 400 | Duplicate name |
| Invalid platform id | 400 | Malformed ObjectID |
| Platform not found | 404 | Platform does not exist or is deleted |
| Admin access required | 403 | User lacks admin privileges |
| Unauthorized | 401 | Missing or invalid authentication token |

---

## Integration with Transactions

Platforms are referenced by transactions via the `platform_id` field:

```json
{
  "id": "507f1f77bcf86cd799439013",
  "user_id": "507f1f77bcf86cd799439010",
  "type": "INCOME",
  "amount": 5000.00,
  "pocket_to": "507f1f77bcf86cd799439011",
  "platform_id": "507f1f77bcf86cd799439012",
  "note": "Salary via BCA transfer",
  "date": "2024-01-31T10:30:00Z"
}
```

### Transaction Rules with Platforms
- `platform_id` is optional in transactions
- When provided, platform must exist and not be deleted
- Deleted platforms can still be viewed in historical transactions
- Inactive platforms cannot be selected for new transactions (enforced by transaction service)

---

## Examples

### Example 1: Create Bank Platform

```bash
curl -X POST http://localhost:8080/api/v1/platforms/admin \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "BCA",
    "type": "BANK",
    "is_active": true
  }'
```

### Example 2: Create E-Wallet Platform

```bash
curl -X POST http://localhost:8080/api/v1/platforms/admin \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "GoPay",
    "type": "E_WALLET",
    "is_active": true
  }'
```

### Example 3: Create Cash Platform

```bash
curl -X POST http://localhost:8080/api/v1/platforms/admin \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cash",
    "type": "CASH",
    "is_active": true
  }'
```

### Example 4: List Active Platforms

```bash
curl -X GET "http://localhost:8080/api/v1/platforms/active?limit=20&skip=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Example 5: List Bank Platforms

```bash
curl -X GET "http://localhost:8080/api/v1/platforms/type/BANK?limit=20&skip=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Example 6: Update Platform

```bash
curl -X PUT http://localhost:8080/api/v1/platforms/admin/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "BCA Transfer",
    "is_active": false
  }'
```

### Example 7: Deactivate Platform

```bash
curl -X PUT http://localhost:8080/api/v1/platforms/admin/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": false
  }'
```

### Example 8: Delete Platform

```bash
curl -X DELETE http://localhost:8080/api/v1/platforms/admin/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## Implementation Notes

### Soft Deletes
- Platforms use soft deletes (setting `deleted_at` timestamp)
- Deleted platforms are excluded from all queries
- Historical transactions can still reference deleted platforms
- Soft-deleted platforms remain in database for audit purposes

### Name Uniqueness
- Platform names must be globally unique
- Uniqueness is enforced at the database level with a unique index
- Case-sensitive comparison

### Admin-Only Operations
- Create, Update, Delete operations require admin role
- List and Get operations are available to authenticated users
- Admin middleware enforces role-based access control

### Soft Delete Behavior
- When a platform is deleted, it becomes invisible in:
  - List all platforms queries
  - List active platforms queries
  - List by type queries
- Deleted platforms remain visible in transaction history
- Transactions can still reference deleted platforms

### Default Values
- `is_active` defaults to `true` when creating platforms
- `created_at` and `updated_at` are set automatically
- `deleted_at` is null until platform is deleted

### Performance Considerations
- Indexes on `name`, `type`, `is_active`, and `deleted_at` for fast queries
- Unique index on `name` prevents duplicates
- Composite indexes optimize filtered queries
