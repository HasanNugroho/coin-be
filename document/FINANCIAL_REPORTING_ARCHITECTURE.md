# Financial Reporting Architecture
## Production-Ready MongoDB Design for High-Scale Finance System

---

## 1. DAILY FINANCIAL REPORT LAYER

### 1.1 Collection Schema: `daily_financial_reports`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,                    // Indexed
  report_date: Date,                    // Indexed (start of day UTC)
  
  // Opening/Closing Balances
  opening_balance: Decimal128,          // Sum of all pocket balances at 00:00 UTC
  closing_balance: Decimal128,          // Sum of all pocket balances at 23:59:59 UTC
  
  // Daily Aggregates
  total_income: Decimal128,             // Sum of all income transactions
  total_expense: Decimal128,            // Sum of all expense transactions
  total_transfer_in: Decimal128,        // Sum of transfer-in transactions
  total_transfer_out: Decimal128,       // Sum of transfer-out transactions
  
  // Expense Breakdown by Category
  expense_by_category: [
    {
      category_id: ObjectId,
      category_name: String,
      amount: Decimal128,
      transaction_count: Int32
    }
  ],
  
  // Transactions Grouped by Pocket
  transactions_by_pocket: [
    {
      pocket_id: ObjectId,
      pocket_name: String,
      pocket_type: String,              // "main", "allocation", "saving", "debt", "system"
      
      income: Decimal128,
      expense: Decimal128,
      transfer_in: Decimal128,
      transfer_out: Decimal128,
      
      opening_balance: Decimal128,
      closing_balance: Decimal128,
      
      transaction_count: Int32,
      transactions: [
        {
          transaction_id: ObjectId,
          type: String,
          amount: Decimal128,
          category_id: ObjectId,
          category_name: String,
          note: String,
          timestamp: Date
        }
      ]
    }
  ],
  
  // Metadata
  generated_at: Date,                   // When this report was generated
  is_final: Boolean,                    // true = locked, false = draft/in-progress
  
  created_at: Date,
  updated_at: Date
}
```

### 1.2 Indexing Strategy

```javascript
// Primary lookup index
db.daily_financial_reports.createIndex(
  { user_id: 1, report_date: -1 },
  { unique: true }
);

// Range queries for date filtering
db.daily_financial_reports.createIndex(
  { user_id: 1, report_date: 1 }
);

// Efficient sorting for dashboard
db.daily_financial_reports.createIndex(
  { user_id: 1, report_date: -1, is_final: 1 }
);

// TTL index for auto-deletion of draft reports after 30 days
db.daily_financial_reports.createIndex(
  { created_at: 1 },
  { 
    expireAfterSeconds: 2592000,
    partialFilterExpression: { is_final: false }
  }
);
```

### 1.3 Generation Strategy

**Trigger:** Daily cron job at 23:59:59 UTC (or configurable timezone offset)

**Method:** Batch aggregation pipeline

```javascript
// Pseudo-code for generation
async function generateDailyReport(userId, reportDate) {
  const startOfDay = new Date(reportDate);
  startOfDay.setUTCHours(0, 0, 0, 0);
  
  const endOfDay = new Date(reportDate);
  endOfDay.setUTCHours(23, 59, 59, 999);
  
  // 1. Get opening balance (snapshot from previous day or calculate)
  const openingBalance = await getOpeningBalance(userId, startOfDay);
  
  // 2. Aggregate transactions for the day
  const dailyData = await db.transactions.aggregate([
    {
      $match: {
        user_id: userId,
        date: { $gte: startOfDay, $lte: endOfDay },
        deleted_at: null
      }
    },
    {
      $facet: {
        byType: [
          {
            $group: {
              _id: "$type",
              amount: { $sum: "$amount" },
              count: { $sum: 1 }
            }
          }
        ],
        byCategory: [
          {
            $match: { type: "expense", category_id: { $ne: null } }
          },
          {
            $lookup: {
              from: "user_categories",
              localField: "category_id",
              foreignField: "_id",
              as: "category"
            }
          },
          {
            $group: {
              _id: "$category_id",
              category_name: { $first: "$category.name" },
              amount: { $sum: "$amount" },
              count: { $sum: 1 }
            }
          }
        ],
        byPocket: [
          {
            $group: {
              _id: "$pocket_from",
              transactions: { $push: "$$ROOT" }
            }
          }
        ]
      }
    }
  ]).toArray();
  
  // 3. Get closing balance (current pocket balances)
  const closingBalance = await getClosingBalance(userId);
  
  // 4. Construct and insert report
  const report = {
    user_id: userId,
    report_date: startOfDay,
    opening_balance: openingBalance,
    closing_balance: closingBalance,
    total_income: dailyData[0].byType.find(t => t._id === "income")?.amount || 0,
    total_expense: dailyData[0].byType.find(t => t._id === "expense")?.amount || 0,
    total_transfer_in: dailyData[0].byType.find(t => t._id === "transfer")?.amount || 0,
    total_transfer_out: 0, // Calculated separately
    expense_by_category: dailyData[0].byCategory,
    transactions_by_pocket: dailyData[0].byPocket,
    generated_at: new Date(),
    is_final: true,
    created_at: new Date(),
    updated_at: new Date()
  };
  
  await db.daily_financial_reports.updateOne(
    { user_id: userId, report_date: startOfDay },
    { $set: report },
    { upsert: true }
  );
}
```

### 1.4 API Response Format

```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "report_date": "2024-01-15T00:00:00Z",
    "opening_balance": 5000.50,
    "closing_balance": 4850.75,
    "summary": {
      "total_income": 2000.00,
      "total_expense": 1500.00,
      "total_transfer_in": 500.00,
      "total_transfer_out": 150.25,
      "net_change": -150.25
    },
    "expense_by_category": [
      {
        "category_id": "507f1f77bcf86cd799439012",
        "category_name": "Food & Dining",
        "amount": 450.00,
        "transaction_count": 8
      },
      {
        "category_id": "507f1f77bcf86cd799439013",
        "category_name": "Transportation",
        "amount": 200.00,
        "transaction_count": 3
      }
    ],
    "transactions_by_pocket": [
      {
        "pocket_id": "507f1f77bcf86cd799439014",
        "pocket_name": "Main Wallet",
        "pocket_type": "main",
        "opening_balance": 3000.00,
        "closing_balance": 2850.75,
        "income": 1000.00,
        "expense": 1200.00,
        "transfer_in": 200.00,
        "transfer_out": 150.25,
        "transaction_count": 12,
        "transactions": [
          {
            "id": "507f1f77bcf86cd799439015",
            "type": "expense",
            "amount": 50.00,
            "category_name": "Food & Dining",
            "note": "Lunch",
            "timestamp": "2024-01-15T12:30:00Z"
          }
        ]
      }
    ],
    "generated_at": "2024-01-15T23:59:59Z",
    "is_final": true
  }
}
```

---

## 2. DASHBOARD DATA (PERFORMANCE-ORIENTED)

### 2.1 KPI Cards - Aggregation Pipelines

#### 2.1.1 Total Balance (All Pockets)

**Real-time (no pre-aggregation needed)**

```javascript
db.pockets.aggregate([
  {
    $match: {
      user_id: ObjectId("userId"),
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
]).toArray();
```

**Index:**
```javascript
db.pockets.createIndex({ user_id: 1, is_active: 1, deleted_at: 1 });
```

#### 2.1.2 Total Income (Current Month)

**Pre-aggregated (use daily_financial_reports)**

```javascript
db.daily_financial_reports.aggregate([
  {
    $match: {
      user_id: ObjectId("userId"),
      report_date: {
        $gte: new Date("2024-01-01"),
        $lte: new Date("2024-01-31")
      },
      is_final: true
    }
  },
  {
    $group: {
      _id: null,
      total_income: { $sum: "$total_income" }
    }
  }
]).toArray();
```

**Index:**
```javascript
db.daily_financial_reports.createIndex({
  user_id: 1,
  report_date: 1,
  is_final: 1
});
```

#### 2.1.3 Total Expense (Current Month)

**Pre-aggregated (use daily_financial_reports)**

```javascript
db.daily_financial_reports.aggregate([
  {
    $match: {
      user_id: ObjectId("userId"),
      report_date: {
        $gte: new Date("2024-01-01"),
        $lte: new Date("2024-01-31")
      },
      is_final: true
    }
  },
  {
    $group: {
      _id: null,
      total_expense: { $sum: "$total_expense" }
    }
  }
]).toArray();
```

#### 2.1.4 Free Money Total (main + allocation pockets)

**Real-time**

```javascript
db.pockets.aggregate([
  {
    $match: {
      user_id: ObjectId("userId"),
      type: { $in: ["main", "allocation"] },
      is_active: true,
      deleted_at: null
    }
  },
  {
    $group: {
      _id: null,
      free_money: { $sum: "$balance" }
    }
  }
]).toArray();
```

**Index:**
```javascript
db.pockets.createIndex({
  user_id: 1,
  type: 1,
  is_active: 1,
  deleted_at: 1
});
```

### 2.2 Chart Aggregations

#### 2.2.1 Monthly Income vs Expense (Last 12 Months)

**Pre-aggregated (use monthly_financial_summaries)**

```javascript
db.monthly_financial_summaries.aggregate([
  {
    $match: {
      user_id: ObjectId("userId"),
      month: {
        $gte: new Date("2023-01-01"),
        $lte: new Date("2024-01-31")
      }
    }
  },
  {
    $sort: { month: 1 }
  },
  {
    $project: {
      _id: 0,
      month: {
        $dateToString: { format: "%Y-%m", date: "$month" }
      },
      income: 1,
      expense: 1
    }
  }
]).toArray();
```

#### 2.2.2 Balance Distribution Per Pocket

**Real-time (cached for 5 minutes)**

```javascript
db.pockets.aggregate([
  {
    $match: {
      user_id: ObjectId("userId"),
      is_active: true,
      deleted_at: null
    }
  },
  {
    $group: {
      _id: "$type",
      pockets: {
        $push: {
          pocket_id: "$_id",
          pocket_name: "$name",
          balance: "$balance"
        }
      },
      total_by_type: { $sum: "$balance" }
    }
  },
  {
    $sort: { total_by_type: -1 }
  }
]).toArray();
```

**Index:**
```javascript
db.pockets.createIndex({
  user_id: 1,
  is_active: 1,
  type: 1,
  deleted_at: 1
});
```

#### 2.2.3 Expense Distribution Per Category (Current Month)

**Pre-aggregated (use daily_financial_reports)**

```javascript
db.daily_financial_reports.aggregate([
  {
    $match: {
      user_id: ObjectId("userId"),
      report_date: {
        $gte: new Date("2024-01-01"),
        $lte: new Date("2024-01-31")
      },
      is_final: true
    }
  },
  {
    $unwind: "$expense_by_category"
  },
  {
    $group: {
      _id: "$expense_by_category.category_id",
      category_name: { $first: "$expense_by_category.category_name" },
      total_amount: { $sum: "$expense_by_category.amount" },
      transaction_count: { $sum: "$expense_by_category.transaction_count" }
    }
  },
  {
    $sort: { total_amount: -1 }
  }
]).toArray();
```

### 2.3 Indexing Summary for Dashboard

```javascript
// KPI Cards
db.pockets.createIndex({ user_id: 1, is_active: 1, deleted_at: 1 });
db.pockets.createIndex({ user_id: 1, type: 1, is_active: 1, deleted_at: 1 });

// Charts
db.daily_financial_reports.createIndex({
  user_id: 1,
  report_date: 1,
  is_final: 1
});
db.monthly_financial_summaries.createIndex({
  user_id: 1,
  month: 1
});
```

### 2.4 Caching Strategy

| Data | Source | Cache TTL | Reason |
|------|--------|-----------|--------|
| Total Balance | pockets (real-time) | 30s | Balances change frequently |
| Monthly Income/Expense | daily_financial_reports | 1h | Locked reports, no change |
| Pocket Distribution | pockets (real-time) | 5m | Changes with transactions |
| Category Distribution | daily_financial_reports | 1h | Locked reports, no change |

---

## 3. AI-READY DATA LAYER (CHATBOT PREPARATION)

### 3.1 Data Gaps for AI Usage

**Current gaps:**
- No enriched transaction context (merchant, category confidence, tags)
- No spending patterns/trends pre-calculated
- No budget vs actual comparison
- No anomaly flags
- No transaction descriptions for NLP

### 3.2 AI-Optimized Collections

#### 3.2.1 Collection: `ai_transaction_enrichment`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  transaction_id: ObjectId,              // Reference to original transaction
  
  // Original transaction data (denormalized for LLM)
  transaction_type: String,
  amount: Decimal128,
  date: Date,
  
  // Enrichment
  merchant_name: String,                 // Extracted/inferred from note
  merchant_category: String,             // Standardized category
  confidence_score: Float,               // 0-1, ML confidence
  
  // NLP-friendly description
  description: String,                   // "Lunch at Starbucks" vs "expense"
  tags: [String],                        // ["food", "coffee", "weekday"]
  
  // Anomaly detection
  is_anomaly: Boolean,
  anomaly_reason: String,                // "3x higher than usual", "unusual time", etc.
  anomaly_score: Float,                  // 0-1
  
  // Spending pattern context
  category_avg_amount: Decimal128,       // User's avg for this category
  category_avg_frequency: Int32,         // Transactions per month
  is_recurring: Boolean,
  
  // Budget context
  budget_category: String,
  budget_remaining: Decimal128,
  budget_utilization_percent: Float,
  
  // Temporal context
  day_of_week: String,
  is_weekend: Boolean,
  is_holiday: Boolean,
  
  created_at: Date,
  updated_at: Date
}
```

**Index:**
```javascript
db.ai_transaction_enrichment.createIndex({
  user_id: 1,
  date: -1
});

db.ai_transaction_enrichment.createIndex({
  user_id: 1,
  transaction_id: 1
}, { unique: true });

db.ai_transaction_enrichment.createIndex({
  user_id: 1,
  is_anomaly: 1,
  date: -1
});
```

#### 3.2.2 Collection: `ai_spending_patterns`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  
  // Pattern identification
  category: String,
  merchant: String,
  
  // Statistics
  avg_amount: Decimal128,
  median_amount: Decimal128,
  std_dev: Decimal128,
  min_amount: Decimal128,
  max_amount: Decimal128,
  
  // Frequency
  frequency_per_month: Float,
  frequency_per_week: Float,
  last_transaction_date: Date,
  
  // Trend
  trend: String,                         // "increasing", "decreasing", "stable"
  trend_percent_change: Float,           // Month-over-month
  
  // Seasonality
  is_seasonal: Boolean,
  seasonal_months: [Int32],              // [1, 12] for Jan, Dec
  
  // Behavioral insights
  preferred_day_of_week: String,
  preferred_time_of_day: String,
  
  // AI context
  spending_category_rank: Int32,         // Rank among user's categories
  is_essential: Boolean,                 // Inferred: food, utilities, etc.
  
  data_points: Int32,                    // Number of transactions in pattern
  confidence: Float,                     // 0-1
  
  created_at: Date,
  updated_at: Date
}
```

**Index:**
```javascript
db.ai_spending_patterns.createIndex({
  user_id: 1,
  category: 1
});

db.ai_spending_patterns.createIndex({
  user_id: 1,
  frequency_per_month: -1
});
```

#### 3.2.3 Collection: `ai_financial_insights`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  
  // Insight metadata
  insight_type: String,                  // "anomaly", "trend", "opportunity", "warning"
  category: String,                      // "spending", "income", "savings", "budget"
  severity: String,                      // "info", "warning", "critical"
  
  // Content for chatbot
  title: String,                         // "Unusual spending detected"
  description: String,                   // Human-readable explanation
  recommendation: String,                // Actionable advice
  
  // Data backing the insight
  metric_name: String,
  metric_value: Decimal128,
  metric_baseline: Decimal128,
  metric_change_percent: Float,
  
  // Context
  affected_categories: [String],
  affected_pockets: [ObjectId],
  date_range: {
    start: Date,
    end: Date
  },
  
  // Engagement
  is_actionable: Boolean,
  action_url: String,                    // Deep link to relevant UI
  
  // Lifecycle
  created_at: Date,
  expires_at: Date,                      // Auto-delete old insights
  is_read: Boolean,
  read_at: Date
}
```

**Index:**
```javascript
db.ai_financial_insights.createIndex({
  user_id: 1,
  created_at: -1
});

db.ai_financial_insights.createIndex({
  user_id: 1,
  is_read: 1,
  severity: 1
});

db.ai_financial_insights.createIndex(
  { expires_at: 1 },
  { expireAfterSeconds: 0 }
);
```

### 3.3 LLM-Optimized Document Format

**Example for chatbot context:**

```json
{
  "user_financial_context": {
    "period": "last_30_days",
    "summary": {
      "total_income": 5000.00,
      "total_expense": 3200.00,
      "net_savings": 1800.00,
      "savings_rate": 0.36
    },
    "top_spending_categories": [
      {
        "name": "Food & Dining",
        "amount": 850.00,
        "percent_of_total": 0.265,
        "trend": "increasing",
        "vs_average": "+15%"
      },
      {
        "name": "Transportation",
        "amount": 420.00,
        "percent_of_total": 0.131,
        "trend": "stable"
      }
    ],
    "anomalies": [
      {
        "date": "2024-01-15",
        "type": "unusual_spending",
        "description": "Spent $250 on electronics (3x your average)",
        "category": "Shopping",
        "severity": "info"
      }
    ],
    "budget_status": {
      "food": { "allocated": 1000, "spent": 850, "remaining": 150 },
      "transport": { "allocated": 500, "spent": 420, "remaining": 80 }
    },
    "insights": [
      {
        "type": "opportunity",
        "message": "Your coffee spending is up 40% this month. Consider reducing frequency.",
        "potential_savings": 60.00
      }
    ]
  }
}
```

### 3.4 Finance Intent Mapping

| User Intent | Data Source | Query Type |
|-------------|-------------|-----------|
| "How much did I spend on food?" | ai_transaction_enrichment + daily_financial_reports | Aggregation |
| "What are my spending trends?" | ai_spending_patterns | Lookup |
| "Am I on budget?" | daily_financial_reports + budget collection | Comparison |
| "Unusual transactions?" | ai_transaction_enrichment (is_anomaly=true) | Filter |
| "Where can I save money?" | ai_financial_insights (type=opportunity) | Lookup |
| "What's my income this month?" | daily_financial_reports | Aggregation |
| "Compare this month vs last?" | monthly_financial_summaries | Comparison |
| "Recurring expenses?" | ai_spending_patterns (frequency > threshold) | Filter |

---

## 4. SNAPSHOT / CUT-OFF COLLECTIONS

### 4.1 Collection: `daily_financial_snapshots`

**Purpose:** Reduce transaction aggregation load by pre-computing daily state

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  snapshot_date: Date,                  // Start of day UTC
  
  // Pocket balances at end of day
  pocket_balances: [
    {
      pocket_id: ObjectId,
      pocket_name: String,
      pocket_type: String,
      balance: Decimal128,
      balance_change: Decimal128        // vs previous day
    }
  ],
  
  // Daily totals
  total_balance: Decimal128,
  total_income: Decimal128,
  total_expense: Decimal128,
  total_transfer_in: Decimal128,
  total_transfer_out: Decimal128,
  
  // Cumulative (year-to-date)
  ytd_income: Decimal128,
  ytd_expense: Decimal128,
  ytd_net: Decimal128,
  
  // Metadata
  transaction_count: Int32,
  is_complete: Boolean,                 // false = partial, true = locked
  generated_at: Date,
  
  created_at: Date,
  updated_at: Date
}
```

**Indexing:**
```javascript
db.daily_financial_snapshots.createIndex(
  { user_id: 1, snapshot_date: -1 },
  { unique: true }
);

db.daily_financial_snapshots.createIndex({
  user_id: 1,
  snapshot_date: 1
});

db.daily_financial_snapshots.createIndex({
  user_id: 1,
  is_complete: 1,
  snapshot_date: -1
});
```

**Generation Method:**

```javascript
// Cron: Daily at 23:59:59 UTC
async function generateDailySnapshot(userId, snapshotDate) {
  const startOfDay = new Date(snapshotDate);
  startOfDay.setUTCHours(0, 0, 0, 0);
  
  const endOfDay = new Date(snapshotDate);
  endOfDay.setUTCHours(23, 59, 59, 999);
  
  // Get current pocket balances
  const pockets = await db.pockets.find({
    user_id: userId,
    deleted_at: null
  }).toArray();
  
  // Get previous day snapshot for comparison
  const prevSnapshot = await db.daily_financial_snapshots.findOne({
    user_id: userId,
    snapshot_date: new Date(snapshotDate.getTime() - 86400000)
  });
  
  // Aggregate daily transactions
  const dailyTotals = await db.transactions.aggregate([
    {
      $match: {
        user_id: userId,
        date: { $gte: startOfDay, $lte: endOfDay },
        deleted_at: null
      }
    },
    {
      $group: {
        _id: "$type",
        amount: { $sum: "$amount" }
      }
    }
  ]).toArray();
  
  const snapshot = {
    user_id: userId,
    snapshot_date: startOfDay,
    pocket_balances: pockets.map(p => ({
      pocket_id: p._id,
      pocket_name: p.name,
      pocket_type: p.type,
      balance: p.balance,
      balance_change: prevSnapshot 
        ? p.balance - (prevSnapshot.pocket_balances.find(pb => pb.pocket_id === p._id)?.balance || 0)
        : 0
    })),
    total_balance: pockets.reduce((sum, p) => sum + p.balance, 0),
    total_income: dailyTotals.find(t => t._id === "income")?.amount || 0,
    total_expense: dailyTotals.find(t => t._id === "expense")?.amount || 0,
    total_transfer_in: dailyTotals.find(t => t._id === "transfer")?.amount || 0,
    total_transfer_out: 0,
    ytd_income: (prevSnapshot?.ytd_income || 0) + (dailyTotals.find(t => t._id === "income")?.amount || 0),
    ytd_expense: (prevSnapshot?.ytd_expense || 0) + (dailyTotals.find(t => t._id === "expense")?.amount || 0),
    ytd_net: 0, // Calculated
    transaction_count: await db.transactions.countDocuments({
      user_id: userId,
      date: { $gte: startOfDay, $lte: endOfDay },
      deleted_at: null
    }),
    is_complete: true,
    generated_at: new Date(),
    created_at: new Date(),
    updated_at: new Date()
  };
  
  snapshot.ytd_net = snapshot.ytd_income - snapshot.ytd_expense;
  
  await db.daily_financial_snapshots.updateOne(
    { user_id: userId, snapshot_date: startOfDay },
    { $set: snapshot },
    { upsert: true }
  );
}
```

### 4.2 Collection: `monthly_financial_summaries`

**Purpose:** Pre-aggregated monthly data for charts and reports

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  month: Date,                          // First day of month
  
  // Monthly totals
  income: Decimal128,
  expense: Decimal128,
  transfer_in: Decimal128,
  transfer_out: Decimal128,
  net: Decimal128,                      // income - expense
  
  // Opening/Closing
  opening_balance: Decimal128,
  closing_balance: Decimal128,
  
  // By category
  expense_by_category: [
    {
      category_id: ObjectId,
      category_name: String,
      amount: Decimal128,
      percent_of_total: Float,
      transaction_count: Int32
    }
  ],
  
  // By pocket
  by_pocket: [
    {
      pocket_id: ObjectId,
      pocket_name: String,
      pocket_type: String,
      income: Decimal128,
      expense: Decimal128,
      net: Decimal128
    }
  ],
  
  // Cumulative
  ytd_income: Decimal128,
  ytd_expense: Decimal128,
  ytd_net: Decimal128,
  
  // Metadata
  day_count: Int32,                     // Days with transactions
  transaction_count: Int32,
  is_complete: Boolean,
  
  created_at: Date,
  updated_at: Date
}
```

**Indexing:**
```javascript
db.monthly_financial_summaries.createIndex(
  { user_id: 1, month: -1 },
  { unique: true }
);

db.monthly_financial_summaries.createIndex({
  user_id: 1,
  month: 1
});
```

**Generation Method:**

```javascript
// Cron: First day of month at 00:00 UTC (generates previous month)
async function generateMonthlySummary(userId, month) {
  const startOfMonth = new Date(month);
  startOfMonth.setUTCHours(0, 0, 0, 0);
  
  const endOfMonth = new Date(month);
  endOfMonth.setUTCMonth(endOfMonth.getUTCMonth() + 1);
  endOfMonth.setUTCHours(0, 0, 0, 0);
  endOfMonth.setUTCMilliseconds(-1);
  
  const aggregation = await db.transactions.aggregate([
    {
      $match: {
        user_id: userId,
        date: { $gte: startOfMonth, $lte: endOfMonth },
        deleted_at: null
      }
    },
    {
      $facet: {
        byType: [
          {
            $group: {
              _id: "$type",
              amount: { $sum: "$amount" }
            }
          }
        ],
        byCategory: [
          {
            $match: { type: "expense" }
          },
          {
            $lookup: {
              from: "user_categories",
              localField: "category_id",
              foreignField: "_id",
              as: "category"
            }
          },
          {
            $group: {
              _id: "$category_id",
              category_name: { $first: "$category.name" },
              amount: { $sum: "$amount" },
              count: { $sum: 1 }
            }
          }
        ],
        byPocket: [
          {
            $group: {
              _id: "$pocket_from",
              income: {
                $sum: {
                  $cond: [{ $eq: ["$type", "income"] }, "$amount", 0]
                }
              },
              expense: {
                $sum: {
                  $cond: [{ $eq: ["$type", "expense"] }, "$amount", 0]
                }
              }
            }
          }
        ],
        count: [
          {
            $count: "total"
          }
        ]
      }
    }
  ]).toArray();
  
  const summary = {
    user_id: userId,
    month: startOfMonth,
    income: aggregation[0].byType.find(t => t._id === "income")?.amount || 0,
    expense: aggregation[0].byType.find(t => t._id === "expense")?.amount || 0,
    transfer_in: aggregation[0].byType.find(t => t._id === "transfer")?.amount || 0,
    transfer_out: 0,
    net: 0,
    opening_balance: await getOpeningBalance(userId, startOfMonth),
    closing_balance: await getClosingBalance(userId, endOfMonth),
    expense_by_category: aggregation[0].byCategory,
    by_pocket: aggregation[0].byPocket,
    ytd_income: 0,
    ytd_expense: 0,
    ytd_net: 0,
    day_count: 0,
    transaction_count: aggregation[0].count[0]?.total || 0,
    is_complete: true,
    created_at: new Date(),
    updated_at: new Date()
  };
  
  summary.net = summary.income - summary.expense;
  
  // Calculate YTD
  const ytdData = await db.monthly_financial_summaries.aggregate([
    {
      $match: {
        user_id: userId,
        month: {
          $gte: new Date(month.getUTCFullYear(), 0, 1),
          $lte: endOfMonth
        }
      }
    },
    {
      $group: {
        _id: null,
        ytd_income: { $sum: "$income" },
        ytd_expense: { $sum: "$expense" }
      }
    }
  ]).toArray();
  
  if (ytdData.length > 0) {
    summary.ytd_income = ytdData[0].ytd_income;
    summary.ytd_expense = ytdData[0].ytd_expense;
    summary.ytd_net = summary.ytd_income - summary.ytd_expense;
  }
  
  await db.monthly_financial_summaries.updateOne(
    { user_id: userId, month: startOfMonth },
    { $set: summary },
    { upsert: true }
  );
}
```

### 4.3 Collection: `pocket_balance_snapshots`

**Purpose:** Point-in-time balance records for audit trail and historical queries

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  pocket_id: ObjectId,
  
  // Balance at snapshot time
  balance: Decimal128,
  balance_before: Decimal128,           // Previous snapshot
  change: Decimal128,
  
  // Snapshot metadata
  snapshot_time: Date,                  // Timestamp of snapshot
  snapshot_type: String,                // "hourly", "daily", "transaction"
  
  // Transaction that caused change (if applicable)
  transaction_id: ObjectId,
  transaction_type: String,
  transaction_amount: Decimal128,
  
  // Context
  pocket_name: String,
  pocket_type: String,
  
  created_at: Date
}
```

**Indexing:**
```javascript
db.pocket_balance_snapshots.createIndex({
  user_id: 1,
  pocket_id: 1,
  snapshot_time: -1
});

db.pocket_balance_snapshots.createIndex({
  user_id: 1,
  snapshot_time: -1
});

db.pocket_balance_snapshots.createIndex(
  { snapshot_time: 1 },
  { 
    expireAfterSeconds: 7776000  // 90 days
  }
);
```

**Generation Method:**

**Option A: Transaction Hook (Real-time)**
```javascript
// After every balance update
async function recordBalanceSnapshot(userId, pocketId, transaction) {
  const pocket = await db.pockets.findOne({ _id: pocketId });
  const prevSnapshot = await db.pocket_balance_snapshots.findOne(
    { user_id: userId, pocket_id: pocketId },
    { sort: { snapshot_time: -1 } }
  );
  
  await db.pocket_balance_snapshots.insertOne({
    user_id: userId,
    pocket_id: pocketId,
    balance: pocket.balance,
    balance_before: prevSnapshot?.balance || 0,
    change: pocket.balance - (prevSnapshot?.balance || 0),
    snapshot_time: new Date(),
    snapshot_type: "transaction",
    transaction_id: transaction._id,
    transaction_type: transaction.type,
    transaction_amount: transaction.amount,
    pocket_name: pocket.name,
    pocket_type: pocket.type,
    created_at: new Date()
  });
}
```

**Option B: Hourly Batch (Lower write volume)**
```javascript
// Cron: Every hour
async function recordHourlySnapshots(userId) {
  const pockets = await db.pockets.find({
    user_id: userId,
    deleted_at: null
  }).toArray();
  
  const snapshots = pockets.map(pocket => ({
    user_id: userId,
    pocket_id: pocket._id,
    balance: pocket.balance,
    snapshot_time: new Date(),
    snapshot_type: "hourly",
    pocket_name: pocket.name,
    pocket_type: pocket.type,
    created_at: new Date()
  }));
  
  await db.pocket_balance_snapshots.insertMany(snapshots);
}
```

### 4.4 Trade-offs Summary

| Collection | Consistency | Performance | Write Volume | Use Case |
|------------|-------------|-------------|--------------|----------|
| daily_financial_snapshots | Strong (daily lock) | Excellent | Low (1/day) | Reports, historical queries |
| monthly_financial_summaries | Strong (monthly lock) | Excellent | Very Low (1/month) | Charts, long-term analysis |
| pocket_balance_snapshots | Eventual (hourly) | Good | Medium (hourly) | Audit trail, balance history |
| daily_financial_reports | Strong (daily lock) | Excellent | Low (1/day) | Detailed daily reports |

---

## 5. IMPLEMENTATION ROADMAP

### Phase 1: Foundation (Week 1)
1. Create `daily_financial_snapshots` collection + indexes
2. Create `monthly_financial_summaries` collection + indexes
3. Create `pocket_balance_snapshots` collection + indexes
4. Deploy cron jobs for snapshot generation

### Phase 2: Reporting (Week 2)
1. Create `daily_financial_reports` collection + indexes
2. Implement daily report generation pipeline
3. Build report API endpoints
4. Add caching layer (Redis)

### Phase 3: Dashboard (Week 3)
1. Implement KPI aggregation pipelines
2. Implement chart aggregation pipelines
3. Add query optimization and caching
4. Performance testing and tuning

### Phase 4: AI Layer (Week 4)
1. Create `ai_transaction_enrichment` collection
2. Create `ai_spending_patterns` collection
3. Create `ai_financial_insights` collection
4. Implement enrichment pipeline
5. Build chatbot data endpoints

---

## 6. PERFORMANCE BENCHMARKS (Expected)

| Query | Collection | Expected Time | Notes |
|-------|-----------|----------------|-------|
| Total balance | pockets | <10ms | Real-time, indexed |
| Monthly KPIs | daily_financial_reports | <50ms | Pre-aggregated |
| 12-month chart | monthly_financial_summaries | <100ms | Pre-aggregated |
| Category breakdown | daily_financial_reports | <50ms | Pre-aggregated |
| Anomalies | ai_transaction_enrichment | <100ms | Indexed filter |
| Spending patterns | ai_spending_patterns | <50ms | Indexed lookup |

---

## 7. MIGRATION STRATEGY

**Do NOT modify existing collections.** All new collections are additive:

1. Deploy new collections in parallel
2. Start populating snapshots from current data
3. Backfill historical data (last 12 months)
4. Enable cron jobs
5. Migrate dashboard queries incrementally
6. Monitor and validate before full cutover

---

## 8. MONITORING & ALERTS

```javascript
// Monitor snapshot generation lag
db.daily_financial_snapshots.aggregate([
  {
    $match: { is_complete: false }
  },
  {
    $group: {
      _id: null,
      count: { $sum: 1 },
      oldest: { $min: "$created_at" }
    }
  }
]);

// Monitor collection sizes
db.daily_financial_snapshots.stats().size;
db.monthly_financial_summaries.stats().size;
db.ai_transaction_enrichment.stats().size;

// Monitor query performance
db.setProfilingLevel(1, { slowms: 100 });
```

---

## 9. COST OPTIMIZATION

- **Snapshot collections:** Reduce raw transaction aggregation by 95%
- **TTL indexes:** Auto-delete old snapshots (configurable retention)
- **Denormalization:** Trade storage for query speed (acceptable for financial data)
- **Batch operations:** Cron jobs run during off-peak hours
- **Indexing:** Composite indexes reduce scan overhead

---

## 10. SECURITY CONSIDERATIONS

- All collections include `user_id` in primary index (data isolation)
- Snapshots are immutable after `is_complete: true`
- Balance snapshots provide audit trail for disputes
- AI enrichment excludes sensitive data (PII)
- All timestamps in UTC for consistency

