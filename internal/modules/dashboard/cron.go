package dashboard

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

func (c *CronJob) Start() {
	c.cron.AddFunc("1 0 * * *", func() {
		ctx := context.Background()
		yesterday := time.Now().AddDate(0, 0, -1)
		
		log.Printf("Starting daily summary generation for date: %s", yesterday.Format("2006-01-02"))
		
		err := c.service.GenerateDailySummariesForAllUsers(ctx, yesterday)
		if err != nil {
			log.Printf("Error generating daily summaries: %v", err)
		} else {
			log.Printf("Daily summaries generated successfully for date: %s", yesterday.Format("2006-01-02"))
		}
	})

	c.cron.Start()
	log.Println("Dashboard cron job started - Daily summaries will be generated at 00:01 every day")
}

func (c *CronJob) Stop() {
	c.cron.Stop()
	log.Println("Dashboard cron job stopped")
}
