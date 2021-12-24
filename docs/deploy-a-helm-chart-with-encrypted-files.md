# Deploy a helm chart with encrypted files

Lets deploy a simple [nginx helm chart with basic auth](./chart/).

The helm chart requires a username and password value. These credentials shouldnt be stored
as cleartext in a git repository.

If you want to store the credentials in git you can use Azure Keyvault keys to safely encrypt and decrypt files.
The encrypted files can safely be stored in git as decryption is only possible with access to the Azure keyvault and the key stored in it.

```bash
# create yaml file containing the credentials
$ cat <<'EOF'>/tmp/credentials.yaml
---
htpasswd:
  username: myawesomeuser
  password: myevenmoreawesomepassword
EOF

# next create a 4096 bit key in azure keyvault which can be used to encrypt the file
$ helm keyvault key create --keyvault helm-keyvault-test --key htpasswd-credentials
{"kid":"https://helm-keyvault-test.vault.azure.net/keys/htpasswd-credentials/ba28ad7ebb7f4f668a0d4561d9e40e02","name":"htpasswd-credentials","keyvault":"helm-keyvault-test","version":"ba28ad7ebb7f4f668a0d4561d9e40e02"}

# now use the key to encrypt the credentials.yaml file
$ helm keyvault file encrypt --keyvault helm-keyvault-test --key htpasswd-credentials --file /tmp/credentials.yaml
```

The `file encrypt` command creates a new file besides the credentials.yaml file, suffixed with `.enc`. This file contains the encrypted data and the key information to decrypt the file again.
The encrypted file can be safely stored in git.

```bash
$ cat /tmp/credentials.yaml.enc 
{
 "kid": "https://helm-keyvault-test.vault.azure.net/keys/htpasswd-credentials/ba28ad7ebb7f4f668a0d4561d9e40e02",
 "chunks": [
  "nhMVxN2tRzzmOSHXX-yh580ZoYUYKmlADpQjXvXI94VbLBkzn8Ap2_ft3ZbxIjC9U_TcQ15-SC7pLf5441j3sUGPQKbysmvevjJ_yDS5ZpvD_tuTNtPAlZvsVYNBXBr6N6ClorLRr8VXAgc4zHV7flGndTVImjyR35qdtINqDuxoobpT5TjZfRxRf5Dgxt3GqkrqaJxCxv1TkFL_9goOg3yBXMDFKor7AucAAZ-Rqo9LsqVwKcoKjUAHW939lH6fG7AuaFIy_owv4_86KYr6zxuNp2PqeJbjyeNCn-cBY3reMFHNcnBVKwzUOd_nCf-EB_iaVtpo8ZOECjPglxcWKaIX5M1cylUAFgQ-7q_YBpQqc0IQKN7m6ki9dThdZEDWhdLsTu0VLzG-6dswmYkpFK7K35qJOzH2AEolxUoXi57eBZ5lcCwasQN4DO_ojXRSq-T-8PQPU9S1WWpBAbopK_kEEgbm-JYWJeSRRTo1x_LRoY74xg7zVIVwcBmBNDuowQ3GvqhW-vb3TjwhrEUGEDbK0TDGdE817CQvER7yR_1vPhmGeIkOEqn3XG4wJNv1NeCqz56QiTllSLANMvvKU5bDFfnK5WOGcB7LEhWpxprDsKwb5Z_ayFSF_A7r6fwGqHPHNW4tR3xhlVq2YTDZI8w1xRbXlk4CUdDD4RjDtCY"
 ],
 "lastmodified": "2021-12-24T05:30:27+01:0"
} 
```

To render and deploy the helm chart with the encrypted file you can use the helm downloader plugin. It supports the `keyvault+file://` uri type.


```bash
# render the helm chart with values from the encrypted file. The file is decrypted during execution and print to stdout.
helm template \
  --values keyvault+file:///tmp/credentials.yaml.enc \
  example \
  chart/

# deploy the chart to a kubernetes cluster
helm upgrade --install \
  --values keyvault+file:///tmp/credentials.yaml.enc \
  example \
  chart/

# setup a local port forwarding to test the deployment
kubectl port-forward service/example-helm-vault-example 8080:80 
```

Open a web browser and open http://localhost:8080.
Try to login with the username `myawesomeuser` and the password `myevenmoreawesomepassword`.

