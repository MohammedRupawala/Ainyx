# Ainyx Backend API

A RESTful API built with Go, GoFiber, SQLC, Neon PostgreSQL, and Uber Zap. It manages users with names and dates of birth, dynamically calculating and returning their ages in user retrieval operations.

## Features
- Dynamic Age Calculation
- PostgreSQL database integration
- SQLC generated queries
- Fully validated inputs with `go-playground/validator`
- Middleware support for Request ID and request logging
- Clean separation of routes, config, logger, handlers, repositories, and services
- Multi-stage Docker build and Docker Compose configuration
- Unit testing for key service computations

## Tech Stack
- **Go** (1.25.5)
- **GoFiber** (v2)
- **SQLC** (v2)
- **pgx/v5** (PostgreSQL driver)
- **Uber Zap** (Logger)
- **go-playground/validator** (Validation)

## Directory Structure
- `/cmd/server/main.go` - Entry point
- `/config/` - App configuration package
- `/db/migrations/` - Database migrations
- `/db/sqlc/` - Generated database queries
- `/internal/` - Internal core implementation
  - `handler/` - HTTP request handlers
  - `logger/` - Zap logger helper
  - `middleware/` - Custom middleware
  - `models/` - Data structures and request/response models
  - `repository/` - Database access repository
  - `routes/` - API routes configuration
  - `service/` - Business logic and dynamic calculations

## API Documentation

### Create User
- **Method**: `POST`
- **Path**: `/users`
- **Request Body**:
```json
{
  "name": "Alice",
  "dob": "1990-05-10"
}
```
- **Response**:
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Update User
- **Method**: `PUT`
- **Path**: `/users/:id`
- **Request Body**:
```json
{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```
- **Response**:
```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### Get User by ID
- **Method**: `GET`
- **Path**: `/users/:id`
- **Response**:
```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15",
  "age": 35
}
```

### Delete User
- **Method**: `DELETE`
- **Path**: `/users/:id`
- **Response**: `204 No Content`

### List All Users
- **Method**: `GET`
- **Path**: `/users`
- **Query Parameters**:
  - `limit` (default: 10, max: 100)
  - `offset` (default: 0)
- **Response**:
```json
[
  {
    "id": 1,
    "name": "Alice Updated",
    "dob": "1991-03-15",
    "age": 35
  }
]
```

### Health Check
- **Method**: `GET`
- **Path**: `/health`
- **Response**:
```json
{
  "status": "ok"
}
```

## Running the Application

### Local Setup

1. Install Go (1.25.5 or higher).
2. Configure `.env` file at the root:
```env
DB_URL="postgresql://neondb_owner:npg_snUV5PdaLZ7r@ep-lively-wildflower-ao5dth45-pooler.c-2.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
PORT=3000
```
3. Run the project:
```bash
go run ./cmd/server
```

### Running with Docker

Run with Docker Compose:
```bash
docker-compose up --build
```

### Running Tests

Run all unit tests:
```bash
go test ./...
```
