# 🪚 URL-Shortener

a URL-shortener that automatically generates a short url and adds it in database.

[![Go CI](https://github.com/Zapi-web/url-shortener/actions/workflows/test.yml/badge.svg)](https://github.com/Zapi-web/url-shortener/actions/workflows/test.yml)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

## 📦 Tech Stack
This project requires:
* **Go 1.24+** (Security verified with 'govulncheck')
* **Redis** (Persistent storage with Docker Volumes)
* **Docker & Docker Compose** (Containerization)

## 🚀 Installation
1. Clone the repo
```bash
    git clone https://github.com/Zapi-web/url-shortener.git
    cd url-shortener
```
2. Start the docker compose file
```bash
    docker compose up --build
```
**That's all!**

## Usage
1. Once Docker is running you can now use curl to push a link and short it
```bash
    curl -X POST -H "Content-Type: application/json" -d '{"url": "<type a url here"}' http://localhost:8282/save
```
You will receive
```bash
{"short_url":"<short_url>"}
```
2. After first step you will have a short link you need to go on adress `http://localhost:8282/YOUR_SHORT_LINK`
3. After that you will be redirected to the your web page
