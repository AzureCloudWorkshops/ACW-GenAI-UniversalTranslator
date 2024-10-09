package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"bytes"
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

type TranslateResponse struct {
	DetectedLanguage struct {
		Language string `json:"language"`
		Score    float64 `json:"score"`
	} `json:"detectedLanguage"`
	Translations []struct {
		Text string `json:"text"`
		To   string `json:"to"`
	} `json:"translations"`
}

func readConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func makePostRequest(url string, body []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func main() {
	config, err := readConfig("example-appsettings.json")
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": config.OcpApimSubscriptionKey,
		"Ocp-Apim-Subscription-Region": config.OcpApimSubscriptionRegion,
		"Content-Type": "application/json",
	}

	phrase := "Boutons et voyants du panneau de commande\n"
	body := fmt.Sprintf(`[{"Text":"%s"}]`, phrase)

	// FIRST API CALL
	response, err := makePostRequest(config.AzureTranslateURL+"?to=en", []byte(body), headers)
	if err != nil {
		fmt.Println("Error making first API call:", err)
		return
	}

	var firstResponse TranslateResponse
	err = json.Unmarshal(response, &firstResponse)
	if err != nil {
		fmt.Println("Error parsing first response:", err)
		return
	}

	translatedText := firstResponse.Translations[0].Text
	detectedLanguage := firstResponse.DetectedLanguage.Language

	// SECOND API CALL
	openAiUrl := fmt.Sprintf("%s/openai/deployments/%s/extensions/chat/completions?api-version=2023-06-01-preview", config.YourResourceName, config.YourDeploymentName)
	body2 := fmt.Sprintf(`{
		"temperature": 0,
		"max_tokens": 1000,
		"top_p": 1.0,
		"dataSources": [{
			"type": "AzureCognitiveSearch",
			"parameters": {
				"endpoint": "%s",
				"key": "%s",
				"indexName": "%s"
			}
		}],
		"messages": [{
			"role": "user",
			"content": "%s"
		}]
	}`, config.Endpoint, config.Key, config.IndexName, translatedText)

	headers2 := map[string]string{
		"api-key": config.ApiKey,
		"Content-Type": "application/json",
	}

	response2, err := makePostRequest(openAiUrl, []byte(body2), headers2)
	if err != nil {
		fmt.Println("Error making second API call:", err)
		return
	}

	var secondResponse map[string]interface{}
	err = json.Unmarshal(response2, &secondResponse)
	if err != nil {
		fmt.Println("Error parsing second response:", err)
		return
	}

	choices := secondResponse["choices"].([]interface{})
	aiResponse := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	// THIRD API CALL
	body3 := fmt.Sprintf(`[{"Text":"%s"}]`, aiResponse)
	response3, err := makePostRequest(config.AzureTranslateURL+"?to="+detectedLanguage, []byte(body3), headers)
	if err != nil {
		fmt.Println("Error making third API call:", err)
		return
	}

	var thirdResponse TranslateResponse
	err = json.Unmarshal(response3, &thirdResponse)
	if err != nil {
		fmt.Println("Error parsing third response:", err)
		return
	}

	finalTranslatedText := thirdResponse.Translations[0].Text
	fmt.Println(finalTranslatedText)
}
