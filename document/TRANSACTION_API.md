# Transaction API Documentation

## Overview

The Transaction module records all financial activities for a user. It integrates with the Pocket module to maintain accurate balance tracking and supports multiple transaction types.

## Transaction Types

- **INCOME**: Money entering a pocket
  - `pocket_to`: Required
  - `pocket_from`: Must be null
  - Effect: Increases `pocket_to` balance

- **EXPENSE**: Money leaving a pocket
  - `pocket_from`: Required
  - `pocket_to`: Must be null
  - Effect: Decreases `pocket_from` balance

- **TRANSFER**: Money moving between pockets
  - `pocket_from`: Required
  - `pocket_to`: Required (must differ from `pocket_from`)
  - Effect: Decreases `pocket_from`, increases `pocket_to`

- **DEBT_PAYMENT**: Payment towards debt
  - `pocket_from`: Required
  - `pocket_to`: Optional
  - Effect: Decreases `pocket_from` balance

- **WITHDRAW**: Cash withdrawal
  - `pocket_from`: Required
  - `pocket_to`: Must be null
  - Effect: Decreases `pocket_from` balance

## Transaction Fields

```json
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "type": "INCOME|EXPENSE|TRANSFER|DEBT_PAYMENT|WITHDRAW",
  "amount": 0.0,
  "pocket_from": "ObjectID (nullable)",
  "pocket_to": "ObjectID (nullable)",
  "category_id": "ObjectID (nullable)",
  "platform_id": "ObjectID (nullable)",
  "note": "string (nullable)",
  "date": "ISO8601 datetime",
  "ref": "string (nullable, external reference)",
  "created_at": "ISO8601 datetime",
  "updated_at": "ISO8601 datetime",
  "deleted_at": "ISO8601 datetime (nullable, soft delete)"
}
```

## API Endpoints

### Create Transaction

**POST** `/api/v1/transactions`

Creates a new transaction and updates pocket balances atomically.

**Request Body:**
```json
{
  "type": "INCOME",
  "amount": 100.50,
  "pocket_to": "507f1f77bcf86cd799439011",
  "category_id": "507f1f77bcf86cd799439012",
  "note": "Salary payment",
  "date": "2024-01-31T10:30:00Z",
  "ref": "SAL-2024-001"
}
```

**Response (201):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Transaction created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "user_id": "507f1f77bcf86cd799439010",
    "type": "INCOME",
    "amount": 100.50,
    "pocket_from": null,
    "pocket_to": "507f1f77bcf86cd799439011",
    "category_id": "507f1f77bcf86cd799439012",
    "platform_id": null,
    "note": "Salary payment",
    "date": "2024-01-31T10:30:00Z",
    "ref": "SAL-2024-001",
    "created_at": "2024-01-31T10:35:00Z",
    "updated_at": "2024-01-31T10:35:00Z",
    "deleted_at": null
  }
}
```

**Validation Rules:**
- `type`: Required, must be valid transaction type
- `amount`: Required, must be > 0
- `date`: Required, must be valid ISO8601 datetime
- `pocket_from`/`pocket_to`: Must be valid ObjectID if provided
- Type-specific rules enforced (see Transaction Types)
- Pockets must exist and belong to user
- Pockets must not be locked or inactive
- Sufficient balance required for debit operations

**Error Responses:**
- 400: Invalid request, validation failed, insufficient balance
- 401: Unauthorized
- 404: Pocket not found

---

### Get Transaction

**GET** `/api/v1/transactions/{id}`

Retrieves a specific transaction by ID.

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Transaction retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "user_id": "507f1f77bcf86cd799439010",
    "type": "INCOME",
    "amount": 100.50,
    "pocket_from": null,
    "pocket_to": "507f1f77bcf86cd799439011",
    "category_id": "507f1f77bcf86cd799439012",
    "platform_id": null,
    "note": "Salary payment",
    "date": "2024-01-31T10:30:00Z",
    "ref": "SAL-2024-001",
    "created_at": "2024-01-31T10:35:00Z",
    "updated_at": "2024-01-31T10:35:00Z",
    "deleted_at": null
  }
}
```

**Error Responses:**
- 401: Unauthorized
- 404: Transaction not found

---

### List User Transactions

**GET** `/api/v1/transactions?limit=10&skip=0`

Retrieves all transactions for the authenticated user with pagination.

**Query Parameters:**
- `limit`: Number of results (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Transactions retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439013",
      "user_id": "507f1f77bcf86cd799439010",
      "type": "INCOME",
      "amount": 100.50,
      "pocket_from": null,
      "pocket_to": "507f1f77bcf86cd799439011",
      "category_id": "507f1f77bcf86cd799439012",
      "platform_id": null,
      "note": "Salary payment",
      "date": "2024-01-31T10:30:00Z",
      "ref": "SAL-2024-001",
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

### List Pocket Transactions

**GET** `/api/v1/transactions/pocket/{pocket_id}?limit=10&skip=0`

Retrieves all transactions for a specific pocket (both incoming and outgoing).

**Path Parameters:**
- `pocket_id`: ID of the pocket

**Query Parameters:**
- `limit`: Number of results (default: 10, max: 1000)
- `skip`: Number of results to skip (default: 0)

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Transactions retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439013",
      "user_id": "507f1f77bcf86cd799439010",
      "type": "INCOME",
      "amount": 100.50,
      "pocket_from": null,
      "pocket_to": "507f1f77bcf86cd799439011",
      "category_id": "507f1f77bcf86cd799439012",
      "platform_id": null,
      "note": "Salary payment",
      "date": "2024-01-31T10:30:00Z",
      "ref": "SAL-2024-001",
      "created_at": "2024-01-31T10:35:00Z",
      "updated_at": "2024-01-31T10:35:00Z",
      "deleted_at": null
    }
  ]
}
```

**Error Responses:**
- 401: Unauthorized
- 404: Pocket not found

---

### Delete Transaction

**DELETE** `/api/v1/transactions/{id}`

Soft deletes a transaction (sets `deleted_at` timestamp).

**Response (200):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Transaction deleted successfully",
  "data": null
}
```

**Error Responses:**
- 401: Unauthorized
- 404: Transaction not found

---

## Validation Rules

### Type-Specific Validation

**INCOME:**
- `pocket_to` is required
- `pocket_from` must be null
- Amount must be > 0

**EXPENSE:**
- `pocket_from` is required
- `pocket_to` must be null
- Amount must be > 0
- Sufficient balance required

**TRANSFER:**
- Both `pocket_from` and `pocket_to` are required
- `pocket_from` ≠ `pocket_to`
- Amount must be > 0
- Sufficient balance in `pocket_from` required

**DEBT_PAYMENT:**
- `pocket_from` is required
- `pocket_to` is optional
- Amount must be > 0
- Sufficient balance required

**WITHDRAW:**
- `pocket_from` is required
- `pocket_to` must be null
- Amount must be > 0
- Sufficient balance required

### General Validation

- All pocket IDs must be valid ObjectIDs
- Pockets must exist and belong to the authenticated user
- Pockets must be active (`is_active: true`)
- Pockets must not be locked (`is_locked: false`)
- Category and Platform IDs must be valid ObjectIDs if provided
- Date must be valid ISO8601 format
- Note must not exceed 500 characters
- Ref must not exceed 100 characters

---

## Error Handling

### Common Error Codes

| Error | HTTP Status | Description |
|-------|------------|-------------|
| Invalid transaction type | 400 | Transaction type not in allowed list |
| Amount must be greater than 0 | 400 | Amount is zero or negative |
| Invalid date format | 400 | Date is not valid ISO8601 |
| Pocket not found | 400 | Referenced pocket does not exist |
| Unauthorized: pocket does not belong to user | 400 | User attempting to access another user's pocket |
| Pocket is locked | 400 | Pocket is locked and cannot be used |
| Pocket is not active | 400 | Pocket is inactive |
| Insufficient balance | 400 | Not enough funds in source pocket |
| pocket_to is required for INCOME transactions | 400 | Type-specific validation failed |
| pocket_from and pocket_to cannot be the same | 400 | Transfer to same pocket |
| Unauthorized | 401 | Missing or invalid authentication token |
| Transaction not found | 404 | Transaction does not exist |

---

## Database Schema

### Transactions Collection

```javascript
db.createCollection("transactions", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["user_id", "type", "amount", "date"],
      properties: {
        _id: { bsonType: "objectId" },
        user_id: { bsonType: "objectId" },
        type: { enum: ["INCOME", "EXPENSE", "TRANSFER", "DEBT_PAYMENT", "WITHDRAW"] },
        amount: { bsonType: "double", minimum: 0 },
        pocket_from: { bsonType: ["objectId", "null"] },
        pocket_to: { bsonType: ["objectId", "null"] },
        category_id: { bsonType: ["objectId", "null"] },
        platform_id: { bsonType: ["objectId", "null"] },
        note: { bsonType: ["string", "null"] },
        date: { bsonType: "date" },
        ref: { bsonType: ["string", "null"] },
        created_at: { bsonType: "date" },
        updated_at: { bsonType: "date" },
        deleted_at: { bsonType: ["date", "null"] }
      }
    }
  }
});

// Indexes for performance
db.transactions.createIndex({ user_id: 1, date: -1 });
db.transactions.createIndex({ pocket_from: 1, date: -1 });
db.transactions.createIndex({ pocket_to: 1, date: -1 });
db.transactions.createIndex({ user_id: 1, type: 1, date: -1 });
db.transactions.createIndex({ deleted_at: 1 });
```

---

## Examples

### Example 1: Record Income

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "INCOME",
    "amount": 5000.00,
    "pocket_to": "507f1f77bcf86cd799439011",
    "category_id": "507f1f77bcf86cd799439012",
    "note": "Monthly salary",
    "date": "2024-01-31T10:30:00Z",
    "ref": "SAL-2024-01"
  }'
```

### Example 2: Record Expense

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "EXPENSE",
    "amount": 50.00,
    "pocket_from": "507f1f77bcf86cd799439011",
    "category_id": "507f1f77bcf86cd799439013",
    "note": "Grocery shopping",
    "date": "2024-01-31T14:20:00Z"
  }'
```

### Example 3: Transfer Between Pockets

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "TRANSFER",
    "amount": 1000.00,
    "pocket_from": "507f1f77bcf86cd799439011",
    "pocket_to": "507f1f77bcf86cd799439014",
    "note": "Transfer to savings",
    "date": "2024-01-31T15:00:00Z"
  }'
```

### Example 4: List User Transactions

```bash
curl -X GET "http://localhost:8080/api/v1/transactions?limit=20&skip=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Example 5: List Pocket Transactions

```bash
curl -X GET "http://localhost:8080/api/v1/transactions/pocket/507f1f77bcf86cd799439011?limit=20&skip=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Implementation Notes

### Atomic Operations

Transaction creation updates pocket balances atomically. If balance update fails, the transaction is rolled back (soft deleted).

### Soft Deletes

Transactions use soft deletes (setting `deleted_at` timestamp). Deleted transactions are excluded from queries but remain in the database for audit purposes.

### Balance Consistency

- Pocket balances are updated immediately upon transaction creation
- Each transaction type follows specific balance update rules
- Insufficient balance validation prevents negative balances

### Authorization

- Users can only access their own transactions
- Users can only use their own pockets
- Unauthorized access attempts return 401 or 400 errors

### Pagination

- Default limit: 10 records
- Maximum limit: 1000 records
- Results sorted by date (newest first)
