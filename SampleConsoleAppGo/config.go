package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	OcpApimSubscriptionKey string `json:"Ocp-Apim-Subscription-Key"`
	OcpApimSubscriptionRegion string `json:"Ocp-Apim-Subscription-Region"`
	AzureTranslateURL string `json:"AzureTranslateURL"`
	YourDeploymentName string `json:"YOUR_DEPLOYMENT_NAME"`
	YourResourceName string `json:"YOUR_RESOURCE_NAME"`
	ApiKey string `json:"api-key"`
	Endpoint string `json:"endpoint"`
	Key string `json:"key"`
	IndexName string `json:"indexName"`
}

func readConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if (err != nil) {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
