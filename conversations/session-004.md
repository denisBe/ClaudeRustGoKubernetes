# Session 4 — 2026-02-14

## Summary

Chose project: Retro Image Filter Service. Defined v1 scope and created PROJECT.md.

## Changes

- Created `PROJECT.md` with architecture, API design, filter specs, and tech stack

## Discussion

- Explored project ideas fitting Go + Rust + Redis + K8s autoscaling requirements
- Denis asked how Go would know when Rust finished processing (Redis patterns: polling, pub/sub, WebSocket, blocking pop)
- Chose **Retro Image Filter Service** — fun, visual results, good fit for all technologies
- v1 scope: grayscale filter, PNG only, API-only (no web UI), async polling
- Grayscale chosen as first filter for mathematical simplicity (per-pixel operation)

---

## Full Transcript

**Denis:**
> Let's talk about the project itself. Like I said my goal is to learn go, rust and kubernetes. I also want to use redis. I want to test how the scalability of kubernetes work to spawn new docker for the rust component when load is too high. Do you have suggestion for a project that could fit this need?

**Claude:**
> Suggested 3 project ideas fitting Go + Rust + Redis + K8s autoscaling:
> 1. Image processing pipeline
> 2. Log aggregation / text analysis
> 3. Prime number / hash computation service
>
> Recommended #1 (image processing) for visible results and genuine CPU load.

**Denis:**
> If redis is used as a buffer between go and rust, how will go know when the processing has been completed to reply to the request?

**Claude:**
> Explained 4 patterns:
> 1. Async with polling/webhook (simplest, recommended to start)
> 2. Redis Pub/Sub (real-time but holds connection)
> 3. WebSocket/SSE (real-time, more complex)
> 4. Redis blocking pop (synchronous feel)

**Denis:**
> I was just curious. I'm still thinking about what project to implement. Would you have a bit more fun or usefull suggestions?

**Claude:**
> Suggested more creative options:
> - Multiplayer game server
> - Retro image filter service (like Prisma/Instagram)
> - Personal web crawler / bookmark enricher
> - PDF / report generator
> - Music audio processor

**Denis:**
> let's go with the retro image filter service

**Claude:**
> Great choice! Asked about filters, input/output, formats, and result delivery.

**Denis:**
> for the filter I'd like something mathematically basic for a first implementation. We can add more complex filters later

**Claude:**
> Suggested simple per-pixel filters: grayscale, sepia, invert, brightness/contrast, posterize.

**Denis:**
> Let's go with grayscale for a first filter. For the rest, let's go with just api first, a PNG, and async polling

**Claude:**
> Summarized v1 scope: Go API + Rust worker + Redis queue, PNG upload, grayscale filter, async polling.

**Denis:**
> Create a project description file (markdown) to store this, commit and push
