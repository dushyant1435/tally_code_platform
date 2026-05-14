# tally_code_platform

A small competitive-coding / online-judge platform.

- **Backend**: Go + `gorilla/mux` + Postgres. Executes user-submitted Python in a sandboxed subprocess with a per-run timeout.
- **Frontend**: React (Create React App) + MUI + Monaco editor.
- **Database**: Postgres, schema in [`db/init.sql`](db/init.sql).

Three modes (visible on the home page):

1. **Playground** — free-form editor, run code with custom stdin.
2. **Coding Arena** — list of problems, write code, run against sample test cases, submit against the full suite.
3. **Code Battle** — placeholder for contests.

---

## One-click start (Docker)

Requires **Docker Desktop** (Windows / macOS) or **Docker Engine** (Linux). Nothing else needs to be installed — Go, Node, Python and Postgres all live inside containers.

### Windows

Double-click `start.bat`.

### macOS / Linux

```bash
chmod +x start.sh stop.sh
./start.sh
```

Either way, after the build finishes:

- Frontend: <http://localhost:3000>
- API:      <http://localhost:8080/api/v1/health>

To shut everything down:

- Windows: `stop.bat`
- macOS / Linux: `./stop.sh`

To wipe the database too:

```bash
docker compose down -v
```

### Custom config

Copy `.env.example` to `.env` at the repo root and edit any of:

```
POSTGRES_DB / POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_PORT
SERVER_PORT
CLIENT_PORT
CORS_ALLOWED_ORIGINS
REACT_APP_API_URL
```

`docker compose` picks `.env` up automatically.

---

## Local development (without Docker)

You only need this if you want hot-reloading. The Docker setup is more than enough for end-to-end testing.

### Prereqs

- Go 1.19+
- Node 18+
- Python 3 on `PATH` (the server shells out to `python3`)
- A running Postgres instance

### Database

```bash
psql -U postgres -c "CREATE DATABASE codedb;"
psql -U postgres -d codedb -f db/init.sql
```

### Server

```bash
cd server
cp .env.example .env       # edit POSTGRES_URL if needed
go run main.go
```

The API listens on `:8080`.

### Client

```bash
cd client
npm install
# Point the client at your local API:
echo "REACT_APP_API_URL=http://localhost:8080" > .env
npm start
```

The dev server listens on `:3000` with hot-reload.

---

## Project layout

```
.
├── client/                  # React app
│   ├── Dockerfile           # multi-stage build -> nginx
│   ├── nginx.conf           # SPA fallback + caching
│   └── src/
│       ├── config.js        # API_BASE + CURRENT_USER_ID
│       ├── components/      # CodeEditor, Textbox, NavBar
│       └── pages/           # Home, Problems, Problem, Playground, ...
├── server/                  # Go API
│   ├── Dockerfile           # multi-stage build, runtime has python3
│   ├── handler/             # HTTP handlers + DB pool + Python runner
│   ├── models/              # struct types
│   ├── router/              # routes + CORS middleware
│   └── main.go
├── db/init.sql              # schema + seed data
├── docker-compose.yml
├── .env.example
├── start.bat / start.sh     # one-click launchers
└── stop.bat  / stop.sh
```

---

## API

| Method | Path | Body / Query | Description |
| ------ | ---- | ------------ | ----------- |
| GET    | `/api/v1/health`                             | —                                 | Liveness probe. |
| GET    | `/api/v1/problems?user_id=<id>`              | —                                 | List problems with per-user `status`. |
| GET    | `/api/v1/problem/{id}`                       | —                                 | Single problem; 404 if missing. |
| POST   | `/api/v1/newproblem`                         | `Problem`                         | Create a problem. |
| GET    | `/api/v1/testcases/{id}`                     | —                                 | All test cases for a problem. |
| GET    | `/api/v1/problem/{id}/sampleTestCases`       | —                                 | Sample-only test cases. |
| POST   | `/api/v1/createTestCase`                     | `TestCase`                        | Add a test case to a problem. |
| POST   | `/api/v1/runCode`                            | `{id, code, user_id}`             | Run against all test cases; records a submission on success. |
| POST   | `/api/v1/runSampleCode`                      | `{id, code, user_id}`             | Run against sample test cases only; returns per-case results. |
| POST   | `/api/v1/runCustomCode`                      | `{code, input}`                   | Run with user-supplied stdin (Playground). |

---

## Hosting

The whole thing ships as three containers (`postgres`, `server`, `client`) with no host-specific glue, so any Docker-friendly host works. Two easy paths:

### Option A: any VPS (DigitalOcean, Hetzner, EC2, Linode, …)

1. SSH into a fresh Linux box and install Docker + Docker Compose.
2. `git clone <your-repo>` and `cd` into it.
3. Create a `.env` and set at least:
   ```
   POSTGRES_PASSWORD=<strong-random-string>
   CORS_ALLOWED_ORIGINS=https://<your-domain>
   REACT_APP_API_URL=https://<your-domain>/api
   ```
4. Put nginx / Caddy / Traefik in front to terminate TLS and reverse-proxy:
   - `https://<your-domain>/`        → `tally-client:80`
   - `https://<your-domain>/api/`    → `tally-server:8080`
5. `docker compose up -d --build`.

A minimal Caddyfile is enough:

```
your-domain.com {
    reverse_proxy /api/* tally-server:8080
    reverse_proxy *      tally-client:80
}
```

### Option B: Render / Railway / Fly.io

These all detect the per-service `Dockerfile`s. Deploy each service from the same repo:

1. **Postgres** — use the platform's managed Postgres (free tiers on Render and Railway). Run `db/init.sql` once against it.
2. **Server** — root directory `server/`, build via Dockerfile, env vars:
   - `POSTGRES_URL`
   - `PORT` (Render injects this — the server already honors it)
   - `CORS_ALLOWED_ORIGINS=https://<client-url>`
3. **Client** — root directory `client/`, build via Dockerfile, build arg:
   - `REACT_APP_API_URL=https://<server-url>`

After deploy, open the client URL.

### Production safety checklist

The code-execution endpoint runs arbitrary Python in the server container. The container isolation buys you a lot, but for anything past a demo you should also:

- Lower the per-run timeout in `server/handler/runCode.go` (`executionTimeout`).
- Run user code in a separate, throw-away container per request (Docker-in-Docker or a sidecar runner).
- Drop network access from the runner.
- Add rate limiting on the `/run*` endpoints.
- Replace the hardcoded `user_id: 123` in `client/src/config.js` with real auth.

---

## License

MIT.
