# Azure Container Registry name (change this to your ACR)
ACR_NAME=meuacr
ACR_LOGIN=$(ACR_NAME).azurecr.io
RESOURCE_GROUP=meu-rg
CLUSTER_NAME=pingpong-cluster
NAMESPACE=pingpong

# Services to build and push
SERVICES=webservice pingservice pongservice

# Build and push all Docker images to ACR
docker-build-push:
	@for service in $(SERVICES); do \
		echo "🚀 Building and pushing $$service..."; \
		docker build -t $(ACR_LOGIN)/$$service:latest ./$$service; \
		docker push $(ACR_LOGIN)/$$service:latest; \
	done

# Create namespace (ignore error if it already exists)
create-namespace:
	-kubectl create namespace $(NAMESPACE)

# Apply Kubernetes manifests
kubectl-apply: create-namespace
	kubectl apply -f k8s/ -n $(NAMESPACE)

# Connect kubectl to AKS cluster
aks-connect:
	az aks get-credentials --resource-group $(RESOURCE_GROUP) --name $(CLUSTER_NAME)

# Full deploy: build + push + apply
deploy: docker-build-push kubectl-apply

# Clean up namespace
clean:
	kubectl delete namespace $(NAMESPACE) --ignore-not-found=true

# Local Usage
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down




	