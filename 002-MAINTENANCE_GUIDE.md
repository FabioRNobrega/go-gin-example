# Kubernetes Maintenance Guide

This document provides step-by-step instructions for maintaining the **PingPong** Kubernetes cluster after deployment.  
It covers accessing the cluster, working with the `pingpong` namespace, monitoring pods, checking logs, and scaling workloads.

---

## 1. Access the Kubernetes Cluster

First, connect your local `kubectl` to the Azure Kubernetes Service (AKS) cluster:

```bash
az aks get-credentials --resource-group ping-pong-resource-group --name ping-pong-cluster-name

```

✅ This merges the AKS context into your local `~/.kube/config`.

Check cluster nodes:

```bash
kubectl get nodes
```

---

## 2. Access Namespace

List all available namespaces:

```bash
kubectl get ns
```

Switch context to the `pingpong` namespace:

```bash
kubectl config set-context --current --namespace=pingpong
```

✅ From now on, all commands will run inside the `pingpong` namespace unless specified otherwise.

---

## 3. Access Pods

List all pods in the `pingpong` namespace:

```bash
kubectl get pods -n pingpong
```

Get detailed information about a pod (example: `webservice`):

```bash
kubectl describe pod -n pingpong -l app=webservice
```

---

## 4. View Logs

Check logs of the `webservice` pod:

```bash
kubectl logs -n pingpong -l app=webservice
```

Follow logs in real time:

```bash
kubectl logs -f -n pingpong -l app=webservice
```

For `pingservice`:

```bash
kubectl logs -f -n pingpong -l app=pingservice
```

For `pongservice`:

```bash
kubectl logs -f -n pingpong -l app=pongservice
```

---

## 5. Scale Deployments

Scale a deployment up or down.  
Example: scale `webservice` to 3 replicas:

```bash
kubectl scale deployment webservice --replicas=3 -n pingpong
```

Scale `pingservice` to 2 replicas:

```bash
kubectl scale deployment pingservice --replicas=2 -n pingpong
```

Scale `pongservice` to 2 replicas:

```bash
kubectl scale deployment pongservice --replicas=2 -n pingpong
```

Check status of deployments:

```bash
kubectl get deployments -n pingpong
```

---
# Accessing the Frontend via Ingress

After the services and pods are running, you can use **Ingress** to expose the frontend.  
Follow these steps to check the Ingress configuration and find the public URL:

---

### 1. List Ingress resources

```sh
kubectl get ingress -n pingpong
```

Expected output:

```sh
NAME              CLASS   HOSTS   ADDRESS         PORTS   AGE
pingpong-ingress  <none>  *       48.217.134.70   80      2h
```

👉 The **ADDRESS** field (`48.217.134.70` in this example) is the public IP you can use to access the frontend in your browser.

---

### 2. Describe the Ingress for more details

```sh
kubectl describe ingress pingpong-ingress -n pingpong
```

This command will show the configured rules and the backend services they point to, e.g.:

- `/` → `webservice`
- `/ping` → `pingservice`
- `/pong` → `pongservice`

---

### 3. Verify the Services connected to the Ingress

```sh
kubectl get svc -n pingpong
```

Expected output:

```sh
NAME          TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
pingservice   ClusterIP   10.0.11.225    <none>        8081/TCP   2h
pongservice   ClusterIP   10.0.226.67    <none>        8082/TCP   2h
webservice    ClusterIP   10.0.205.123   <none>        80/TCP     2h
```

This confirms that all services are available internally and that the Ingress is correctly routing external traffic to them.

---
