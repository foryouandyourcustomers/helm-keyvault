# helm-keyvault

A small little downloader plugin to retrieve values files from azure keyvault.

This allows to render helm charts with secret values if no other possibility to retrieve secrets
are given, e.g. when bootstrapping a cluster and no secret csi is available.

Here a short example to retrieve the secret `argocd-yaml` in the keyvault `helm-keyvault-test`.

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

Next use the keyvault plugin to retrieve the generated secret during helm execution by defining
keyvault://` as url for the values file.

```bash
# downlaod the latest secret
--values keyvault://helm-keyvault-test.vault.azure.net/secrets/argocd-yaml
# downlaod the secret with a specific verion
--values keyvault://helm-keyvault-test.vault.azure.net/secrets/argocd-yaml/2d6e0430c0724ad1bdc277af8b549c57
```
