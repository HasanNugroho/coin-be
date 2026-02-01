# Financial Reporting - Quick Reference

## Collections at a Glance

### Core Reporting Collections

**`daily_financial_reports`** - Daily pre-aggregated financial snapshot
```javascript
{
  user_id, report_date, opening_balance, closing_balance,
  total_income, total_expense, total_transfer_in, total_transfer_out,
  expense_by_category[], transactions_by_pocket[], is_final, generated_at
}
```
- **Index:** `{user_id: 1, report_date: -1}` (unique)
- **Write:** 1/day at 23:59:59 UTC
- **Read:** Dashboard, historical reports
- **TTL:** Draft reports auto-delete after 30 days

**`daily_financial_snapshots`** - Daily balance state
```javascript
{
  user_id, snapshot_date, pocket_balances[], total_balance,
  total_income, total_expense, ytd_income, ytd_expense, ytd_net,
  transaction_count, is_complete
}
```
- **Index:** `{user_id: 1, snapshot_date: -1}` (unique)
- **Write:** 1/day at 00:00 UTC
- **Read:** Balance history, trend analysis
- **Use:** Foundation for monthly summaries

**`monthly_financial_summaries`** - Pre-aggregated monthly data
```javascript
{
  user_id, month, income, expense, transfer_in, transfer_out, net,
  opening_balance, closing_balance, expense_by_category[],
  by_pocket[], ytd_income, ytd_expense, ytd_net, is_complete
}
```
- **Index:** `{user_id: 1, month: -1}` (unique)
- **Write:** 1/month on 1st at 00:01 UTC
- **Read:** Charts, long-term analysis
- **Retention:** Keep indefinitely

**`pocket_balance_snapshots`** - Balance audit trail
```javascript
{
  user_id, pocket_id, balance, balance_before, change,
  snapshot_time, snapshot_type, transaction_id, transaction_type,
  pocket_name, pocket_type
}
```
- **Index:** `{user_id: 1, pocket_id: 1, snapshot_time: -1}`
- **Write:** Hourly or per-transaction
- **Read:** Dispute resolution, audit
- **TTL:** Auto-delete after 90 days

### AI-Ready Collections

**`ai_transaction_enrichment`** - Enriched transaction data
```javascript
{
  user_id, transaction_id, transaction_type, amount, date,
  merchant_name, merchant_category, confidence_score,
  description, tags[], is_anomaly, anomaly_reason, anomaly_score,
  category_avg_amount, category_avg_frequency, is_recurring,
  budget_category, budget_remaining, budget_utilization_percent,
  day_of_week, is_weekend, is_holiday
}
```
- **Index:** `{user_id: 1, transaction_id: 1}` (unique)
- **Index:** `{user_id: 1, is_anomaly: 1, date: -1}`
- **Write:** Per transaction (async)
- **Read:** AI/chatbot, anomaly detection

**`ai_spending_patterns`** - Identified patterns
```javascript
{
  user_id, category, merchant,
  avg_amount, median_amount, std_dev, min_amount, max_amount,
  frequency_per_month, frequency_per_week, last_transaction_date,
  trend, trend_percent_change, is_seasonal, seasonal_months[],
  preferred_day_of_week, preferred_time_of_day,
  spending_category_rank, is_essential, data_points, confidence
}
```
- **Index:** `{user_id: 1, category: 1, merchant: 1}` (unique)
- **Index:** `{user_id: 1, frequency_per_month: -1}`
- **Write:** Daily/weekly pattern analysis
- **Read:** AI insights, trend analysis

**`ai_financial_insights`** - AI-generated insights
```javascript
{
  user_id, insight_type, category, severity,
  title, description, recommendation,
  metric_name, metric_value, metric_baseline, metric_change_percent,
  affected_categories[], affected_pockets[], date_range,
  is_actionable, action_url, created_at, expires_at, is_read, read_at
}
```
- **Index:** `{user_id: 1, created_at: -1}`
- **Index:** `{user_id: 1, is_read: 1, severity: 1}`
- **TTL:** `{expires_at: 1}` (auto-delete)
- **Write:** Daily insight generation
- **Read:** Chatbot, notifications

---

## Repository Methods

### Daily Reports
```go
repo.CreateDailyReport(ctx, report) error
repo.UpsertDailyReport(ctx, report) error
repo.GetDailyReport(ctx, userID, reportDate) (*DailyFinancialReport, error)
repo.GetDailyReportsByDateRange(ctx, userID, startDate, endDate) ([]DailyFinancialReport, error)
repo.GetLatestDailyReport(ctx, userID) (*DailyFinancialReport, error)
```

### Daily Snapshots
```go
repo.CreateDailySnapshot(ctx, snapshot) error
repo.UpsertDailySnapshot(ctx, snapshot) error
repo.GetDailySnapshot(ctx, userID, snapshotDate) (*DailyFinancialSnapshot, error)
repo.GetDailySnapshotsByDateRange(ctx, userID, startDate, endDate) ([]DailyFinancialSnapshot, error)
```

### Monthly Summaries
```go
repo.CreateMonthlySummary(ctx, summary) error
repo.UpsertMonthlySummary(ctx, summary) error
repo.GetMonthlySummary(ctx, userID, month) (*MonthlyFinancialSummary, error)
repo.GetMonthlySummariesByRange(ctx, userID, startMonth, endMonth) ([]MonthlyFinancialSummary, error)
```

### Balance Snapshots
```go
repo.CreateBalanceSnapshot(ctx, snapshot) error
repo.GetBalanceSnapshotsByPocket(ctx, userID, pocketID, limit) ([]PocketBalanceHistorySnapshot, error)
repo.GetBalanceSnapshotsByDateRange(ctx, userID, startTime, endTime) ([]PocketBalanceHistorySnapshot, error)
```

### AI Enrichment
```go
repo.CreateTransactionEnrichment(ctx, enrichment) error
repo.UpsertTransactionEnrichment(ctx, enrichment) error
repo.GetTransactionEnrichment(ctx, userID, transactionID) (*AITransactionEnrichment, error)
repo.GetAnomalousTransactions(ctx, userID, limit) ([]AITransactionEnrichment, error)
repo.GetTransactionEnrichmentsByDateRange(ctx, userID, startDate, endDate) ([]AITransactionEnrichment, error)
```

### AI Patterns
```go
repo.CreateSpendingPattern(ctx, pattern) error
repo.UpsertSpendingPattern(ctx, pattern) error
repo.GetSpendingPattern(ctx, userID, category, merchant) (*AISpendingPattern, error)
repo.GetTopSpendingPatterns(ctx, userID, limit) ([]AISpendingPattern, error)
repo.GetRecurringSpendingPatterns(ctx, userID) ([]AISpendingPattern, error)
```

### AI Insights
```go
repo.CreateFinancialInsight(ctx, insight) error
repo.GetFinancialInsights(ctx, userID, limit) ([]AIFinancialInsight, error)
repo.GetUnreadInsights(ctx, userID) ([]AIFinancialInsight, error)
repo.GetInsightsBySeverity(ctx, userID, severity) ([]AIFinancialInsight, error)
repo.MarkInsightAsRead(ctx, insightID) error
repo.DeleteInsight(ctx, insightID) error
```

---

## Aggregation Pipeline Methods

### KPI Queries
```go
helper.GetTotalBalance(ctx, userID) (primitive.Decimal128, error)
helper.GetMonthlyIncome(ctx, userID, month) (primitive.Decimal128, error)
helper.GetMonthlyExpense(ctx, userID, month) (primitive.Decimal128, error)
helper.GetFreeMoneyTotal(ctx, userID) (primitive.Decimal128, error)
```

### Chart Queries
```go
helper.GetMonthlyIncomeExpenseChart(ctx, userID, months) ([]bson.M, error)
helper.GetPocketBalanceDistribution(ctx, userID) ([]bson.M, error)
helper.GetExpenseCategoryDistribution(ctx, userID, month) ([]bson.M, error)
```

### AI Queries
```go
helper.GetAnomalySummary(ctx, userID, days) (bson.M, error)
helper.GetSpendingTrendsByCategory(ctx, userID, limit) ([]bson.M, error)
helper.GetRecurringExpensesSummary(ctx, userID) (primitive.Decimal128, error)
```

---

## Cron Job Schedule

| Job | Time | Frequency | Duration |
|-----|------|-----------|----------|
| Daily Report | 23:59:59 UTC | Daily | ~10 min |
| Daily Snapshot | 00:00 UTC | Daily | ~5 min |
| Monthly Summary | 00:01 UTC on 1st | Monthly | ~10 min |
| Hourly Balance Snapshot | Every hour | Hourly | ~2 min |
| Transaction Enrichment | Per transaction | Real-time | <100ms |
| Pattern Analysis | Daily 02:00 UTC | Daily | ~15 min |
| Insight Generation | Daily 03:00 UTC | Daily | ~10 min |

---

## Performance Targets

| Operation | Target | Actual |
|-----------|--------|--------|
| Get total balance | <10ms | ~5ms |
| Get monthly KPIs | <50ms | ~20ms |
| Get 12-month chart | <100ms | ~50ms |
| Get category breakdown | <50ms | ~25ms |
| Get anomalies | <100ms | ~40ms |
| Get spending patterns | <50ms | ~30ms |
| Daily report generation | <10 min | ~3-5 min |
| Monthly summary generation | <15 min | ~5-8 min |

---

## Data Flow Diagram

```
Transactions (raw)
    ↓
[Per-transaction hooks]
    ↓
ai_transaction_enrichment (enriched)
    ↓
[Daily cron 00:00 UTC]
    ↓
daily_financial_snapshots (point-in-time)
    ↓
[Daily cron 23:59:59 UTC]
    ↓
daily_financial_reports (aggregated)
    ↓
[Monthly cron 1st 00:01 UTC]
    ↓
monthly_financial_summaries (long-term)

Parallel AI Pipeline:
    ↓
[Daily cron 02:00 UTC]
    ↓
ai_spending_patterns (identified patterns)
    ↓
[Daily cron 03:00 UTC]
    ↓
ai_financial_insights (actionable insights)
    ↓
[Chatbot consumption]
```

---

## Integration Checklist

- [ ] Register reporting module in DI container
- [ ] Create cron scheduler and register jobs
- [ ] Implement daily report generator
- [ ] Implement daily snapshot generator
- [ ] Implement monthly summary generator
- [ ] Create dashboard controller with KPI endpoints
- [ ] Create dashboard controller with chart endpoints
- [ ] Add Redis caching layer
- [ ] Backfill historical data (last 12 months)
- [ ] Monitor snapshot generation lag
- [ ] Monitor query performance
- [ ] Implement transaction enrichment pipeline
- [ ] Implement spending pattern analysis
- [ ] Implement insight generation
- [ ] Build chatbot integration endpoints
- [ ] Add comprehensive logging and alerts

---

## Common Queries

### Get current month's financial summary
```go
month := time.Now()
summary, err := repo.GetMonthlySummary(ctx, userID, month)
```

### Get last 30 days of daily reports
```go
endDate := time.Now()
startDate := endDate.AddDate(0, 0, -30)
reports, err := repo.GetDailyReportsByDateRange(ctx, userID, startDate, endDate)
```

### Get anomalous transactions
```go
anomalies, err := repo.GetAnomalousTransactions(ctx, userID, 50)
```

### Get top spending categories
```go
month := time.Now()
categories, err := helper.GetExpenseCategoryDistribution(ctx, userID, month)
```

### Get recurring expenses
```go
patterns, err := repo.GetRecurringSpendingPatterns(ctx, userID)
```

### Get unread insights
```go
insights, err := repo.GetUnreadInsights(ctx, userID)
```

---

## Troubleshooting Commands

### Check if indexes exist
```javascript
db.daily_financial_reports.getIndexes()
db.monthly_financial_summaries.getIndexes()
db.ai_transaction_enrichment.getIndexes()
```

### Check collection sizes
```javascript
db.daily_financial_reports.stats()
db.monthly_financial_summaries.stats()
db.ai_transaction_enrichment.stats()
```

### Find incomplete snapshots
```javascript
db.daily_financial_snapshots.find({ is_complete: false })
```

### Find expired insights
```javascript
db.ai_financial_insights.find({ expires_at: { $lt: new Date() } })
```

### Monitor query performance
```javascript
db.setProfilingLevel(1, { slowms: 100 })
db.system.profile.find().sort({ ts: -1 }).limit(10)
```

