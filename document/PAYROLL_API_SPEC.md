# Payroll API Specification

## Overview

This document specifies the API endpoints and changes related to the payroll auto-input system. The payroll system automatically processes salary distribution on configured salary days using allocation rules.

---

## Updated Endpoints

### 1. Update User Profile with Default UserPlatform

**Endpoint**: `PUT /v1/users/profile`

**Purpose**: Update user profile including setting the default UserPlatform for payroll income

**HTTP Method**: PUT

**Authentication**: Required (Bearer token)

**Request Body**:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+6281234567890",
  "telegramId": "john_doe_123",
  "currency": "IDR",
  "baseSalary": 10000000,
  "salaryCycle": "monthly",
  "salaryDay": 25,
  "language": "id",
  "autoInputPayroll": true,
  "defaultUserPlatformId": "507f1f77bcf86cd799439011"
}
```

**Request Fields**:
- `name` (string, optional): User's full name
- `email` (string, optional): User's email address
- `phone` (string, optional): User's phone number
- `telegramId` (string, optional): Telegram ID
- `currency` (string, optional): Pay currency (IDR, USD)
- `baseSalary` (number, optional): Base salary amount
- `salaryCycle` (string, optional): Salary cycle (monthly, weekly, biweekly)
- `salaryDay` (integer, optional): Day of month for salary (1-28)
- `language` (string, optional): Preferred language (id, en)
- `autoInputPayroll` (boolean, optional): Enable automatic payroll processing
- `defaultUserPlatformId` (string, optional): ID of UserPlatform to receive payroll income

**Validation Rules**:
- `defaultUserPlatformId` must be a valid 24-character hexadecimal string
- Referenced UserPlatform must exist and be active
- Referenced UserPlatform must belong to the authenticated user
- `salaryDay` must be between 1 and 28
- `baseSalary` must be >= 0
- `autoInputPayroll` can only be enabled if `baseSalary > 0` and `defaultUserPlatformId` is set

**Success Response** (200 OK):
```json
{
  "status": "success",
  "message": "Profile updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439010",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+6281234567890",
    "telegramId": "john_doe_123",
    "currency": "IDR",
    "baseSalary": 10000000,
    "salaryCycle": "monthly",
    "salaryDay": 25,
    "language": "id",
    "autoInputPayroll": true,
    "defaultUserPlatformId": "507f1f77bcf86cd799439011",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:45:00Z"
  }
}
```

**Validation Error Response** (400 Bad Request):
```json
{
  "status": "error",
  "message": "invalid default_user_platform_id"
}
```

**Unauthorized Response** (401 Unauthorized):
```json
{
  "status": "error",
  "message": "unauthorized"
}
```

**Side Effects**:
- Setting `defaultUserPlatformId` enables payroll income to be deposited to that specific UserPlatform
- Changing `defaultUserPlatformId` affects future payroll processing
- Disabling `autoInputPayroll` stops automatic payroll processing immediately
- Enabling `autoInputPayroll` requires valid `defaultUserPlatformId` and `baseSalary > 0`

---

### 2. Get User Profile

**Endpoint**: `GET /v1/users/profile`

**Purpose**: Retrieve authenticated user's profile including payroll configuration

**HTTP Method**: GET

**Authentication**: Required (Bearer token)

**Request Body**: None

**Success Response** (200 OK):
```json
{
  "status": "success",
  "message": "Profile retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439010",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+6281234567890",
    "telegramId": "john_doe_123",
    "currency": "IDR",
    "baseSalary": 10000000,
    "salaryCycle": "monthly",
    "salaryDay": 25,
    "language": "id",
    "autoInputPayroll": true,
    "defaultUserPlatformId": "507f1f77bcf86cd799439011",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:45:00Z"
  }
}
```

**Response Fields**:
- `defaultUserPlatformId` (string, optional): ID of UserPlatform receiving payroll income
- `autoInputPayroll` (boolean): Whether automatic payroll processing is enabled
- All other standard user profile fields

**Unauthorized Response** (401 Unauthorized):
```json
{
  "status": "error",
  "message": "unauthorized"
}
```

---

## Payroll Processing Flow

### Automatic Payroll Execution

**Trigger**: Daily at 00:01 AM (1 minute after midnight)

**Processing Steps**:
1. System fetches all users where:
   - `user_profile.auto_input_payroll = true`
   - `user_profile.salary_day = today's date`
   - `user_profile.base_salary > 0`

2. For each eligible user:
   - Validates `default_user_platform_id` exists and is active
   - Checks if payroll already processed today (idempotency)
   - Creates INCOME transaction to default UserPlatform
   - Executes allocation rules in priority order
   - Records payroll execution status

3. All operations within a database transaction:
   - Rollback on any failure
   - Ensures data consistency

### Idempotency

Payroll is processed **once per user per day** using a unique key:
```
payroll_YYYY_MM_DD_user_id
```

If payroll already executed for today, the system skips processing and logs the event.

### Allocation Execution

After income transaction is created:

1. **Load Allocations**: Fetch all active allocations for user
2. **Sort by Priority**: HIGH (1) → MEDIUM (2) → LOW (3)
3. **Execute Sequentially**:
   - Calculate allocation amount (percentage or nominal)
   - Validate sufficient balance in default UserPlatform
   - Create TRANSFER transaction
   - Update balances
4. **Remaining Balance**: Stays in default UserPlatform as free cash

---

## Data Models

### UserProfile (Updated)

```json
{
  "id": "507f1f77bcf86cd799439010",
  "user_id": "507f1f77bcf86cd799439001",
  "phone": "+6281234567890",
  "telegram_id": "john_doe_123",
  "base_salary": 10000000,
  "salary_cycle": "monthly",
  "salary_day": 25,
  "pay_currency": "IDR",
  "lang": "id",
  "auto_input_payroll": true,
  "default_user_platform_id": "507f1f77bcf86cd799439011",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:45:00Z"
}
```

**New Field**:
- `default_user_platform_id` (ObjectID, optional): Reference to UserPlatform receiving payroll income

### PayrollRecord (New)

Tracks payroll execution for idempotency and auditing:

```json
{
  "id": "507f1f77bcf86cd799439020",
  "user_id": "507f1f77bcf86cd799439001",
  "year": 2024,
  "month": 1,
  "day": 25,
  "amount": 10000000,
  "status": "SUCCESS",
  "error": null,
  "created_at": "2024-01-25T00:01:30Z"
}
```

**Fields**:
- `user_id`: User who received payroll
- `year`, `month`, `day`: Date of payroll processing
- `amount`: Salary amount processed
- `status`: SUCCESS or FAILED
- `error`: Error message if failed

---

## Error Handling

### Payroll Processing Errors

**Scenario**: Default UserPlatform not found

**Behavior**:
- Payroll processing skipped for user
- PayrollRecord created with status=FAILED
- Error logged: "default user platform not found"
- Next day's processing will retry

**Scenario**: Insufficient balance for allocation

**Behavior**:
- Allocation skipped
- Processing continues with next allocation
- Remaining balance stays in default UserPlatform

**Scenario**: Database transaction failure

**Behavior**:
- Entire payroll transaction rolled back
- PayrollRecord created with status=FAILED
- Error logged with details
- User's balances unchanged

---

## Validation Rules

### Setting Default UserPlatform

1. **Must Exist**: Referenced UserPlatform must exist in database
2. **Must Be Active**: `user_platform.is_active = true`
3. **Must Belong to User**: `user_platform.user_id = authenticated_user_id`
4. **Valid Format**: Must be 24-character hexadecimal string

### Enabling Auto Payroll

1. **Base Salary Required**: `base_salary > 0`
2. **Default Platform Required**: `default_user_platform_id` must be set
3. **Salary Day Required**: `salary_day` must be configured (1-28)

---

## Examples

### Example 1: Enable Payroll with Default Platform

**Request**:
```bash
curl -X PUT http://localhost:8080/v1/users/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "baseSalary": 10000000,
    "salaryDay": 25,
    "autoInputPayroll": true,
    "defaultUserPlatformId": "507f1f77bcf86cd799439011"
  }'
```

**Response**:
```json
{
  "status": "success",
  "message": "Profile updated successfully",
  "data": {
    "baseSalary": 10000000,
    "salaryDay": 25,
    "autoInputPayroll": true,
    "defaultUserPlatformId": "507f1f77bcf86cd799439011"
  }
}
```

### Example 2: Change Default Platform

**Request**:
```bash
curl -X PUT http://localhost:8080/v1/users/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "defaultUserPlatformId": "507f1f77bcf86cd799439012"
  }'
```

**Response**:
```json
{
  "status": "success",
  "message": "Profile updated successfully",
  "data": {
    "defaultUserPlatformId": "507f1f77bcf86cd799439012"
  }
}
```

### Example 3: Disable Payroll

**Request**:
```bash
curl -X PUT http://localhost:8080/v1/users/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "autoInputPayroll": false
  }'
```

**Response**:
```json
{
  "status": "success",
  "message": "Profile updated successfully",
  "data": {
    "autoInputPayroll": false
  }
}
```

---

## Important Notes

### Atomicity

All payroll operations are atomic:
- Income transaction creation
- Balance updates
- Allocation execution
- All within a single database transaction

If any step fails, the entire payroll is rolled back and user's balances remain unchanged.

### Idempotency

Payroll cannot be processed twice on the same day for the same user:
- System checks `payroll_records` collection
- If record exists for today, processing is skipped
- Prevents duplicate income entries

### Balance Updates

Balance updates occur in this order:
1. Income transaction created and persisted
2. Default UserPlatform balance increased
3. Main pocket balance increased
4. Allocation transactions created and persisted
5. Allocation balances updated

All within a single transaction for consistency.

### Allocation Execution

Allocations are executed in priority order:
- **Priority 1 (HIGH)**: Essential allocations (e.g., emergency fund)
- **Priority 2 (MEDIUM)**: Important allocations (e.g., savings)
- **Priority 3 (LOW)**: Discretionary allocations (e.g., entertainment)

Remaining balance stays in default UserPlatform as free cash.

---

## Backward Compatibility

- `defaultUserPlatformId` is optional (null by default)
- Existing users without this field continue to work
- Payroll only processes if field is set and valid
- No breaking changes to existing endpoints

---

## Monitoring & Debugging

### Check Payroll Status

Query `payroll_records` collection:
```json
{
  "user_id": "507f1f77bcf86cd799439001",
  "year": 2024,
  "month": 1
}
```

### Verify Transactions

Check `transactions` collection for:
- Type: INCOME (payroll income)
- Ref: "payroll_YYYY_MM_DD"
- Type: TRANSFER (allocations)
- Ref: "alloc_<allocation_id>"

### Monitor Logs

Look for:
- "Starting daily payroll processing..."
- "successfully processed payroll for user..."
- "failed to process payroll for user..."
- "Daily payroll processing completed"

---

## Summary

The payroll API enables:

✓ **Explicit payroll destination** via `defaultUserPlatformId`  
✓ **Automatic payroll processing** on configured salary day  
✓ **Atomic transactions** ensuring data consistency  
✓ **Idempotent execution** preventing duplicates  
✓ **Allocation-based distribution** with priority rules  
✓ **Comprehensive error handling** with rollback  
✓ **Audit trail** via payroll records  

All changes are backward compatible and optional.
