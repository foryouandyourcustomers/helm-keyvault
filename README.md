# helm-keyvault

A small little downloader plugin to retrieve values files from azure keyvault.

This allows to render helm charts with secret values if no other possibility to retrieve secrets
are given, e.g. when bootstrapping a cluster and no secret csi is available.

Here a short example to retrieve the secret `argocd-yaml` in the keyvault `helm-keyvault-test`.

## Installation

Install the plugin with `helm plugin install`

```bash
# install with direct download
helm plugin install https://github.com/foryouandyourcustomers/helm-keyvault/releases/download/0.1.0/helm-keyvault_Linux_x86_64.tar.gz

# if you receive an error message 'fatal: repository github.com/foryouandyourcustomers/helm-keyvault/releases/download/0.0.2/helm-keyvault_Linux_x86_64.tar.gz/ not found 
# your helm cli can't handle tar downloads (should be fixed in helm cli v3.8!). You need to install the plugin manually
```

## Authentication

The plugin supports three authentication methods against azure.


### auth.json
First it checks for the env var `AZURE_AUTH_LOCATION`. If the env var exists it will try to load the
authentication from the given [file](https://docs.microsoft.com/en-us/dotnet/azure/sdk/authentication#mgmt-file). 

### environment variables
If no auth file is given it will try to setup the authentication via environment variables.
The supported environment variables are:

Client Credentials - Specify the env vars:

    AZURE_CLIENT_ID
    AZURE_CLIENT_SECRET
    AZURE_TENANT_ID

Client Certificate - Specify the env vars:

    AZURE_CERTIFICATE_PATH
    AZURE_CERTIFICATE_PASSWORD
    AZURE_TENANT_ID

Username Password - Specify the env vars:

    AZURE_USERNAME
    AZURE_PASSWORD
    AZURE_CLIENT_ID
    AZURE_TENANT_ID

MSI - specify the env vars:

    AZURE_AD_RESOURCE
    AZURE_CLIENT_ID

### azure cli
Last but not least it will try to login with the local azure cli credentials.

## Usage

First, create a base64 encoded values.yaml file in the azure keyvault
```bash
# first create an azure keyvault secret
yaml=$(cat <<'EOF' | base64
argocd:
  git:
    sshkey: ssh-rsa mysupersecretprivatersarepositorykey
  EOF
)

az keyvault secret set --name argocd.yaml --vault-name helm-keyvault-test --value $yaml --encoding base64
```

You can also use the helm-keyvault utility to write the secret
```bash
cat <<'EOF' > /tmp/values.yaml
argocd:
  git:
    sshkey: ssh-rsa mysupersecretprivatersarepositorykey
EOF

helm keyvault secret put --file /tmp/values.yaml --id keyvault+secret://helm-keyvault-test/argocd-yaml
```

Next use the keyvault plugin to retrieve the generated secret during helm execution by defining
keyvault+secret://` as url for the values file.

```bash
# downlaod the latest secret
--values keyvault+secret://helm-keyvault-test/argocd-yaml
# downlaod the secret with a specific verion
--values keyvault+secret://helm-keyvault-test/argocd-yaml/2d6e0430c0724ad1bdc277af8b549c57
```
