package jobs

import (
	"log"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/files"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/sisconf"
)

func SaveAllAvailableFoods() {
	sisconfApiClient := sisconf.NewSisconfAPIClient()
	foodsList, err := sisconfApiClient.GetAllAvailableFoods()
	if err != nil {
		log.Printf("Could not retrieve foods: %s\n", err.Error())
		return
	}
	err = files.WriteApiResponseToJSON(foodsList, "./internal/data/foods.json")
	if err != nil {
		log.Printf("Could not save JSON file %s\n", err.Error())
	}
}
