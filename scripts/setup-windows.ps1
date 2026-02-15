#Requires -RunAsAdministrator
<#
.SYNOPSIS
    Installs all dependencies for the Retro Image Filter Service on Windows.
.DESCRIPTION
    Uses Chocolatey to install Docker Desktop, Go, Rust, minikube, and kubectl.
    Installs Go-based tools (golangci-lint, kubeconform).
    Fixes common Git Bash issues.
    Builds Docker images and deploys to minikube.
#>

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Write-Step { param([string]$Message) Write-Host "`n==> $Message" -ForegroundColor Cyan }
function Write-Ok   { param([string]$Message) Write-Host "    [OK] $Message" -ForegroundColor Green }
function Write-Skip { param([string]$Message) Write-Host "    [SKIP] $Message" -ForegroundColor Yellow }

# --- Chocolatey ---
Write-Step "Checking Chocolatey"
if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "Chocolatey not found. Install it from https://chocolatey.org/install and re-run this script." -ForegroundColor Red
    exit 1
}
Write-Ok "Chocolatey found"

# --- Core dependencies via Chocolatey ---
$chocoPackages = @("docker-desktop", "golang", "rustup.install", "minikube", "kubernetes-cli")

foreach ($pkg in $chocoPackages) {
    Write-Step "Installing $pkg"
    $installed = choco list --local-only $pkg 2>$null | Select-String $pkg
    if ($installed) {
        Write-Skip "$pkg already installed"
    } else {
        choco install $pkg -y
        Write-Ok "$pkg installed"
    }
}

# Refresh PATH so new tools are available
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
            [System.Environment]::GetEnvironmentVariable("Path", "User")

# --- Go tools ---
Write-Step "Installing golangci-lint"
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
Write-Ok "golangci-lint installed"

Write-Step "Installing kubeconform"
go install github.com/yannh/kubeconform/cmd/kubeconform@latest
Write-Ok "kubeconform installed"

# --- Git Bash fix ---
Write-Step "Fixing Git Bash post-install scripts"
$gitBashPostInstall = "C:\Program Files\Git\etc\post-install"
if (Test-Path $gitBashPostInstall) {
    $postFiles = Get-ChildItem "$gitBashPostInstall\*.post" -ErrorAction SilentlyContinue
    if ($postFiles) {
        $postFiles | Remove-Item -Force
        Write-Ok "Removed $($postFiles.Count) .post file(s)"
    } else {
        Write-Skip "No .post files found"
    }
} else {
    Write-Skip "Git Bash post-install directory not found"
}

# --- Docker credential fix ---
Write-Step "Checking Docker config for credsStore"
$dockerConfig = "$env:USERPROFILE\.docker\config.json"
if (Test-Path $dockerConfig) {
    $config = Get-Content $dockerConfig -Raw | ConvertFrom-Json
    if ($config.PSObject.Properties.Name -contains "credsStore") {
        $config.PSObject.Properties.Remove("credsStore")
        $config | ConvertTo-Json -Depth 10 | Set-Content $dockerConfig
        Write-Ok "Removed credsStore from Docker config"
    } else {
        Write-Skip "No credsStore key found"
    }
} else {
    Write-Skip "Docker config not found (Docker Desktop may not have been run yet)"
}

# --- Build Docker images ---
Write-Step "Building Docker images"
docker compose build
Write-Ok "Docker images built"

# --- Minikube ---
Write-Step "Starting minikube"
minikube start
Write-Ok "minikube started"

Write-Step "Deploying Kubernetes manifests"
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/
Write-Ok "Manifests applied"

# --- Verification summary ---
Write-Step "Verification"
Write-Host ""
$tools = @(
    @{ Name = "docker";         Cmd = "docker --version" },
    @{ Name = "docker compose"; Cmd = "docker compose version" },
    @{ Name = "go";             Cmd = "go version" },
    @{ Name = "rustc";          Cmd = "rustc --version" },
    @{ Name = "cargo";          Cmd = "cargo --version" },
    @{ Name = "minikube";       Cmd = "minikube version" },
    @{ Name = "kubectl";        Cmd = "kubectl version --client" },
    @{ Name = "golangci-lint";  Cmd = "golangci-lint --version" },
    @{ Name = "kubeconform";    Cmd = "kubeconform -v" }
)

foreach ($tool in $tools) {
    try {
        $output = Invoke-Expression $tool.Cmd 2>&1 | Select-Object -First 1
        Write-Ok "$($tool.Name): $output"
    } catch {
        Write-Host "    [FAIL] $($tool.Name): not found" -ForegroundColor Red
    }
}

Write-Host "`nSetup complete!" -ForegroundColor Green
