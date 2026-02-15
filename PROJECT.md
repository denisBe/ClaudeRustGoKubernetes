# Retro Image Filter Service

A microservices-based image processing pipeline built to learn Go, Rust, Redis, and Kubernetes.

## Architecture

- **Go API** — HTTP service that accepts image uploads, queues processing jobs, and serves results
- **Rust Worker** — Pulls jobs from Redis, applies image filters, writes results back
- **Redis** — Job queue between Go and Rust, result storage
- **Kubernetes** — Orchestration with autoscaling of Rust workers under load

## API

### POST /jobs
- Upload a PNG image
- Returns a job ID

### GET /jobs/{id}
- Poll for job status and result
- Returns status (`pending`, `processing`, `completed`, `failed`) and result when ready

## Filters

### v1
- **Grayscale** — per-pixel weighted average: `0.299*R + 0.587*G + 0.114*B`

### Future
- Sepia, invert, brightness/contrast, posterize
- Kernel-based filters: blur, edge detection, emboss

## Result delivery

Async polling — client submits a job, receives a job ID, polls until completion.

## Tech stack

| Component | Technology |
|-----------|------------|
| API | Go |
| Worker | Rust |
| Queue / Cache | Redis |
| Container | Docker |
| Orchestration | Kubernetes (with HPA for Rust workers) |
| CI | GitHub Actions |
| Image format | PNG (v1) |
