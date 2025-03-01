package sisconf

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (sisconfApiClient SisconfAPIClient) Login(loginData LoginData) (string, error) {
	loginUrl := fmt.Sprintf("%s/people/login", sisconfApiClient.baseUrl)
	reqBody, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(reqBody)
	resp, err := http.Post(loginUrl, "application/json", reader)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("couldn't log user in")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var authenticationData LoginResponse
	err = json.Unmarshal(body, &authenticationData)
	if err != nil {
		return "", err
	}

	return authenticationData.AuthenticationToken, nil
}

func (sisconfApiClient SisconfAPIClient) UpdateOrdersGroupSheetFileURL(ordersGroupId int, sheetUrl string) error {
	godotenv.Load()

	ordersGroupPartialUpdateUrl := fmt.Sprintf(
		"%s/orders-group/%d/doc-url",
		sisconfApiClient.baseUrl,
		ordersGroupId,
	)
	requestBody, err := json.Marshal(map[string]any{
		"docUrl": sheetUrl,
	})
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		http.MethodPatch,
		ordersGroupPartialUpdateUrl,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}

	authenticationToken, err := sisconfApiClient.Login(LoginData{
		Email:    os.Getenv("SISCONF_MS_USER_EMAIL"),
		Password: os.Getenv("SISCONF_MS_USER_PASSWORD"),
	})
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authenticationToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("couldn't update orders group of %d id", ordersGroupId)
		return errors.New(errMsg)
	}

	return nil
}
