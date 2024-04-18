#!/bin/bash
az group create --location eastus --resource-group MyResourceGroup

az deployment group create \
  --resource-group MyResourceGroup \
  --template-file Infra/infra.bicep \
  --verbose

