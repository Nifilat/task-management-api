# Taskify: Task Management API

Backend for Railtronics' task management platform, written in Go (Gin + GORM) on PostgreSQL. Users register and log in, manage their own tasks and profile; admins manage all users and tasks. Access is controlled with JWT access tokens, role-based access control (RBAC) and permission-level attribute checks (ABAC).

> Status: **skeleton**. Search for the TODOs to find what to build next: `grep -rn "TODO(" backend database-migrations`

## Project layout

```
backend/
  main.go                       router and wiring
  internal/handlers/            HTTP handlers (register, auth, refresh, task, users)
  internal/services/            business logic and database queries
  internal/models/              GORM models (user, token, task, role/permission)
  internal/middleware/          auth middleware (JWT verification, role/permission checks)
  internal/repositories/        database connection
  internal/utils/               env helpers, JWT and UUID helpers
database-migrations/migrations/ numbered SQL migrations (000001 to 000008)
frontend/                       React UI from the starter
docker-compose.yml              Postgres, migrations, backend and frontend
```

## Setup

Requires Go 1.23+, PostgreSQL 13+ (or Docker), and optionally Postman.

1. Create your local config (`.env` is git-ignored, never commit it):

```bash
cp .env.example .env
# set DB_PASSWORD, and a JWT_SECRET from:  openssl rand -hex 32
```

### Option A: everything with Docker Compose

```bash
docker compose up -d
```

The frontend is at http://localhost and the API at http://localhost:8080.

### Option B: local backend + Postgres

Start Postgres (Docker example from the course; use the `DB_PASSWORD` from `.env`):

```bash
set -a; source .env; set +a
docker run -itd --name postgres --restart=always -p 5432:5432 \
  -e POSTGRES_PASSWORD=$DB_PASSWORD -e POSTGRES_USER=taskmanager -e POSTGRES_DB=taskmanager postgres:17
```

Run the migrations with [golang-migrate](https://github.com/golang-migrate/migrate) (`brew install golang-migrate`):

```bash
migrate -path database-migrations/migrations \
  -database "postgres://taskmanager:$DB_PASSWORD@localhost:5432/taskmanager?sslmode=disable" -verbose up
```

> The new migrations (000003 to 000008) start as TODO comments. Fill them in **before** the first `up`; if you already ran them, `migrate ... down 6` then `up` after editing.

Run the API:

```bash
cd backend
go run main.go
```

Connect with psql: `psql -h localhost -p 5432 -U taskmanager -d taskmanager`

## Endpoints

| Method | Path | Access |
|---|---|---|
| POST | `/api/v1/auth/register` | public |
| POST | `/api/v1/auth/login` | public |
| POST | `/api/v1/auth/refresh` | public (valid refresh_token) |
| GET / POST | `/api/v1/tasks` | user (own) / admin (all) |
| GET / PUT | `/api/v1/tasks/:id` | owner or admin |
| DELETE | `/api/v1/tasks/:id` | admin |
| GET | `/api/v1/users/:user_id/tasks` | owner or admin |
| GET | `/api/v1/users/profile` | the caller |
| GET | `/api/v1/users/profile/:user_id` | owner or admin |
| GET | `/api/v1/users` | admin |
| DELETE | `/api/v1/users/:user_id` | admin |

## Roadmap (TODO numbers match the code comments)

1. [ ] **Registration**: validate input, bcrypt-hash the password, store in `users`
2. [ ] **Login**: verify credentials, issue a JWT access_token and refresh_token (1 hour each), store the refresh_token in `tokens`
3. [ ] **Refresh token** endpoint: validate against `tokens`, return a new access_token and refresh_token
4. [ ] **RBAC/ABAC schema** (migrations 000004 to 000008): `roles`, `user_roles`, `permissions`, `role_permissions`, plus default data (roles `user` and `admin`, an admin user, role-permission mappings)
5. [ ] **Wire roles and permissions** into registration (default `user` role) and login (roles and permissions in the token)
6. [ ] **Auth middleware**: verify the access_token on each request, enforce role and permission checks
7. [ ] **Tasks**: `tasks` migration (000003) and CRUD handlers and services
8. [ ] **Fix the SQL injection** in the Get User Profile service (`GetUserProfileMalicious`) and finish the user endpoints
9. [ ] Optional extras: Swagger, Postman collection, unit and end-to-end tests, pagination and sorting on tasks, README design notes, rate limiting on login

## Security notes to carry through

- Passwords are stored only as bcrypt hashes and are never returned by the API.
- Secrets (`DB_PASSWORD`, `JWT_SECRET`) come from the environment, not from source.
- Always use bound query parameters; never build SQL with string formatting.
- Return the same error for an unknown user and a wrong password.
