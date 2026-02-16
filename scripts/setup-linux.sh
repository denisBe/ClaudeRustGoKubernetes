#!/usr/bin/env bash
#
# Installs all dependencies for the Retro Image Filter Service on Linux.
# Supports Debian/Ubuntu and Fedora/RHEL.
# Run with: sudo ./scripts/setup-linux.sh
#

set -euo pipefail

# --- Root check ---
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run with sudo: sudo $0" >&2
    exit 1
fi

# --- Helpers ---
step()  { echo -e "\n\033[36m==> $1\033[0m"; }
ok()    { echo -e "    \033[32m[OK]\033[0m $1"; }
skip()  { echo -e "    \033[33m[SKIP]\033[0m $1"; }
fail()  { echo -e "    \033[31m[FAIL]\033[0m $1"; }

command_exists() { command -v "$1" &>/dev/null; }

# --- Detect distro ---
step "Detecting Linux distribution"
if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO_FAMILY=""
    case "$ID" in
        ubuntu|debian|pop|linuxmint) DISTRO_FAMILY="debian" ;;
        fedora|rhel|centos|rocky|alma) DISTRO_FAMILY="fedora" ;;
        *) ;;
    esac
fi

if [ -z "${DISTRO_FAMILY:-}" ]; then
    echo "Unsupported distribution. This script supports Debian/Ubuntu and Fedora/RHEL."
    echo "You can still install the dependencies manually — see SETUP.md."
    exit 1
fi
ok "Detected: $PRETTY_NAME ($DISTRO_FAMILY family)"

# --- Docker ---
step "Installing Docker"
if command_exists docker; then
    skip "Docker already installed: $(docker --version)"
else
    if [ "$DISTRO_FAMILY" = "debian" ]; then
        apt-get update
        apt-get install -y docker.io docker-compose-v2
    else
        dnf install -y docker docker-compose-plugin
    fi
    systemctl enable --now docker
    ok "Docker installed"
fi

# Add current user to docker group (takes effect on next login)
if ! groups "${SUDO_USER:-$USER}" | grep -q docker; then
    usermod -aG docker "${SUDO_USER:-$USER}"
    ok "Added ${SUDO_USER:-$USER} to docker group (log out and back in to take effect)"
fi

# --- Go ---
step "Installing Go"
GO_VERSION="1.26.0"
install_go() {
    rm -f /tmp/go.tar.gz
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
}

if [ -x /usr/local/go/bin/go ]; then
    CURRENT_GO=$(/usr/local/go/bin/go version | grep -oP 'go\K[0-9]+\.[0-9]+(\.[0-9]+)?')
    if printf '%s\n%s\n' "$GO_VERSION" "$CURRENT_GO" | sort -V | head -1 | grep -qx "$GO_VERSION"; then
        skip "Go already up to date: go$CURRENT_GO"
    else
        echo "    Go $CURRENT_GO found, upgrading to $GO_VERSION..."
        install_go
        ok "Go upgraded from $CURRENT_GO to $GO_VERSION"
    fi
else
    install_go
    ok "Go $GO_VERSION installed"
fi

# Ensure Go is on PATH for the installing user (both .bashrc and .profile)
REAL_USER="${SUDO_USER:-$USER}"
REAL_HOME="/home/$REAL_USER"
for rcfile in "$REAL_HOME/.bashrc" "$REAL_HOME/.profile"; do
    if ! grep -q '/usr/local/go/bin' "$rcfile" 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> "$rcfile"
        ok "Added Go to PATH in $(basename "$rcfile")"
    fi
done
export PATH=$PATH:/usr/local/go/bin

# --- Rust ---
step "Installing Rust"
if command_exists rustc; then
    skip "Rust already installed: $(rustc --version)"
else
    # Install as the real user, not root
    sudo -u "${SUDO_USER:-$USER}" bash -c \
        'curl --proto "=https" --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y'
    ok "Rust installed"
fi

# --- kubectl ---
step "Installing kubectl"
if command_exists kubectl; then
    skip "kubectl already installed"
else
    KUBECTL_VERSION=$(curl -Ls https://dl.k8s.io/release/stable.txt)
    curl -LO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl"
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    rm kubectl
    ok "kubectl $KUBECTL_VERSION installed"
fi

# --- minikube ---
step "Installing minikube"
if command_exists minikube; then
    skip "minikube already installed"
else
    curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    install minikube-linux-amd64 /usr/local/bin/minikube
    rm minikube-linux-amd64
    ok "minikube installed"
fi

# --- Go tools (run as real user) ---
step "Installing Go tools"
export GOPATH="/home/${SUDO_USER:-$USER}/go"
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

sudo -u "${SUDO_USER:-$USER}" bash -c \
    'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin && \
     go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
     go install github.com/yannh/kubeconform/cmd/kubeconform@latest'
ok "golangci-lint and kubeconform installed"

# --- Build Docker images ---
step "Building Docker images"
cd "$(dirname "$0")/.."
docker compose build
ok "Docker images built"

# --- Minikube (run as real user) ---
step "Starting minikube"
sudo -u "${SUDO_USER:-$USER}" minikube delete 2>/dev/null || true
sudo -u "${SUDO_USER:-$USER}" minikube start
ok "minikube started"

step "Deploying Kubernetes manifests"
sudo -u "${SUDO_USER:-$USER}" kubectl apply -f k8s/namespace.yaml
sudo -u "${SUDO_USER:-$USER}" kubectl apply -f k8s/
ok "Manifests applied"

# --- Verification ---
step "Verification"
echo ""
# Run verification as the real user so tools in ~/.cargo/bin, ~/go/bin, etc. are found
REAL_USER="${SUDO_USER:-$USER}"
REAL_HOME="/home/$REAL_USER"
USER_PATH="/usr/local/go/bin:$REAL_HOME/go/bin:$REAL_HOME/.cargo/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

verify() {
    local name="$1"; shift
    if full_output=$(sudo -u "$REAL_USER" env "PATH=$USER_PATH" "HOME=$REAL_HOME" "$@" 2>&1); then
        output=$(echo "$full_output" | head -1)
        ok "$name: $output"
    else
        fail "$name: not found"
    fi
}

verify "docker"         docker --version
verify "docker compose" docker compose version
verify "go"             go version
verify "rustc"          rustc --version
verify "cargo"          cargo --version
verify "minikube"       minikube version
verify "kubectl"        kubectl version --client
verify "golangci-lint"  golangci-lint --version
verify "kubeconform"    kubeconform -v

echo ""
echo -e "\033[32mSetup complete!\033[0m"
echo "NOTE: Log out and back in for docker group membership to take effect."
