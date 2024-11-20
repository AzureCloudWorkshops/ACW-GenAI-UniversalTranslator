require 'json'
require 'net/http'
require 'uri'

def load_config
  JSON.parse(File.read('config.json'))
end

def translate_text(text, to_language, config)
  uri = URI("#{config['AzureTranslateURL']}?to=#{to_language}")
  request = Net::HTTP::Post.new(uri)
  request['Ocp-Apim-Subscription-Key'] = config['Ocp-Apim-Subscription-Key']
  request['Ocp-Apim-Subscription-Region'] = config['Ocp-Apim-Subscription-Region']
  request['Content-Type'] = 'application/json'
  request.body = [{ 'Text' => text }].to_json

  response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: uri.scheme == 'https') do |http|
    http.request(request)
  end

  JSON.parse(response.body)
end

def call_openai(content, config)
  uri = URI("#{config['YOUR_RESOURCE_NAME']}/openai/deployments/#{config['YOUR_DEPLOYMENT_NAME']}/extensions/chat/completions?api-version=2023-06-01-preview")
  request = Net::HTTP::Post.new(uri)
  request['api-key'] = config['api-key']
  request['Content-Type'] = 'application/json'
  request.body = {
    'temperature' => 0,
    'max_tokens' => 1000,
    'top_p' => 1.0,
    'dataSources' => [
      {
        'type' => 'AzureCognitiveSearch',
        'parameters' => {
          'endpoint' => config['endpoint'],
          'key' => config['key'],
          'indexName' => config['indexName']
        }
      }
    ],
    'messages' => [
      {
        'role' => 'user',
        'content' => content
      }
    ]
  }.to_json

  response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: uri.scheme == 'https') do |http|
    http.request(request)
  end

  JSON.parse(response.body)
end

def main
  config = load_config

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

  puts final_translated_text
end

main if __FILE__ == $PROGRAM_NAME
