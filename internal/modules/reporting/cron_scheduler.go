package reporting

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type CronScheduler struct {
	db                      *mongo.Database
	repo                    *Repository
	dailyReportGenerator    *DailyReportGenerator
	dailySnapshotGenerator  *DailySnapshotGenerator
	monthlySummaryGenerator *MonthlySummaryGenerator
	stopChan                chan struct{}
	wg                      sync.WaitGroup
	running                 bool
	mu                      sync.Mutex
}

func NewCronScheduler(db *mongo.Database, repo *Repository) *CronScheduler {
	return &CronScheduler{
		db:                      db,
		repo:                    repo,
		dailyReportGenerator:    NewDailyReportGenerator(db, repo),
		dailySnapshotGenerator:  NewDailySnapshotGenerator(db, repo),
		monthlySummaryGenerator: NewMonthlySummaryGenerator(db, repo),
		stopChan:                make(chan struct{}),
	}
}

// Start initializes and starts all cron jobs
func (cs *CronScheduler) Start() error {
	cs.mu.Lock()
	if cs.running {
		cs.mu.Unlock()
		return nil
	}
	cs.running = true
	cs.mu.Unlock()

	log.Println("[CronScheduler] Starting reporting cron jobs")

	// Daily Snapshot: 00:00 UTC every day
	cs.wg.Add(1)
	go cs.runDailySnapshotJob()
	log.Println("[CronScheduler] Started: Daily Snapshot job (runs at 00:00 UTC)")

	// Daily Report: 23:59 UTC every day
	cs.wg.Add(1)
	go cs.runDailyReportJob()
	log.Println("[CronScheduler] Started: Daily Report job (runs at 23:59 UTC)")

	// Monthly Summary: 00:01 UTC on 1st of month
	cs.wg.Add(1)
	go cs.runMonthlySummaryJob()
	log.Println("[CronScheduler] Started: Monthly Summary job (runs at 00:01 UTC on 1st)")

	log.Println("[CronScheduler] All cron jobs started successfully")
	return nil
}

// Stop gracefully stops the cron scheduler
func (cs *CronScheduler) Stop() {
	cs.mu.Lock()
	if !cs.running {
		cs.mu.Unlock()
		return
	}
	cs.running = false
	cs.mu.Unlock()

	log.Println("[CronScheduler] Stopping cron scheduler")
	close(cs.stopChan)
	cs.wg.Wait()
	log.Println("[CronScheduler] Cron scheduler stopped")
}

// runDailySnapshotJob runs the daily snapshot generation at 00:00 UTC
func (cs *CronScheduler) runDailySnapshotJob() {
	defer cs.wg.Done()

	for {
		now := time.Now().UTC()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		if nextRun.Before(now) {
			nextRun = nextRun.AddDate(0, 0, 1)
		}

		duration := nextRun.Sub(now)
		log.Printf("[CronScheduler] Daily snapshot scheduled in %v", duration)

		select {
		case <-cs.stopChan:
			return
		case <-time.After(duration):
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			snapshotDate := time.Now().UTC().AddDate(0, 0, -1)
			if err := cs.dailySnapshotGenerator.GenerateForAllUsers(ctx, snapshotDate); err != nil {
				log.Printf("[CronScheduler] Daily snapshot error: %v", err)
			}
			cancel()
		}
	}
}

// runDailyReportJob runs the daily report generation at 23:59 UTC
func (cs *CronScheduler) runDailyReportJob() {
	defer cs.wg.Done()

	for {
		now := time.Now().UTC()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, time.UTC)
		if nextRun.Before(now) {
			nextRun = nextRun.AddDate(0, 0, 1)
		}

		duration := nextRun.Sub(now)
		log.Printf("[CronScheduler] Daily report scheduled in %v", duration)

		select {
		case <-cs.stopChan:
			return
		case <-time.After(duration):
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			reportDate := time.Now().UTC()
			if err := cs.dailyReportGenerator.GenerateForAllUsers(ctx, reportDate); err != nil {
				log.Printf("[CronScheduler] Daily report error: %v", err)
			}
			cancel()
		}
	}
}

// runMonthlySummaryJob runs the monthly summary generation at 00:01 UTC on 1st
func (cs *CronScheduler) runMonthlySummaryJob() {
	defer cs.wg.Done()

	for {
		now := time.Now().UTC()
		nextRun := time.Date(now.Year(), now.Month(), 1, 0, 1, 0, 0, time.UTC)
		if nextRun.Before(now) {
			nextRun = nextRun.AddDate(0, 1, 0)
		}

		duration := nextRun.Sub(now)
		log.Printf("[CronScheduler] Monthly summary scheduled in %v", duration)

		select {
		case <-cs.stopChan:
			return
		case <-time.After(duration):
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
			now := time.Now().UTC()
			previousMonth := now.AddDate(0, -1, 0)
			if err := cs.monthlySummaryGenerator.GenerateForAllUsers(ctx, previousMonth); err != nil {
				log.Printf("[CronScheduler] Monthly summary error: %v", err)
			}
			cancel()
		}
	}
}

// TriggerDailyReportNow triggers daily report generation immediately (for testing/backfill)
func (cs *CronScheduler) TriggerDailyReportNow(ctx context.Context, reportDate time.Time) error {
	log.Printf("[CronScheduler] Triggering daily report for %s", reportDate.Format("2006-01-02"))
	return cs.dailyReportGenerator.GenerateForAllUsers(ctx, reportDate)
}

// TriggerDailySnapshotNow triggers daily snapshot generation immediately (for testing/backfill)
func (cs *CronScheduler) TriggerDailySnapshotNow(ctx context.Context, snapshotDate time.Time) error {
	log.Printf("[CronScheduler] Triggering daily snapshot for %s", snapshotDate.Format("2006-01-02"))
	return cs.dailySnapshotGenerator.GenerateForAllUsers(ctx, snapshotDate)
}

// TriggerMonthlySummaryNow triggers monthly summary generation immediately (for testing/backfill)
func (cs *CronScheduler) TriggerMonthlySummaryNow(ctx context.Context, month time.Time) error {
	log.Printf("[CronScheduler] Triggering monthly summary for %s", month.Format("2006-01"))
	return cs.monthlySummaryGenerator.GenerateForAllUsers(ctx, month)
}
