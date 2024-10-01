import requests
import json

def load_config():
    with open('config.json') as config_file:
        return json.load(config_file)

def translate_text(text, to_language, config):
    url = f"{config['AzureTranslateURL']}?to={to_language}"
    headers = {
        'Ocp-Apim-Subscription-Key': config['Ocp-Apim-Subscription-Key'],
        'Ocp-Apim-Subscription-Region': config['Ocp-Apim-Subscription-Region'],
        'Content-Type': 'application/json'
    }
    body = [{'Text': text}]
    response = requests.post(url, headers=headers, json=body)
    return response.json()

def call_openai(content, config):
    url = f"{config['YOUR_RESOURCE_NAME']}/openai/deployments/{config['YOUR_DEPLOYMENT_NAME']}/extensions/chat/completions?api-version=2023-06-01-preview"
    headers = {
        'api-key': config['api-key'],
        'Content-Type': 'application/json'
    }
    body = {
        'temperature': 0,
        'max_tokens': 1000,
        'top_p': 1.0,
        'dataSources': [
            {
                'type': 'AzureCognitiveSearch',
                'parameters': {
                    'endpoint': config['endpoint'],
                    'key': config['key'],
                    'indexName': config['indexName']
                }
            }
        ],
        'messages': [
            {
                'role': 'user',
                'content': content
            }
        ]
    }
    response = requests.post(url, headers=headers, json=body)
    return response.json()

def main():
    config = load_config()

    # Change me to ask a new question!
    input_text = "Boutons et voyants du panneau de commande\n"

    # FIRST API CALL
    first_response = translate_text(input_text, 'en', config)
    translated_text = first_response[0]['translations'][0]['text']
    detected_language = first_response[0]['detectedLanguage']['language']

    # SECOND API CALL
    openai_response = call_openai(translated_text, config)
    ai_response = openai_response['choices'][0]['messages'][1]['content']

    # THIRD API CALL
    final_response = translate_text(ai_response, detected_language, config)
    final_translated_text = final_response[0]['translations'][0]['text']

    print(final_translated_text)

if __name__ == "__main__":
    main()
