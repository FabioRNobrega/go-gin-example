# 🔑 Configuring Azure Key Vault with Secrets in AKS via CSI Driver (env)

## 1. Prerequisites

- AKS cluster created ✅  
- Azure CLI configured ✅ (`az login`)  
- Key Vault already exists with secrets (`rsa-public-key`, `rsa-private-key`) ✅  

---

## 2. Enable the CSI Driver in the cluster

```bash
az aks enable-addons   --addons azure-keyvault-secrets-provider   --name ping-pong-cluster-name   --resource-group ping-pong-resource-group
```

Check if CSI driver pods are running:

```bash
kubectl get pods -n kube-system | grep csi-secrets-store
```

---

## 3. Managed Identity

Identify the Managed Identity used by the nodepool:

```bash
az vmss list   --resource-group MC_ping-pong-resource-group_ping-pong-cluster-name_eastus   -o table
```

Check assigned IDs:

```bash
az vmss identity show   --resource-group MC_ping-pong-resource-group_ping-pong-cluster-name_eastus   --name aks-agentpool-<id>
```

---

## 4. Grant Key Vault access

Since the Key Vault uses RBAC, grant access via role assignment:

```bash
MSYS_NO_PATHCONV=1 az role assignment create   --role "Key Vault Secrets User"   --assignee <USER_ASSIGNED_IDENTITY_CLIENT_ID>   --scope /subscriptions/<SUBSCRIPTION_ID>/resourceGroups/ping-pong-resource-group/providers/Microsoft.KeyVault/vaults/ping-pong-vault-keys
```

---

## 5. Create SecretProviderClass

File: `secretproviderclass.yaml`

```yaml
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: pingpong-keys
  namespace: pingpong
spec:
  provider: azure
  parameters:
    usePodIdentity: "false"
    useVMManagedIdentity: "true"
    userAssignedIdentityID: "<USER_ASSIGNED_IDENTITY_CLIENT_ID>"
    keyvaultName: ping-pong-vault-keys
    tenantId: <TENANT_ID>
    cloudName: ""
    objects: |
      array:
        - |
          objectName: rsa-public-key
          objectType: secret
        - |
          objectName: rsa-private-key
          objectType: secret
  secretObjects:
  - secretName: pingpong-keys-secret
    type: Opaque
    data:
    - objectName: rsa-public-key
      key: RSA_PUBLIC_KEY
    - objectName: rsa-private-key
      key: RSA_PRIVATE_KEY
```

Apply to the cluster:

```bash
kubectl apply -f secretproviderclass.yaml -n pingpong
```

---

## 6. Mount secrets in Deployment

In the service deployment (e.g., `webservice`), add:

```yaml
        env:
        - name: RSA_PUBLIC_KEY
          valueFrom:
            secretKeyRef:
              name: pingpong-keys-secret
              key: RSA_PUBLIC_KEY
        - name: RSA_PRIVATE_KEY
          valueFrom:
            secretKeyRef:
              name: pingpong-keys-secret
              key: RSA_PRIVATE_KEY
```

---

## 7. Validate

Check if secrets are created:

```bash
kubectl get secrets -n pingpong
```

Check pod logs:

```bash
kubectl logs -n pingpong <webservice-pod>
```

Open secrets on bash

```bash
KUBE_EDITOR=nano kubectl edit secret pingpong-keys-secret -n pingpong
```

---

✅ Now your application receives the keys from Azure Key Vault via **env vars** injected by the **CSI Driver**.
