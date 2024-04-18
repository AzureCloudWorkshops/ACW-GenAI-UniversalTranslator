using System;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.Extensions.Configuration;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Text;
using ACW_GenAI_UniversalTranslator;
using System.Text.Json.Nodes;

class Program
{
    static readonly HttpClient client = new();
    static readonly HttpClient client2 = new();
    static readonly HttpClient client3 = new();
    static IConfiguration? Configuration { get; set; }

    static async Task Main()
    {
        var builder = new ConfigurationBuilder()
            .SetBasePath(AppDomain.CurrentDomain.BaseDirectory)
            .AddJsonFile("appsettings.json");

        Configuration = builder.Build();

        // Add headers
        client.DefaultRequestHeaders.Add("Ocp-Apim-Subscription-Key", Configuration["Ocp-Apim-Subscription-Key"]);
        client.DefaultRequestHeaders.Add("Ocp-Apim-Subscription-Region", Configuration["Ocp-Apim-Subscription-Region"]);
        // Create body
        // Change me to ask a new question!
        var body = new[] { new { Text = "Boutons et voyants du panneau de commande\n" } };
        var content = new StringContent(JsonSerializer.Serialize(body), Encoding.UTF8, "application/json");

        // FIRST API CALL
        HttpResponseMessage response = await client.PostAsync($"{Configuration["AzureTranslateURL"]}?to=en", content);
        // Read and deserialize the response
        var responseString = await response.Content.ReadAsStringAsync();
        responseString = responseString.Substring(1, responseString.Length - 2); 
        TranslateResponseModel? firstResponseBody = JsonSerializer.Deserialize<TranslateResponseModel>(responseString);

        // SECOND API CALL
        var endpoint = Configuration["endpoint"];
        var key = Configuration["key"];
        var indexName = Configuration["indexName"];
        var openAiUrl =
            $"{Configuration["YOUR_RESOURCE_NAME"]}/openai/deployments/{Configuration["YOUR_DEPLOYMENT_NAME"]}/extensions/chat/completions?api-version=2023-06-01-preview";

        var body2 = new
        {
            temperature = 0,
            max_tokens = 1000,
            top_p = 1.0,
            dataSources = new[]
            {
                new
                {
                    type = "AzureCognitiveSearch",
                    parameters = new
                    {
                        endpoint = endpoint,
                        key = key,
                        indexName = indexName
                    }
                }
            },
            messages = new[]
            {
                new
                {
                    role = "user",
                    content = firstResponseBody?.translations?[0].text
                }
            }
        };
        
        
        client2.DefaultRequestHeaders.Add("api-key", Configuration["api-key"]);
        var content2 = new StringContent(JsonSerializer.Serialize(body2), Encoding.UTF8, "application/json");

        HttpResponseMessage response2 = await client2.PostAsync(openAiUrl, content2);
       
        string responseString2 = await response2.Content.ReadAsStringAsync();
        var choices = ParseMessages(responseString2);
        var answers = JsonNode.Parse(choices);
        var aiResponse = answers?[1]?["content"];
        
        
        // Add headers
        client3.DefaultRequestHeaders.Add("Ocp-Apim-Subscription-Key", Configuration["Ocp-Apim-Subscription-Key"]);
        client3.DefaultRequestHeaders.Add("Ocp-Apim-Subscription-Region", Configuration["Ocp-Apim-Subscription-Region"]);
        // Create body
        
        var body3 = new[] { new { Text = $"{aiResponse}" } };
        var content3 = new StringContent(JsonSerializer.Serialize(body3), Encoding.UTF8, "application/json");

        // THIRD API CALL
        HttpResponseMessage response3 = await client3.PostAsync($"{Configuration["AzureTranslateURL"]}?to={firstResponseBody?.detectedLanguage?.language}", content3);
        // Read and deserialize the response
        var responseString3 = await response3.Content.ReadAsStringAsync();
        responseString = responseString3.Substring(1, responseString3.Length - 2); 
        TranslateResponseModel? thirdResponseBody = JsonSerializer.Deserialize<TranslateResponseModel>(responseString);
    
        Console.WriteLine(thirdResponseBody?.translations?[0].text);
    }


    // Using this to parse the json and kick out the messages objects
    private static string ParseMessages(string json)
    {
        string keyToFind = "\"messages\":";
        int startIndex = json.IndexOf(keyToFind);

        if (startIndex == -1)
        {
            return "No 'messages' key found in the JSON string.";
        }

        // Move the start index to the start of the array
        startIndex += keyToFind.Length;
        int arrayDepth = 0;
        int i = startIndex;

        // Determine the array bounds
        for (; i < json.Length; i++)
        {
            if (json[i] == '[')
            {
                arrayDepth++;
            }
            else if (json[i] == ']')
            {
                arrayDepth--;
                if (arrayDepth == 0) // Found the end of the array
                {
                    break;
                }
            }
        }

        if (i == json.Length)
        {
            return "Malformed JSON: Array not properly closed.";
        }

        return json.Substring(startIndex, i - startIndex + 1);
    }
}