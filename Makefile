# Azure Container Registry config
ACR_NAME=pingpongacr333
ACR_LOGIN=$(ACR_NAME).azurecr.io
RESOURCE_GROUP=ping-pong-resource-group
CLUSTER_NAME=ping-pong-cluster-name
NAMESPACE=pingpong

# Services to build and push
SERVICES=webservice pingservice pongservice

# Build and push all Docker images (production) to ACR
docker-build-push:
	@for service in $(SERVICES); do \
		echo "Building and pushing $$service..."; \
		docker build -t $(ACR_LOGIN)/$$service:latest ./$$service; \
		docker push $(ACR_LOGIN)/$$service:latest; \
	done

# Create namespace (ignore error if it already exists)
create-namespace:
	-kubectl create namespace $(NAMESPACE)

# Apply Kubernetes manifests
kubectl-apply: 
	-kubectl apply -f k8s/ -n $(NAMESPACE)

# Connect kubectl to AKS cluster
aks-connect:
	az aks get-credentials --resource-group $(RESOURCE_GROUP) --name $(CLUSTER_NAME)

# Full deploy: build + push + apply + rollout
deploy: docker-build-push kubectl-apply rollout-restart

# Restart all deployments in the namespace
rollout-restart:
	kubectl rollout restart deployment -n $(NAMESPACE)

# Clean up namespace
clean:
	kubectl delete namespace $(NAMESPACE) --ignore-not-found=true

# Local development
docker-run:
	docker-compose up --build

docker-down:
	docker-compose down
