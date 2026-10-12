# SubNode Kubernetes Deployment

## Prerequisites

- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl configured
- Docker or container runtime for building images

## Quick Start

### 1. Build and Load Images

```bash
# Using minikube
minikube start
eval $(minikube docker-env)

# Build images
make build-all

# For kind
kind load docker-image subnode-backend:latest
kind load docker-image subnode-frontend:latest
```

### 2. Create Namespace and Secrets

```bash
# Create namespace
kubectl create namespace subnode

# Create database secret
kubectl create secret generic subnode-db \
  --namespace=subnode \
  --from-literal=POSTGRES_USER=subnode \
  --from-literal=POSTGRES_PASSWORD=change-me-in-production \
  --from-literal=POSTGRES_DB=subnode

# Create JWT secret
kubectl create secret generic subnode-secrets \
  --namespace=subnode \
  --from-literal=JWT_SECRET=your-secret-key-here
```

### 3. Deploy

```bash
# Apply manifests
kubectl apply -f k8s/deployment.yaml

# Check deployment
kubectl get pods -n subnode
kubectl get services -n subnode
```

### 4. Access the Application

```bash
# For NodePort (local dev)
minikube service frontend -n subnode

# For Ingress (production)
# Add entry to /etc/hosts:
# <cluster-ip> subnode.example.com
```

## Production Deployment

### 1. Configure TLS

```bash
# Generate self-signed certificate (for testing)
openssl req -x509 -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365 -nodes

# Create TLS secret
kubectl create secret tls subnode-tls \
  --namespace=subnode \
  --key=tls.key \
  --cert=tls.crt
```

### 2. Deploy Production Manifest

```bash
kubectl apply -f k8s/deployment.prod.yaml
```

### 3. Configure Ingress

Update `k8s/deployment.prod.yaml` with your domain and TLS certificate.

## Scaling

```bash
# Scale backend
kubectl scale deployment backend -n subnode --replicas=3

# Scale frontend
kubectl scale deployment frontend -n subnode --replicas=3
```

## Monitoring

```bash
# View logs
kubectl logs -n subnode -l app=backend -f
kubectl logs -n subnode -l app=frontend -f

# View events
kubectl get events -n subnode --sort-by='.lastTimestamp'
```

## Cleanup

```bash
# Delete everything
kubectl delete -f k8s/deployment.yaml
kubectl delete namespace subnode
```
