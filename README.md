# helm-keyvault

Helm Plugin to manage [Azure Keyvault](https://azure.microsoft.com/en-us/services/key-vault/) secrets and keys. It allows to safely store chart files either as a secret inside an Azure Keyvault
or as an encrypted file in git.

## Installation

```bash
# download the latest release from github 
mkdir .helm-plugins
cd .helm-plugins
curl -L https://github.com/foryouandyourcustomers/helm-keyvault/releases/latest/download/helm-keyvault_linux_amd64.tar.gz | tar -xzf -
helm plugin install ./keyvault 
```

## Usage

The plugin can be used to manage keys and secrets and it can be used as a downloader plugin for helm charts.

Have a look at the two examples:
- [Deploy a helm chart with keyvault secrets](./docs/deploy-a-helm-chart-with-keyvault-secrets.md)
- [Deploy a helm chart with encrypted files](./docs/deploy-a-helm-chart-with-encrypted-files.md)

## Authentication

The plugin requires to authenticate with Azure. The user, service principal or managed identity used by the plugin needs permissions
on the keyvault(s) to read and/or manage Keyvault secrets and keys.

The plugin supports three authentication methods against azure.

### auth.json
First it checks for the env var `AZURE_AUTH_LOCATION`. If the env var exists it will try to load the
authentication from the given [file](https://docs.microsoft.com/en-us/dotnet/azure/sdk/authentication#mgmt-file). 

### environment variables
If no auth file is given it will try to setup the authentication via environment variables.

The supported environment variables are:

**Client Credentials**

    AZURE_CLIENT_ID
    AZURE_CLIENT_SECRET
    AZURE_TENANT_ID

**Client Certificate**

    AZURE_CERTIFICATE_PATH
    AZURE_CERTIFICATE_PASSWORD
    AZURE_TENANT_ID

**Username Password**:

    AZURE_USERNAME
    AZURE_PASSWORD
    AZURE_CLIENT_ID
    AZURE_TENANT_ID

**MSI**

    AZURE_AD_RESOURCE
    AZURE_CLIENT_ID

### azure cli
Last but not least it will try to login with the local azure cli credentials.
