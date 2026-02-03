# Allocation API Specification

## Overview
The Allocation API manages salary allocation rules for users. Allocations define how income should be distributed across different pockets and platforms. Each allocation can be scheduled for execution on a specific day of the month using the `execute_day` field (1-31). The system intelligently handles month-end dates: if an allocation is scheduled for day 31 but the month has fewer days, it executes on the last day of that month.

## Base URL
```
/v1/allocations
```

## Authentication
All endpoints require Bearer token authentication via the `Authorization` header.

## Data Models

### Allocation Object
```json
{
  "id": "507f1f77bcf86cd799439011",
  "user_id": "507f1f77bcf86cd799439012",
  "pocket_id": "507f1f77bcf86cd799439013",
  "user_platform_id": "507f1f77bcf86cd799439014",
  "priority": 1,
  "allocation_type": "PERCENTAGE",
  "nominal": 50.0,
  "is_active": true,
  "execute_day": 15,
  "created_at": "2026-02-03T11:14:00Z",
  "updated_at": "2026-02-03T11:14:00Z",
  "deleted_at": null
}
```

### Field Descriptions
- **id**: Unique identifier for the allocation (ObjectID)
- **user_id**: ID of the user who owns this allocation
- **pocket_id**: Optional target pocket for allocation (either pocket_id or user_platform_id must be provided)
- **user_platform_id**: Optional target platform for allocation
- **priority**: Priority level (1=HIGH, 2=MEDIUM, 3=LOW)
- **allocation_type**: Type of allocation - "PERCENTAGE" or "NOMINAL"
- **nominal**: Amount or percentage value
  - For PERCENTAGE: 0-100
  - For NOMINAL: positive number
- **is_active**: Whether the allocation is active
- **execute_day**: Optional day of month (1-31) for scheduled execution. If null, allocation is not automatically executed. Month-end handling: if set to 31 and month has fewer days, executes on the last day of that month (e.g., Feb 28/29)
- **created_at**: Creation timestamp
- **updated_at**: Last update timestamp
- **deleted_at**: Soft delete timestamp (null if not deleted)

## Endpoints

### 1. Create Allocation
Create a new allocation rule for the authenticated user.

**Request**
```
POST /v1/allocations
Content-Type: application/json
Authorization: Bearer {token}

{
  "pocket_id": "507f1f77bcf86cd799439013",
  "user_platform_id": "",
  "priority": 1,
  "allocation_type": "PERCENTAGE",
  "nominal": 50.0,
  "execute_day": 15
}
```

**Request Body Parameters**
- **pocket_id** (string, optional): Target pocket ID (24-char hex string)
- **user_platform_id** (string, optional): Target platform ID (24-char hex string)
- **priority** (integer, required): Priority level (1-3)
- **allocation_type** (string, required): "PERCENTAGE" or "NOMINAL"
- **nominal** (number, required): Amount or percentage (> 0)
- **execute_day** (integer, optional): Day of month (1-31) for scheduled execution

**Validation Rules**
- At least one of `pocket_id` or `user_platform_id` must be provided
- `priority` must be between 1 and 3
- `allocation_type` must be "PERCENTAGE" or "NOMINAL"
- `nominal` must be greater than 0
- If `allocation_type` is "PERCENTAGE", `nominal` cannot exceed 100
- Target pocket/platform must belong to the user and be active

**Response (201 Created)**
```json
{
  "success": true,
  "message": "Allocation created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439012",
    "pocket_id": "507f1f77bcf86cd799439013",
    "user_platform_id": null,
    "priority": 1,
    "allocation_type": "PERCENTAGE",
    "nominal": 50.0,
    "is_active": true,
    "execute_day": 15,
    "created_at": "2026-02-03T11:14:00Z",
    "updated_at": "2026-02-03T11:14:00Z",
    "deleted_at": null
  }
}
```

**Error Responses**
- 400: Bad request (validation error)
- 401: Unauthorized (missing/invalid token)

---

### 2. Get Allocation by ID
Retrieve a specific allocation by its ID.

**Request**
```
GET /v1/allocations/{id}
Authorization: Bearer {token}
```

**Path Parameters**
- **id** (string, required): Allocation ID (24-char hex string)

**Response (200 OK)**
```json
{
  "success": true,
  "message": "Allocation retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439012",
    "pocket_id": "507f1f77bcf86cd799439013",
    "user_platform_id": null,
    "priority": 1,
    "allocation_type": "PERCENTAGE",
    "nominal": 50.0,
    "is_active": true,
    "execute_date": "2026-02-15T00:00:00Z",
    "created_at": "2026-02-03T11:14:00Z",
    "updated_at": "2026-02-03T11:14:00Z",
    "deleted_at": null
  }
}
```

**Error Responses**
- 400: Bad request (invalid ID format)
- 401: Unauthorized
- 404: Allocation not found

---

### 3. List All Allocations
Retrieve all allocations for the authenticated user.

**Request**
```
GET /v1/allocations
Authorization: Bearer {token}
```

**Response (200 OK)**
```json
{
  "success": true,
  "message": "Allocations retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "user_id": "507f1f77bcf86cd799439012",
      "pocket_id": "507f1f77bcf86cd799439013",
      "user_platform_id": null,
      "priority": 1,
      "allocation_type": "PERCENTAGE",
      "nominal": 50.0,
      "is_active": true,
      "execute_date": "2026-02-15T00:00:00Z",
      "created_at": "2026-02-03T11:14:00Z",
      "updated_at": "2026-02-03T11:14:00Z",
      "deleted_at": null
    }
  ]
}
```

**Error Responses**
- 401: Unauthorized

---

### 4. Update Allocation
Update an existing allocation.

**Request**
```
PUT /v1/allocations/{id}
Content-Type: application/json
Authorization: Bearer {token}

{
  "priority": 2,
  "nominal": 75.0,
  "is_active": true,
  "execute_day": 20
}
```

**Path Parameters**
- **id** (string, required): Allocation ID (24-char hex string)

**Request Body Parameters** (all optional)
- **pocket_id** (string, optional): New target pocket ID
- **user_platform_id** (string, optional): New target platform ID
- **priority** (integer, optional): New priority level (1-3)
- **allocation_type** (string, optional): "PERCENTAGE" or "NOMINAL"
- **nominal** (number, optional): New amount or percentage (> 0)
- **is_active** (boolean, optional): Active status
- **execute_day** (integer, optional): New day of month (1-31) for scheduled execution

**Validation Rules**
- Same as create endpoint for provided fields
- At least one field must be provided for update

**Response (200 OK)**
```json
{
  "success": true,
  "message": "Allocation updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439012",
    "pocket_id": "507f1f77bcf86cd799439013",
    "user_platform_id": null,
    "priority": 2,
    "allocation_type": "PERCENTAGE",
    "nominal": 75.0,
    "is_active": true,
    "execute_day": 20,
    "created_at": "2026-02-03T11:14:00Z",
    "updated_at": "2026-02-03T12:00:00Z",
    "deleted_at": null
  }
}
```

**Error Responses**
- 400: Bad request (validation error)
- 401: Unauthorized
- 404: Allocation not found

---

### 5. Delete Allocation
Soft delete an allocation (marks as deleted without removing from database).

**Request**
```
DELETE /v1/allocations/{id}
Authorization: Bearer {token}
```

**Path Parameters**
- **id** (string, required): Allocation ID (24-char hex string)

**Response (200 OK)**
```json
{
  "success": true,
  "message": "Allocation deleted successfully",
  "data": null
}
```

**Error Responses**
- 400: Bad request (invalid ID format)
- 401: Unauthorized
- 404: Allocation not found

---

## Allocation Execution

### Scheduled Execution
Allocations with an `execute_day` are automatically processed by the allocation cron job at 01:00 AM (Asia/Jakarta timezone) daily. The system intelligently handles month-end dates:
- If `execute_day` is set to 31 and the current month has fewer than 31 days, the allocation executes on the last day of that month
- For example, an allocation scheduled for day 31 will execute on February 28 (or 29 in leap years)

### Execution Process
1. The cron job checks the current day of the month
2. If today is the last day of the month and it's less than 31, it also fetches allocations scheduled for day 31
3. For each matching allocation, it:
   - Validates the allocation is still active
   - Validates the target pocket/platform is valid and active
   - Creates a TRANSFER transaction from the main pocket to the target
   - Updates balances for both source and target entities
   - All operations are wrapped in a database transaction for atomicity

### Execution Requirements
- Allocation must be `is_active = true`
- Target pocket/platform must be active and belong to the user
- User must have an active default user platform configured
- User must have an active main pocket

### Error Handling
If execution fails for an allocation:
- The transaction is rolled back
- The allocation remains in the database for retry
- Error is logged for monitoring
- Processing continues with other allocations

---

## Timezone Information
All timestamps in the API use ISO 8601 format with UTC timezone. However, the allocation cron job and execution logic operates in **Asia/Jakarta** timezone (UTC+7).

When scheduling allocations with `execute_day`, the day is evaluated in the Asia/Jakarta timezone context.

---

## Example Workflows

### Workflow 1: Percentage-based Allocation
Create an allocation that distributes 30% of income to a savings pocket on the 15th of each month:

```bash
curl -X POST http://localhost:8080/v1/allocations \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "pocket_id": "507f1f77bcf86cd799439013",
    "priority": 1,
    "allocation_type": "PERCENTAGE",
    "nominal": 30.0,
    "execute_day": 15
  }'
```

### Workflow 2: Fixed Amount Allocation
Create an allocation that transfers a fixed amount to an investment platform on the last day of each month:

```bash
curl -X POST http://localhost:8080/v1/allocations \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "user_platform_id": "507f1f77bcf86cd799439014",
    "priority": 2,
    "allocation_type": "NOMINAL",
    "nominal": 500000.0,
    "execute_day": 31
  }'
```
Note: Since day 31 is specified, this will execute on the last day of every month (28/29 for February, 30 for April/June/September/November, 31 for other months).

### Workflow 3: Update Execution Day
Update an existing allocation to execute on a different day of the month:

```bash
curl -X PUT http://localhost:8080/v1/allocations/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "execute_day": 20
  }'
```

---

## Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Validation error |
| 401 | Unauthorized - Missing or invalid token |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

---

## Rate Limiting
No rate limiting is currently implemented. Subject to change in future versions.

---

## Changelog

### Version 1.0.0 (2026-02-03)
- Initial API specification
- Support for percentage and nominal allocations
- Scheduled execution via `execute_date` field
- Asia/Jakarta timezone support for cron jobs
