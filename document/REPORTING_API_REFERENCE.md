# Reporting API Reference

## Base URL
```
http://localhost:8080/api/v1/reports
```

All endpoints require `Authorization: Bearer <token>` header.

---

## Daily Reports

### Get Daily Report
Retrieve a precomputed daily financial report for a specific date.

**Endpoint**: `GET /daily`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| date | string | Yes | Report date in YYYY-MM-DD format |
| include_details | boolean | No | Include detailed breakdowns (default: true) |

**Example Request**:
```bash
GET /api/v1/reports/daily?date=2024-01-15&include_details=true
Authorization: Bearer eyJhbGc...
```

**Success Response** (200 OK):
```json
{
  "id": "507f1f77bcf86cd799439001",
  "user_id": "507f1f77bcf86cd799439002",
  "report_date": "2024-01-15T00:00:00Z",
  "opening_balance": 5000.00,
  "closing_balance": 5250.50,
  "total_income": 1500.00,
  "total_expense": 1249.50,
  "total_transfer_in": 0.00,
  "total_transfer_out": 0.00,
  "expense_by_category": [
    {
      "category_id": "507f1f77bcf86cd799439011",
      "category_name": "Food & Dining",
      "amount": 450.00,
      "transaction_count": 5
    }
  ],
  "transactions_grouped_by_pocket": [
    {
      "pocket_id": "507f1f77bcf86cd799439012",
      "pocket_name": "Main Wallet",
      "pocket_type": "main",
      "pocket_balance": 5250.50,
      "income": 1500.00,
      "expense": 1249.50,
      "transfer_in": 0.00,
      "transfer_out": 0.00,
      "transaction_count": 10
    }
  ],
  "generated_at": "2024-01-15T23:59:59Z",
  "is_final": true
}
```

**Pending Response** (202 Accepted):
```json
{
  "message": "report queued for generation"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid date format
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: Report date is in the future

---

### Generate Daily Report
Manually trigger generation of a daily report.

**Endpoint**: `POST /daily/generate`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| date | string | Yes | Report date in YYYY-MM-DD format |

**Example Request**:
```bash
POST /api/v1/reports/daily/generate?date=2024-01-15
Authorization: Bearer eyJhbGc...
```

**Success Response** (200 OK):
```json
{
  "id": "507f1f77bcf86cd799439001",
  "report_date": "2024-01-15T00:00:00Z",
  "closing_balance": 5250.50,
  "total_income": 1500.00,
  "total_expense": 1249.50,
  "generated_at": "2024-01-15T23:59:59Z",
  "is_final": false
}
```

---

## Dashboard

### Get KPI Cards
Retrieve key performance indicators for dashboard display.

**Endpoint**: `GET /dashboard/kpis`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| month | string | No | Month in YYYY-MM format (default: current month) |

**Example Request**:
```bash
GET /api/v1/reports/dashboard/kpis?month=2024-01
Authorization: Bearer eyJhbGc...
```

**Response** (200 OK):
```json
{
  "total_balance": 15750.50,
  "total_income_current_month": 5000.00,
  "total_expense_current_month": 3200.00,
  "free_money_total": 8500.00,
  "net_change_current_month": 1800.00
}
```

---

### Get Monthly Trend Chart
Retrieve 12-month income vs expense trend data.

**Endpoint**: `GET /dashboard/charts/monthly-trend`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| months | integer | No | Number of months to return (default: 12, max: 36) |

**Example Request**:
```bash
GET /api/v1/reports/dashboard/charts/monthly-trend?months=12
Authorization: Bearer eyJhbGc...
```

**Response** (200 OK):
```json
{
  "data": [
    {
      "month": "2023-01",
      "income": 4500.00,
      "expense": 3000.00,
      "net": 1500.00
    },
    {
      "month": "2023-02",
      "income": 4800.00,
      "expense": 3100.00,
      "net": 1700.00
    },
    {
      "month": "2024-01",
      "income": 5000.00,
      "expense": 3200.00,
      "net": 1800.00
    }
  ]
}
```

---

### Get Pocket Distribution Chart
Retrieve balance distribution across all pockets.

**Endpoint**: `GET /dashboard/charts/pocket-distribution`

**Example Request**:
```bash
GET /api/v1/reports/dashboard/charts/pocket-distribution
Authorization: Bearer eyJhbGc...
```

**Response** (200 OK):
```json
{
  "data": [
    {
      "pocket_id": "507f1f77bcf86cd799439012",
      "pocket_name": "Main Wallet",
      "pocket_type": "main",
      "balance": 8500.00,
      "percentage": 54.0
    },
    {
      "pocket_id": "507f1f77bcf86cd799439013",
      "pocket_name": "Savings",
      "pocket_type": "saving",
      "balance": 7250.50,
      "percentage": 46.0
    }
  ]
}
```

---

### Get Expense by Category Chart
Retrieve expense distribution by category for a specific month.

**Endpoint**: `GET /dashboard/charts/expense-by-category`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| month | string | No | Month in YYYY-MM format (default: current month) |

**Example Request**:
```bash
GET /api/v1/reports/dashboard/charts/expense-by-category?month=2024-01
Authorization: Bearer eyJhbGc...
```

**Response** (200 OK):
```json
{
  "data": [
    {
      "category_id": "507f1f77bcf86cd799439011",
      "category_name": "Food & Dining",
      "amount": 900.00,
      "percentage": 28.13,
      "transaction_count": 15
    },
    {
      "category_id": "507f1f77bcf86cd799439014",
      "category_name": "Transportation",
      "amount": 600.00,
      "percentage": 18.75,
      "transaction_count": 8
    },
    {
      "category_id": "507f1f77bcf86cd799439015",
      "category_name": "Entertainment",
      "amount": 400.00,
      "percentage": 12.50,
      "transaction_count": 5
    }
  ]
}
```

---

## AI Integration

### Get AI Financial Context
Retrieve precomputed financial context optimized for AI chatbot consumption.

**Endpoint**: `GET /ai/financial-context`

**Example Request**:
```bash
GET /api/v1/reports/ai/financial-context
Authorization: Bearer eyJhbGc...
```

**Response** (200 OK):
```json
{
  "id": "507f1f77bcf86cd799439001",
  "user_id": "507f1f77bcf86cd799439002",
  "context_date": "2024-01-15T23:59:59Z",
  "current_balance": 15750.50,
  "free_money": 8500.00,
  "last_30_days": {
    "total_income": 5000.00,
    "total_expense": 3200.00,
    "net_change": 1800.00,
    "average_daily_expense": 106.67,
    "transaction_count": 45,
    "top_expense_categories": [
      {
        "category_name": "Food & Dining",
        "amount": 900.00,
        "percentage": 28.13
      },
      {
        "category_name": "Transportation",
        "amount": 600.00,
        "percentage": 18.75
      }
    ]
  },
  "year_to_date": {
    "total_income": 15000.00,
    "total_expense": 9600.00,
    "net_change": 5400.00,
    "average_monthly_expense": 3200.00
  },
  "pockets": [
    {
      "pocket_id": "507f1f77bcf86cd799439012",
      "pocket_name": "Main Wallet",
      "pocket_type": "main",
      "balance": 8500.00,
      "percentage_of_total": 54.0,
      "monthly_trend": [
        {
          "month": "2023-12",
          "balance": 8200.00
        },
        {
          "month": "2024-01",
          "balance": 8500.00
        }
      ]
    }
  ],
  "spending_patterns": {
    "highest_expense_day_of_week": "Friday",
    "highest_expense_category": "Food & Dining",
    "average_transaction_amount": 71.11,
    "largest_transaction": 500.00,
    "smallest_transaction": 5.00
  },
  "alerts": [
    {
      "type": "high_spending",
      "message": "Spending 15% above monthly average",
      "severity": "warning"
    }
  ],
  "updated_at": "2024-01-15T00:30:00Z"
}
```

**Pending Response** (202 Accepted):
```json
{
  "message": "context being generated"
}
```

---

## Health Check

### Service Health
Check if the reporting service is operational.

**Endpoint**: `GET /health`

**Example Request**:
```bash
GET /api/v1/reports/health
```

**Response** (200 OK):
```json
{
  "status": "healthy",
  "service": "reporting"
}
```

---

## Error Handling

All errors follow this format:

```json
{
  "error": "error message describing what went wrong"
}
```

### Common HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 202 | Accepted (report/context queued for generation) |
| 400 | Bad Request (invalid parameters) |
| 401 | Unauthorized (missing/invalid token) |
| 404 | Not Found (resource doesn't exist) |
| 500 | Internal Server Error |

---

## Rate Limiting

No rate limiting is currently applied. For production, consider implementing:
- Per-user rate limits: 100 requests/minute
- Per-endpoint rate limits for expensive operations

---

## Caching Recommendations

**Cache these endpoints** (5-minute TTL):
- `GET /dashboard/kpis`
- `GET /dashboard/charts/monthly-trend`
- `GET /dashboard/charts/pocket-distribution`
- `GET /dashboard/charts/expense-by-category`
- `GET /ai/financial-context`

**Don't cache** (real-time):
- `GET /daily` (may be regenerated)
- `POST /daily/generate`

---

## Pagination

Currently, all endpoints return complete datasets. For future pagination:
- Monthly trend: Limited to 36 months max
- Category breakdown: Top 10 categories returned

---

## Data Freshness

| Endpoint | Data Freshness |
|----------|-----------------|
| Daily report | Updated at 23:59:59 UTC daily |
| Monthly trend | Updated at 00:00:01 UTC on 1st of month |
| Pocket distribution | Real-time (from live pockets collection) |
| Expense by category | Updated at 00:00:01 UTC on 1st of month |
| AI context | Updated at 00:30:00 UTC daily |

---

## Integration Examples

### React Dashboard Component
```javascript
// Fetch KPIs
const response = await fetch('/api/v1/reports/dashboard/kpis', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const kpis = await response.json();

// Display in dashboard
setTotalBalance(kpis.total_balance);
setTotalIncome(kpis.total_income_current_month);
setTotalExpense(kpis.total_expense_current_month);
```

### AI Chatbot Integration
```python
# Fetch AI context
response = requests.get(
    'http://localhost:8080/api/v1/reports/ai/financial-context',
    headers={'Authorization': f'Bearer {token}'}
)
context = response.json()

# Pass to LLM
prompt = f"""
User's current balance: {context['current_balance']}
Last 30 days expense: {context['last_30_days']['total_expense']}
Top spending category: {context['spending_patterns']['highest_expense_category']}
...
"""
```

