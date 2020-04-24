package main

import (
	"log"
	"time"

	"github.com/VISIndonesia/conveyor"
	"github.com/VISIndonesia/conveyor/config"
)

func main() {
	psProjID := config.Configuration.PubSub.ProjectID
	psCredFile := config.Configuration.PubSub.CredFile
	subID := config.Configuration.PubSub.Subscriber
	timeout := config.Configuration.PubSub.Timeout

	bqProjID := config.Configuration.BigQuery.ProjectID
	bqCredFile := config.Configuration.BigQuery.CredFile
	bqDataset := config.Configuration.BigQuery.Dataset

	gap := config.Configuration.Gap

	timeoutDuration := time.Duration(timeout) * time.Second

	log.Printf("Running with configs: %+v\n", config.Configuration)

	for x := range time.Tick(time.Duration(gap) * time.Second) {
		log.Println("Starting job:", x)
		if events, err := conveyor.Consume(subID, timeoutDuration, psProjID, psCredFile); err != nil {
			log.Println("Error fetching messages", err)
		} else {
			if errs := conveyor.UploadEvents(events, bqDataset, bqProjID, bqCredFile); len(errs) != 0 {
				for _, err := range errs {
					log.Println("Error uploading data", err)
				}
			}
		}
		log.Println("run completed at:", time.Now())
	}
}
