package payroll

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
		cron:    cron.New(),
	}
}

// Start begins the daily payroll processing cron job
// Runs every day at 00:01 AM (1 minute after midnight)
func (c *CronJob) Start() error {
	_, err := c.cron.AddFunc("1 0 * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		log.Println("Starting daily payroll processing...")
		if err := c.service.ProcessDailyPayroll(ctx); err != nil {
			log.Printf("Error processing daily payroll: %v", err)
		}
		log.Println("Daily payroll processing completed")
	})

	if err != nil {
		return err
	}

	c.cron.Start()
	log.Println("Payroll cron job started")
	return nil
}

// Stop stops the cron job
func (c *CronJob) Stop() {
	c.cron.Stop()
	log.Println("Payroll cron job stopped")
}
