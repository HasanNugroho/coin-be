package dashboard

import (
	"context"
	"log"
	"time"

	"github.com/HasanNugroho/coin-be/internal/modules/daily_summary"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	service             *Service
	dailySummaryService *daily_summary.Service
	cron                *cron.Cron
}

func NewCronJob(service *Service, dss *daily_summary.Service) *CronJob {
	return &CronJob{
		service:             service,
		dailySummaryService: dss,
		cron:                cron.New(),
	}
}

func (c *CronJob) Start() {
	c.cron.AddFunc("1 0 * * *", func() {
		ctx := context.Background()
		yesterday := time.Now().AddDate(0, 0, -1)

		log.Printf("Starting daily summary generation for date: %s", yesterday.Format("2006-01-02"))

		err := c.dailySummaryService.GenerateDailySummariesForAllUsers(ctx, yesterday)
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
