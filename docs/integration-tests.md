# Integration Tests

To ensure a working keyvault plugin an integration test suite can be executed with go test.

## Requirements

The azure service principal or identity used to execute the test suite needs to have permissions to create and delete
keyvault resources in the specified resource group.

Create a service principal with Contributor rights on the resource group. 

```bash
# create a new service principal with contributor rights
az ad sp create-for-rbac --name "helm-keyvault-integration-tests" --role contributor \
  --scopes /subscriptions/{subscription-id}/resourceGroups/{resource-group}
# ensure the password doesnt expire anytime soon
az ad sp credential reset --name "helm-keyvault-integration-tests" --years 100 
```

See https://github.com/marketplace/actions/azure-login#configure-a-service-principal-with-a-secret for more details.

To execute the integration tests a few environment variables are required
- INTEGRATION: Env var needs to be set to enable the integration test suite
- AZURE_SUBSCRIPTION: The azure subscription containing the resource group
- AZURE_RESOURCE_GROUP: The name of the resource group in which to create the keyvaults for the integration tests
- AZURE_KEYVAULT_NAME: The name of the azure keyvault - will be suffixed with a random string
- AZURE_TENANT_ID: The azure tenant id of the azure ad account or service principal used for the integration tests
- AZURE_CLIENT_ID: The service principal id, if a service principal
- AZURE_CLIENT_SECRET: The service principal client secret
