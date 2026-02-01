# Dashboard API Documentation

## Base URL
```
/api/v1/dashboard
```

## Authentication
All endpoints require `user_id` in the request context (typically from JWT token).

---

## 1. KPI Cards Endpoint

### Get Dashboard KPIs
Returns key performance indicators for the current month.

**Endpoint:**
```
GET /api/v1/dashboard/kpis
```

**Query Parameters:**
None

**Response:**
```json
{
  "success": true,
  "data": {
    "total_balance": "5000.50",
    "monthly_income": "3000.00",
    "monthly_expense": "1500.00",
    "free_money_total": "2500.00",
    "monthly_net_change": "1500.00"
  }
}
```

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/kpis" \
  -H "Authorization: Bearer <token>"
```

---

## 2. Charts Endpoint

### Get Dashboard Charts
Returns all chart data (income/expense, pocket distribution, category distribution).

**Endpoint:**
```
GET /api/v1/dashboard/charts
```

**Query Parameters:**
None

**Response:**
```json
{
  "success": true,
  "data": {
    "income_expense_chart": [
      {
        "month": "2024-01",
        "income": "3000.00",
        "expense": "1500.00"
      },
      {
        "month": "2024-02",
        "income": "3500.00",
        "expense": "1800.00"
      }
    ],
    "pocket_distribution": [
      {
        "_id": "main",
        "pockets": [
          {
            "pocket_id": "507f1f77bcf86cd799439014",
            "pocket_name": "Main Wallet",
            "balance": "2500.00"
          }
        ],
        "total_by_type": "2500.00"
      }
    ],
    "category_distribution": [
      {
        "_id": "507f1f77bcf86cd799439012",
        "category_name": "Food & Dining",
        "total_amount": "450.00",
        "transaction_count": 8
      }
    ]
  }
}
```

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/charts" \
  -H "Authorization: Bearer <token>"
```

---

## 3. Daily Reports with Date Range Filter

### Get Daily Reports by Date Range
Returns daily financial reports for a specified date range.

**Endpoint:**
```
GET /api/v1/dashboard/reports/daily
```

**Query Parameters:**
| Parameter | Type | Required | Format | Description |
|-----------|------|----------|--------|-------------|
| start_date | string | Yes | YYYY-MM-DD | Start date (inclusive) |
| end_date | string | Yes | YYYY-MM-DD | End date (inclusive) |

**Response:**
```json
{
  "success": true,
  "data": {
    "start_date": "2024-01-01",
    "end_date": "2024-01-31",
    "count": 31,
    "reports": [
      {
        "id": "507f1f77bcf86cd799439011",
        "report_date": "2024-01-31T00:00:00Z",
        "opening_balance": "5000.50",
        "closing_balance": "4850.75",
        "total_income": "2000.00",
        "total_expense": "1500.00",
        "total_transfer_in": "500.00",
        "total_transfer_out": "150.25",
        "expense_by_category": [
          {
            "category_id": "507f1f77bcf86cd799439012",
            "category_name": "Food & Dining",
            "amount": "450.00",
            "transaction_count": 8
          }
        ],
        "transactions_by_pocket": [
          {
            "pocket_id": "507f1f77bcf86cd799439014",
            "pocket_name": "Main Wallet",
            "pocket_type": "main",
            "income": "1000.00",
            "expense": "1200.00",
            "transfer_in": "200.00",
            "transfer_out": "150.25",
            "opening_balance": "3000.00",
            "closing_balance": "2850.75",
            "transaction_count": 12
          }
        ],
        "is_final": true,
        "generated_at": "2024-01-31T23:59:59Z"
      }
    ]
  }
}
```

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/daily?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer <token>"
```

**Example JavaScript/Fetch:**
```javascript
const startDate = '2024-01-01';
const endDate = '2024-01-31';

const response = await fetch(
  `/api/v1/dashboard/reports/daily?start_date=${startDate}&end_date=${endDate}`,
  {
    headers: { 'Authorization': `Bearer ${token}` }
  }
);

const data = await response.json();
console.log(data.data.reports);
```

---

## 4. Monthly Reports with Date Range Filter

### Get Monthly Reports by Date Range
Returns monthly financial summaries for a specified month range.

**Endpoint:**
```
GET /api/v1/dashboard/reports/monthly
```

**Query Parameters:**
| Parameter | Type | Required | Format | Description |
|-----------|------|----------|--------|-------------|
| start_month | string | Yes | YYYY-MM | Start month (inclusive) |
| end_month | string | Yes | YYYY-MM | End month (inclusive) |

**Response:**
```json
{
  "success": true,
  "data": {
    "start_month": "2024-01",
    "end_month": "2024-03",
    "count": 3,
    "summaries": [
      {
        "id": "507f1f77bcf86cd799439020",
        "month": "2024-01-01T00:00:00Z",
        "income": "10000.00",
        "expense": "5000.00",
        "transfer_in": "1000.00",
        "transfer_out": "500.00",
        "net": "5000.00",
        "opening_balance": "15000.00",
        "closing_balance": "20000.00",
        "expense_by_category": [
          {
            "category_id": "507f1f77bcf86cd799439012",
            "category_name": "Food & Dining",
            "amount": "1500.00",
            "percent_of_total": 0.30,
            "transaction_count": 45
          }
        ],
        "by_pocket": [
          {
            "pocket_id": "507f1f77bcf86cd799439014",
            "pocket_name": "Main Wallet",
            "pocket_type": "main",
            "income": "8000.00",
            "expense": "4000.00",
            "net": "4000.00"
          }
        ],
        "ytd_income": "10000.00",
        "ytd_expense": "5000.00",
        "ytd_net": "5000.00",
        "transaction_count": 150,
        "is_complete": true
      }
    ]
  }
}
```

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/monthly?start_month=2024-01&end_month=2024-12" \
  -H "Authorization: Bearer <token>"
```

---

## 5. Income vs Expense Chart

### Get Monthly Income/Expense Chart
Returns monthly income vs expense data for the last N months.

**Endpoint:**
```
GET /api/v1/dashboard/charts/income-expense
```

**Query Parameters:**
| Parameter | Type | Required | Default | Range | Description |
|-----------|------|----------|---------|-------|-------------|
| months | integer | No | 12 | 1-60 | Number of months to include |

**Response:**
```json
{
  "success": true,
  "data": {
    "months": 12,
    "chart": [
      {
        "month": "2023-02",
        "income": "2500.00",
        "expense": "1200.00"
      },
      {
        "month": "2023-03",
        "income": "3000.00",
        "expense": "1500.00"
      },
      {
        "month": "2024-01",
        "income": "3500.00",
        "expense": "1800.00"
      }
    ]
  }
}
```

**Example cURL:**
```bash
# Get last 12 months (default)
curl -X GET "http://localhost:8080/api/v1/dashboard/charts/income-expense" \
  -H "Authorization: Bearer <token>"

# Get last 24 months
curl -X GET "http://localhost:8080/api/v1/dashboard/charts/income-expense?months=24" \
  -H "Authorization: Bearer <token>"
```

---

## 6. Category Distribution Chart

### Get Expense Distribution by Category
Returns expense breakdown by category for a specific month.

**Endpoint:**
```
GET /api/v1/dashboard/charts/category-distribution
```

**Query Parameters:**
| Parameter | Type | Required | Format | Default | Description |
|-----------|------|----------|--------|---------|-------------|
| month | string | No | YYYY-MM | Current month | Month to analyze |

**Response:**
```json
{
  "success": true,
  "data": {
    "month": "2024-01",
    "distribution": [
      {
        "_id": "507f1f77bcf86cd799439012",
        "category_name": "Food & Dining",
        "total_amount": "1500.00",
        "transaction_count": 45
      },
      {
        "_id": "507f1f77bcf86cd799439013",
        "category_name": "Transportation",
        "total_amount": "800.00",
        "transaction_count": 12
      },
      {
        "_id": "507f1f77bcf86cd799439015",
        "category_name": "Entertainment",
        "total_amount": "500.00",
        "transaction_count": 8
      }
    ]
  }
}
```

**Example cURL:**
```bash
# Current month
curl -X GET "http://localhost:8080/api/v1/dashboard/charts/category-distribution" \
  -H "Authorization: Bearer <token>"

# Specific month
curl -X GET "http://localhost:8080/api/v1/dashboard/charts/category-distribution?month=2024-01" \
  -H "Authorization: Bearer <token>"
```

---

## 7. Pocket Distribution Chart

### Get Balance Distribution by Pocket Type
Returns current balance breakdown by pocket type.

**Endpoint:**
```
GET /api/v1/dashboard/charts/pocket-distribution
```

**Query Parameters:**
None

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "_id": "main",
      "pockets": [
        {
          "pocket_id": "507f1f77bcf86cd799439014",
          "pocket_name": "Main Wallet",
          "balance": "2500.00"
        }
      ],
      "total_by_type": "2500.00"
    },
    {
      "_id": "allocation",
      "pockets": [
        {
          "pocket_id": "507f1f77bcf86cd799439016",
          "pocket_name": "Savings Allocation",
          "balance": "1500.00"
        }
      ],
      "total_by_type": "1500.00"
    },
    {
      "_id": "saving",
      "pockets": [
        {
          "pocket_id": "507f1f77bcf86cd799439017",
          "pocket_name": "Emergency Fund",
          "balance": "5000.00"
        }
      ],
      "total_by_type": "5000.00"
    }
  ]
}
```

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/charts/pocket-distribution" \
  -H "Authorization: Bearer <token>"
```

---

## 8. Complete Dashboard Summary

### Get Complete Dashboard Summary
Returns all dashboard data (KPIs, charts, and optional date range reports) in one request.

**Endpoint:**
```
GET /api/v1/dashboard/summary
```

**Query Parameters:**
| Parameter | Type | Required | Format | Description |
|-----------|------|----------|--------|-------------|
| start_date | string | No | YYYY-MM-DD | Start date for daily reports |
| end_date | string | No | YYYY-MM-DD | End date for daily reports |

**Response:**
```json
{
  "success": true,
  "data": {
    "kpis": {
      "total_balance": "5000.50",
      "monthly_income": "3000.00",
      "monthly_expense": "1500.00",
      "free_money_total": "2500.00",
      "monthly_net_change": "1500.00"
    },
    "charts": {
      "income_expense_chart": [...],
      "pocket_distribution": [...],
      "category_distribution": [...]
    },
    "date_range_reports": [
      {
        "id": "507f1f77bcf86cd799439011",
        "report_date": "2024-01-31T00:00:00Z",
        ...
      }
    ]
  }
}
```

**Example cURL:**
```bash
# Get dashboard summary with optional date range
curl -X GET "http://localhost:8080/api/v1/dashboard/summary?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer <token>"
```

---

## Error Responses

### Unauthorized
```json
{
  "error": "unauthorized"
}
```
**Status:** 401

### Bad Request
```json
{
  "error": "start_date and end_date are required"
}
```
**Status:** 400

### Invalid Date Format
```json
{
  "error": "invalid start_date format (use YYYY-MM-DD)"
}
```
**Status:** 400

### Invalid Date Range
```json
{
  "error": "start_date must be before or equal to end_date"
}
```
**Status:** 400

### Server Error
```json
{
  "error": "internal server error message"
}
```
**Status:** 500

---

## Date Range Filter Examples

### Get reports for a specific month
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/daily?start_date=2024-01-01&end_date=2024-01-31"
```

### Get reports for a quarter
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/daily?start_date=2024-01-01&end_date=2024-03-31"
```

### Get reports for a year
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/daily?start_date=2024-01-01&end_date=2024-12-31"
```

### Get last 7 days
```bash
# Calculate dates in your client
const endDate = new Date();
const startDate = new Date(endDate.getTime() - 7 * 24 * 60 * 60 * 1000);

curl -X GET "http://localhost:8080/api/v1/dashboard/reports/daily?start_date=${startDate}&end_date=${endDate}"
```

### Get monthly summaries for a year
```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/monthly?start_month=2024-01&end_month=2024-12"
```

---

## Frontend Integration Examples

### React Hook for Dashboard Data
```typescript
import { useQuery } from '@tanstack/react-query';

export function useDashboardData(startDate?: string, endDate?: string) {
  const params = new URLSearchParams();
  if (startDate) params.append('start_date', startDate);
  if (endDate) params.append('end_date', endDate);

  return useQuery({
    queryKey: ['dashboard', startDate, endDate],
    queryFn: async () => {
      const response = await fetch(
        `/api/v1/dashboard/summary?${params}`,
        {
          headers: { 'Authorization': `Bearer ${token}` }
        }
      );
      return response.json();
    }
  });
}

// Usage
function Dashboard() {
  const { data, isLoading } = useDashboardData('2024-01-01', '2024-01-31');
  
  return (
    <div>
      <KPICards kpis={data?.data?.kpis} />
      <Charts charts={data?.data?.charts} />
      <DailyReports reports={data?.data?.date_range_reports} />
    </div>
  );
}
```

### Vue 3 Composable
```typescript
import { ref, computed } from 'vue';

export function useDashboard(startDate: string, endDate: string) {
  const data = ref(null);
  const loading = ref(false);
  const error = ref(null);

  const fetchDashboard = async () => {
    loading.value = true;
    try {
      const response = await fetch(
        `/api/v1/dashboard/summary?start_date=${startDate}&end_date=${endDate}`,
        {
          headers: { 'Authorization': `Bearer ${token}` }
        }
      );
      data.value = await response.json();
    } catch (e) {
      error.value = e;
    } finally {
      loading.value = false;
    }
  };

  return { data, loading, error, fetchDashboard };
}
```

---

## Performance Notes

- **KPI queries:** <50ms (cached for 30 seconds)
- **Chart queries:** <100ms (cached for 1 hour)
- **Daily reports:** <100ms per 30 days
- **Monthly summaries:** <50ms per 12 months

For best performance, cache responses on the client side and use the date range filters to limit data returned.

