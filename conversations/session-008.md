# Session 8 — First Go API: Health Check Endpoint

**Date:** 2026-02-14
**Commit:** (this commit)

## Summary

Denis implemented his first Go HTTP API with three endpoints (`/healthz`, `POST /jobs`, `GET /jobs`). Claude explained Go HTTP concepts comparing to C++, reviewed code iteratively, and helped debug issues.

## What was done

- **go-service/main.go** — Implemented REST API with `net/http`:
  - `GET /healthz` — health check (matching K8s deployment manifest)
  - `POST /jobs` — placeholder for job creation
  - `GET /jobs` — placeholder for job status
- **.gitignore** — Added `__debug_bin*` and `*.exe*` patterns for VS Code debug binaries

## Key concepts covered

- Go HTTP handler signature: `func(w http.ResponseWriter, r *http.Request)`
- `ServeMux` routing with Go 1.22+ method patterns
- JSON responses with struct tags
- Error handling as return values (no exceptions)
- `fmt.Println` vs `fmt.Printf` (and the `f` suffix convention)
- Import syntax (quotes required, grouped imports)
- Exported vs unexported (uppercase first letter = public)
- `gofmt` enforces tabs, no style debates
- `/healthz` K8s naming convention
- Browsers default to `https://`, `curl` doesn't

## Full Transcript

### User
> Let's start with the GO api. Can you show me how to create a REST endpoint in golang ?

### Claude
Great, let's get into Go! Since you're a C++ developer, here's how REST endpoints work in Go.

Go has a built-in HTTP server in the standard library (`net/http`) — no framework needed. This is very different from C++ where you'd typically reach for a library like Boost.Beast, cpp-httplib, or a framework like Crow.

Key concepts explained:
1. Handler functions — `func(w http.ResponseWriter, r *http.Request)` — like a callback
2. Routing — `http.NewServeMux()` with Go 1.22+ method patterns like `"GET /health"`
3. Starting the server — `http.ListenAndServe(":8080", mux)` — one line
4. JSON responses — `json.NewEncoder(w).Encode(data)` with struct tags like `json:"status"`

Provided a C++ comparison table (HTTP server, routing, JSON, concurrency, error handling).

Mentioned Go's HTTP server spawns a goroutine per request automatically — no thread pool management.

Suggested starting with `/health` since it's the simplest, then `POST /jobs`, then `GET /jobs/{id}`.

### User
> Thank you for the explanation they're great! I think you're mistaken about the suggested routes though, I'm not working on a gaming server but on a image filter service

### Claude
You're right, I apologize! Let me re-read PROJECT.md to get the actual routes.

[Read PROJECT.md]

The actual API is:
- `POST /jobs` — Upload a PNG image, returns a job ID
- `GET /jobs/{id}` — Poll for job status and result

Suggested starting with a health check first since it's simplest and needed for K8s probes.

### User
> yes let's start with a health check, show me the current main.go

### Claude
[Read go-service/main.go — just `package main` and empty `func main() {}`]

A blank canvas! Explained what Denis needs to do:
1. Import `net/http` and `encoding/json`
2. Create a handler function for `GET /health`
3. Set up a `ServeMux` and register the route
4. Start the server with `http.ListenAndServe`

Hint: Define a struct for the response and use `json.NewEncoder(w).Encode(...)`.

### User
> What's the equivalent of std::cout in go ?

### Claude
`fmt.Println()` — from the `fmt` package. `fmt.Printf` for formatted output. But for the health endpoint, you'll write to `http.ResponseWriter` instead — it's a stream back to the client, not stdout. `fmt.Println` is handy for debug logging though.

### User
> [First review request]

### Claude
[Read main.go — first attempt with several bugs]

Bugs found:
1. `mux.handleFunc` → `mux.HandleFunc` (uppercase = exported)
2. `http.ResponseWrite` → `http.ResponseWriter` (typo)
3. `fmt.PrintLn` → `fmt.Println` (lowercase `l`)
4. `{"Okayyy!!!"}` — need struct name: `StatusResponse{Status: "Okayyy!!!"}`

Design notes: use grouped imports, register both handlers, handle error from `ListenAndServe` with `log.Fatal`.

### User
> what does an error value looks like in go ?

### Claude
Explained Go error interface (`Error() string`), multiple return values pattern `(result, err)`, `if err != nil` convention. Compared to C++ exceptions. Mentioned the `f` suffix convention (`Fatal`/`Fatalf`, `Print`/`Printf`).

### User
> [Second review request]

### Claude
[Read main.go — improved but still had issues]

Remaining bugs:
1. `log.Fatal("Error: %s", err)` → `log.Fatalf("Error: %s", err)` (need `f` suffix for formatting)
2. `http.ResponseWriterr` → `http.ResponseWriter` (extra `r`)
3. `{Status:"Okayyy!!!"}` → `StatusResponse{Status: "Okayyy!!!"}` (still need type name)
4. Import block indented with spaces — Go uses tabs

Still missing: `handlePostJob` not registered, no `/health` endpoint yet.

### User
> what's the GET /health endpoint supposed to do ? Is it a common pattern for rest api ?

### Claude
Explained health endpoints: used for K8s liveness/readiness probes and load balancer checks. Simplest version returns 200 OK with `{"status": "ok"}`. `/healthz` is the K8s convention (Google tradition). Mentioned our `k8s/go-deployment.yaml` likely already references it.

### User
> When checking the go-deployment.yaml, the livenessProbe path is defined to /healthz. Should it be /health instead ?

### Claude
Good catch! Either way is fine as long as they match. `/healthz` is the K8s convention. Suggested Denis pick whichever he prefers and make them consistent.

### User
> let's go with /healthz

### Claude
Go ahead and update your `HandleFunc` route from `"GET /health"` to `"GET /healthz"`.

### User
> When sending a curl command to localhost:8081/healthz it's working, but not when opening a browser (edge) to the same address. Why ?

### Claude
Likely the browser is trying `https://` by default. Suggested typing `http://localhost:8081/healthz` explicitly.

### User
> it's working in the browser, I had to add http://

### Claude
That's it! Browsers default to `https://` these days, `curl` doesn't. Asked if Denis is ready to write tests.

### User
> It's ok for a first time working with go. Please create a commit and push it

### Claude
[Staged go-service/main.go and .gitignore, created session-008.md, committed and pushed]
