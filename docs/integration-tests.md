# Integration Tests

To ensure a working keyvault plugin an integration test suite can be executed with go test.

## Requirements

The azure service principal or identity used to execute the test suite needs to have permissions to create and delete
keyvault resources in the specified resource group.

Create a service principal with Contributor rights on the resource group.
Store the output as a secret in the github repository.

```bash
# create a new service principal with contributor rights
az ad sp create-for-rbac --name "helm-keyvault-integration-tests" --role contributor \
  --scopes /subscriptions/{subscription-id}/resourceGroups/{resource-group} \
  --sdk-auth
```

See https://github.com/marketplace/actions/azure-login#configure-a-service-principal-with-a-secret for more details/.