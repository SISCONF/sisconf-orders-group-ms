package jobs

import (
	"log"
	"os"

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

	err = os.Mkdir("./internal/data", 0777)
	if err != nil {
		log.Printf("Failed to create foods.json dir: %s\n", err.Error())
	}
	err = files.WriteApiResponseToJSON(foodsList, "./internal/data/foods.json")
	if err != nil {
		log.Printf("Could not save JSON file %s\n", err.Error())
	}
}
