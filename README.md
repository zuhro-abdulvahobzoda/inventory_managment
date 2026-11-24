# inventory-api

A simple inventory system: items, their quantity and location, and a change history.
Go + Gin + PostgreSQL, JWT authentication, two roles: `admin` and `employee`.

## Running

1. Copy `.env.example` to `.env`, set `JWT_SECRET` and fill in the `POSTGRES_*` values.
2. Create the database (`createdb inventory`). Migrations are applied **automatically on startup**
   (see [Migrations](#migrations)), so there is nothing else to run by hand.
3. `go run ./cmd/api`
4. Tests: `go test ./...`

API documentation — `docs/v1/openapi.yaml` (can be opened in Swagger Editor).

## Migrations

Schema changes live in `migrations/` as numbered pairs of files (`NNN_name.up.sql` / `NNN_name.down.sql`)
and are applied with [golang-migrate](https://github.com/golang-migrate/migrate). The files are embedded into the
binary (`migrations/embed.go`), and `postgres.Migrate` runs all pending "up" migrations on every start;
when the database is already up to date it does nothing. The applied version is stored in the `schema_migrations` table.

To add a migration, create the next pair of files, e.g. `007_something.up.sql` and `007_something.down.sql`.

If your database was created earlier by running the `.sql` files with `psql`, that is fine: migrations 001–004
use `IF NOT EXISTS`, so they are safely re-applied and only the new ones take effect.

## Roles and authorization

The role is stored in the `users` table and checked on the server: on every request the user and their role are loaded from the database by the `user_id` from the token; the token itself carries no role.

- The role is chosen at sign-up: `POST /sign-up` accepts an optional `"role"` field (`"admin"` or `"employee"`, `employee` by default). The first registered user is **not** made admin automatically.
- **Only one user can be `admin`.** This is enforced by the database (a partial unique index on `users`), so signing up as `admin` while an admin already exists returns `409`.
- To hand the admin role to somebody else, first change the current admin: `UPDATE users SET role = 'employee' WHERE role = 'admin';` and then `UPDATE users SET role = 'admin' WHERE username = '...';`

| Request | Who can do it |
|---------|---------------|
| `POST /sign-up`, `POST /login` | everyone, no token needed |
| `GET /users/me`, `GET /items`, `GET /items/{id}`, `GET /items/{id}/history`, `GET /locations` | `admin` and `employee` |
| `POST /items`, `PUT /items/{id}`, `DELETE /items/{id}`, `POST /locations` | `admin` only |

Status codes: no token or an invalid one — `401`; valid token but the role is not allowed — `403`.

All paths start with `/api/v1`.

## Tokens (JWT)

`/login` returns a pair of tokens signed with HS256 using `JWT_SECRET`:
- `access_token` — sent as `Authorization: Bearer <token>`, lives for `ACCESS_TOKEN_TTL` (default `15m`);
- `refresh_token` — lives for `REFRESH_TOKEN_TTL` (default `720h`) and is not accepted as an access token.

A token contains `user_id`, `type` (`access`/`refresh`) and the standard `iat`, `exp`. Tokens are not stored anywhere, so there is no logout endpoint: the client simply discards the token. Passwords are stored only as bcrypt hashes.

## Items, locations and history

- `quantity` — the amount, cannot be negative; `location_id` — an optional reference to a place from `/locations`.
- A location has a `name` and its exact position on the map: `latitude` (-90..90) and `longitude` (-180..180) in decimal degrees. Both are required when creating a location. Locations created before coordinates were introduced return `null` for both.
- `DELETE` is a soft delete: the item disappears from the list but stays in the database (`items.deleted_at`), so its history is preserved.
- Every create, update and delete is written to `history` in the same transaction as the item change itself: who (`user_id`), when (`changed_at`), what they did (`action`) and the state before and after (`old_data`, `new_data`).

## Example

```
# sign up as the admin and log in
curl -X POST localhost:8080/api/v1/sign-up -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"12345678","role":"admin"}'
curl -X POST localhost:8080/api/v1/login -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"12345678"}'

# then use the token from the response
TOKEN=...
curl -X POST localhost:8080/api/v1/locations -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Warehouse A","latitude":38.5598,"longitude":68.787}'
curl -X POST localhost:8080/api/v1/items -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' -d '{"name":"Laptop","quantity":5,"location_id":1}'
curl localhost:8080/api/v1/items -H "Authorization: Bearer $TOKEN"
curl localhost:8080/api/v1/items/1/history -H "Authorization: Bearer $TOKEN"
```

## Logging

Logs are written with `slog` to stdout at the `LOG_LEVEL` level (default `info`):
- `ENVIRONMENT=local` — human-readable text;
- otherwise — JSON.
