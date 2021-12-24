# Deploy a helm chart with keyvault secrets

Lets deploy a simple [nginx helm chart with basic auth](./chart/).

The helm chart requires a username and password value. These credentials shouldnt be stored
as cleartext in a git repository.

We can use Azure Keyvault with the helm-keyvault downloader plugin to safely store the credentials as a secret.

```bash
# create yaml file containing the credentials
$ cat <<'EOF'>/tmp/credentials.yaml
---
htpasswd:
  username: mysupersecretuser
  password: mysupersecretpassword
EOF

# next create a secret with the files content in azure keyvault
$ helm keyvault secret put --keyvault helm-keyvault-test --secret htpasswd-credentials --file /tmp/credentials.yaml
{"id":"https://helm-keyvault-test.vault.azure.net/secrets/htpasswd-credentials/0f219949d08b459b80c7fcdaf2d56abd","name":"htpasswd-credentials","keyvault":"helm-keyvault-test","version":"0f219949d08b459b80c7fcdaf2d56abd","value":"LS0tCmh0cGFzc3dkOgogIHVzZXJuYW1lOiBteXN1cGVyc2VjcmV0dXNlcgogIHBhc3N3b3JkOiBteXN1cGVyc2VjcmV0cGFzc3dvcmQK"}
```

With the secret in place we can now render the helm chart with the values retrieved from keyvault. To do this we can use the downloader plugin in helm-keyvault. The downloader plugin supports the `keyvault+secret://` uri type.

```bash
# render the helm chart with values retrieved from the Azure Keyvault secret
# the keyvault uri equals the id of the generated secret minus the https.
# if no version is given the downloader plugin will retrieve the latest version from the azure keyvault
$ helm template \
  --values keyvault+secret://helm-keyvault-test.vault.azure.net/secrets/htpasswd-credentials/0f219949d08b459b80c7fcdaf2d56abd \
  example \
  chart/

# deploy the chart to a kubernetes cluster
$ helm upgrade --install \
  --values keyvault+secret://helm-keyvault-test.vault.azure.net/secrets/htpasswd-credentials/0f219949d08b459b80c7fcdaf2d56abd \
  example \
  chart/

# setup a local port forwarding to test the deployment
$ kubectl port-forward service/example-helm-vault-example 8080:80 
```

Open a web browser and open http://localhost:8080.
Try to login with the username `mysupersecretuser` and the password `mysupersecretpassword`.

