package report

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"os"
	"time"
)

const storagePrefix = "https://storage.cloud.google.com"

// StorageSession holds Cloud client session
type StorageSession struct {
	client *storage.Client
}

var objectStorage *StorageSession

// GetStorageObject returns object objectStorage reference
func GetStorageObject() (*StorageSession, error) {
	if objectStorage == nil {
		err := getClient()
		if err != nil {
			return nil, err
		}
	}

	return &StorageSession{
		client: objectStorage.client,
	}, nil
}

// Get storage client once as the application starts and reuse it.
func getClient() error {
	ss := StorageSession{}
	ctx := context.Background()
	googleCred := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if googleCred == "" {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS environment variable must be set.")
	}
	opt := option.WithCredentialsFile(googleCred)
	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		return err
	}
	ss.client = client
	objectStorage = &ss
	return nil
}

// UploadFile to cloud storage.
// bucketName - cloud bucket name
// object - object name for cloud
// filePath - filePath for uploading file
func (storage *StorageSession) UploadFile(bucketName, object, filePath string) (string, error) {
	ctx := context.Background()
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := storage.client.Bucket(bucketName).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	url := storagePrefix + "/" + bucketName + "/" + object
	return url, nil
}
