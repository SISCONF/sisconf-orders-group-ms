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
