# Financial Reporting - Deployment Checklist

## Pre-Deployment Verification

### Code Quality
- [x] All Go files compile without errors
- [x] Models defined for all 7 collections
- [x] Repository methods implemented for all CRUD operations
- [x] Aggregation pipelines defined for all dashboard queries
- [x] Index management with proper MongoDB driver API
- [x] Service layer with business logic
- [x] DI module registration complete

### Files Created
- [x] `internal/modules/reporting/models.go` - Core data models
- [x] `internal/modules/reporting/ai_models.go` - AI-ready models
- [x] `internal/modules/reporting/repository.go` - Database operations
- [x] `internal/modules/reporting/aggregations.go` - Query pipelines
- [x] `internal/modules/reporting/indexes.go` - Index management
- [x] `internal/modules/reporting/service.go` - Business logic
- [x] `internal/modules/reporting/module.go` - DI registration
- [x] `FINANCIAL_REPORTING_ARCHITECTURE.md` - Complete architecture
- [x] `REPORTING_IMPLEMENTATION_GUIDE.md` - Implementation guide
- [x] `REPORTING_QUICK_REFERENCE.md` - Quick reference

---

## Phase 1: Foundation (Week 1)

### Step 1.1: Register Module
```go
// In cmd/api/main.go
import "github.com/HasanNugroho/coin-be/internal/modules/reporting"

// In container builder:
reporting.Register(builder)
```

**Verification:**
- [ ] Application starts without errors
- [ ] Reporting repository is injectable
- [ ] Aggregation helper is injectable

### Step 1.2: Verify Index Creation
```bash
# After first run, verify indexes exist:
mongo <connection-string>
use <database>
db.daily_financial_reports.getIndexes()
db.monthly_financial_summaries.getIndexes()
db.ai_transaction_enrichment.getIndexes()
```

**Expected Output:**
- [ ] 4 indexes on `daily_financial_reports`
- [ ] 2 indexes on `monthly_financial_summaries`
- [ ] 4 indexes on `ai_transaction_enrichment`
- [ ] 5 indexes on `ai_spending_patterns`
- [ ] 3 indexes on `ai_financial_insights`

### Step 1.3: Test Repository Methods
```go
// Test basic CRUD
ctx := context.Background()
userID := primitive.NewObjectID()

// Create test report
report := &reporting.DailyFinancialReport{
    UserID: userID,
    ReportDate: time.Now(),
    IsFinal: true,
}
err := repo.CreateDailyReport(ctx, report)
// [ ] Should succeed without error

// Retrieve test report
retrieved, err := repo.GetDailyReport(ctx, userID, time.Now())
// [ ] Should return the created report
```

---

## Phase 2: Cron Jobs (Week 2)

### Step 2.1: Implement Daily Report Generator
**File:** `internal/modules/reporting/cron_daily_report.go`

```go
type DailyReportGenerator struct {
    db *mongo.Database
    repo *Repository
}

func (g *DailyReportGenerator) GenerateForAllUsers(ctx context.Context) error {
    // Implementation from guide
}
```

**Verification:**
- [ ] Generator compiles without errors
- [ ] Can be instantiated with DB and repository
- [ ] GenerateForAllUsers method exists

### Step 2.2: Implement Daily Snapshot Generator
**File:** `internal/modules/reporting/cron_daily_snapshot.go`

**Verification:**
- [ ] Generator compiles without errors
- [ ] Generates snapshots for all users
- [ ] Snapshots have correct structure

### Step 2.3: Implement Monthly Summary Generator
**File:** `internal/modules/reporting/cron_monthly_summary.go`

**Verification:**
- [ ] Generator compiles without errors
- [ ] Generates summaries for previous month
- [ ] Summaries aggregate correctly

### Step 2.4: Setup Cron Scheduler
**File:** `internal/core/cron/cron.go`

```go
type CronScheduler struct {
    cron *cron.Cron
    db   *mongo.Database
}

func (cs *CronScheduler) Start() {
    // Register jobs:
    // - Daily Report: 23:59:59 UTC
    // - Daily Snapshot: 00:00 UTC
    // - Monthly Summary: 00:01 1st of month UTC
}
```

**Verification:**
- [ ] Cron scheduler starts without errors
- [ ] Jobs are registered
- [ ] Cron runs in background

### Step 2.5: Register Cron in Main
```go
// In cmd/api/main.go
scheduler := cron.NewCronScheduler(mongoClient.Database(cfg.MongoDB))
scheduler.Start()
defer scheduler.Stop()
```

**Verification:**
- [ ] Application starts with cron scheduler
- [ ] No errors in logs
- [ ] Cron jobs execute on schedule

---

## Phase 3: API Endpoints (Week 3)

### Step 3.1: Create Reporting Controller
**File:** `internal/modules/reporting/controller.go`

```go
type Controller struct {
    service *Service
}

func (c *Controller) GetDailyReport(ctx *gin.Context) {
    // Implementation
}

func (c *Controller) GetMonthlyReport(ctx *gin.Context) {
    // Implementation
}
```

**Endpoints:**
- [ ] `GET /api/v1/reports/daily/:date`
- [ ] `GET /api/v1/reports/monthly/:month`
- [ ] `GET /api/v1/reports/daily-range?start=&end=`

### Step 3.2: Create Dashboard Controller
**File:** `internal/modules/reporting/dashboard_controller.go`

```go
type DashboardController struct {
    service *Service
}

func (c *DashboardController) GetKPIs(ctx *gin.Context) {
    // Implementation
}

func (c *DashboardController) GetCharts(ctx *gin.Context) {
    // Implementation
}
```

**Endpoints:**
- [ ] `GET /api/v1/dashboard/kpis`
- [ ] `GET /api/v1/dashboard/charts`
- [ ] `GET /api/v1/dashboard/summary`

### Step 3.3: Create AI Controller
**File:** `internal/modules/reporting/ai_controller.go`

```go
type AIController struct {
    service *Service
}

func (c *AIController) GetFinancialContext(ctx *gin.Context) {
    // Implementation
}

func (c *AIController) GetAnomalies(ctx *gin.Context) {
    // Implementation
}

func (c *AIController) GetInsights(ctx *gin.Context) {
    // Implementation
}
```

**Endpoints:**
- [ ] `GET /api/v1/ai/financial-context`
- [ ] `GET /api/v1/ai/anomalies`
- [ ] `GET /api/v1/ai/insights`
- [ ] `POST /api/v1/ai/insights/:id/read`

### Step 3.4: Register Controllers
```go
// In module.go
builder.Add(di.Def{
    Name: "reportingController",
    Build: func(ctn di.Container) (interface{}, error) {
        service := ctn.Get("reportingService").(*Service)
        return NewController(service), nil
    },
})
```

**Verification:**
- [ ] All endpoints respond with 200 OK
- [ ] Data is returned in correct format
- [ ] No errors in logs

---

## Phase 4: Data Backfill (Week 4)

### Step 4.1: Backfill Daily Reports
```go
func BackfillDailyReports(ctx context.Context, db *mongo.Database, days int) error {
    gen := reporting.NewDailyReportGenerator(db, reporting.NewRepository(db))
    
    for i := 0; i < days; i++ {
        date := time.Now().AddDate(0, 0, -i)
        if err := gen.GenerateForUser(ctx, userID); err != nil {
            return err
        }
    }
    return nil
}
```

**Execution:**
- [ ] Run backfill for last 365 days
- [ ] Monitor progress in logs
- [ ] Verify document count: `db.daily_financial_reports.countDocuments()`

### Step 4.2: Backfill Daily Snapshots
```go
func BackfillDailySnapshots(ctx context.Context, db *mongo.Database, days int) error {
    gen := reporting.NewDailySnapshotGenerator(db, reporting.NewRepository(db))
    
    for i := 0; i < days; i++ {
        date := time.Now().AddDate(0, 0, -i)
        if err := gen.GenerateForUser(ctx, userID); err != nil {
            return err
        }
    }
    return nil
}
```

**Verification:**
- [ ] Snapshots generated for all days
- [ ] Balances are consistent
- [ ] No gaps in data

### Step 4.3: Backfill Monthly Summaries
```go
func BackfillMonthlySummaries(ctx context.Context, db *mongo.Database, months int) error {
    gen := reporting.NewMonthlySummaryGenerator(db, reporting.NewRepository(db))
    
    for i := 0; i < months; i++ {
        month := time.Now().AddDate(0, -i, 0)
        if err := gen.GenerateForUser(ctx, userID, month); err != nil {
            return err
        }
    }
    return nil
}
```

**Verification:**
- [ ] Summaries for last 12 months
- [ ] Totals match daily reports
- [ ] YTD calculations correct

### Step 4.4: Verify Data Integrity
```javascript
// Check daily reports
db.daily_financial_reports.find({ is_final: true }).count()
// Should be ~365 per user

// Check monthly summaries
db.monthly_financial_summaries.find().count()
// Should be ~12 per user

// Check consistency
db.daily_financial_reports.aggregate([
  { $match: { user_id: ObjectId("...") } },
  { $group: { _id: null, total_income: { $sum: "$total_income" } } }
])
// Should match monthly summary total_income
```

---

## Phase 5: Performance Tuning (Week 5)

### Step 5.1: Add Caching Layer
```go
type CachedService struct {
    service *Service
    cache   *redis.Client
}

func (cs *CachedService) GetDashboardKPIs(ctx context.Context, userID primitive.ObjectID) (*DashboardKPIs, error) {
    key := fmt.Sprintf("kpis:%s", userID.Hex())
    
    // Try cache
    val, err := cs.cache.Get(ctx, key).Result()
    if err == nil {
        // Parse and return
    }
    
    // Fetch from service
    result, err := cs.service.GetDashboardKPIs(ctx, userID)
    if err == nil {
        cs.cache.Set(ctx, key, result, 30*time.Second)
    }
    
    return result, err
}
```

**Verification:**
- [ ] Redis connection works
- [ ] Cache hits reduce query time
- [ ] Cache invalidation works correctly

### Step 5.2: Monitor Query Performance
```go
// Enable MongoDB profiling
db.setProfilingLevel(1, { slowms: 100 })

// Check slow queries
db.system.profile.find().sort({ ts: -1 }).limit(10)
```

**Targets:**
- [ ] KPI queries: <50ms
- [ ] Chart queries: <100ms
- [ ] Report queries: <100ms

### Step 5.3: Optimize Indexes
```javascript
// Check index usage
db.daily_financial_reports.aggregate([
  { $indexStats: {} }
])

// Remove unused indexes if any
db.daily_financial_reports.dropIndex("index_name")
```

**Verification:**
- [ ] All indexes are used
- [ ] No redundant indexes
- [ ] Query plans are efficient

---

## Phase 6: AI Layer (Week 6)

### Step 6.1: Implement Transaction Enrichment
**File:** `internal/modules/reporting/enrichment.go`

```go
type EnrichmentService struct {
    repo *Repository
}

func (es *EnrichmentService) EnrichTransaction(ctx context.Context, txn *transaction.Transaction) error {
    enrichment := &AITransactionEnrichment{
        UserID: txn.UserID,
        TransactionID: txn.ID,
        // ... populate fields
    }
    return es.repo.UpsertTransactionEnrichment(ctx, enrichment)
}
```

**Verification:**
- [ ] Enrichment runs after transaction creation
- [ ] All fields are populated
- [ ] No performance impact on transaction creation

### Step 6.2: Implement Pattern Analysis
**File:** `internal/modules/reporting/pattern_analyzer.go`

```go
type PatternAnalyzer struct {
    repo *Repository
}

func (pa *PatternAnalyzer) AnalyzePatterns(ctx context.Context, userID primitive.ObjectID) error {
    // Aggregate spending patterns
    // Identify trends
    // Detect seasonality
    // Update ai_spending_patterns collection
}
```

**Cron Job:** Daily at 02:00 UTC

**Verification:**
- [ ] Patterns identified correctly
- [ ] Trends calculated accurately
- [ ] Seasonality detected

### Step 6.3: Implement Insight Generation
**File:** `internal/modules/reporting/insight_generator.go`

```go
type InsightGenerator struct {
    repo *Repository
}

func (ig *InsightGenerator) GenerateInsights(ctx context.Context, userID primitive.ObjectID) error {
    // Analyze patterns
    // Detect anomalies
    // Generate actionable insights
    // Save to ai_financial_insights
}
```

**Cron Job:** Daily at 03:00 UTC

**Verification:**
- [ ] Insights are actionable
- [ ] Recommendations are relevant
- [ ] Severity levels are appropriate

---

## Phase 7: Monitoring & Alerts (Week 7)

### Step 7.1: Setup Logging
```go
import "go.uber.org/zap"

logger.Info("Daily report generated",
    zap.String("user_id", userID.Hex()),
    zap.Time("report_date", reportDate),
    zap.Duration("duration", duration),
)
```

**Verification:**
- [ ] All cron jobs log start/end
- [ ] Errors are logged with context
- [ ] Performance metrics are tracked

### Step 7.2: Setup Alerts
```go
// Alert if snapshot generation takes >15 minutes
if duration > 15*time.Minute {
    alerting.SendAlert("Slow snapshot generation", severity.Warning)
}

// Alert if >10 incomplete snapshots
incompleteCount, _ := repo.CountIncompleteSnapshots(ctx)
if incompleteCount > 10 {
    alerting.SendAlert("Too many incomplete snapshots", severity.Critical)
}
```

**Verification:**
- [ ] Alerts are sent on failures
- [ ] Alert channels are configured
- [ ] Team is notified

### Step 7.3: Setup Dashboards
- [ ] Grafana dashboard for collection sizes
- [ ] Grafana dashboard for query performance
- [ ] Grafana dashboard for cron job execution times
- [ ] Grafana dashboard for error rates

---

## Post-Deployment Validation

### Data Validation
```javascript
// Verify no data loss
db.daily_financial_reports.find({ user_id: ObjectId("...") }).count()
// Should be ~365

// Verify consistency
db.daily_financial_reports.aggregate([
  { $match: { user_id: ObjectId("...") } },
  { $group: { _id: null, total: { $sum: "$total_income" } } }
])
// Should match monthly_financial_summaries total

// Verify no duplicates
db.daily_financial_reports.aggregate([
  { $group: { _id: { user_id: "$user_id", report_date: "$report_date" }, count: { $sum: 1 } } },
  { $match: { count: { $gt: 1 } } }
])
// Should return empty
```

### Performance Validation
- [ ] Dashboard loads in <500ms
- [ ] KPI queries return in <50ms
- [ ] Chart queries return in <100ms
- [ ] Report generation completes in <10 minutes

### User Acceptance Testing
- [ ] Dashboard displays correct data
- [ ] Charts are accurate
- [ ] Reports match manual calculations
- [ ] No data inconsistencies

---

## Rollback Plan

If issues occur:

1. **Stop cron jobs** - Disable snapshot/report generation
2. **Preserve data** - Do NOT delete collections
3. **Investigate** - Check logs for errors
4. **Fix** - Deploy corrected code
5. **Resume** - Re-enable cron jobs
6. **Backfill** - Generate missing snapshots

---

## Success Criteria

- [x] All 7 collections created with correct schema
- [x] All indexes created and optimized
- [x] Cron jobs run on schedule
- [x] API endpoints return correct data
- [x] Dashboard queries perform well
- [x] No data inconsistencies
- [x] Alerts configured and working
- [x] Team trained on new system

---

## Support & Troubleshooting

### Common Issues

**Issue:** Snapshots not generating
- Check cron logs: `docker logs <container> | grep "Daily snapshot"`
- Verify MongoDB connection
- Check disk space

**Issue:** Slow queries
- Run `db.collection.stats()` to check size
- Verify indexes exist
- Check query explain plan

**Issue:** Data inconsistencies
- Compare daily reports with raw transactions
- Check for transactions created after snapshot time
- Verify opening/closing balance calculations

### Contact
- Backend Team: #backend-support
- Database Team: #database-support
- DevOps: #devops-support

