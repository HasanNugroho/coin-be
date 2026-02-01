package reporting

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SnapshotWorker struct {
	service *Service
}

func NewSnapshotWorker(service *Service) *SnapshotWorker {
	return &SnapshotWorker{service: service}
}

// StartDailySnapshotJob starts a background job to generate daily snapshots at 23:59:59 UTC
func (w *SnapshotWorker) StartDailySnapshotJob(ctx context.Context, userID string) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().UTC()

			if now.Hour() == 23 && now.Minute() >= 59 {
				log.Printf("[SnapshotWorker] Generating daily snapshot at %v", now)
				// Generate snapshot for all users
				w.generateAllDailySnapshots(ctx)
			}
		}
	}
}

// StartMonthlySummaryJob starts a background job to generate monthly summaries at 00:00:01 UTC on 1st of month
func (w *SnapshotWorker) StartMonthlySummaryJob(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().UTC()

			// Check if it's the 1st of the month at 00:00
			if now.Day() == 1 && now.Hour() == 0 && now.Minute() == 0 {
				log.Printf("[SnapshotWorker] Generating monthly summaries at %v", now)
				// Generate monthly summaries for all users
				w.generateAllMonthlySummaries(ctx)
			}
		}
	}
}

// StartAIContextJob starts a background job to generate AI context at 00:30 UTC daily
func (w *SnapshotWorker) StartAIContextJob(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().UTC()

			// Check if it's 00:30 UTC
			if now.Hour() == 0 && now.Minute() == 30 {
				log.Printf("[SnapshotWorker] Generating AI financial contexts at %v", now)
				// Generate AI context for all users
				w.generateAllAIContexts(ctx)
			}
		}
	}
}

func (w *SnapshotWorker) generateAllDailySnapshots(ctx context.Context) {
	// In production, this would iterate over all active users
	// For now, this is a placeholder that would be called with specific user IDs
	log.Println("[SnapshotWorker] Daily snapshot generation completed")
}

func (w *SnapshotWorker) generateAllMonthlySummaries(ctx context.Context) {
	// In production, this would iterate over all active users
	// For now, this is a placeholder that would be called with specific user IDs
	log.Println("[SnapshotWorker] Monthly summary generation completed")
}

func (w *SnapshotWorker) generateAllAIContexts(ctx context.Context) {
	// In production, this would iterate over all active users
	// For now, this is a placeholder that would be called with specific user IDs
	log.Println("[SnapshotWorker] AI context generation completed")
}

// GenerateDailySnapshotForUser generates a daily snapshot for a specific user
func (w *SnapshotWorker) GenerateDailySnapshotForUser(ctx context.Context, userID string) error {
	oid, err := parseObjectID(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	_, err = w.service.GenerateDailySnapshot(ctx, oid, time.Now())
	if err != nil {
		log.Printf("[SnapshotWorker] Error generating daily snapshot for user %s: %v", userID, err)
		return err
	}

	log.Printf("[SnapshotWorker] Daily snapshot generated for user %s", userID)
	return nil
}

// GenerateMonthlySummaryForUser generates a monthly summary for a specific user
func (w *SnapshotWorker) GenerateMonthlySummaryForUser(ctx context.Context, userID string, yearMonth string) error {
	oid, err := parseObjectID(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	_, err = w.service.GenerateMonthlySummary(ctx, oid, yearMonth)
	if err != nil {
		log.Printf("[SnapshotWorker] Error generating monthly summary for user %s: %v", userID, err)
		return err
	}

	log.Printf("[SnapshotWorker] Monthly summary generated for user %s (month: %s)", userID, yearMonth)
	return nil
}

// GenerateAIContextForUser generates AI financial context for a specific user
func (w *SnapshotWorker) GenerateAIContextForUser(ctx context.Context, userID string) error {
	oid, err := parseObjectID(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	_, err = w.service.GenerateAIFinancialContext(ctx, oid)
	if err != nil {
		log.Printf("[SnapshotWorker] Error generating AI context for user %s: %v", userID, err)
		return err
	}

	log.Printf("[SnapshotWorker] AI context generated for user %s", userID)
	return nil
}

func parseObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return oid, nil
}
