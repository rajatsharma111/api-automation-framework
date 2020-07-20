package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// JSONConfig  consists json object
type JSONConfig struct {
	ConfigJSON map[string]interface{}
}

var instance *JSONConfig

// GetConfig to get CloudConfig singelton session
func GetConfig(filePath string, encrypted bool, encrytionFilePath ...string) *JSONConfig {
	if encrypted {
		decryptConfig(encrytionFilePath[0], filePath)
	}
	configJSON := loadConfig(filePath)
	return &JSONConfig{configJSON}
}

func decryptConfig(encrytionFilePath string, decrytionFilePath string) {
	cipher, _ := os.Open(encrytionFilePath)
	data, _ := ioutil.ReadAll(cipher)
	kms := os.Getenv("GCP_KMS_NAME")
	if kms == "" {
		LogError("GCP_KMS_NAME environment variable must be set.")
	}
	KmsClient().DecryptSymmetric(decrytionFilePath, kms, data)
}

// LoadConfig reads json file path and return the json object
func loadConfig(filePath string) map[string]interface{} {
	var configJSON map[string]interface{}
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		LogError("Error loading " + filePath + " file")
	}
	json.Unmarshal([]byte(data), &configJSON)
	return configJSON
}
