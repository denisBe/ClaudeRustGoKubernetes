# Developer Setup Guide

Step-by-step instructions to get the Retro Image Filter Service running locally.

## Dependencies

| Tool | Minimum Version | Purpose |
|------|----------------|---------|
| [Docker Desktop](https://www.docker.com/products/docker-desktop/) | Latest | Containers (includes Docker Compose v2) |
| [Go](https://go.dev/dl/) | 1.22+ | Go API service |
| [Rust / Cargo](https://rustup.rs/) | Latest stable | Rust worker service |
| [minikube](https://minikube.sigs.k8s.io/docs/start/) | Latest | Local Kubernetes cluster |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | Latest | Kubernetes CLI |
| [golangci-lint](https://golangci-lint.run/welcome/install/) | Latest | Go linting (see QUALITY.md) |
| [kubeconform](https://github.com/yannh/kubeconform) | Latest | K8s manifest validation (see QUALITY.md) |

## Automated Setup

Scripts that install all dependencies and build the project:

- **Windows (PowerShell, admin):** `.\scripts\setup-windows.ps1`
- **Linux (bash, sudo):** `./scripts/setup-linux.sh`

If you prefer manual installation, follow the sections below.

## Windows Installation

### Using Chocolatey

Install [Chocolatey](https://chocolatey.org/install) first, then run in an **admin** PowerShell:

```powershell
choco install docker-desktop golang rustup.install minikube kubernetes-cli -y
```

After Rust installs, open a **new** terminal so `cargo` is on your PATH.

### Go tools

```powershell
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### kubeconform

Download the latest release from [GitHub](https://github.com/yannh/kubeconform/releases) and place the binary somewhere on your PATH, or:

```powershell
go install github.com/yannh/kubeconform/cmd/kubeconform@latest
```

### Git Bash fix (important)

Git Bash on Windows may hang or break on interactive use due to post-install scripts. To fix, open an **admin** Git Bash and run:

```bash
rm -f /etc/post-install/*.post
```

This removes one-time post-install hooks that can interfere with Git Bash sessions.

### Docker credential fix

If you see `docker-credential-desktop` errors in Git Bash, edit `~/.docker/config.json` and remove the `"credsStore"` key:

```jsonc
{
  // Remove this line:
  // "credsStore": "desktop"
}
```

## Linux Installation

### Debian / Ubuntu

```bash
# Docker
sudo apt-get update
sudo apt-get install -y docker.io docker-compose-v2
sudo usermod -aG docker $USER

# Go (check https://go.dev/dl/ for latest version)
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc

# Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
source "$HOME/.cargo/env"

# kubectl
curl -LO "https://dl.k8s.io/release/$(curl -Ls https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

### Fedora / RHEL

```bash
# Docker
sudo dnf install -y docker docker-compose-plugin
sudo systemctl enable --now docker
sudo usermod -aG docker $USER

# Go, Rust, kubectl, minikube — same as Debian instructions above
```

### Go tools (all distros)

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/yannh/kubeconform/cmd/kubeconform@latest
```

## Verification

Run these commands to confirm everything is installed:

```bash
docker --version
docker compose version
go version
rustc --version && cargo --version
minikube version
kubectl version --client
golangci-lint --version
kubeconform -v
```

All commands should print version info without errors.

## Build and Run

### Docker Compose (local development)

```bash
docker compose build
docker compose up
```

This starts the Go API on `localhost:8080`, the Rust worker, and Redis on `localhost:6379`.

### Kubernetes (minikube)

```bash
minikube start
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/
minikube service go-api -n retro-filter
```
