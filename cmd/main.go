package main

import (
	"fmt"
	"time"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/jobs"
	"github.com/go-co-op/gocron/v2"
)

func main() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println("Error creating scheduler")
		return
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			jobs.SaveAllAvailableFoods,
		),
	)
	if err != nil {
		fmt.Println("Error creating job!")
		return
	}
	scheduler.Start()

	time.Sleep(time.Minute)

	err = scheduler.Shutdown()
	if err != nil {
		fmt.Println("Error during shutdown!")
	}
}
