# Quality Standards

## Test-Driven Development (TDD)

This project follows a strict TDD workflow:

1. **Red** — Write a failing test that defines the expected behavior
2. **Green** — Write the minimum code to make the test pass
3. **Refactor** — Clean up the code while keeping tests green

### Rules

- No production code is written without a failing test first
- Tests must be run and passing before every commit
- Each test should test one behavior and have a clear name describing that behavior
- Tests are first-class code — same quality standards as production code

## Go Testing

- Use the standard `testing` package
- Table-driven tests where appropriate
- Use `t.Run()` for subtests
- Run with `go test ./...`

## Rust Testing

- Use built-in `#[cfg(test)]` module tests for unit tests
- Integration tests in the `tests/` directory
- Use descriptive test names: `test_<function>_<scenario>_<expected>`
- Run with `cargo test`

## Linting

### Go

- Run `golangci-lint run` before every commit
- Enabled linters (configured in `.golangci.yml`):
  - `govet` — correctness checks
  - `staticcheck` — advanced static analysis
  - `errcheck` — ensure errors are handled
  - `gosimple` — simplification suggestions
  - `ineffassign` — detect unused assignments
  - `gofmt` — enforce standard formatting
- All code must be formatted with `gofmt`

### Rust

- Run `cargo clippy -- -D warnings` before every commit (warnings treated as errors)
- Run `cargo fmt --check` to enforce formatting
- All code must be formatted with `rustfmt`

## Kubernetes Manifests

- Validate manifests with `kubectl --dry-run=client` before applying
- Use `kubeval` or `kubeconform` for schema validation
- All manifests must specify:
  - `resources.requests` and `resources.limits` for CPU and memory
  - `livenessProbe` and `readinessProbe` for containers
  - Explicit `namespace` (no relying on default)
- Use labels consistently: `app`, `version`, `component`
- No `latest` image tags — always pin to a specific version
- Secrets must never be committed in plain text — use sealed secrets or external secret managers
- Prefer declarative management (`kubectl apply`) over imperative commands

## Continuous Integration (GitHub Actions)

CI runs on every push to `master` and on all pull requests. See `.github/workflows/ci.yml`.

### Pipeline jobs

- **Go** — format check (`gofmt`), lint (`golangci-lint`), test with race detector + coverage
- **Rust** — format check (`cargo fmt`), lint (`cargo clippy`), test (`cargo test`)
- **Kubernetes** — validate manifests with `kubeconform -strict`

### Rules

- All CI checks must pass before merging
- PRs require a green CI status
- No force-pushing to `master`

## General

- All code must compile with no warnings
- No commented-out code in commits
- Keep functions small and focused
- Linters and formatters must pass before every commit
