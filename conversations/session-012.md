# Session 12

## Summary
Redis job persistence and queue: added `JobQueueMessage`, fixed `CreateJob` to store job status with `Set` and push to `jobs:queue` with `RPush`. Updated tests to use `miniredis` and added three new Redis-verifying tests.

## Date
2026-02-21

## Commit
ee74e1c

## What was done
- Added `JobQueueMessage` struct with `ImageData []byte` (auto base64 in JSON)
- Rewrote `CreateJob` in `db.go`: stores `JobInfo` via `Set("job:<id>")` and pushes `JobQueueMessage` via `RPush("jobs:queue")`
- Added `miniredis` test dependency
- Updated `newTestJobsContext` to use a real in-process Redis (miniredis) instead of nil client
- Added tests: `TestPostJob_JobStatusStoredInRedis`, `TestPostJob_JobPushedToQueue`, `TestPostJob_QueueMessageImageIsBase64Encoded`
- All 14 tests pass

## Key decisions
- Separate structs for DB record (`JobInfo`) and queue message (`JobQueueMessage`) — different consumers need different data
- `[]byte` fields in Go JSON are automatically base64-encoded, matching what Rust's `serde_json` expects
- `miniredis` used for testing — in-process, no real Redis server needed, shuts down automatically per test
- Denis was tired so Claude wrote the production code this session (exception to the learning rules)

## Full Transcript

> Hi Claude! How are you doing today ?

Doing well, thanks for asking! Ready to help with your learning project.

I can see there are some modified files in your Go service (`jobs.go`, `jobs_test.go`, `routes.go`, and a new `db.go`). Looks like you've been busy! What are you working on today?

---

> What's the difference between Sonnet and Opus ?

[Used claude-code-guide agent to research the answer]

Here's the short version:

| | Sonnet | Opus |
|---|---|---|
| **Intelligence** | Strong, general-purpose | Highest, frontier-level |
| **Speed** | Faster | Slower |
| **Cost** | Lower | Higher |
| **Best for** | Most everyday tasks | Complex reasoning, deep analysis |

Think of it like CPU tiers — Sonnet is a fast, efficient mid-range; Opus is the high-end for demanding workloads.

**For this project:** You're currently talking to Sonnet 4.6. Opus 4.6 would give deeper reasoning but costs more. For most of our Go/Rust learning sessions, Sonnet handles it well — though if you want to tackle something particularly complex, Opus would shine there.

You can switch models in Claude Code settings if you want to experiment.

---

> With redis in golang, how can I serialize a struct that already have json annotations ? Do i need to serialize it manually ? Knowing that it needs to be read from rust

Good question — and the answer is: **you need to call `json.Marshal()` explicitly before storing to Redis**.

[Explained the flow: struct → json.Marshal() → []byte → Redis, and Rust reads with serde_json]

---

> what's the difference between HSet and Set in the go redis library ? And what are their expected arguments ?

[Explained Set (Redis STRING) vs HSet (Redis HASH), and recommended Set+JSON for the queue use case]

---

> I need to do both. I want to store in queue for the rust to consume the job, and I want to add an entry in the DB (JobInfo) to persist the job status. Can you update the code of CreateJob please ?

[Read db.go and jobs.go]

I'd love to help, but per our learning agreement — I guide, you write! Let me read the current code first so I can give you useful direction.

[Explained: Bug 1 (wrong HSet usage), Bug 2 (image not stored), Missing piece (RPush for queue). Denis was already mid-edit on the file.]

---

> I'm very tired. My plan is to have one entry in redis (the Set command I think) to persist the job and its status, so that it can be queried in the handleGetJob method, and to push the image with the job uuid in a queue so that the rust component can access it. I know it's against the rules, but can you do the update for me please ?

[Read db.go again — Denis had already partially edited it. Wrote complete updated db.go with JobQueueMessage struct, Set for status, RPush for queue.]

Done. Here's what changed and why: [explanation of JobQueueMessage, Set with 0 TTL, RPush FIFO ordering]

---

> Thank you! Can you update the unit tests ?

[Read jobs_test.go and go.mod]
[Added miniredis dependency via go get]
[Wrote updated jobs_test.go with miniredis and 3 new Redis-verifying tests]
[Fixed LLen call (mr.LLen → redisClient.LLen)]
[All 14 tests pass]

---

> commit and push

[This commit]
