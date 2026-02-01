# Financial Management System Documentation

## Overview

A comprehensive financial management system built with Go, featuring income/expense tracking, smart allocation engine, saving targets, and detailed financial reports.

## Features

### 1. Category Management
- **Income Categories**: Salary, Bonus, Investment, etc.
- **Expense Categories**: Food, Transportation, Shopping, Bills, etc.
- **Custom Categories**: Users can create their own categories
- **Default Categories**: Auto-created on user registration

### 2. Transaction Management (CORE)
- **Income Tracking**: Record income with automatic allocation distribution
- **Expense Tracking**: Record expenses with allocation tracking
- **CRUD Operations**: Create, read, update, delete transactions
- **Advanced Filtering**:
  - By date range
  - By category
  - By transaction type (income/expense)
  - By allocation

### 3. Allocation System (ADVANCED)
The allocation system automatically distributes income based on priorities and percentages.

#### Allocation Rules:
1. **Priority-Based**: Allocations are processed in priority order (1 = highest)
2. **Percentage-Based**: Each allocation has an ideal percentage
3. **Target-Aware**: Allocations with targets stop when target is reached
4. **Flexible**: Total percentage can exceed 100%

#### Example:
```
Income: Rp 5,000,000

Allocation          Priority    %      Target          Amount Allocated
Bills & Utilities   1           40%    -               Rp 2,000,000
Emergency Fund      2           10%    Rp 10,000,000   Rp 500,000
Investment          3           30%    -               Rp 1,500,000
Savings             4           20%    Rp 5,000,000    Rp 1,000,000

Free Cash: Rp 0
```

### 4. Saving Target
- **Target Creation**: Link targets to allocations
- **Progress Tracking**: Automatic progress calculation
- **Deadline Management**: Track days remaining
- **Auto-Sync**: Targets sync with allocation balance
- **Status Management**: Active → Completed when target reached

### 5. Financial Reports
- **Dashboard Summary**:
  - Total balance across all allocations
  - Free cash available
  - Income this month
  - Expense this month
  - Remaining balance

- **Income Report**:
  - Total income
  - Breakdown by category
  - Breakdown by month

- **Expense Report**:
  - Total expense
  - Breakdown by category
  - Breakdown by allocation
  - Breakdown by month

- **Allocation Report**:
  - Current balance per allocation
  - Distribution history
  - Progress vs target
  - Total allocated vs spent

- **Target Progress**:
  - Active targets
  - Progress percentage
  - Days remaining
  - Estimated completion

## Business Rules

### 1. User Scope Rule
- All data MUST have `user_id`
- Complete data isolation per user

### 2. Transaction Rule
- **Income** → Triggers allocation engine
- **Expense** → Reduces allocation or free cash
- **Validation**: Cannot expense more than available balance

### 3. Allocation Rule
- Ordered by priority (ASC)
- Percentage = ideal target
- Target reached → allocation skipped
- Total percentage can be > 100%

### 4. Saving Target Rule
- Target linked to allocation
- Allocation increases → target progress increases
- Target completed → status changes to "completed"

## API Endpoints

### Authentication
```
POST   /api/auth/register          - Register new user (creates default categories & allocations)
POST   /api/auth/login             - User login
POST   /api/auth/refresh           - Refresh access token
POST   /api/auth/logout            - User logout
```

### Categories (Protected)
```
POST   /api/categories             - Create category
GET    /api/categories             - Get all categories (filter by type with ?type=income|expense)
GET    /api/categories/:id         - Get category by ID
PUT    /api/categories/:id         - Update category
DELETE /api/categories/:id         - Delete category
```

### Transactions (Protected)
```
POST   /api/transactions/income    - Create income (auto-distributes to allocations)
POST   /api/transactions/expense   - Create expense
GET    /api/transactions           - Get all transactions (paginated)
GET    /api/transactions/filter    - Filter transactions
GET    /api/transactions/:id       - Get transaction by ID
PUT    /api/transactions/:id       - Update transaction
DELETE /api/transactions/:id       - Delete transaction
```

### Allocations (Protected)
```
POST   /api/allocations            - Create allocation
GET    /api/allocations            - Get all allocations
GET    /api/allocations/:id        - Get allocation by ID
PUT    /api/allocations/:id        - Update allocation
DELETE /api/allocations/:id        - Delete allocation
GET    /api/allocations/logs       - Get allocation distribution logs
GET    /api/allocations/:id/logs   - Get logs for specific allocation
```

### Saving Targets (Protected)
```
POST   /api/targets                - Create saving target
GET    /api/targets                - Get all targets
GET    /api/targets/:id            - Get target by ID
PUT    /api/targets/:id            - Update target
DELETE /api/targets/:id            - Delete target
```

### Reports (Protected)
```
GET    /api/reports/dashboard      - Get dashboard summary
GET    /api/reports/income         - Get income report (with date filters)
GET    /api/reports/expense        - Get expense report (with date filters)
GET    /api/reports/allocation     - Get allocation report
GET    /api/reports/target-progress - Get target progress
```

## Data Models

### Category
```go
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "name": "string",
  "type": "income|expense",
  "icon": "string",
  "color": "string",
  "is_default": "boolean",
  "created_at": "timestamp"
}
```

### Transaction
```go
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "type": "income|expense",
  "amount": "float64",
  "category_id": "ObjectID",
  "allocation_id": "ObjectID (optional)",
  "description": "string",
  "transaction_date": "timestamp",
  "is_distributed": "boolean",
  "created_at": "timestamp"
}
```

### Allocation
```go
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "name": "string",
  "priority": "int",
  "percentage": "float64",
  "current_amount": "float64",
  "target_amount": "float64 (optional)",
  "is_active": "boolean",
  "created_at": "timestamp"
}
```

### Allocation Log
```go
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "allocation_id": "ObjectID",
  "transaction_id": "ObjectID",
  "income_amount": "float64",
  "allocated_amount": "float64",
  "percentage": "float64",
  "priority": "int",
  "created_at": "timestamp"
}
```

### Saving Target
```go
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "allocation_id": "ObjectID",
  "name": "string",
  "target_amount": "float64",
  "current_amount": "float64",
  "deadline": "timestamp",
  "status": "active|completed",
  "created_at": "timestamp"
}
```

## System Flow

### 1. Registration Flow
```
User registers
  → Create user account
  → Create default categories (8 categories)
  → Create default allocations (4 allocations)
  → Return authentication token
```

### 2. Income Flow
```
User inputs income
  → Save transaction
  → Run allocation engine:
      - Get active allocations (sorted by priority)
      - For each allocation:
          * Check if target reached (skip if yes)
          * Calculate ideal amount (income × percentage)
          * Allocate (min of ideal and remaining)
          * Save allocation log
          * Update allocation balance
      - Calculate free cash (remaining amount)
  → Update saving targets
  → Return distribution summary
```

### 3. Expense Flow
```
User inputs expense
  → Validate balance (allocation or free cash)
  → Save transaction
  → Reduce allocation balance (if allocation specified)
  → Return expense confirmation
```

## Allocation Engine Algorithm

```go
income = 5,000,000
remaining = income

allocations = getActiveAllocations() // sorted by priority ASC

for allocation in allocations:
    // Skip if target reached
    if allocation.target_amount != null 
       and allocation.current_amount >= allocation.target_amount:
       continue

    // Calculate ideal allocation
    ideal = income × allocation.percentage / 100

    // Allocate (limited by remaining)
    allocated = min(ideal, remaining)

    // Respect target cap
    if allocation.target_amount != null:
        max_allowed = allocation.target_amount - allocation.current_amount
        allocated = min(allocated, max_allowed)

    // Save log and update balance
    save_allocation_log(...)
    allocation.current_amount += allocated
    remaining -= allocated

    if remaining <= 0:
        break

free_cash = remaining
```

## Project Structure

```
internal/modules/
├── category/
│   ├── models.go           - Category data model
│   ├── repository.go       - Database operations
│   ├── service.go          - Business logic
│   ├── controller.go       - HTTP handlers
│   ├── routes.go           - Route definitions
│   ├── module.go           - DI registration
│   └── dto/
│       ├── request.go      - Request DTOs
│       └── response.go     - Response DTOs
│
├── transaction/
│   ├── models.go
│   ├── repository.go
│   ├── service.go
│   ├── controller.go
│   ├── routes.go
│   ├── module.go
│   └── dto/
│
├── allocation/
│   ├── models.go
│   ├── repository.go
│   ├── service.go
│   ├── controller.go
│   ├── routes.go
│   ├── module.go
│   ├── engine.go           - Allocation distribution engine
│   └── dto/
│
├── target/
│   ├── models.go
│   ├── repository.go
│   ├── service.go
│   ├── controller.go
│   ├── routes.go
│   ├── module.go
│   └── dto/
│
└── report/
    ├── models.go
    ├── service.go
    ├── controller.go
    ├── routes.go
    └── module.go
```

## Usage Examples

### 1. Create Income Transaction
```bash
POST /api/transactions/income
{
  "type": "income",
  "amount": 5000000,
  "category_id": "category_id_here",
  "description": "Monthly salary",
  "transaction_date": "2024-01-15T00:00:00Z"
}

Response:
{
  "status": "success",
  "message": "Income created and distributed successfully",
  "data": {
    "transaction": {...},
    "total_income": 5000000,
    "distributed": 5000000,
    "free_cash": 0,
    "distributions": [
      {
        "allocation_id": "...",
        "allocation_name": "Bills & Utilities",
        "amount": 2000000,
        "percentage": 40,
        "priority": 1
      },
      ...
    ]
  }
}
```

### 2. Create Expense Transaction
```bash
POST /api/transactions/expense
{
  "type": "expense",
  "amount": 150000,
  "category_id": "category_id_here",
  "allocation_id": "allocation_id_here",
  "description": "Grocery shopping",
  "transaction_date": "2024-01-16T00:00:00Z"
}
```

### 3. Get Dashboard Summary
```bash
GET /api/reports/dashboard

Response:
{
  "status": "success",
  "data": {
    "total_balance": 4850000,
    "free_cash": 0,
    "income_this_month": 5000000,
    "expense_this_month": 150000,
    "remaining_this_month": 4850000
  }
}
```

### 4. Create Saving Target
```bash
POST /api/targets
{
  "allocation_id": "allocation_id_here",
  "name": "Vacation Fund",
  "target_amount": 10000000,
  "deadline": "2024-12-31T00:00:00Z"
}
```

## Edge Cases Handled

1. **Total Percentage > 100%**: System handles gracefully by distributing based on priority until income is exhausted
2. **Total Percentage < 100%**: Remaining amount becomes free cash
3. **Target Already Reached**: Allocation is skipped in distribution
4. **Insufficient Balance for Expense**: Transaction is rejected with error message
5. **Concurrent Transactions**: MongoDB handles with atomic operations

## Database Collections

- `users` - User accounts
- `categories` - Income/expense categories
- `transactions` - All financial transactions
- `allocations` - User allocation configurations
- `allocation_logs` - Distribution history (audit trail)
- `saving_targets` - Saving goals

## Security

- JWT-based authentication
- All routes (except auth) are protected
- User data isolation via `user_id`
- Password hashing with bcrypt

## Performance Considerations

- Indexed fields: `user_id`, `transaction_date`, `priority`
- Pagination support on list endpoints
- Aggregation pipelines for reports
- Efficient allocation engine (O(n) complexity)

## Future Enhancements

1. Recurring transactions
2. Budget planning
3. Financial forecasting
4. Multi-currency support
5. Export to CSV/PDF
6. Mobile app integration
7. Notification system
8. Shared allocations (family accounts)
