.PHONY: help dev backend-frontend frontend backend k8s-apply k8s-delete clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dev: ## Start local development with docker-compose
	cd compose && docker-compose up --build

dev-backend: ## Start backend only (with SQLite for quick testing)
	cd backend && go run ./cmd/server

dev-frontend: ## Start frontend dev server
	cd frontend && npm install && npm run dev

backend: ## Build and run backend
	cd backend && go build -o ../bin/subnode ./cmd/server && ../bin/subnode

frontend-build: ## Build frontend for production
	cd frontend && npm install && npm run build

k8s-apply: ## Apply Kubernetes manifests
	kubectl apply -f k8s/deployment.yaml

k8s-delete: ## Delete all Kubernetes resources
	kubectl delete -f k8s/deployment.yaml

k8s-port-forward: ## Port forward to frontend
	kubectl port-forward -n subnode svc/frontend 30080:80

k8s-logs-backend: ## View backend logs
	kubectl logs -n subnode -l app=backend -f

k8s-logs-frontend: ## View frontend logs
	kubectl logs -n subnode -l app=frontend -f

test: ## Run tests
	cd backend && go test ./...

lint: ## Run linters
	cd backend && golangci-lint run
	cd frontend && npm run lint

build-all: ## Build both backend and frontend images
	docker build -t subnode-backend:latest -f backend/Dockerfile backend
	docker build -t subnode-frontend:latest -f frontend/Dockerfile frontend

push-all: build-all ## Build and push images
	docker push subnode-backend:latest
	docker push subnode-frontend:latest

clean: ## Clean build artifacts
	rm -rf bin/
	cd frontend && rm -rf node_modules/ dist/
