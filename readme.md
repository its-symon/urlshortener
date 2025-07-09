# 🔗 URL Shortener API

A scalable, secure, and extensible URL Shortening API built with Go using Gin, GORM, PostgreSQL, Redis, RabbitMQ, and more. Includes rate limiting, token-based authentication, click tracking via message queues, and Swagger documentation.

---

## 📦 Tech Stack

| Layer           | Tool/Library                     |
|----------------|----------------------------------|
| Web Framework  | [Gin](https://github.com/gin-gonic/gin)           |
| ORM            | [GORM](https://gorm.io/)                      |
| Database       | PostgreSQL (Production), SQLite (Test) |
| Caching        | [Redis](https://redis.io/)                    |
| Queue          | [RabbitMQ](https://www.rabbitmq.com/)         |
| Auth           | JWT-based API key/token auth       |
| Docs           | Swagger (Go Swagger v2 via `swag`)|
| Tests          | Go `testing` + [Testify](https://github.com/stretchr/testify) |
| CI/CD          | GitHub Actions                    |

---

## Features

- URL shortening (with optional custom alias)
- URL redirection with click tracking
- Auth system with token issuance
- API key validation for protected endpoints
- Rate limiting using Redis
- Asynchronous click logging via RabbitMQ
- Swagger documentation
- Dockerized + GitHub CI/CD setup

---

## Getting Started

### 1. Clone & Run

```bash
git clone https://github.com/its-symon/urlshortener.git
cd urlshortener
go mod tidy
```

### 2. Setup .env
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=12345
DB_NAME=urlshortener

REDIS_ADDR=localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

PORT=8080
```

### 3. Run with Docker Compose

```bash
docker-compose up -d
```

### 4. Running Tests
```bash
go test ./tests -v
```

### API Documentation
- Navigate to: http://localhost:8080/swagger/index.html

- Built using Swagger + Go annotations (swag init)


### Contributing
- Pull requests and suggestions welcome! Please fork and submit a PR.

---

## Author

**Symon Barua**  
[![GitHub](https://img.shields.io/badge/GitHub-its--symon-black?logo=github)](https://github.com/its-symon)  
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Symon%20Barua-blue?logo=linkedin)](https://linkedin.com/in/SymonBarua)
