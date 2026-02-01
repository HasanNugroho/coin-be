# Financial Management System - Quick Start Guide

## Setup

1. **Start MongoDB and Redis**
```bash
# MongoDB
mongod --dbpath /path/to/data

# Redis
redis-server
```

2. **Configure Environment**
```bash
cp .env.example .env
# Edit .env with your settings
```

3. **Run the Application**
```bash
make dev
# or
air -c .air.toml
```

## Step-by-Step Usage

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "phone": "081234567890",
    "password": "securepassword",
    "name": "John Doe"
  }'
```

**What happens:**
- User account created
- 8 default categories created (Salary, Bonus, Investment, Food, Transportation, etc.)
- 4 default allocations created (Bills 40%, Emergency Fund 10%, Investment 30%, Savings 20%)

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "user": {...},
    "token_pair": {
      "access_token": "eyJhbGc...",
      "refresh_token": "eyJhbGc..."
    }
  }
}
```

Save the `access_token` for subsequent requests.

### 3. View Your Categories

```bash
curl -X GET http://localhost:8080/api/categories \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 4. View Your Allocations

```bash
curl -X GET http://localhost:8080/api/allocations \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5. Create Income Transaction

```bash
curl -X POST http://localhost:8080/api/transactions/income \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "income",
    "amount": 5000000,
    "category_id": "SALARY_CATEGORY_ID",
    "description": "Monthly salary",
    "transaction_date": "2024-01-15T00:00:00Z"
  }'
```

**What happens:**
- Transaction saved
- Allocation engine runs automatically
- Income distributed to allocations based on priority:
  - Bills & Utilities (Priority 1): Rp 2,000,000 (40%)
  - Emergency Fund (Priority 2): Rp 500,000 (10%)
  - Investment (Priority 3): Rp 1,500,000 (30%)
  - Savings (Priority 4): Rp 1,000,000 (20%)
- Allocation logs created for audit trail

**Response:**
```json
{
  "status": "success",
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

### 6. Create Expense Transaction

```bash
curl -X POST http://localhost:8080/api/transactions/expense \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "expense",
    "amount": 150000,
    "category_id": "FOOD_CATEGORY_ID",
    "allocation_id": "BILLS_ALLOCATION_ID",
    "description": "Grocery shopping",
    "transaction_date": "2024-01-16T00:00:00Z"
  }'
```

**What happens:**
- System validates allocation has sufficient balance
- Transaction saved
- Allocation balance reduced by expense amount

### 7. View Dashboard

```bash
curl -X GET http://localhost:8080/api/reports/dashboard \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response:**
```json
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

### 8. Create Saving Target

```bash
curl -X POST http://localhost:8080/api/targets \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "allocation_id": "EMERGENCY_FUND_ALLOCATION_ID",
    "name": "Emergency Fund Goal",
    "target_amount": 10000000,
    "deadline": "2024-12-31T00:00:00Z"
  }'
```

### 9. View Income Report

```bash
curl -X GET "http://localhost:8080/api/reports/income?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 10. View Expense Report

```bash
curl -X GET "http://localhost:8080/api/reports/expense?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 11. View Allocation Report

```bash
curl -X GET http://localhost:8080/api/reports/allocation \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 12. View Target Progress

```bash
curl -X GET http://localhost:8080/api/reports/target-progress \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Advanced Usage

### Create Custom Category

```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Freelance Income",
    "type": "income",
    "icon": "💼",
    "color": "#3498db"
  }'
```

### Create Custom Allocation

```bash
curl -X POST http://localhost:8080/api/allocations \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Vacation Fund",
    "priority": 5,
    "percentage": 15,
    "target_amount": 20000000
  }'
```

### Filter Transactions

```bash
# By type
curl -X GET "http://localhost:8080/api/transactions/filter?type=income" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# By category
curl -X GET "http://localhost:8080/api/transactions/filter?category_id=CATEGORY_ID" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# By date range
curl -X GET "http://localhost:8080/api/transactions/filter?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Combined filters
curl -X GET "http://localhost:8080/api/transactions/filter?type=expense&start_date=2024-01-01&end_date=2024-01-31&limit=20" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### View Allocation Distribution Logs

```bash
# All logs
curl -X GET http://localhost:8080/api/allocations/logs \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Logs for specific allocation
curl -X GET http://localhost:8080/api/allocations/ALLOCATION_ID/logs \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Common Scenarios

### Scenario 1: Monthly Salary Processing

1. Receive salary income → Creates income transaction
2. System automatically distributes to allocations
3. View dashboard to see updated balances
4. Check allocation logs to see distribution details

### Scenario 2: Making a Purchase

1. Create expense transaction with allocation
2. System validates sufficient balance
3. Allocation balance reduced
4. View expense report to track spending

### Scenario 3: Saving for a Goal

1. Create allocation for the goal
2. Create saving target linked to allocation
3. Income automatically distributed to this allocation
4. View target progress to monitor achievement

### Scenario 4: Adjusting Allocations

1. Update allocation percentages based on needs
2. Future income will use new percentages
3. Existing balances remain unchanged
4. View allocation report to see changes

## Testing the Allocation Engine

### Test Case 1: Normal Distribution
```
Income: Rp 5,000,000
Allocations:
- Bills (Priority 1, 40%) → Rp 2,000,000
- Emergency (Priority 2, 10%) → Rp 500,000
- Investment (Priority 3, 30%) → Rp 1,500,000
- Savings (Priority 4, 20%) → Rp 1,000,000
Free Cash: Rp 0
```

### Test Case 2: Target Reached
```
Income: Rp 5,000,000
Emergency Fund already at Rp 10,000,000 (target reached)

Distribution:
- Bills (Priority 1, 40%) → Rp 2,000,000
- Emergency (Priority 2, 10%) → SKIPPED (target reached)
- Investment (Priority 3, 30%) → Rp 1,500,000
- Savings (Priority 4, 20%) → Rp 1,000,000
Free Cash: Rp 500,000 (Emergency's share)
```

### Test Case 3: Insufficient Income
```
Income: Rp 1,000,000
Allocations (same as above)

Distribution:
- Bills (Priority 1, 40%) → Rp 400,000
- Emergency (Priority 2, 10%) → Rp 100,000
- Investment (Priority 3, 30%) → Rp 300,000
- Savings (Priority 4, 20%) → Rp 200,000
Free Cash: Rp 0
```

## Troubleshooting

### Error: "insufficient balance in allocation"
- Check allocation current balance
- Ensure expense amount doesn't exceed available balance
- View allocation report to see balances

### Error: "category not found"
- Verify category ID is correct
- List categories to get valid IDs
- Ensure category belongs to your user

### Error: "allocation not found"
- Verify allocation ID is correct
- List allocations to get valid IDs
- Ensure allocation belongs to your user

### Error: "unauthorized"
- Check access token is valid
- Token may have expired (15 minutes)
- Use refresh token to get new access token

## Next Steps

1. Explore the full API documentation at `/swagger/index.html`
2. Read detailed documentation in `FINANCIAL_SYSTEM.md`
3. Customize categories and allocations for your needs
4. Set up saving targets for your goals
5. Monitor your financial progress with reports

## Tips

- **Default Allocations**: Modify the default allocations in `allocation/repository.go` to match your preferences
- **Default Categories**: Customize default categories in `category/repository.go`
- **Percentages**: Total percentage can exceed 100% - system distributes based on priority until income exhausted
- **Audit Trail**: Allocation logs provide complete history of income distribution
- **Reports**: Use date filters to analyze specific periods
- **Targets**: Link multiple targets to same allocation for different goals
