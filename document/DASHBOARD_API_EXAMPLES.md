# Dashboard Analytics API - Request/Response Examples

## Authentication

All dashboard endpoints require JWT authentication.

```bash
# Login first to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

---

## Endpoint 1: Dashboard Summary

### Request

```bash
curl -X GET http://localhost:8080/api/v1/dashboard/summary \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Response (Success - 200 OK)

```json
{
  "status": "success",
  "message": "Dashboard summary retrieved successfully",
  "data": {
    "total_net_worth": 15750000.50,
    "monthly_income": 8500000.00,
    "monthly_expense": 3200000.00,
    "monthly_net": 5300000.00
  }
}
```

### Response Breakdown

| Field | Description | Calculation |
|-------|-------------|-------------|
| `total_net_worth` | Sum of all active pocket balances | Real-time from `pockets` collection |
| `monthly_income` | Total income this month | Historical (daily_summaries) + Live Delta (today) |
| `monthly_expense` | Total expense this month | Historical (daily_summaries) + Live Delta (today) |
| `monthly_net` | Net cash flow this month | monthly_income - monthly_expense |

### Error Responses

**401 Unauthorized** - Missing or invalid token
```json
{
  "status": "error",
  "message": "unauthorized",
  "data": null
}
```

**400 Bad Request** - Invalid user ID
```json
{
  "status": "error",
  "message": "invalid user id",
  "data": null
}
```

---

## Endpoint 2: Dashboard Charts

### Request - 7 Days

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/charts?range=7d" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Request - 30 Days

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/charts?range=30d" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Request - 90 Days

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/charts?range=90d" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Response (Success - 200 OK)

```json
{
  "status": "success",
  "message": "Dashboard charts retrieved successfully",
  "data": {
    "cash_flow_trend": [
      {
        "date": "2026-01-25",
        "income": 500000.00,
        "expense": 150000.00
      },
      {
        "date": "2026-01-26",
        "income": 0.00,
        "expense": 75000.00
      },
      {
        "date": "2026-01-27",
        "income": 8000000.00,
        "expense": 450000.00
      },
      {
        "date": "2026-01-28",
        "income": 0.00,
        "expense": 200000.00
      },
      {
        "date": "2026-01-29",
        "income": 0.00,
        "expense": 325000.00
      },
      {
        "date": "2026-01-30",
        "income": 0.00,
        "expense": 180000.00
      },
      {
        "date": "2026-01-31",
        "income": 0.00,
        "expense": 95000.00
      }
    ],
    "income_breakdown": [
      {
        "category_id": "507f1f77bcf86cd799439011",
        "category_name": "Salary",
        "amount": 8000000.00,
        "percentage": 94.12
      },
      {
        "category_id": "507f1f77bcf86cd799439012",
        "category_name": "Freelance",
        "amount": 500000.00,
        "percentage": 5.88
      }
    ],
    "expense_breakdown": [
      {
        "category_id": "507f1f77bcf86cd799439013",
        "category_name": "Food & Dining",
        "amount": 650000.00,
        "percentage": 44.83
      },
      {
        "category_id": "507f1f77bcf86cd799439014",
        "category_name": "Transportation",
        "amount": 400000.00,
        "percentage": 27.59
      },
      {
        "category_id": "507f1f77bcf86cd799439015",
        "category_name": "Shopping",
        "amount": 275000.00,
        "percentage": 18.97
      },
      {
        "category_id": "",
        "category_name": "Uncategorized",
        "amount": 125000.00,
        "percentage": 8.62
      }
    ]
  }
}
```

### Response Breakdown

#### `cash_flow_trend`
Array of daily income/expense data points for the specified range.

| Field | Type | Description |
|-------|------|-------------|
| `date` | string | Date in YYYY-MM-DD format |
| `income` | number | Total income for that day |
| `expense` | number | Total expense for that day |

#### `income_breakdown`
Array of income categories with amounts and percentages.

| Field | Type | Description |
|-------|------|-------------|
| `category_id` | string | Category ID (empty for uncategorized) |
| `category_name` | string | Category name |
| `amount` | number | Total amount for category |
| `percentage` | number | Percentage of total income |

#### `expense_breakdown`
Array of expense categories with amounts and percentages.

| Field | Type | Description |
|-------|------|-------------|
| `category_id` | string | Category ID (empty for uncategorized) |
| `category_name` | string | Category name |
| `amount` | number | Total amount for category |
| `percentage` | number | Percentage of total expense |

### Error Responses

**400 Bad Request** - Invalid range parameter
```json
{
  "status": "error",
  "message": "invalid range parameter. Valid values: 7d, 30d, 90d",
  "data": null
}
```

**401 Unauthorized** - Missing or invalid token
```json
{
  "status": "error",
  "message": "unauthorized",
  "data": null
}
```

---

## Frontend Integration Examples

### React/Next.js Example

```typescript
// api/dashboard.ts
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export const dashboardAPI = {
  getSummary: async (token: string) => {
    const response = await axios.get(`${API_BASE_URL}/dashboard/summary`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    return response.data.data;
  },

  getCharts: async (token: string, range: '7d' | '30d' | '90d' = '7d') => {
    const response = await axios.get(`${API_BASE_URL}/dashboard/charts`, {
      params: { range },
      headers: { Authorization: `Bearer ${token}` }
    });
    return response.data.data;
  }
};

// components/Dashboard.tsx
import { useEffect, useState } from 'react';
import { dashboardAPI } from '@/api/dashboard';

export default function Dashboard() {
  const [summary, setSummary] = useState(null);
  const [charts, setCharts] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      const token = localStorage.getItem('token');
      
      const [summaryData, chartsData] = await Promise.all([
        dashboardAPI.getSummary(token),
        dashboardAPI.getCharts(token, '30d')
      ]);
      
      setSummary(summaryData);
      setCharts(chartsData);
      setLoading(false);
    };

    fetchData();
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Dashboard</h1>
      
      {/* Summary Cards */}
      <div className="grid grid-cols-4 gap-4">
        <Card title="Net Worth" value={summary.total_net_worth} />
        <Card title="Monthly Income" value={summary.monthly_income} />
        <Card title="Monthly Expense" value={summary.monthly_expense} />
        <Card title="Monthly Net" value={summary.monthly_net} />
      </div>

      {/* Cash Flow Chart */}
      <LineChart data={charts.cash_flow_trend} />

      {/* Category Breakdowns */}
      <div className="grid grid-cols-2 gap-4">
        <PieChart 
          title="Income Breakdown" 
          data={charts.income_breakdown} 
        />
        <PieChart 
          title="Expense Breakdown" 
          data={charts.expense_breakdown} 
        />
      </div>
    </div>
  );
}
```

### Vue.js Example

```vue
<template>
  <div class="dashboard">
    <h1>Dashboard</h1>
    
    <!-- Summary Cards -->
    <div class="summary-grid">
      <SummaryCard 
        title="Net Worth" 
        :value="summary.total_net_worth" 
      />
      <SummaryCard 
        title="Monthly Income" 
        :value="summary.monthly_income" 
      />
      <SummaryCard 
        title="Monthly Expense" 
        :value="summary.monthly_expense" 
      />
      <SummaryCard 
        title="Monthly Net" 
        :value="summary.monthly_net" 
      />
    </div>

    <!-- Charts -->
    <CashFlowChart :data="charts.cash_flow_trend" />
    <CategoryCharts 
      :income="charts.income_breakdown" 
      :expense="charts.expense_breakdown" 
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const summary = ref({});
const charts = ref({});

onMounted(async () => {
  const token = localStorage.getItem('token');
  const headers = { Authorization: `Bearer ${token}` };

  const [summaryRes, chartsRes] = await Promise.all([
    axios.get('/api/v1/dashboard/summary', { headers }),
    axios.get('/api/v1/dashboard/charts?range=30d', { headers })
  ]);

  summary.value = summaryRes.data.data;
  charts.value = chartsRes.data.data;
});
</script>
```

---

## Testing with Postman

### 1. Create Environment Variables

```
base_url: http://localhost:8080
token: (will be set after login)
```

### 2. Login Request

**POST** `{{base_url}}/api/v1/auth/login`

**Body (JSON):**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Tests Script:**
```javascript
pm.test("Login successful", function () {
    pm.response.to.have.status(200);
    var jsonData = pm.response.json();
    pm.environment.set("token", jsonData.data.token);
});
```

### 3. Dashboard Summary Request

**GET** `{{base_url}}/api/v1/dashboard/summary`

**Headers:**
```
Authorization: Bearer {{token}}
```

**Tests Script:**
```javascript
pm.test("Summary retrieved", function () {
    pm.response.to.have.status(200);
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('total_net_worth');
    pm.expect(jsonData.data).to.have.property('monthly_income');
});
```

### 4. Dashboard Charts Request

**GET** `{{base_url}}/api/v1/dashboard/charts?range=30d`

**Headers:**
```
Authorization: Bearer {{token}}
```

**Tests Script:**
```javascript
pm.test("Charts retrieved", function () {
    pm.response.to.have.status(200);
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('cash_flow_trend');
    pm.expect(jsonData.data.cash_flow_trend).to.be.an('array');
});
```

---

## Performance Benchmarks

### Expected Response Times

| Endpoint | Range | Expected Time | Max Acceptable |
|----------|-------|---------------|----------------|
| `/summary` | N/A | < 50ms | 100ms |
| `/charts` | 7d | < 50ms | 100ms |
| `/charts` | 30d | < 75ms | 150ms |
| `/charts` | 90d | < 100ms | 200ms |

### Load Testing with Apache Bench

```bash
# Test summary endpoint (100 requests, 10 concurrent)
ab -n 100 -c 10 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/dashboard/summary

# Test charts endpoint
ab -n 100 -c 10 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/api/v1/dashboard/charts?range=30d"
```

---

## Common Use Cases

### 1. Dashboard Page Load
```javascript
// Fetch both endpoints in parallel
const [summary, charts] = await Promise.all([
  fetch('/api/v1/dashboard/summary'),
  fetch('/api/v1/dashboard/charts?range=7d')
]);
```

### 2. Range Selector
```javascript
const handleRangeChange = async (range) => {
  setLoading(true);
  const charts = await fetch(`/api/v1/dashboard/charts?range=${range}`);
  setCharts(charts.data);
  setLoading(false);
};
```

### 3. Auto-Refresh
```javascript
// Refresh every 5 minutes
useEffect(() => {
  const interval = setInterval(() => {
    fetchDashboardData();
  }, 5 * 60 * 1000);
  
  return () => clearInterval(interval);
}, []);
```

---

## Troubleshooting

### Issue: "unauthorized" error
**Solution:** Check if token is included in Authorization header

### Issue: Empty data arrays
**Solution:** Ensure user has transactions in the specified date range

### Issue: Slow response times
**Solution:** Verify indexes are created in MongoDB

### Issue: Incorrect totals
**Solution:** Check if cron job ran successfully for historical data
