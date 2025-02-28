package files

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
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

func UploadToS3Bucket(filename string) (string, error) {
	godotenv.Load()

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AMAZON_S3_BUCKET_REGION")),
	}))

	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AMAZON_S3_BUCKET_NAME")),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
