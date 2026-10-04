# URL Shortener

A REST API for shortening URLs, built in Go using only the standard library — no third-party web framework. The storage layer sits behind an interface, with both an in-memory and a PostgreSQL implementation provided; swapping between them is a one-line change in `main.go`.

## Features

- Shorten a URL with a randomly generated 6-character code
- Redirect through a short link, with per-link access tracking
- Fetch a link's details without following the redirect
- List all stored links
- Update a link's target URL
- Delete a link
- Access statistics (hit count, last accessed time)
- URL validation (scheme and host required)
- Request logging middleware (method, path, duration)

## Stack

- **Go 1.22+** — standard library only for HTTP: `net/http`, `encoding/json`, `net/url`, `database/sql`, `log/slog`
- **PostgreSQL** via `jackc/pgx/v5/stdlib`
- **godotenv** for loading local environment variables
- Routing on `net/http`'s built-in `ServeMux` (Go 1.22+ method-aware patterns, e.g. `"GET /shorten/{id}"`) — no router dependency
- Storage behind a `Storage` interface: `MemoryStorage` (map-backed) and `PostgresStorage` (parameterized SQL queries)

## Architecture

```
url-shortener/
├── main.go
├── init.sql                        # table schema
├── .env.example
├── models/
│   └── url.go                      # URLRecord
├── dto/
│   └── dto.go                      # request/response payloads
└── internal/
    ├── storage/
    │   ├── storage.go              # Storage interface
    │   ├── memory_storage.go       # in-memory implementation
    │   ├── postgres_storage.go     # PostgreSQL implementation
    │   └── id_generator.go         # random short-code generator
    ├── handlers/
    │   └── handlers.go             # HTTP handlers
    └── middleware/
        └── logging.go              # request logging
```

The `Storage` interface is the backbone of the design: handlers depend only on the interface, never on a concrete storage type. Switching from an in-memory store to PostgreSQL required a single change in `main.go` and zero changes to `handlers`.

## Getting started

### 1. Database

```bash
psql -h localhost -p 5432 -U <username>
CREATE DATABASE urlshortener;
\c urlshortener
\i init.sql
```

### 2. Environment

```bash
cp .env.example .env
```
Fill in your local credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=urlshortener
```

### 3. Run

```bash
git clone https://github.com/flecc1/URL-Shortener.git
cd URL-Shortener
go run main.go
```
Server starts on `http://localhost:8080`.

## API

| Method | Path                   | Description                  |
|--------|------------------------|-------------------------------|
| POST   | `/shorten`             | Create a short link           |
| GET    | `/{id}`                | Redirect to the original URL  |
| GET    | `/shorten/all`         | List all links                |
| GET    | `/shorten/{id}`        | Get link details              |
| PUT    | `/shorten/{id}`        | Update the target URL         |
| DELETE | `/shorten/{id}`        | Delete a link                 |
| GET    | `/shorten/{id}/stats`  | Access statistics             |

## Example

```bash
curl -X POST http://localhost:8080/shorten \
  -d '{"url": "https://google.com"}'
```
```json
{
  "id": "AULm6Z",
  "short_url": "http://localhost:8080/AULm6Z",
  "original_url": "https://google.com",
  "created_at": "2026-09-27T11:49:42+03:00"
}
```

```bash
curl http://localhost:8080/shorten/AULm6Z/stats
```

## Design notes

- **Interface-driven storage.** `Storage` defines the contract (`Create`, `Get`, `GetAll`, `UpdateById`, `DeleteById`, `IncrementAccess`); `handlers` is written against the interface, not a concrete type. This is what made the in-memory → PostgreSQL migration a no-op for the handler layer.
- **SQL injection safety.** All PostgreSQL queries use parameterized placeholders (`$1`, `$2`, …) — no string concatenation with user input, anywhere.
- **Routing.** Built on Go 1.22's enhanced `ServeMux`: HTTP method and path parameters are declared directly in the pattern (`"PUT /shorten/{id}"`), extracted via `r.PathValue(...)`. No manual path parsing, no external router.
- **Error handling.** A shared `ErrNotFound` sentinel is used by both storage implementations, so callers can check for it with `errors.Is` regardless of which backend is active. `database/sql` errors are distinguished by cause: `sql.ErrNoRows` → 404, everything else → 500; `RowsAffected()` is checked on `UPDATE`/`DELETE` to detect a no-op write.
- **Dependency injection.** Both `Handler` and each storage implementation are constructed explicitly (`NewHandler`, `NewMemoryStorage`, `NewPostgresStorage`) rather than relying on global state.

## Known limitations

- `MemoryStorage` has a `sync.RWMutex` field declared but not yet wired in — concurrency safety is scheduled once goroutines/channels are covered
- No automated tests yet
- `GetAll` has no pagination
- Connection pool uses `database/sql` defaults (no `SetMaxOpenConns` tuning)

## What this project covers

Built while working through *The Go Programming Language* (Donovan & Kernighan), chapters 1–7, alongside Go by Example and hands-on practice. Topics applied here: structs and methods, interfaces, closures, recursion, `net/http` from first principles (including the pre- and post-1.22 routing styles), `database/sql`, middleware, and package-level architecture (`internal/`, dependency injection via constructors).

## Author

flecc1