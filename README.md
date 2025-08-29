# Go Gin Ping Pong | Fábio R. Nóbrega  

This project is a simple **microservice playground** built with [Gin Web Framework](https://github.com/gin-gonic/gin) in Go 1.23.  
The goal is to demonstrate how to structure multiple services (Ping, Pong, and Web UI) with **Docker Compose**, **htmx**, **Shoelace**, and **Air** (for hot reload).  

We use:  
- Go: 1.23  
- Gin: v1.10.1  
- Docker + Docker Compose  
- Air for live reload  
- htmx + Shoelace for frontend  

---

## Table of contents

* [Install](#install)  
* [Usage](#usage)  
* [Architecture](#architecture)  
* [Troubleshooting](#troubleshooting)  
* [Git Guideline](#git-guideline)  
* [AZURE Kubernetes with ACR](#azure-kubernetes-with-acr)

---

## Install

Clone the repo and cd into the project:

```bash
git clone https://github.com/FabioRNobrega/go-gin-example.git
cd go-gin-example
```

Make sure Docker and Docker Compose are installed.  

Build and run the stack with:

```bash
make docker-up
```

This will start **3 services** with hot reload:  
- **webservice** → UI + API Gateway (port 8080)  
- **pingservice** → returns 🏓 "Hello World, I'm the ping service" (port 8081)  
- **pongservice** → returns 🏓 "Hello World, I'm the pong service" (port 8082)  

Stop everything with:

```bash
make docker-down
```

---

## Usage

Access the UI at:

```
http://localhost:8080/
```

On the screen you’ll see:
- **Blue side** → Ping button  
- **Green side** → Pong button  
- A central ball 🏐 that moves left/right depending on which button you click.  
- The ball displays the message returned by the microservice.  

Example endpoints (if called directly):  
- `GET http://localhost:8080/call-ping` → forwards request to `pingservice`  
- `GET http://localhost:8080/call-pong` → forwards request to `pongservice`  

---

## Architecture

### Service Relationship Diagram  

```mermaid
flowchart LR
    user([🌍 User<br/>Browser])

    subgraph Azure[Azure AKS Cluster]
        subgraph NS[Namespace: pingpong]
            ingress[Ingress Controller<br/>NGINX]
            websvc[(Service: webservice<br/>ClusterIP:80→Pod:8080)]
            pingsvc[(Service: pingservice<br/>ClusterIP:8081)]
            pongsvc[(Service: pongservice<br/>ClusterIP:8082)]

            web([Pod: WebService - GO APP])
            ping([Pod: PingService - GO APP])
            pong([Pod: PongService - GO APP])
        end
    end

    acr[(Azure Container Registry<br/>pingpongacr333)]

    %% Connections
    user --> ingress --> websvc --> web
    web --> pingsvc --> ping
    web --> pongsvc --> pong
    acr --> Azure
```

+ Ingress: Entry point that exposes the application to the internet with a stable external IP.

+ Service: Provides a stable network endpoint inside the cluster, forwarding traffic to the right Pods.

+ Pod: The actual running instance of your application (container with your Go app).

### Flow:

```mermaid

flowchart LR
    user([User<br/>Browser])
    ingress([Ingress<br/>Public Entry Point])
    websvc([WebService Service])
    webpod([WebService Pod<br/>UI + API Gateway])
    pingsvc([PingService Service])
    pingpod([PingService Pod])
    pongsvc([PongService Service])
    pongpod([PongService Pod])

    user --> ingress --> websvc --> webpod
    webpod --> pingsvc --> pingpod
    webpod --> pongsvc --> pongpod
``` 

---

## Troubleshooting

- If you cannot access the UI via `localhost:8080`, verify that all containers are up with `docker ps`.  
- Inside Docker, services resolve each other by **service name** (`pingservice:8081`, `pongservice:8082`).  
- If hot reload fails, ensure `air` is installed and available in the container.  

---

## Git Guideline

Create your branches and commits using English and follow this guideline:

### Branches
- Feature:  `feat/branch-name`  
- Hotfix: `hotfix/branch-name`  
- POC: `poc/branch-name`  

### Commit prefixes
- Chore: `chore(context): message`  
- Feat: `feat(context): message`  
- Fix: `fix(context): message`  
- Refactor: `refactor(context): message`  
- Tests: `tests(context): message`  
- Docs: `docs(context): message`  

## AZURE Kubernetes with ACR

We have prepared detailed documentation to guide you:  

- [Deployment Guide](./001-DEPLOYMENT_AZURE_GUIDE.md) – Step-by-step instructions on how to deploy Kubernetes on Azure with ACR.  
- [Maintenance Guide](./002-MAINTENANCE_GUIDE.md) – How to access the cluster, manage namespaces, view logs, and scale your services.  
- [Nodes Explanation](./003-NODES.md)  - Understanding how nodes and nodes pools works on ASK.
- [AKS Key Vault Setup](./004-AKS_KEYVAULT_SETUP.md) - How to implement secret with AKS Key Vault on K8S.
