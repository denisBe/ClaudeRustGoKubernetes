# Session 13

## Summary
Implemented `GET /jobs/{id}` endpoint: Denis wrote `handleGetJob` and `GetJobInfo`, Claude added tests and fixed 404 vs 500 error handling for missing jobs.

## Date
2026-02-21

## Commit
(to be filled)

## What was done
- Denis implemented `handleGetJob` using `r.PathValue("id")` and `GetJobInfo`
- Denis implemented `GetJobInfo` in `db.go` with real Redis `Get` + JSON unmarshal
- Claude added `GET /jobs/{id}` route in `routes.go`
- Claude added test helpers `newTestMux` and `createJob`
- Claude added 4 new tests: 200, JSON content-type, correct job info, 404 for missing ID
- Fixed `handleGetJob` to return 404 (not 500) for missing jobs using `errors.Is(err, redis.Nil)`
- Fixed button text mismatch in `app.js` ("Apply Retro Filter" → "Apply Filter")
- All 18 tests pass

## Key decisions
- `redis.Nil` bubbles up from `GetJobInfo` unwrapped; handler distinguishes it with `errors.Is`
- Custom error type (`ErrJobNotFound`) discussed but deferred as overkill for now
- Tests written against correct behavior (404) before fix, causing intentional failure that guided Denis to the fix

## Full Transcript

> With golang. how Can i create a route handler (mux.HandleFunc) with a get argument ?

[Explained Go 1.22 path parameters: `"GET /jobs/{id}"` pattern and `r.PathValue("id")`]

---

> I've implemented the handleGetJob method. Can you update the unittests and the front end please ?

[Read routes.go, jobs.go, db.go, index.html, app.js]

Good work! Let me also check the frontend before writing anything.

[Noted code review issue: handleGetJob returns 500 for all errors including "not found" — wrote tests against correct behavior (404) to catch it]

[Added miniredis dependency, wrote 4 new GET /jobs/{id} tests, fixed app.js button text]

17/18 pass. Only `TestGetJob_NonExistentID_Returns404` fails — still returning 500 instead of 404.

[Explained the fix: `errors.Is(err, redis.Nil)` check needed in handleGetJob]

---

> How can I know if db.redisClient.Get return an error because the entry (id) was not found ?

The go-redis library uses a sentinel error value called `redis.Nil` when a key doesn't exist. You check for it with `errors.Is`. [Explained analogy to C++ `errno == ENOENT`]

---

> Feels like it should be the responsability of GetJobInfo don't you think? But if it means that I need to define a custom error for it it might be overkill at this point of the project ?

Exactly right on both counts. `GetJobInfo` already returns `redis.Nil` unwrapped — the fix only needs to go in `handleGetJob`. Custom error type discussed but deferred as overkill.

---

> Please do this update, we already discussed it now it's boilerplate

[Added `errors` and `github.com/redis/go-redis/v9` imports to jobs.go, added `errors.Is(err, redis.Nil)` check in handleGetJob]

All 18 tests pass.

---

> run the tests

[Ran go test ./... — 18/18 pass]

---

> commit and push

[This commit]
