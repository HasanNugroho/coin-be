# Performance-Optimized Reporting & Analytics Layer

## Executive Summary

This document defines a production-ready, performance-first reporting and analytics architecture for the Coin Finance System. All designs prioritize O(1) reads through precomputed data, snapshot collections, and optimized aggregation pipelines.

**Core Principle**: NO aggregation on raw transactions at read time. All reads consume precomputed snapshots or lightweight aggregations.

---

## 1. DAILY FINANCIAL REPORT LAYER

### 1.1 Collection Schema: `daily_financial_reports`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,                    // Index: compound
  report_date: Date,                    // Index: compound (user_id, report_date)
  
  // Opening state
  opening_balance: Decimal128,          // Sum of all pocket balances at 00:00
  
  // Closing state
  closing_balance: Decimal128,          // Sum of all pocket balances at 23:59:59
  
  // Daily aggregates
  total_income: Decimal128,             // Sum of all income transactions
  total_expense: Decimal128,            // Sum of all expense transactions
  total_transfer_in: Decimal128,        // Sum of transfer-in transactions
  total_transfer_out: Decimal128,       // Sum of transfer-out transactions
  
  // Breakdown by category
  expense_by_category: [
    {
      category_id: ObjectId,
      category_name: String,
      amount: Decimal128,
      transaction_count: Int32
    }
  ],
  
  // Breakdown by pocket
  transactions_grouped_by_pocket: [
    {
      pocket_id: ObjectId,
      pocket_name: String,
      pocket_type: String,              // main, allocation, saving, debt, system
      pocket_balance: Decimal128,       // Balance at end of day
      income: Decimal128,
      expense: Decimal128,
      transfer_in: Decimal128,
      transfer_out: Decimal128,
      transaction_count: Int32
    }
  ],
  
  // Metadata
  generated_at: Date,                   // Timestamp when report was generated
  is_final: Boolean,                    // true = locked, false = preliminary
  created_at: Date,
  updated_at: Date
}
```

### 1.2 Index Strategy

```javascript
// Primary lookup index
db.daily_financial_reports.createIndex({ user_id: 1, report_date: -1 })

// Range queries for date ranges
db.daily_financial_reports.createIndex({ user_id: 1, report_date: 1 })

// Unique constraint (one report per user per day)
db.daily_financial_reports.createIndex(
  { user_id: 1, report_date: 1 },
  { unique: true }
)

// Support for "latest report" queries
db.daily_financial_reports.createIndex({ user_id: 1, is_final: 1, report_date: -1 })
```

### 1.3 Generation Strategy

**Trigger**: Scheduled cron job at 23:59:59 UTC daily + on-demand via API

**Generation Process**:
1. Query all transactions for user on that date (filtered by `date` field)
2. Fetch all pocket balances at end of day
3. Calculate aggregates in memory
4. Upsert into `daily_financial_reports`
5. Mark `is_final: true` after 24 hours

**Pseudocode**:
```go
func GenerateDailyReport(userID ObjectID, reportDate time.Time) {
  startOfDay := reportDate.StartOfDay()
  endOfDay := reportDate.EndOfDay()
  
  // Fetch transactions for the day
  txns := db.transactions.find({
    user_id: userID,
    date: { $gte: startOfDay, $lt: endOfDay },
    deleted_at: null
  })
  
  // Calculate aggregates
  report := calculateAggregates(txns, userID, reportDate)
  
  // Upsert report
  db.daily_financial_reports.updateOne(
    { user_id: userID, report_date: startOfDay },
    { $set: report },
    { upsert: true }
  )
}
```

**Concurrency Safety**:
- Unique index prevents duplicate reports
- `is_final` flag prevents overwrites after 24h
- Idempotent: safe to re-run multiple times

### 1.4 API Endpoint: Get Daily Report

**Endpoint**: `GET /api/v1/reports/daily/:date`

**Query Parameters**:
- `date` (required): ISO 8601 date string (YYYY-MM-DD)
- `include_details` (optional): boolean, default=true. If false, returns summary only.

**Request**:
```
GET /api/v1/reports/daily/2024-01-15?include_details=true
Authorization: Bearer <token>
```

**Response Schema (OpenAPI)**:
```yaml
DailyFinancialReport:
  type: object
  required:
    - report_date
    - opening_balance
    - closing_balance
    - total_income
    - total_expense
    - total_transfer_in
    - total_transfer_out
  properties:
    report_date:
      type: string
      format: date
      example: "2024-01-15"
    opening_balance:
      type: number
      format: decimal
      example: 5000.00
    closing_balance:
      type: number
      format: decimal
      example: 5250.50
    total_income:
      type: number
      format: decimal
      example: 1500.00
    total_expense:
      type: number
      format: decimal
      example: 1249.50
    total_transfer_in:
      type: number
      format: decimal
      example: 0.00
    total_transfer_out:
      type: number
      format: decimal
      example: 0.00
    expense_by_category:
      type: array
      items:
        type: object
        properties:
          category_id:
            type: string
            format: ObjectId
          category_name:
            type: string
            example: "Food & Dining"
          amount:
            type: number
            format: decimal
          transaction_count:
            type: integer
    transactions_grouped_by_pocket:
      type: array
      items:
        type: object
        properties:
          pocket_id:
            type: string
            format: ObjectId
          pocket_name:
            type: string
          pocket_type:
            type: string
            enum: [main, allocation, saving, debt, system]
          pocket_balance:
            type: number
            format: decimal
          income:
            type: number
            format: decimal
          expense:
            type: number
            format: decimal
          transfer_in:
            type: number
            format: decimal
          transfer_out:
            type: number
            format: decimal
          transaction_count:
            type: integer
    generated_at:
      type: string
      format: date-time
    is_final:
      type: boolean
```

**Response Example**:
```json
{
  "report_date": "2024-01-15",
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

**HTTP Status Codes**:
- `200 OK`: Report found and returned
- `202 Accepted`: Report not yet generated, queued for generation
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: Report date is in the future

---

## 2. DASHBOARD DATA LAYER

### 2.1 Dashboard Components & Data Sources

#### KPI Cards (Real-time + Precomputed)

| KPI | Data Source | Query Cost | Refresh Rate |
|-----|-------------|-----------|--------------|
| Total Balance (all pockets) | Pocket collection | O(1) | Real-time |
| Total Income (current month) | Monthly snapshot | O(1) | Daily |
| Total Expense (current month) | Monthly snapshot | O(1) | Daily |
| Free Money Total | Pocket collection | O(1) | Real-time |

#### Charts (Precomputed)

| Chart | Data Source | Query Cost | Refresh Rate |
|-------|-------------|-----------|--------------|
| Monthly income vs expense (12mo) | Monthly summaries | O(1) | Daily |
| Balance distribution per pocket | Pocket collection | O(1) | Real-time |
| Expense distribution per category | Monthly summaries | O(1) | Daily |

### 2.2 Optimized MongoDB Aggregation Pipelines

#### Pipeline 1: KPI - Total Balance (All Pockets)

**Collection**: `pockets`

```javascript
db.pockets.aggregate([
  {
    $match: {
      user_id: ObjectId("..."),
      is_active: true,
      deleted_at: null
    }
  },
  {
    $group: {
      _id: null,
      total_balance: { $sum: "$balance" }
    }
  },
  {
    $project: {
      _id: 0,
      total_balance: 1
    }
  }
])
```

**Indexes**: `{ user_id: 1, is_active: 1, deleted_at: 1 }`

**Expected Result**: O(1) - single document aggregation

---

#### Pipeline 2: KPI - Income/Expense (Current Month)

**Collection**: `monthly_financial_summaries` (precomputed)

```javascript
db.monthly_financial_summaries.findOne({
  user_id: ObjectId("..."),
  year_month: "2024-01"
})
```

**Returns**:
```javascript
{
  total_income: 5000.00,
  total_expense: 3200.00
}
```

**Indexes**: `{ user_id: 1, year_month: 1 }`

**Expected Result**: O(1) - direct lookup

---

#### Pipeline 3: Chart - Monthly Income vs Expense (Last 12 Months)

**Collection**: `monthly_financial_summaries` (precomputed)

```javascript
db.monthly_financial_summaries.find({
  user_id: ObjectId("..."),
  year_month: {
    $gte: "2023-01",
    $lte: "2024-01"
  }
}).sort({ year_month: 1 }).limit(12)
```

**Returns**:
```javascript
[
  {
    year_month: "2023-01",
    total_income: 4500.00,
    total_expense: 3000.00
  },
  // ... 11 more months
]
```

**Indexes**: `{ user_id: 1, year_month: 1 }`

**Expected Result**: O(1) - 12 document lookup

---

#### Pipeline 4: Chart - Balance Distribution Per Pocket

**Collection**: `pockets`

```javascript
db.pockets.aggregate([
  {
    $match: {
      user_id: ObjectId("..."),
      is_active: true,
      deleted_at: null
    }
  },
  {
    $project: {
      _id: 1,
      name: 1,
      type: 1,
      balance: 1
    }
  },
  {
    $sort: { balance: -1 }
  }
])
```

**Indexes**: `{ user_id: 1, is_active: 1, deleted_at: 1 }`

**Expected Result**: O(n) where n = number of pockets (typically < 20)

---

#### Pipeline 5: Chart - Expense Distribution Per Category (Current Month)

**Collection**: `monthly_financial_summaries` (precomputed)

```javascript
db.monthly_financial_summaries.findOne({
  user_id: ObjectId("..."),
  year_month: "2024-01"
})
```

**Extract from document**:
```javascript
{
  expense_by_category: [
    {
      category_id: ObjectId("..."),
      category_name: "Food & Dining",
      amount: 450.00
    },
    // ... more categories
  ]
}
```

**Indexes**: `{ user_id: 1, year_month: 1 }`

**Expected Result**: O(1) - direct lookup, array extraction

---

### 2.3 Index Recommendations

```javascript
// Pockets collection
db.pockets.createIndex({ user_id: 1, is_active: 1, deleted_at: 1 })
db.pockets.createIndex({ user_id: 1, type: 1 })

// Monthly summaries collection
db.monthly_financial_summaries.createIndex({ user_id: 1, year_month: 1 })
db.monthly_financial_summaries.createIndex({ user_id: 1, year_month: -1 })
```

### 2.4 Data Separation: Precomputed vs Real-time

**Precomputed (from snapshots)**:
- Monthly income/expense totals
- Monthly expense by category
- 12-month historical trends

**Real-time (from live collections)**:
- Total balance across pockets
- Pocket-level balance distribution
- Pocket type filtering (main, allocation)

### 2.5 API Endpoints

#### Endpoint 1: GET /api/v1/dashboard/kpis

**Query Parameters**:
- `month` (optional): YYYY-MM format, defaults to current month

**Response Schema**:
```yaml
DashboardKPIs:
  type: object
  properties:
    total_balance:
      type: number
      format: decimal
      description: Sum of all active pocket balances
    total_income_current_month:
      type: number
      format: decimal
    total_expense_current_month:
      type: number
      format: decimal
    free_money_total:
      type: number
      format: decimal
      description: Sum of main + allocation pocket balances
    net_change_current_month:
      type: number
      format: decimal
      description: income - expense
```

**Response Example**:
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

#### Endpoint 2: GET /api/v1/dashboard/charts/monthly-trend

**Query Parameters**:
- `months` (optional): number of months to return, default=12, max=36

**Response Schema**:
```yaml
MonthlyTrendChart:
  type: object
  properties:
    data:
      type: array
      items:
        type: object
        properties:
          month:
            type: string
            format: YYYY-MM
          income:
            type: number
            format: decimal
          expense:
            type: number
            format: decimal
          net:
            type: number
            format: decimal
```

**Response Example**:
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
      "month": "2024-01",
      "income": 5000.00,
      "expense": 3200.00,
      "net": 1800.00
    }
  ]
}
```

---

#### Endpoint 3: GET /api/v1/dashboard/charts/pocket-distribution

**Response Schema**:
```yaml
PocketDistributionChart:
  type: object
  properties:
    data:
      type: array
      items:
        type: object
        properties:
          pocket_id:
            type: string
            format: ObjectId
          pocket_name:
            type: string
          pocket_type:
            type: string
            enum: [main, allocation, saving, debt, system]
          balance:
            type: number
            format: decimal
          percentage:
            type: number
            format: float
            description: Percentage of total balance
```

**Response Example**:
```json
{
  "data": [
    {
      "pocket_id": "507f1f77bcf86cd799439012",
      "pocket_name": "Main Wallet",
      "pocket_type": "main",
      "balance": 8500.00,
      "percentage": 54.0
    }
  ]
}
```

---

#### Endpoint 4: GET /api/v1/dashboard/charts/expense-by-category

**Query Parameters**:
- `month` (optional): YYYY-MM format, defaults to current month

**Response Schema**:
```yaml
ExpenseByCategoryChart:
  type: object
  properties:
    data:
      type: array
      items:
        type: object
        properties:
          category_id:
            type: string
            format: ObjectId
          category_name:
            type: string
          amount:
            type: number
            format: decimal
          percentage:
            type: number
            format: float
          transaction_count:
            type: integer
```

**Response Example**:
```json
{
  "data": [
    {
      "category_id": "507f1f77bcf86cd799439011",
      "category_name": "Food & Dining",
      "amount": 450.00,
      "percentage": 14.06,
      "transaction_count": 15
    }
  ]
}
```

---

## 3. SNAPSHOT / CUT-OFF COLLECTIONS

### 3.1 Collection: `daily_financial_snapshots`

**Purpose**: Point-in-time snapshot of all pocket balances at end of day

**Schema**:
```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  snapshot_date: Date,                  // End of day timestamp
  
  pocket_balances: [
    {
      pocket_id: ObjectId,
      pocket_name: String,
      pocket_type: String,
      balance: Decimal128,
      currency: String                  // e.g., "IDR"
    }
  ],
  
  total_balance: Decimal128,
  free_money_total: Decimal128,         // main + allocation pockets
  
  created_at: Date,
  updated_at: Date
}
```

**Indexes**:
```javascript
db.daily_financial_snapshots.createIndex({ user_id: 1, snapshot_date: -1 })
db.daily_financial_snapshots.createIndex(
  { user_id: 1, snapshot_date: 1 },
  { unique: true }
)
```

**Generation Strategy**:
- Triggered at 23:59:59 UTC daily
- Fetches all active pockets for user
- Stores snapshot of balances
- Idempotent: safe to re-run

**Consistency vs Performance**:
- **Consistency**: Snapshots are taken at fixed time (23:59:59)
- **Performance**: O(1) lookup for historical balance queries
- **Trade-off**: 24-hour delay in snapshot availability (acceptable for historical data)

---

### 3.2 Collection: `monthly_financial_summaries`

**Purpose**: Precomputed monthly aggregates to avoid expensive month-range queries

**Schema**:
```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  year_month: String,                   // Format: "2024-01"
  
  // Monthly totals
  total_income: Decimal128,
  total_expense: Decimal128,
  total_transfer_in: Decimal128,
  total_transfer_out: Decimal128,
  
  // Breakdown by category
  expense_by_category: [
    {
      category_id: ObjectId,
      category_name: String,
      amount: Decimal128,
      transaction_count: Int32
    }
  ],
  
  income_by_category: [
    {
      category_id: ObjectId,
      category_name: String,
      amount: Decimal128,
      transaction_count: Int32
    }
  ],
  
  // Breakdown by pocket
  pocket_summary: [
    {
      pocket_id: ObjectId,
      pocket_name: String,
      pocket_type: String,
      income: Decimal128,
      expense: Decimal128,
      transfer_in: Decimal128,
      transfer_out: Decimal128,
      transaction_count: Int32
    }
  ],
  
  // Metadata
  transaction_count: Int32,
  generated_at: Date,
  is_final: Boolean,                    // true after month ends
  created_at: Date,
  updated_at: Date
}
```

**Indexes**:
```javascript
db.monthly_financial_summaries.createIndex({ user_id: 1, year_month: -1 })
db.monthly_financial_summaries.createIndex(
  { user_id: 1, year_month: 1 },
  { unique: true }
)
db.monthly_financial_summaries.createIndex({ user_id: 1, is_final: 1 })
```

**Generation Strategy**:
- Triggered at 00:00:00 UTC on the 1st of each month
- Aggregates all transactions from previous month
- Calculates category breakdowns
- Marks `is_final: true` after month ends
- Can be regenerated if transactions are corrected

**Consistency vs Performance**:
- **Consistency**: Finalized after month-end (no retroactive changes)
- **Performance**: O(1) lookup for monthly data
- **Trade-off**: 1-day delay in monthly summary availability

---

### 3.3 Collection: `pocket_balance_snapshots`

**Purpose**: Track balance history per pocket for trend analysis

**Schema**:
```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  pocket_id: ObjectId,
  snapshot_date: Date,                  // End of day
  
  balance: Decimal128,
  balance_change: Decimal128,           // Change from previous day
  
  daily_income: Decimal128,
  daily_expense: Decimal128,
  daily_transfer_in: Decimal128,
  daily_transfer_out: Decimal128,
  
  created_at: Date
}
```

**Indexes**:
```javascript
db.pocket_balance_snapshots.createIndex({ user_id: 1, pocket_id: 1, snapshot_date: -1 })
db.pocket_balance_snapshots.createIndex({ pocket_id: 1, snapshot_date: -1 })
```

**Generation Strategy**:
- Triggered at 23:59:59 UTC daily
- Fetches pocket balance and daily transaction totals
- Stores snapshot
- Enables trend analysis without querying raw transactions

**Consistency vs Performance**:
- **Consistency**: Snapshots taken at fixed time
- **Performance**: O(1) lookup for pocket balance history
- **Trade-off**: 24-hour delay, but acceptable for trend analysis

---

## 4. AI-READY DATA LAYER

### 4.1 AI-Friendly Aggregated Collections

#### Collection: `ai_financial_context`

**Purpose**: Pre-aggregated, denormalized data optimized for LLM consumption

**Schema**:
```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  context_date: Date,                   // Latest snapshot date
  
  // Current state
  current_balance: Decimal128,
  free_money: Decimal128,
  
  // Recent performance (last 30 days)
  last_30_days: {
    total_income: Decimal128,
    total_expense: Decimal128,
    net_change: Decimal128,
    average_daily_expense: Decimal128,
    transaction_count: Int32,
    top_expense_categories: [
      {
        category_name: String,
        amount: Decimal128,
        percentage: Float
      }
    ]
  },
  
  // Year-to-date
  year_to_date: {
    total_income: Decimal128,
    total_expense: Decimal128,
    net_change: Decimal128,
    average_monthly_expense: Decimal128
  },
  
  // Pocket breakdown
  pockets: [
    {
      pocket_id: ObjectId,
      pocket_name: String,
      pocket_type: String,
      balance: Decimal128,
      percentage_of_total: Float,
      monthly_trend: [
        {
          month: String,
          balance: Decimal128
        }
      ]
    }
  ],
  
  // Spending patterns
  spending_patterns: {
    highest_expense_day_of_week: String,
    highest_expense_category: String,
    average_transaction_amount: Decimal128,
    largest_transaction: Decimal128,
    smallest_transaction: Decimal128
  },
  
  // Alerts & insights
  alerts: [
    {
      type: String,                     // e.g., "high_spending", "low_balance"
      message: String,
      severity: String                  // "info", "warning", "critical"
    }
  ],
  
  updated_at: Date
}
```

**Indexes**:
```javascript
db.ai_financial_context.createIndex({ user_id: 1 })
db.ai_financial_context.createIndex({ user_id: 1, updated_at: -1 })
```

**Generation Strategy**:
- Triggered daily at 00:30 UTC (after all snapshots are generated)
- Aggregates data from snapshots, summaries, and daily reports
- Denormalizes for easy LLM consumption
- Single document per user for O(1) retrieval

---

### 4.2 Example LLM-Optimized Documents

**Example 1: Current Financial Status**
```json
{
  "user_id": "507f1f77bcf86cd799439001",
  "context_date": "2024-01-15T23:59:59Z",
  "current_balance": 15750.50,
  "free_money": 8500.00,
  "last_30_days": {
    "total_income": 15000.00,
    "total_expense": 9600.00,
    "net_change": 5400.00,
    "average_daily_expense": 320.00,
    "transaction_count": 85,
    "top_expense_categories": [
      {
        "category_name": "Food & Dining",
        "amount": 2400.00,
        "percentage": 25.0
      },
      {
        "category_name": "Transportation",
        "amount": 1800.00,
        "percentage": 18.75
      }
    ]
  },
  "spending_patterns": {
    "highest_expense_day_of_week": "Friday",
    "highest_expense_category": "Food & Dining",
    "average_transaction_amount": 112.94,
    "largest_transaction": 850.00,
    "smallest_transaction": 5.50
  }
}
```

---

### 4.3 Mapping: User Intent → Data Source

| User Intent | Query | Data Source | Response Time |
|-------------|-------|-------------|----------------|
| "What's my current balance?" | Get total balance | `ai_financial_context.current_balance` | O(1) |
| "How much did I spend last month?" | Get monthly expense | `monthly_financial_summaries` | O(1) |
| "What are my top spending categories?" | Get category breakdown | `ai_financial_context.last_30_days.top_expense_categories` | O(1) |
| "Show me my balance trend" | Get 12-month trend | `monthly_financial_summaries` (12 docs) | O(1) |
| "Am I on budget?" | Compare expense vs threshold | `monthly_financial_summaries` + rules | O(1) |
| "What's my average daily spending?" | Calculate average | `ai_financial_context.last_30_days.average_daily_expense` | O(1) |
| "Which pocket has the most money?" | Get pocket breakdown | `ai_financial_context.pockets` | O(1) |
| "When do I usually spend the most?" | Get spending patterns | `ai_financial_context.spending_patterns` | O(1) |

**Key Principle**: All AI queries are O(1) lookups. NO aggregation on raw transactions.

---

## 5. IMPLEMENTATION CHECKLIST

- [ ] Create `daily_financial_reports` collection with indexes
- [ ] Create `daily_financial_snapshots` collection with indexes
- [ ] Create `monthly_financial_summaries` collection with indexes
- [ ] Create `pocket_balance_snapshots` collection with indexes
- [ ] Create `ai_financial_context` collection with indexes
- [ ] Implement snapshot generation service
- [ ] Implement daily report generation service
- [ ] Implement monthly summary generation service
- [ ] Implement AI context generation service
- [ ] Implement reporting API endpoints
- [ ] Implement dashboard API endpoints
- [ ] Implement health check endpoints
- [ ] Add service registration to main.go
- [ ] Create Swagger/OpenAPI documentation
- [ ] Add integration tests
- [ ] Set up cron jobs for scheduled generation

---

## 6. PERFORMANCE GUARANTEES

| Operation | Query Cost | Latency | Notes |
|-----------|-----------|---------|-------|
| Get daily report | O(1) | <50ms | Direct lookup |
| Get dashboard KPIs | O(1) | <50ms | Snapshot lookup |
| Get 12-month trend | O(1) | <50ms | 12 document lookup |
| Get pocket distribution | O(n) | <100ms | n = # pockets (typically <20) |
| Get AI context | O(1) | <50ms | Single document lookup |

**No query will scan raw transactions at read time.**

---

## 7. MONITORING & ALERTS

- Monitor snapshot generation latency (target: <5s)
- Monitor monthly summary generation latency (target: <30s)
- Alert if report generation fails
- Alert if snapshot is missing for a day
- Track collection sizes and growth rates

