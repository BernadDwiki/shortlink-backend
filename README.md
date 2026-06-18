# Shortlink Backend

[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-v1.x-00ADD8?logo=gin&logoColor=white)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![JWT](https://img.shields.io/badge/JWT-Authentication-orange)](https://jwt.io/)

Backend service for the Shortlink application built with Go, Gin, PostgreSQL, and JWT Authentication.

---

## Project Description

Shortlink Backend is a RESTful API that allows users to:

- Register and login securely using JWT authentication
- Create shortened URLs
- Redirect short URLs to their original destinations
- Manage created links
- Track link statistics such as click counts
- View user dashboard data

The application follows a layered architecture consisting of:

- Controller Layer
- Service Layer
- Repository Layer
- Middleware Layer

This separation improves maintainability, scalability, and testability.

---

## Technology Stack

### Backend

- Go
- Gin Framework
- JWT Authentication
- PostgreSQL
- Swagger/OpenAPI

### Development Tools

- Docker
- Docker Compose
- Makefile
- golang-migrate

---

## Project Structure

```text
shortlink-backend/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── controller/
│   ├── service/
│   ├── repository/
│   ├── middleware/
│   └── model/
│
├── db/
│   └── migrations/
│
├── docs/
│   ├── swagger.json
│   └── swagger.yaml
│
├── Dockerfile
├── docker-compose.yml
├── example.env
├── Makefile
├── .env
└── README.md
```

---

# Setup

## 1. Clone Repository

```bash
git clone https://github.com/<your-username>/<your-repository>.git
cd <your-repository>/shortlink-backend
```

If the repository contains both backend and frontend folders in the same root, clone the repo and then change into `shortlink-backend`.

---

## 2. Environment Configuration

Copy `example.env` to `.env` and update values as needed:

```bash
cp example.env .env
```

Example `.env` contents:

```env
APP_HOST=0.0.0.0
APP_PORT=8080

DB_URL=postgres://shortlink:shortlinkpassword@localhost:5432/shortlink?sslmode=disable

JWT_SECRET=your-secret-key
JWT_EXPIRATION_MINUTES=30

BASE_URL=http://localhost:8080
```

> Note: The application reads `DB_URL` directly. If you use a `.env` file with separate `DB_HOST`, `DB_USER`, and friends, make sure `DB_URL` is set or expanded correctly.

### Environment Variables

| Variable | Description |
|-----------|-------------|
| APP_HOST | Application host |
| APP_PORT | Application port |
| DB_URL | PostgreSQL connection string |
| JWT_SECRET | Secret key used for JWT signing |
| JWT_EXPIRATION_MINUTES | JWT expiration time |
| BASE_URL | Base URL for generated short links |

---

## 3. Database Setup

### Option 1 — Docker PostgreSQL (Recommended)

```bash
docker run --name shortlink-db \
  -e POSTGRES_USER=shortlink \
  -e POSTGRES_PASSWORD=shortlinkpassword \
  -e POSTGRES_DB=shortlink \
  -p 5432:5432 \
  -v "$(pwd)/db/migrations:/docker-entrypoint-initdb.d" \
  -d postgres:15-alpine
```

The SQL migration files inside:

```text
db/migrations
```

will be executed automatically during database initialization.

### Option 2 — Local PostgreSQL

Create database:

```sql
CREATE DATABASE shortlink;
```

Update `DB_URL` in `.env` and run:

```bash
make migrate-up
```

Rollback migrations:

```bash
make migrate-down
```

---

## 4. Install Dependencies

```bash
go mod download
```

## 5. Run Application

Load environment variables, then start the server:

```bash
# in bash or zsh
set -a
source .env
set +a

go run cmd/main.go
```

If you prefer, set the environment variables directly in your shell before running the app.

Server will run at:

```text
http://localhost:8080
```

---

# Docker Setup

Run all services using Docker Compose:

```bash
docker compose up --build
```

Run in background:

```bash
docker compose up -d --build
```

Stop services:

```bash
docker compose down
```

---

# API Documentation

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

---

# API Endpoints

## Authentication

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /api/register | POST | Register new user |
| /api/login | POST | User login |

## Links

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /api/links | GET | Get all user links |
| /api/links | POST | Create short link |
| /api/links/:id | DELETE | Delete link |
| /:shortCode | GET | Redirect short URL |

## Dashboard

| Endpoint | Method | Description |
|-----------|----------|-------------|
| /api/dashboard | GET | Dashboard statistics |

---

# Design Decisions

### Layered Architecture

- Controller → HTTP request/response handling
- Service → Business logic
- Repository → Database operations
- Middleware → Authentication and CORS

### JWT Authentication

JWT is used because it provides:

- Stateless authentication
- Easy frontend integration
- Scalability for REST APIs

### PostgreSQL

PostgreSQL was chosen because it provides:

- Reliable relational storage
- Strong indexing support
- Excellent performance

---

# Troubleshooting

## CORS Error

If frontend runs on a different origin (for example `http://localhost:3000`), update the allowed origins inside:

```text
internal/middleware/cors.middleware.go
```

---

## Database Connection Error

Verify PostgreSQL is running:

```bash
docker ps
```

Ensure the `DB_URL` value is correct.

---

# Future Improvements

- Custom alias support
- QR Code generation
- Link expiration
- Advanced analytics
- Unit testing
- Refresh token authentication

---

# Related Project

Frontend Repository:

https://github.com/BernadDwiki/shortlink-frontend

---

# License

This project is licensed under the MIT License.