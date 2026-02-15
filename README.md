# URL Shortener Service (Golang)

A simple URL shortener service built using **Golang** that provides:

- Shorten long URLs using a REST API
- Returns the same short code for the same URL (no duplicates)
- Redirect support (`?code=<code>` → original URL)
- In-memory storage (no DB)
- Metrics API to show top 3 most shortened domains
- Unit tests
- Docker support

---

## Features

### URL Shortening
- Accepts a URL via REST API.
- Generates a short code.
- Stores the mapping in memory.
- If the same URL is sent again, returns the **same short code**.

Example output:
```json
{
    "success": true,
    "data": {
        "short_url": "http://localhost:8080?code=0PFQZkhh",
        "code": "0PFQZkhh"
    }
}
```

### Redirection
- When user visits `?code=<code>`, the service redirects to the original URL.

### Metrics API
- Returns the top 3 domain names that were shortened most frequently.

Example output:
```json
[
  { "domain": "udemy.com", "count": 6 },
  { "domain": "youtube.com", "count": 4 },
  { "domain": "wikipedia.org", "count": 2 }
]
```

## Requirements
- Go 1.22+
- Docker (optional)
- Make (optional)

## Development
This project includes a Makefile to streamline common tasks. Use the following commands:

### URL Shortening
- Show build info: Displays the current Git commit and intended Docker tag.
```bash
make info
```
- Compile the binary: Builds the Go application for Linux/AMD64.
```bash
make build
```
This generates the binary at:
```bash
bin/manager
```

- Run unit tests
```bash
make go-test
```


## Docker

### Build Docker image using Makefile
```bash
make docker-build IMAGE_NAME=<your image name> IMAGE_TAG=<your image tag>
```

### Run Docker Container
```bash
docker run -p 8080:8080 <image>:<tag>
```

## API Endpoints

### 1) Shorten URL

### POST /
Request:
```json
{
  "url": "https://youtube.com/watch?v=abc123"
}
```
Response:
```json
{
  "success": true,
  "data": {
    "short_url": "http://localhost:8080?code=0PFQZkhh",
    "code": "0PFQZkhh"
  }
}
```

### 2) Redirect to Original URL

### GET /?code=<code>
Example:
```bash
curl -v "http://localhost:8080?code=0PFQZkhh"
```

This will redirect to the original URL.

### 3) Metrics API

### GET /metrics
Example:
```bash
curl http://localhost:8080/metrics
```
Response:
```json
[
  { "domain": "udemy.com", "count": 6 },
  { "domain": "youtube.com", "count": 4 },
  { "domain": "wikipedia.org", "count": 2 }
]
```

### Notes
- This project uses in-memory storage, so restarting the server clears all stored URLs.
- The same URL always returns the same short code.
- Thread safety is ensured using locks (sync.RWMutex).

