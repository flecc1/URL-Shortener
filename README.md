# URL Shortener

A REST API for shortening URLs, written in Go using only the standard library for HTTP, with no web framework or router. Links are stored in PostgreSQL. Storage sits behind an interface, and an in-memory implementation is included too; switching between them takes one line in `main.go`. The whole stack starts with a single `docker compose up`.

## Features

- Shorten a URL into a random 6-character code
- Redirect through a short link, counting each visit
- Fetch a link's details without following the redirect
- List, update and delete links
- Per-link statistics: hit count, creation time, last visit
- URL validation (`http`/`https` scheme and a host are required)
- Request logging middleware (method, path, duration)

## Stack

- **Go 1.27**, standard library: `net/http`, `encoding/json`, `net/url`, `database/sql`, `log/slog`, `crypto/rand`
- **PostgreSQL 17** via `jackc/pgx/v5/stdlib`
- **Docker**: multi-stage `Dockerfile` plus `docker-compose.yml` with Postgres
- **godotenv** for loading `.env` when running locally

## Quick start (Docker)

You only need Docker installed.

```bash
git clone https://github.com/flecc1/URL-Shortener.git
cd URL-Shortener
cp .env.example .env      # set DB_USER, DB_PASSWORD (must not be empty), DB_NAME
docker compose up --build
```

The API is available at `http://localhost:8080`.

What happens on startup:
1. Postgres starts and, on first run, creates the user and database from `.env` and applies `scheme/init.sql`.
2. Compose waits for the database healthcheck (`pg_isready`) to pass.
3. The app container starts and connects to the database by service name (`DB_HOST=postgres`).

Data is kept in the `pgdata` volume and survives `docker compose down`. To reset the database completely:

```bash
docker compose down -v
```

## Running locally (without Docker)

Requires Go 1.27+ and a running PostgreSQL.

```bash
createdb urlshortener
psql -d urlshortener -f scheme/init.sql
cp .env.example .env      # DB_HOST=localhost and your local credentials
go run .
```

## API

| Method | Path                  | Description                 | Success |
|--------|-----------------------|-----------------------------|---------|
| POST   | `/shorten`            | Create a short link         | 201     |
| GET    | `/{id}`               | Redirect to the original URL | 302    |
| GET    | `/shorten/all`        | List all links              | 200     |
| GET    | `/shorten/{id}`       | Get link details            | 200     |
| PUT    | `/shorten/{id}`       | Update the target URL       | 200     |
| DELETE | `/shorten/{id}`       | Delete a link               | 204     |
| GET    | `/shorten/{id}/stats` | Access statistics           | 200     |

Invalid input returns `400`, an unknown id returns `404`.

### Examples

Create a link:

```bash
curl -X POST http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url": "https://google.com"}'
```
```json
{
  "id": "AULm6Z",
  "short_url": "http://localhost:8080/AULm6Z",
  "original_url": "https://google.com",
  "created_at": "2026-09-27T11:49:42Z"
}
```

Follow it and check the stats:

```bash
curl -i http://localhost:8080/AULm6Z
curl http://localhost:8080/shorten/AULm6Z/stats
```
```json
{
  "id": "AULm6Z",
  "created_at": "2026-09-27T11:49:42Z",
  "last_access_at": "2026-09-27T11:50:03Z",
  "access_count": 1
}
```

## Project structure

```
URL-Shortener/
├── main.go                         # config, DB connection, routes, server
├── Dockerfile                      # multi-stage build
├── docker-compose.yml              # app + Postgres
├── .env.example
├── scheme/
│   └── init.sql                    # table schema
├── models/
│   └── url.go                      # URLRecord
├── dto/
│   └── Dtos.go                     # request/response payloads
└── internal/
    ├── storage/
    │   ├── Storage.go              # Storage interface + ErrRecordNotFound
    │   ├── memory-storage.go       # in-memory implementation
    │   ├── postgres_storage.go     # PostgreSQL implementation
    │   └── idgen.go                # random short-code generator
    ├── handlers/
    │   └── handlers.go             # HTTP handlers
    └── midleware/
        └── logging.go              # request logging
```

## Design notes

- **Interface-driven storage.** Handlers depend only on the `Storage` interface (`Create`, `Get`, `GetAll`, `UpdateById`, `DeleteById`, `IncrementAccess`), never on a concrete type. Moving from the in-memory store to PostgreSQL changed one line in `main.go` and nothing in `handlers`.
- **Concurrency safety in `MemoryStorage`.** Each HTTP request runs in its own goroutine, so the map is guarded by a `sync.RWMutex`: writes take `Lock`, reads take `RLock`. `Create` generates an id and inserts it under a single lock, so check-then-insert is atomic. Reads return copies, never pointers into the map.
- **Short ids.** Generated with `crypto/rand` from a 62-character alphabet, which gives about 56 billion combinations for 6 characters.
- **SQL injection safety.** Every query uses placeholders (`$1`, `$2`, …); user input is never concatenated into SQL.
- **Routing.** Uses the method-aware `ServeMux` patterns (`"PUT /shorten/{id}"`) with `r.PathValue("id")`, so there is no manual path parsing and no external router.
- **Error handling.** Both storage backends return the shared `ErrRecordNotFound` sentinel, which can be checked with `errors.Is`. For Postgres, `sql.ErrNoRows` and `RowsAffected() == 0` on `UPDATE`/`DELETE` map to "not found".
- **Dependency injection.** `Handler` and the storages are built with explicit constructors (`NewHandler`, `NewMemoryStorage`, `NewPostgresStorage`), with no global state.
- **Containerization.** The Go binary is compiled statically (`CGO_ENABLED=0`) in a `golang` image and copied into a small `alpine` image. `go.mod`/`go.sum` are copied before the source so the dependency layer stays cached. Compose starts the app only after the Postgres healthcheck passes, because the app exits if it can't reach the database at startup.

## Known limitations

- No automated tests yet
- The base URL for short links (`http://localhost:8080`) is hard-coded in `main.go` rather than read from config
- The schema is applied by `init.sql` only when the database is first created; there are no migrations
- `GetAll` has no pagination
- Connection pool uses the `database/sql` defaults

## What this project covers

Built while working through *The Go Programming Language* (Donovan & Kernighan), chapters 1–7, alongside Go by Example and hands-on practice. Topics applied here: structs and methods, interfaces, closures, `net/http` from first principles (both the old and the method-aware routing styles), `database/sql`, middleware, package layout with `internal/`, dependency injection through constructors, and containerizing a Go service with Docker and Compose.

## Author

flecc1
