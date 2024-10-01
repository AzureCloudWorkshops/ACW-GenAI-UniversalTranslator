import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.hc.client5.http.classic.methods.HttpPost;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.CloseableHttpResponse;
import org.apache.hc.client5.http.impl.classic.HttpClients;
import org.apache.hc.core5.http.io.entity.StringEntity;

import java.io.File;
import java.io.IOException;
import java.nio.file.Paths;
import java.util.Scanner;

public class Program {
    private static final ObjectMapper objectMapper = new ObjectMapper();

    public static void main(String[] args) throws IOException {
        // Read configuration settings
        JsonNode config = objectMapper.readTree(new File("example-appsettings.json"));

        // Add headers
        String subscriptionKey = config.get("Ocp-Apim-Subscription-Key").asText();
        String subscriptionRegion = config.get("Ocp-Apim-Subscription-Region").asText();
        String azureTranslateURL = config.get("AzureTranslateURL").asText();
        String deploymentName = config.get("YOUR_DEPLOYMENT_NAME").asText();
        String resourceName = config.get("YOUR_RESOURCE_NAME").asText();
        String apiKey = config.get("api-key").asText();
        String endpoint = config.get("endpoint").asText();
        String key = config.get("key").asText();
        String indexName = config.get("indexName").asText();

        // Create body
        // Change me to ask a new question!
        String body = "[{\"Text\":\"Boutons et voyants du panneau de commande\\n\"}]";

        // FIRST API CALL
        String firstResponseBody = makePostRequest(azureTranslateURL + "?to=en", body, subscriptionKey, subscriptionRegion);

        // Parse the response
        JsonNode firstResponseJson = objectMapper.readTree(firstResponseBody);
        String translatedText = firstResponseJson.get(0).get("translations").get(0).get("text").asText();
        String detectedLanguage = firstResponseJson.get(0).get("detectedLanguage").get("language").asText();

        // SECOND API CALL
        String openAiUrl = resourceName + "/openai/deployments/" + deploymentName + "/extensions/chat/completions?api-version=2023-06-01-preview";
        String body2 = "{ \"temperature\": 0, \"max_tokens\": 1000, \"top_p\": 1.0, \"dataSources\": [{ \"type\": \"AzureCognitiveSearch\", \"parameters\": { \"endpoint\": \"" + endpoint + "\", \"key\": \"" + key + "\", \"indexName\": \"" + indexName + "\" }}], \"messages\": [{ \"role\": \"user\", \"content\": \"" + translatedText + "\" }] }";
        String secondResponseBody = makePostRequest(openAiUrl, body2, apiKey, null);

        // Parse the response
        JsonNode secondResponseJson = objectMapper.readTree(secondResponseBody);
        String aiResponse = secondResponseJson.get("choices").get(0).get("message").get("content").asText();

        // THIRD API CALL
        String body3 = "[{\"Text\":\"" + aiResponse + "\"}]";
        String thirdResponseBody = makePostRequest(azureTranslateURL + "?to=" + detectedLanguage, body3, subscriptionKey, subscriptionRegion);

        // Parse the response
        JsonNode thirdResponseJson = objectMapper.readTree(thirdResponseBody);
        String finalTranslatedText = thirdResponseJson.get(0).get("translations").get(0).get("text").asText();

        System.out.println(finalTranslatedText);
    }

    private static String makePostRequest(String url, String body, String apiKey, String region) throws IOException {
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpPost httpPost = new HttpPost(url);
            httpPost.setEntity(new StringEntity(body));
            httpPost.setHeader("Content-type", "application/json");
            if (apiKey != null) {
                httpPost.setHeader("Ocp-Apim-Subscription-Key", apiKey);
            }
            if (region != null) {
                httpPost.setHeader("Ocp-Apim-Subscription-Region", region);
            }

            try (CloseableHttpResponse response = httpClient.execute(httpPost)) {
                return new Scanner(response.getEntity().getContent()).useDelimiter("\\A").next();
            }
        }
    }
}
