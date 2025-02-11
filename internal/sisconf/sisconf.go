package sisconf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type SisconfAPIClient struct {
	baseUrl string
}

func NewSisconfAPIClient() SisconfAPIClient {
	err := godotenv.Load()
	if err != nil {
		errMsg := fmt.Sprintf("Could not load .env: %s", err.Error())
		panic(errMsg)
	}
	baseApiUrl := fmt.Sprintf(
		"%s/api",
		os.Getenv("SISCONF_SPRING_BACKEND_URL"),
	)
	return SisconfAPIClient{
		baseUrl: baseApiUrl,
	}
}

func (sisconfApiClient SisconfAPIClient) GetAllAvailableFoods() ([]Food, error) {
	getFoodsUrl := fmt.Sprintf("%s/foods", sisconfApiClient.baseUrl)
	resp, err := http.Get(getFoodsUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var foodsList []Food
	err = json.Unmarshal(body, &foodsList)
	if err != nil {
		return nil, err
	}

	return foodsList, err
}
