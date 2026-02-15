# Retro Image Filter Service

A microservices-based image processing pipeline for learning Go, Rust, and Kubernetes. See [PROJECT.md](PROJECT.md) for full architecture and API details.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [minikube](https://minikube.sigs.k8s.io/docs/start/) and [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Go 1.22+](https://go.dev/dl/)
- [Rust / Cargo](https://rustup.rs/)

## Local Development (Docker Compose)

```bash
docker-compose up --build
```

This starts:
- **go-api** on `localhost:8080`
- **rust-worker** (background processor)
- **redis** on `localhost:6379`

## Kubernetes Development (minikube)

```bash
# Start minikube
minikube start

# Deploy all resources
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/

# Access the API
minikube service go-api -n retro-filter
```

## Project Structure

```
.
├── go-service/          # Go HTTP API
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── rust-service/        # Rust image processing worker
│   ├── Dockerfile
│   ├── Cargo.toml
│   └── src/main.rs
├── k8s/                 # Kubernetes manifests
│   ├── namespace.yaml
│   ├── redis-deployment.yaml
│   ├── redis-service.yaml
│   ├── go-deployment.yaml
│   ├── go-service.yaml
│   ├── rust-deployment.yaml
│   └── rust-hpa.yaml
├── docker-compose.yml
├── PROJECT.md           # Architecture and API design
├── QUALITY.md           # Testing, linting, and CI standards
└── CLAUDE.md            # Claude Code interaction rules
```

## Quality Standards

See [QUALITY.md](QUALITY.md) for TDD workflow, linting rules, and CI pipeline details.
