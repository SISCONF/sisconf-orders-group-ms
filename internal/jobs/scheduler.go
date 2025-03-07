package jobs

import (
	"fmt"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/utils"
	"github.com/go-co-op/gocron/v2"
)

func NewSisConfCronScheduler() gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	utils.PanicOnError("Couldn't create scheduler", err)

	job, err := scheduler.NewJob(
		gocron.CronJob(
			"*/5 * * * *",
			false,
		),
		gocron.NewTask(func() {
			SaveAllAvailableFoods()
		}),
	)

	errMsg := fmt.Sprintf("Couldn't created cron job %d", job.ID())
	utils.PanicOnError(errMsg, err)

	return scheduler
}
