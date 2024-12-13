const axios = require('axios');
const fs = require('fs');

const config = JSON.parse(fs.readFileSync('config.json', 'utf8'));

const translateText = async (text, toLanguage) => {
    const url = `${config.AzureTranslateURL}?to=${toLanguage}`;
    const headers = {
        'Ocp-Apim-Subscription-Key': config['Ocp-Apim-Subscription-Key'],
        'Ocp-Apim-Subscription-Region': config['Ocp-Apim-Subscription-Region'],
        'Content-Type': 'application/json'
    };
    const body = [{ Text: text }];
    const response = await axios.post(url, body, { headers });
    return response.data;
};

const callOpenAI = async (content) => {
    const url = `${config['YOUR_RESOURCE_NAME']}/openai/deployments/${config['YOUR_DEPLOYMENT_NAME']}/extensions/chat/completions?api-version=2023-06-01-preview`;
    const headers = {
        'api-key': config['api-key'],
        'Content-Type': 'application/json'
    };
    const body = {
        temperature: 0,
        max_tokens: 1000,
        top_p: 1.0,
        dataSources: [
            {
                type: 'AzureCognitiveSearch',
                parameters: {
                    endpoint: config.endpoint,
                    key: config.key,
                    indexName: config.indexName
                }
            }
        ],
        messages: [
            {
                role: 'user',
                content: content
            }
        ]
    };
    const response = await axios.post(url, body, { headers });
    return response.data;
};

const main = async () => {
    const inputText = process.argv[2] || "Boutons et voyants du panneau de commande\n";

    // FIRST API CALL
    const firstResponse = await translateText(inputText, 'en');
    const translatedText = firstResponse[0].translations[0].text;
    const detectedLanguage = firstResponse[0].detectedLanguage.language;

    // SECOND API CALL
    const openaiResponse = await callOpenAI(translatedText);
    const aiResponse = openaiResponse.choices[0].messages[1].content;

    // THIRD API CALL
    const finalResponse = await translateText(aiResponse, detectedLanguage);
    const finalTranslatedText = finalResponse[0].translations[0].text;

    console.log(finalTranslatedText);
};

main().catch(console.error);
