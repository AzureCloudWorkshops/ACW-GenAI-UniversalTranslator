#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <curl/curl.h>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

struct Config {
    std::string subscriptionKey;
    std::string subscriptionRegion;
    std::string azureTranslateURL;
    std::string deploymentName;
    std::string resourceName;
    std::string apiKey;
    std::string endpoint;
    std::string key;
    std::string indexName;
};

Config loadConfig(const std::string& filename) {
    std::ifstream file(filename);
    json configJson;
    file >> configJson;

    Config config;
    config.subscriptionKey = configJson["Ocp-Apim-Subscription-Key"];
    config.subscriptionRegion = configJson["Ocp-Apim-Subscription-Region"];
    config.azureTranslateURL = configJson["AzureTranslateURL"];
    config.deploymentName = configJson["YOUR_DEPLOYMENT_NAME"];
    config.resourceName = configJson["YOUR_RESOURCE_NAME"];
    config.apiKey = configJson["api-key"];
    config.endpoint = configJson["endpoint"];
    config.key = configJson["key"];
    config.indexName = configJson["indexName"];

    return config;
}

size_t WriteCallback(void* contents, size_t size, size_t nmemb, void* userp) {
    ((std::string*)userp)->append((char*)contents, size * nmemb);
    return size * nmemb;
}

std::string makePostRequest(const std::string& url, const std::string& body, const Config& config, bool isAzure = true) {
    CURL* curl;
    CURLcode res;
    std::string readBuffer;

    curl_global_init(CURL_GLOBAL_DEFAULT);
    curl = curl_easy_init();
    if (curl) {
        struct curl_slist* headers = nullptr;
        headers = curl_slist_append(headers, "Content-Type: application/json");
        if (isAzure) {
            headers = curl_slist_append(headers, ("Ocp-Apim-Subscription-Key: " + config.subscriptionKey).c_str());
            headers = curl_slist_append(headers, ("Ocp-Apim-Subscription-Region: " + config.subscriptionRegion).c_str());
        } else {
            headers = curl_slist_append(headers, ("api-key: " + config.apiKey).c_str());
        }

        curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &readBuffer);

        res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            fprintf(stderr, "curl_easy_perform() failed: %s\n", curl_easy_strerror(res));
        }

        curl_easy_cleanup(curl);
        curl_slist_free_all(headers);
    }

    curl_global_cleanup();
    return readBuffer;
}

int main(int argc, char* argv[]) {
    Config config = loadConfig("config.json");

    std::string inputText = "Boutons et voyants du panneau de commande\n";
    if (argc > 1) {
        inputText = argv[1];
    }

    // FIRST API CALL
    json body = json::array({ {{"Text", inputText}} });
    std::string firstResponseBody = makePostRequest(config.azureTranslateURL + "?to=en", body.dump(), config);

    json firstResponseJson = json::parse(firstResponseBody);
    std::string translatedText = firstResponseJson[0]["translations"][0]["text"];
    std::string detectedLanguage = firstResponseJson[0]["detectedLanguage"]["language"];

    // SECOND API CALL
    std::string openAiUrl = config.resourceName + "/openai/deployments/" + config.deploymentName + "/extensions/chat/completions?api-version=2023-06-01-preview";
    json body2 = {
        {"temperature", 0},
        {"max_tokens", 1000},
        {"top_p", 1.0},
        {"dataSources", json::array({
            {
                {"type", "AzureCognitiveSearch"},
                {"parameters", {
                    {"endpoint", config.endpoint},
                    {"key", config.key},
                    {"indexName", config.indexName}
                }}
            }
        })},
        {"messages", json::array({
            {
                {"role", "user"},
                {"content", translatedText}
            }
        })}
    };
    std::string secondResponseBody = makePostRequest(openAiUrl, body2.dump(), config, false);

    json secondResponseJson = json::parse(secondResponseBody);
    std::string aiResponse = secondResponseJson["choices"][0]["messages"][1]["content"];

    // THIRD API CALL
    json body3 = json::array({ {{"Text", aiResponse}} });
    std::string thirdResponseBody = makePostRequest(config.azureTranslateURL + "?to=" + detectedLanguage, body3.dump(), config);

    json thirdResponseJson = json::parse(thirdResponseBody);
    std::string finalTranslatedText = thirdResponseJson[0]["translations"][0]["text"];

    std::cout << finalTranslatedText << std::endl;

    return 0;
}
