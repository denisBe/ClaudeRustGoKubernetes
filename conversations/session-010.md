# Session 10

**Date**: 2026-02-16
**Focus**: Web frontend + setup-linux.sh fixes

## Summary

Added a web frontend (HTML/CSS/JS) to the Go API service using Go's `embed` package. Fixed multiple issues in `setup-linux.sh`: Go PATH not in `.profile`, Go version check, verification running as wrong user, and SIGPIPE with `pipefail`.

## What was done

1. **Frontend files created** (Claude-authored, per CLAUDE.md rules):
   - `go-service/static/index.html` — SPA with drag-drop upload, job list, result display, health indicator
   - `go-service/static/style.css` — dark retro theme (navy bg, red/teal accents)
   - `go-service/static/app.js` — vanilla JS: file upload, POST /jobs, polling, health check
2. **Go embed integration**:
   - `go-service/embed.go` — `//go:embed static` directive
   - `go-service/routes.go` — static file server at `/static/` and `/`
3. **`docs/FRONTEND.md`** — usage and architecture docs
4. **`~/.profile`** — added `/usr/local/go/bin` to PATH (Go wasn't accessible in non-interactive shells)
5. **`scripts/setup-linux.sh` fixes**:
   - Go version check: compares installed version vs desired, upgrades if older
   - Go detection: uses `[ -x /usr/local/go/bin/go ]` instead of `command_exists go` (which checked root's PATH)
   - PATH setup: always ensures Go is in both `.bashrc` and `.profile`
   - Go version bumped from 1.22.5 to 1.26.0
   - Verification section: runs as real user with proper PATH/HOME
   - Fixed SIGPIPE (exit 141) from `pipefail` + `| head -1` on multi-line output
6. **`go-service/.gitignore`** — ignore compiled binaries

## Key decisions

- Frontend is vanilla JS, no frameworks — keeps it simple and dependency-free
- Static files embedded in Go binary via `//go:embed` — no volume mounts needed
- Go version comparison uses `sort -V` for robust semantic version sorting
- Verification runs commands via `sudo -u` with explicit PATH and HOME

## Full Transcript

> Implement the following plan:
> # Plan: Add Web Frontend to Retro Image Filter Service
> [full plan text]

[Read main.go, routes.go, jobs.go, models.go, health.go to understand current codebase]

Created task list for tracking progress. Created `go-service/static/` directory and three files:
- `index.html` — single-page app with upload zone, job list, result section, health dot
- `style.css` — dark theme with retro colors (#1a1a2e bg, #e94560 primary, #4ecca3 secondary)
- `app.js` — IIFE with drag-drop handling, fetch API calls, polling logic

Created `go-service/embed.go` with `//go:embed static` and `StaticFS()` helper using `fs.Sub`.

Modified `go-service/routes.go` to add file server handler for `/static/` and root `/` serving index.html.

Created `docs/FRONTEND.md` with usage instructions and architecture explanation.

Attempted `go build` — Go not found in PATH.

> It seems that there is still issue with the setup-linux.sh script. It tells that go is installed but it's not accessible from the commandline. Where is it located ? Can you add it to the user path permanently ?

Found Go at `/usr/local/go/bin/go` but PATH didn't include it. The line was in `.bashrc` but not `.profile`, and Claude Code runs non-interactive shells that skip `.bashrc`.

Added Go PATH to `~/.profile`. Fixed `setup-linux.sh`: changed `command_exists go` to `[ -x /usr/local/go/bin/go ]`, moved PATH setup outside if/else to always run, added to both `.bashrc` and `.profile`.

Verified `go version` works with updated PATH. Built Go service successfully. Existing test failures are pre-existing (stub handlers).

> Can you update the setup-linux.sh script to use apt install commands, instead of downloading custom tar files ?

Changed Go to `apt-get install -y golang-go`, kubectl to official K8s apt repo, minikube to `.deb` package install.

> Ok nevermind I need more recent version. Revert this and keep downloading tarball directly

Reverted all three (Go, kubectl, minikube) back to tarball/binary downloads, keeping the bug fixes.

> Can you use the latest version of go ? 1.22 is too old

[Web search found Go 1.26.0 released Feb 10, 2026]

Updated `GO_VERSION="1.26.0"` in setup-linux.sh.

> In the script it just check if go is already installed but it doesn't check the version. Please check the version, and install the new version of go if the currently installed version is older

Added version comparison using `sort -V`: extracts installed version, compares with desired, upgrades if older, skips if up to date or newer.

Tested comparison logic with three cases: older (upgrade), same (skip), newer (skip) — all correct.

> Run the script and verify it's working

[Ran `sudo setup-linux.sh` — Go upgraded 1.22.5→1.26.0, all passed except rustc/cargo/minikube FAIL in verification]

Verification failures because `verify()` ran commands as root but rustc/cargo are in user's `~/.cargo/bin` and minikube needs user's HOME. Added `sudo -u` with explicit PATH and HOME.

Ran again — minikube still failed. Found `REAL_HOME` wasn't defined in verification section. Added it.

Ran again — minikube still failed. Debugged: `set -euo pipefail` + `| head -1` causes SIGPIPE (exit 141) when command outputs multiple lines. Fixed by capturing full output first, then piping to `head -1`.

Tested verification in isolation — all 9 green.

> Run the full script again to verify everything's working

[Ran full script — all 9 verification checks green]

> Run the webserver so that I can verify the front-end you generated is working

[Built and ran `go-service` in background, verified HTTP 200 at localhost:8081]

> commit and push

[Committed and pushed all changes]
