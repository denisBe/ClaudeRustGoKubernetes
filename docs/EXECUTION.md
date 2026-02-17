# Running the Retro Image Filter Service

This document describes every way to run the project, from quick local development to production-like Kubernetes deployment.

## Architecture Overview

The service has three components:

| Component | Language | Role |
|-----------|----------|------|
| **go-api** | Go | HTTP API + web frontend (serves static files via embed) |
| **rust-worker** | Rust | Picks jobs from Redis queue, applies image filters |
| **redis** | — | Message queue between API and worker |

---

## 1. Run Directly (Go API only)

The fastest way to iterate on the Go API. No Docker needed.

```bash
cd go-service
go run .
```

The server starts on **http://localhost:8081** with the web frontend and all API endpoints.

**What works**: API endpoints, frontend UI, health check.
**What doesn't work**: Job processing (no Redis or Rust worker).

To build a binary instead:

```bash
cd go-service
go build -o go-api .
./go-api
```

### Running tests

```bash
cd go-service
go test ./...
```

---

## 2. Docker Compose (Full Stack)

Runs all three components together. Best for integration testing and end-to-end development.

```bash
docker compose up --build
```

| Service | URL |
|---------|-----|
| Go API + Frontend | http://localhost:8081 |
| Redis | localhost:6379 |

To run in the background:

```bash
docker compose up --build -d
```

To stop:

```bash
docker compose down
```

To rebuild a single service:

```bash
docker compose up --build go-api
```

All methods use port 8081 consistently.

---

## 3. Single Docker Container (Go API only)

Useful for testing the Docker image in isolation.

```bash
# Build
docker build -t retro-filter/go-api:dev ./go-service

# Run
docker run --rm -p 8081:8081 retro-filter/go-api:dev
```

Access at **http://localhost:8081**.

For the Rust worker:

```bash
docker build -t retro-filter/rust-worker:dev ./rust-service
docker run --rm retro-filter/rust-worker:dev
```

---

## 4. VS Code

### Run and Debug

You can run the Go service directly from VS Code:

1. Open `go-service/main.go`
2. Click **Run > Run Without Debugging** (Ctrl+F5) or **Run > Start Debugging** (F5)
3. VS Code will use the Go extension's built-in runner

For a custom launch configuration, create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Go API",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/go-service",
            "env": {
                "REDIS_ADDR": "localhost:6379"
            }
        }
    ]
}
```

### Integrated Terminal

Alternatively, use the VS Code integrated terminal to run any of the methods described in this document (go run, docker compose, etc.).

---

## 5. Kubernetes with Minikube (Production-like)

Closest to a production deployment. Runs all components in a Kubernetes cluster.

### Prerequisites

- minikube installed and running (`minikube start`)
- kubectl configured (`kubectl` should point to minikube)

### Deploy

```bash
# Load local Docker images into minikube
minikube image load retro-filter/go-api:0.1.0
minikube image load retro-filter/rust-worker:0.1.0

# Apply manifests
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/
```

### Access the service

```bash
# Get the URL assigned by minikube
minikube service go-api -n retro-filter --url
```

Open the returned URL in your browser.

### Check status

```bash
# Pods
kubectl get pods -n retro-filter

# Logs
kubectl logs -n retro-filter -l component=go-api
kubectl logs -n retro-filter -l component=rust-worker

# Describe a pod (for debugging)
kubectl describe pod -n retro-filter -l component=go-api
```

### Tear down

```bash
kubectl delete -f k8s/
minikube stop
```

### K8s resources

| Manifest | What it creates |
|----------|----------------|
| `namespace.yaml` | `retro-filter` namespace |
| `go-deployment.yaml` | Go API deployment (1 replica, resource limits, health probes) |
| `go-service.yaml` | NodePort service exposing the API |
| `redis-deployment.yaml` | Redis deployment |
| `redis-service.yaml` | ClusterIP service for Redis |
| `rust-deployment.yaml` | Rust worker deployment |
| `rust-hpa.yaml` | Horizontal Pod Autoscaler for the worker |

---

## Quick Reference

| Method | Command | URL | Use case |
|--------|---------|-----|----------|
| Go directly | `cd go-service && go run .` | http://localhost:8081 | Fast API iteration |
| Docker Compose | `docker compose up --build` | http://localhost:8081 | Full stack, integration testing |
| Single container | `docker run -p 8081:8081 ...` | http://localhost:8081 | Test Docker image |
| VS Code | F5 or Ctrl+F5 | http://localhost:8081 | Debug with breakpoints |
| Kubernetes | `kubectl apply -f k8s/` | `minikube service ...` | Production-like deployment |
