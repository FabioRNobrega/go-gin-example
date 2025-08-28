# Kubernetes Deployment on Azure - Step by Step Guide

This document explains how to deploy a Go microservices project on **Azure Kubernetes Service (AKS)** with **Azure Container Registry (ACR)**, Ingress, Secrets/Key Vault, and additional services such as PostgreSQL and Redis.

---

## 1. Create Kubernetes Cluster on Azure

1. Go to the [Azure Portal](https://portal.azure.com/).
2. Search for **Kubernetes Services** and click **Create**.
3. Fill in:
   - **Cluster name** (e.g. `ping-pong-cluster`)
   - **Resource group** (e.g. `ping-pong-resource-group`)
   - **Region**: choose the closest region to your users.
   - **Node pool**: select VM size (start small, e.g. `Standard_B2s`).
   - **Authentication**: enable RBAC, OIDC, and workload identity.
   - **Networking**: default (Azure CNI).
   - **Monitoring**: disable for dev/test, enable for production.
4. Click **Review + Create**.

After provisioning, connect locally:

```sh
az aks get-credentials --resource-group ping-pong-resource-group --name ping-pong-cluster
kubectl get nodes
```

---

## 2. Kubernetes Resources (Manifests)

Inside the project, create a `k8s/` folder with the following manifests:

- **Namespace**

  ```yaml
  apiVersion: v1
  kind: Namespace
  metadata:
    name: pingpong
  ```

- **Deployments & Services**
  - `pingservice-deployment.yaml`
  - `pongservice-deployment.yaml`
  - `webservice-deployment.yaml`

Each includes a `Deployment` (pods) and a `Service` (ClusterIP).

- **Ingress**

  ```yaml
  apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
    name: pingpong-ingress
    namespace: pingpong
    annotations:
      kubernetes.io/ingress.class: nginx
  spec:
    rules:
      - http:
          paths:
            - path: /
              pathType: Prefix
              backend:
                service:
                  name: webservice
                  port:
                    number: 80
  ```

Apply everything:

```sh
kubectl apply -f k8s/ -n pingpong
```

---

## 3. Connect to Azure Container Registry (ACR) & Push Images

1. Create an ACR:

   ```sh
   az acr create --resource-group ping-pong-resource-group --name pingpongacr333 --sku Standard
   ```

2. Log in:

   ```sh
   az acr login --name pingpongacr333
   ```

3. Build & push images for each service:

   ```sh
   docker build -t pingpongacr333.azurecr.io/webservice:latest ./webservice
   docker push pingpongacr333.azurecr.io/webservice:latest

   docker build -t pingpongacr333.azurecr.io/pingservice:latest ./pingservice
   docker push pingpongacr333.azurecr.io/pingservice:latest

   docker build -t pingpongacr333.azurecr.io/pongservice:latest ./pongservice
   docker push pingpongacr333.azurecr.io/pongservice:latest
   ```

4. Verify:

   ```sh
   az acr repository list --name pingpongacr333 --output table
   ```

---

## 4. Secrets & Azure Key Vault

For secrets management you have two options:

### Option A - Kubernetes Secrets

```sh
kubectl create secret generic db-credentials \
  --from-literal=POSTGRES_USER=admin \
  --from-literal=POSTGRES_PASSWORD=secret123 \
  -n pingpong
```

Pods can access via `envFrom.secretRef`.

### Option B - Azure Key Vault + CSI Driver

1. Enable **Key Vault CSI Driver** when creating the AKS cluster.
2. Create a Key Vault:

   ```sh
   az keyvault create --name pingpong-kv --resource-group ping-pong-resource-group --location eastus
   ```

3. Add secrets (DB credentials, API keys, etc).
4. Mount secrets into pods via CSI driver (`SecretProviderClass` manifest).

---

## 5. Databases & Redis

### PostgreSQL (Azure Database for PostgreSQL Flexible Server)

```sh
az postgres flexible-server create \
  --name pingpong-db \
  --resource-group ping-pong-resource-group \
  --location eastus \
  --admin-user adminuser \
  --admin-password SuperSecret123!
```

Store connection string in Kubernetes secret or Key Vault.

### Redis (Azure Cache for Redis)

```sh
az redis create \
  --name pingpong-redis \
  --resource-group ping-pong-resource-group \
  --location eastus \
  --sku Basic --vm-size C0
```

Connection string is stored in secrets.

---

## Final Workflow

1. **Build & push** Docker images → ACR.  
2. **Deploy manifests** → AKS.  
3. **Ingress** provides external IP for the webservice.  
4. **Secrets/Key Vault** store sensitive data.  
5. **Postgres + Redis** provisioned on Azure and connected via environment variables.  

You can now access your app using the ingress public IP or a DNS name if configured.

---
