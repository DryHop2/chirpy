# Chirpy

Chirpy is a secure microblogging API service built in Go, featuring user authentication, chirp posting, and integration with a simulated payment provider (Polka). 

This project was built as part of a guided backend development course to gain hands-on experience with building secure and maintainable HTTP servers in Go. [Boot.dev](https://www.boot.dev/) 

## Features

- Create, retrieve, and delete short-form "chirps"
- Secure user authentication with hashed passwords and JWTs
- Refresh token system for long-lived sessions
- Token revocation support
- Webhook endpoint for upgrading users to "Chirpy Red"
- Admin routes to reset or monitor server state
- Metrics tracking via middleware
- Configurable via `.env`
- Cleanly organized using Go packages and `sqlc`

## API Endpoints

### Auth & Users
- `POST /api/users` – Register
- `POST /api/login` – Login (returns JWT + refresh token)
- `PUT /api/users` – Update user (requires JWT)

### Chirps
- `POST /api/chirps` – Create a chirp (requires JWT)
- `GET /api/chirps` – List chirps (supports `author_id` and `sort=asc|desc`)
- `GET /api/chirps/{chirpID}` – Get single chirp
- `DELETE /api/chirps/{chirpID}` – Delete chirp (author only)

### Admin
- `GET /api/healthz` – Health check
- `GET /admin/metrics` – View file server hit count
- `POST /admin/reset` – Reset all users

### Polka Webhooks
- `POST /api/polka/webhooks` – Marks user as "Chirpy Red" (requires API key)

## Setup

1. Clone the repo
2. Create a `.env` file:
   ```dotenv
   DB_URL=postgres://user:pass@localhost:5432/chirpy
   JWT_SECRET=your_jwt_secret
   POLKA_KEY=polka_key
   PLATFORM=dev
   ```
3. Run migrations with `goose`
4. Start the server: `go run main.go`

## Notes
- All token-protected endpoints use Bearer auth: `Authorization: Bearer <token>`
- Webhooks use API key auth: `Authorization: ApiKey <POLKA_KEY>`

---

Built for hands-on experience with secure Go web development.

