# Session 11

**Date**: 2026-02-16
**Focus**: Frontend, setup-linux.sh fixes, multipart API, Redis connection, Go mentoring

## Summary

Large session covering: web frontend implementation, setup-linux.sh Go version handling and verification fixes, switching POST /jobs to multipart/form-data, Go concepts teaching, Denis implementing handlePostJob, and connecting the Go API to Redis.

## What was done

### Commits in this session:

1. **5757db3** — Web frontend + setup-linux.sh Go version handling
2. **91a8acc** — Execution guide (docs/EXECUTION.md)
3. **6283cd0** — Switch POST /jobs to multipart/form-data with filter field
4. **a0c8c6f** — Denis implements handlePostJob with multipart parsing
5. **d2a739a** — Redis connection, JobsContext DI, port standardization to 8081
6. **dc24d58** — Session 11 conversation log
7. **5834745** — Refactor handlePostJob: extract helpers

### Details:

- Created embedded web frontend (HTML/CSS/JS) with drag-drop upload, filter dropdown, job polling
- Fixed setup-linux.sh: Go version comparison/upgrade, verification as real user, SIGPIPE fix
- Added Go PATH to ~/.profile for non-interactive shells
- Bumped Go to 1.26.0
- Created docs/EXECUTION.md and docs/FRONTEND.md
- Switched POST /jobs from raw binary to multipart/form-data (image + filter)
- Updated openapi.json, frontend, and tests accordingly
- Removed StatusResponse/returnStatus, replaced with setHeader
- Denis implemented handlePostJob (Claude reviewed, suggested fixes)
- Denis connected Go API to Redis with JobsContext dependency injection
- Standardized port to 8081 across all configs
- Refactored handlePostJob: extracted extractImage(), extractFilter(), writeJSON() helpers

## Key decisions

- multipart/form-data chosen over base64-in-JSON for file upload (no encoding overhead, stdlib support)
- JobsContext struct for dependency injection of Redis client into handlers
- Frontend is vanilla JS, no frameworks
- Go version comparison uses `sort -V` for semantic version sorting
- Tests use nil Redis client for unit tests (no Redis needed for validation logic)

## Teaching moments

- Go multipart: `r.FormFile()` returns `(multipart.File, *FileHeader, error)` — use the File directly, don't call `FileHeader.Open()`
- Go enums: `const + iota` pattern vs string constants with validation map
- Go `WriteHeader` must be called before `Write`, and only once — buffer response before committing status code
- Go naming: don't use `ctx` as receiver name (reserved for context.Context), don't use `error` as variable name (shadows builtin)
- `io.ReadAll` preferred over manual `Read` + size check
- `set -euo pipefail` + `| head -1` causes SIGPIPE (exit 141) on multi-line output

## Full Transcript

> Implement the following plan:
> # Plan: Add Web Frontend to Retro Image Filter Service
> [full plan text]

[Read main.go, routes.go, jobs.go, models.go, health.go]

Created static files (index.html, style.css, app.js), embed.go, updated routes.go, created docs/FRONTEND.md.

Attempted `go build` — Go not found in PATH.

> It seems that there is still issue with the setup-linux.sh script. It tells that go is installed but it's not accessible from the commandline. Where is it located ? Can you add it to the user path permanently ?

Found Go at `/usr/local/go/bin/go` but not in PATH. Added to `~/.profile`. Fixed setup-linux.sh: `[ -x /usr/local/go/bin/go ]` instead of `command_exists go`, PATH in both .bashrc and .profile.

> Can you update the setup-linux.sh script to use apt install commands, instead of downloading custom tar files ?

Changed Go/kubectl/minikube to apt packages.

> Ok nevermind I need more recent version. Revert this and keep downloading tarball directly

Reverted to tarball downloads.

> Can you use the latest version of go ? 1.22 is too old

[Web search: Go 1.26.0 released Feb 10, 2026]

Updated GO_VERSION to 1.26.0.

> In the script it just check if go is already installed but it doesn't check the version. Please check the version, and install the new version of go if the currently installed version is older

Added version comparison with `sort -V`. Tested with three cases (older/same/newer).

> Run the script and verify it's working

[Ran sudo setup-linux.sh — Go upgraded, verification had FAIL for rustc/cargo/minikube]

Fixed verification to run as real user with proper PATH/HOME. Still had minikube FAIL.

> yes fix the verification section

Added REAL_HOME, then found SIGPIPE issue: `pipefail` + `| head -1`. Fixed by capturing full output first.

> Run the full script again to verify everything's working

[All 9 checks green]

> Run the webserver so that I can verify the front-end you generated is working

[Built and ran go-service, verified HTTP 200 at localhost:8081]

> commit and push

[Committed 5757db3, pushed]

> in the docs folder, can you generate a EXECUTION.md that describes the different way to execute this project ?

Created docs/EXECUTION.md with 5 execution methods.

> yes, commit and push

[Committed 91a8acc, pushed]

> In test openapi.json and in the test, the content of the endpoint POST /jobs is directly a binary. I'd like to have more information after. Can you update this so it receives a json with the image in it ? What's the best way to send a binary file and a json together ?

Explained three approaches: multipart/form-data (recommended), base64-in-JSON, two requests. Denis chose multipart.

> I'll go with the multipart/from-data solution. Can you update the openapi, the front-end, and the go unit tests please ? I'd like to have only a filter field in the json for now

Updated openapi.json, frontend (filter dropdown), and rewrote jobs_test.go with multipart requests + new test cases.

> I've updated models.go to remove the StatusResponse struct and change the function returnStatus to setHeader. Can you update the openapi and the unittests accordingly ?

Removed models_test.go, updated health_test.go, removed StatusResponse from openapi.json.

> commit and push

[Committed 6283cd0, pushed]

> My next step will be to connect the go component to redis. Can you give me some guidelines ?

Explained Redis data model (LIST for queue, HASH for job state), go-redis library, dependency injection patterns, suggested implementation order.

> I'll start with 1-2

Gave hints: `go get github.com/redis/go-redis/v9`, `redis.NewClient`, `os.Getenv`, `client.Ping`.

> I did the update to connect to redis and have it passed to the jobs handlers. Can you review the code and update the unit test ?

[Read updated main.go, routes.go, jobs.go, docker-compose.yml]

Review: good DI pattern with JobsContext, noted `ctx` receiver name should be `jc` (ctx reserved for context.Context), Redis client unused in handlers yet (expected). Updated tests with `newTestJobsContext()`. Updated ports to 8081 in Dockerfile, docker-compose, K8s manifests, docs.

> commit and push

[Committed d2a739a, pushed]

> Can you generate a curl command that will send the file from images_bank/html5.png ?

```
curl -X POST http://localhost:8081/jobs -F "image=@..." -F "filter=grayscale"
```

> Can you try to run the command ? I see errors when I copy/paste this

[Ran curl successfully, server returned empty body as expected from stub]

> How can I ask curl to display the code returned ?

Explained `-w "\n%{http_code}\n"` flag.

> In Go, what's the equivalent of declaring an enum and their conversion to string values ?

Explained `const + iota` pattern with String() method, and simpler string constants with validation map alternative.

> How do I set the return code on a http.ResponseWriter ?

Explained `w.WriteHeader()` — must call before Write, can only call once.

> Now I'm writing the return code StatusCreated before encoding the body. What happens if there is an error line 68 and I trigger a http.Error ?

Explained: first WriteHeader wins, 500 silently ignored. Showed buffer pattern to delay WriteHeader.

> How can I declare a constant buffer from hexadecimal values in go ?

Explained `var pngMagic = []byte{0x89, ...}` (no const for slices) and `const` string with `\x` escapes alternative.

> I implemented handlePostJob and all tests are passing. Can you review my code ?

[Read jobs.go]

Review: Fix header/WriteHeader order (swap lines 74-75), FilterSepia: false bug, rename `error` to `err`, use `io.ReadAll`, add bounds check for PNG magic.

> I did the fixes you suggested. For 2. I already did the update on the unit test

[Read updated jobs.go, verified sepia removed from tests]

Code looked good, ready to commit.

> commit and push

[Committed a0c8c6f, pushed]

> I did the update to connect to redis and have it passed to the jobs handlers. Can you review the code and update the unit test ? I also updated the port forwarding in the docker-compose, please update the markdown accordingly

[Read all updated files, updated tests with JobsContext, updated ports everywhere]

> commit and push

[Committed d2a739a, pushed]

> Looks like you forgot to update the conversations again... Can you update it and store this directive in MEMORY.md ?

Updated MEMORY.md with critical reminder to update conversation log on every commit. Created session-011.md.

[Committed dc24d58, pushed]

> handlePostJob function is a bit too big. Do you have suggestion to make it shorter ? Like creating sub-functions

Suggested extracting validatePNG(), validateFilter(), and writeJSON() helpers.

> why only validate the filter and PNG ? Wouldn't it be better if the functions extract the data as well ?

Denis was right — better to have functions that extract + validate together, returning clean data or error.

> Could you do the update ? It's boilerplate

Refactored jobs.go: extracted extractImage(), extractFilter(), writeJSON(). Handler went from ~45 to ~17 lines.

[Tests pass]

> commit and push

[Committed 5834745, pushed]

> man you forgot to update the conversations AGAIN... why do you keep forgetting it ?

Updated session-011.md and SUMMARY.md. Strengthened MEMORY.md reminder.
