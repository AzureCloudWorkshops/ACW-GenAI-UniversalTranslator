param location string = 'eastus'
param deploymentId string = toLower(uniqueString(resourceGroup().id))


resource storageAccount 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: 'safor${deploymentId}'
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'Storage'
}

resource cognitiveService 'Microsoft.CognitiveServices/accounts@2021-10-01' = {
  name: 'cog${deploymentId}'
  location: location
  kind: 'CognitiveServices'
  sku: {
    name: 'S0'
  }
  properties: {
    publicNetworkAccess: 'Enabled'
  }
}

resource open_ai 'Microsoft.CognitiveServices/accounts@2022-03-01' = {
  name: 'aiS${deploymentId}'
  location: location
  kind: 'OpenAI'
  sku: {
    name: 'S0'
  }
  properties: {
    publicNetworkAccess: 'Enabled'
  }
}

resource openaiDeployment 'Microsoft.CognitiveServices/accounts/deployments@2023-05-01' = {
  name: 'myDeployment'
  sku: {
    name: 'Standard'
    capacity: 239
  }
  parent: open_ai
  properties: {
    model: {
      format: 'OpenAI'
      name: 'gpt-35-turbo'
    }
    raiPolicyName: 'Microsoft.Default'
    versionUpgradeOption: 'OnceCurrentVersionExpired'
  }
}

resource aiSearch 'Microsoft.Search/searchServices@2023-11-01' = {
  name: 'aisearch${deploymentId}'
  location: location
  sku: {
    name: 'basic'
  }
}

resource storageAccountName_default_containerName 'Microsoft.Storage/storageAccounts/blobServices/containers@2023-01-01' = {
  name: 'safor${deploymentId}/default/testcontainer'
  dependsOn: [
    storageAccount
  ]
}
