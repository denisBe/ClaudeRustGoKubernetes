# Web Frontend

## Accessing the UI

Start the service and open **http://localhost:8081** in a browser.

With Docker Compose:
```bash
docker compose up --build
```

Or run directly:
```bash
cd go-service && go run .
```

## Features

- **Drag-and-drop upload** — drop a PNG image onto the upload zone (or click to browse)
- **Job submission** — sends the image to `POST /jobs` and tracks the returned job
- **Auto-polling** — polls `GET /jobs/{id}` every 2 seconds for pending/processing jobs
- **Result display** — shows the filtered image when a job completes
- **Health indicator** — green/red dot in the header showing API status via `/healthz`

## How It Works

The frontend is a single-page app built with vanilla HTML, CSS, and JavaScript (no frameworks). The static files live in `go-service/static/` and are embedded into the Go binary at compile time using Go's `embed` package.

### File structure

```
go-service/
  embed.go          # //go:embed directive, exports StaticFS()
  routes.go         # Registers file server handlers
  static/
    index.html      # Main page
    style.css       # Styling (dark theme, retro colors)
    app.js          # API calls, polling, UI logic
```

### Embed mechanism

`embed.go` contains:
```go
//go:embed static
var staticFiles embed.FS
```

This tells the Go compiler to include the entire `static/` directory inside the compiled binary. No volume mounts or file copies are needed at runtime — the binary is self-contained.

### Routing

- `GET /` — serves `index.html`
- `GET /static/*` — serves CSS, JS, and other static assets
- API routes (`/healthz`, `/jobs`, etc.) are unaffected

## Modifying the Frontend

Edit files in `go-service/static/` and rebuild:

```bash
cd go-service && go build .
```

Since files are embedded at compile time, you must rebuild after every change. For faster iteration during development, you could temporarily serve from the filesystem instead of the embed (but the embed approach is used for production/Docker).

## Styling

The CSS uses a dark theme with retro-inspired accent colors:
- Background: `#1a1a2e` (dark navy)
- Primary accent: `#e94560` (retro red/pink)
- Secondary accent: `#4ecca3` (teal green)

These can be adjusted in `style.css`.
