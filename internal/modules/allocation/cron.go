package allocation

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type CronJob struct {
	service *Service
	cron    *cron.Cron
}

func NewCronJob(service *Service) *CronJob {
	return &CronJob{
		service: service,
		cron:    cron.New(cron.WithLocation(getJakartaLocation())),
	}
}

// Start begins the daily allocation processing cron job
// Runs every day at 01:00 AM (Asia/Jakarta timezone)
func (c *CronJob) Start() error {
	_, err := c.cron.AddFunc("0 1 * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		log.Println("Starting daily allocation processing...")
		if err := c.service.ProcessDailyAllocations(ctx); err != nil {
			log.Printf("Error processing daily allocations: %v", err)
		}
		log.Println("Daily allocation processing completed")
	})

	if err != nil {
		return err
	}

	c.cron.Start()
	log.Println("Allocation cron job started")
	return nil
}

// Stop stops the cron job
func (c *CronJob) Stop() {
	c.cron.Stop()
	log.Println("Allocation cron job stopped")
}

// getJakartaLocation returns the Asia/Jakarta timezone location
func getJakartaLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Failed to load Asia/Jakarta timezone: %v, using UTC", err)
		return time.UTC
	}
	return loc
}
