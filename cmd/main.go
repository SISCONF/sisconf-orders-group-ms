package main

import (
	"log"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/jobs"
)

func main() {
	scheduler := jobs.NewSisConfCronScheduler()
	scheduler.Start()
	log.Println("Successfully started cron scheduler")

	var forever chan any

	<-forever
}
