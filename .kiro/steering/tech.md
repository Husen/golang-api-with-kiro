# Tech Stack

- **Language**: Go 1.22
- **Framework**: Gin (`github.com/gin-gonic/gin v1.10.0`)
- **API Docs**: Swaggo (`swaggo/swag`, `swaggo/gin-swagger`, `swaggo/files`)
- **Containerization**: Docker + Docker Compose

## Common Commands

```bash
# Generate Swagger docs (required before building if annotations changed)
swag init

# Run locally
go run .

# Build binary
go build -o server .

# Run with Docker Compose
docker-compose up --build
```

## Notes

- Swagger docs are auto-generated into `docs/` via `swag init` — do not edit that folder manually
- The Dockerfile runs `swag init` as part of the build step
- No external database; storage is in-memory using a slice + `sync.Mutex`
- `GIN_MODE=release` is set in Docker Compose for production-like behavior
