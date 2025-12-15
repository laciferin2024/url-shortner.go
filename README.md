# URL Shortener Service

A production-ready backend service to shorten long URLs and retrieve original URLs efficiently. Built with Go, PostgreSQL, and Redis.

## Features
- **Shorten URLs**: Convert long URLs into short, unique codes.
- **Redirection**: Redirect short codes to original URLs with low latency.
- **Caching**: Redis caching for fast lookups.
- **Rate Limiting**: IP-based rate limiting to prevent abuse.
- **Admin Listing**: Paginated list of all shortened URLs.
- **Metrics**: Tracks click counts and last accessed time.
- **Dockerized**: Easy deployment with Docker Compose.

## Setup & Run

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (for local development)

### Running with Docker (Recommended)
The easiest way to run the service is using Docker Compose. This starts the application, PostgreSQL, and Redis.

```bash
docker-compose up --build
```

The service will be available at `http://localhost:8080`.

### Running Locally
1.  Ensure PostgreSQL and Redis are running.
2.  Update `config/config.go` or set environment variables if your DB/Redis credentials differ from defaults.
3.  Run the application:

```bash
go run ./cmd
```

## API Documentation

### 1. Shorten URL
**Endpoint**: `POST /app/url`
**Description**: Submit a long URL to receive a shortened version.
**Rate Limit**: 5 requests/second (Burst 10).

**Request**:
```bash
curl -X POST http://localhost:8080/app/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

**Response**:
```json
{
  "Url": "https://www.google.com",
  "short-url": "http://localhost:8080/app/url/abc123xyz"
}
```

### 2. Redirect
**Endpoint**: `GET /app/url/:code`
**Description**: Redirects to the original URL.

**Request**:
```bash
curl -v http://localhost:8080/app/url/abc123xyz
```

**Response**: `308 Permanent Redirect` to original URL.

### 3. List URLs (Admin)
**Endpoint**: `GET /admin/urls`
**Description**: Retrieve a paginated list of shortened URLs.
**Query Params**:
- `page` (default: 1)
- `page_size` (default: 10)

**Request**:
```bash
curl http://localhost:8080/admin/urls
```

**Response**:
```json
[
  {
    "ID": 1,
    "Url": "https://www.google.com",
    "ShortenedUrl": "abc123xyz",
    "clickCount": 5,
    "createdAt": "..."
  }
]
```

## Architecture & Design Decisions

### Clean Architecture
The project follows a modular structure:
- **Handlers**: HTTP layer, request parsing, and response formatting.
- **Services**: Business logic (shortening, caching, stats).
- **Repositories**: Data access layer (PostgreSQL, Redis).

### Database & Caching
- **PostgreSQL**: Primary storage for durability. Stores metadata like creation time and click counts.
- **Redis**: Caching layer for high-performance redirects.
    - **Write-Through/Look-Aside**: On creation, data is written to DB. On access, Cache is checked first; if miss, DB is queried and Cache is updated.

### Rate Limiting
Implemented a **Token Bucket** algorithm for IP-based rate limiting to protect the service from abuse (DoS).

### URL Validation
Strict validation using `net/url` ensures only valid absolute URLs (http/https) are processed.

### Docker
- **Multi-stage Build**: Uses `golang:alpine` for building and a minimal `alpine` image for the final container to keep image size small.
- **Docker Compose**: Orchestrates App, DB, and Redis for a one-command setup.
