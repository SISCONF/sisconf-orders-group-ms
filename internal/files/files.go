package files

import (
	"bytes"
	"encoding/json"
	"os"
)

func WriteApiResponseToJSON(responseBody any, filePath string) error {
	respBytes, err := json.Marshal(responseBody)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer([]byte{})
	err = json.Indent(buffer, respBytes, " ", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return err
}

func ReadStructJSON[T any](fileName string) (*T, error) {
	dataBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var data T
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, err
}
