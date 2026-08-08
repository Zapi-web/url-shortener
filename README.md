# URL-Shortener

a URL-shortener that automatically generates a short url and adds it in database.

[![Main Pipeline](https://github.com/Zapi-web/url-shortener/actions/workflows/main.yaml/badge.svg)](https://github.com/Zapi-web/url-shortener/actions/workflows/main.yaml)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

## Features
* **Snowflake & Base62 Encoding**: Fast, collision-free short ID generation with $O(1)$ memory encoding.
* **Caching Layer**: Redis cache with automatic fallback to `fakeCache` if Redis is unreachable.
* **PostgreSQL Storage**: Indexed on `expired_at` for efficient database cleanup routines.
* **Observability**: Built-in VictoriaMetrics integration tracking request counts, latencies, and cache stats.

## Tech Stack
* **Go 1.26+**
* **PostgreSQL** (Main database)
* **Redis** (Cache layer)
* **Docker & Docker Compose** (Containerization)

## Installation
1. Clone the repo
```bash
    git clone https://github.com/Zapi-web/url-shortener.git
    cd url-shortener
```
2. Start the docker compose file (**NOTE**:Remember to update default database passwords in docker-compose.yaml before deploying to production.)
```bash
    docker compose up --build
```
**That's all!**

## Usage
1. Once Docker is running you can now use curl to push a link and short it
```bash
    curl -X POST http://localhost:8080/api/v1/ \
    -H "Content-Type: application/json" \
    -d '{"url": "https://example.com", "user_id": 123}'
```
**Response**
```bash
{"short_url":"2u6dLGX2rZK"}
```
2. Access the short URL in a browser or via curl:
``` bash
    curl -i localhost:8080/2u6dLGX2rZK
```
3. After that you will be redirected to the your web page
