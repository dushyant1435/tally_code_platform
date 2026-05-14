# HardCode

> *Code hard. Win harder.*  A modern LeetCode-style online judge. Forged by **Dushyant Kuntal**.

- **Backend**: Go + `gorilla/mux` + Postgres + JWT auth + bcrypt password hashes. Runs user-submitted Python in a subprocess with a per-request timeout.
- **Frontend**: React (Create React App) + MUI + Monaco editor.
- **Database**: Postgres, schema and seed data in [`db/init.sql`](db/init.sql).

## Features

- Username/email + password **signup & login**, JWT-based sessions.
- Two roles: **user** and **admin**. Admins can author problems and test cases and see every user's submissions.
- Problems list with **difficulty** (easy / medium / hard), **tags**, search, filter, and per-user solved status.
- Problem detail with **Description** / **Submissions** tabs, sample test cases, real-time **Run** (sample only) and **Submit** (full judge).
- Verdicts: `accepted`, `wrong_answer`, `time_limit_exceeded`, `runtime_error`, `no_test_cases`, `server_error`.
- Every submission is persisted with status, language, runtime, failing test number, and a short message. Browse them on the **Submissions** page or per-problem tab.
- **Playground** for running custom code with your own stdin.

---

## One-click start (Docker)

Requires **Docker Desktop** (Windows / macOS) or Docker Engine + the compose plugin (Linux). Nothing else needs to be installed — Go, Node, Python and Postgres all live inside containers.

### Windows

Double-click `start.bat`.

### macOS / Linux

```bash
chmod +x start.sh stop.sh
./start.sh
```

After the build finishes:

- Frontend: <http://localhost:3000>
- API:      <http://localhost:8080/api/v1/health>

### Demo accounts (seeded automatically)

| Username | Password   | Role  |
| -------- | ---------- | ----- |
| `admin`  | `admin123` | admin |
| `demo`   | `demo123`  | user  |

You can also create your own account from the **Sign up** page; new accounts are regular users.

### Shutting down / wiping

- Windows: `stop.bat`
- macOS / Linux: `./stop.sh`
- Wipe the database too: `docker compose down -v`

### Custom config

Copy `.env.example` → `.env` in the repo root and edit any of:

```
POSTGRES_DB / POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_PORT
SERVER_PORT
CLIENT_PORT
CORS_ALLOWED_ORIGINS
JWT_SECRET              # change for production
REACT_APP_API_URL       # baked into the React bundle at build time
```

> If you change `JWT_SECRET` after users have signed in, every existing token is invalidated. People simply have to log in again.

---

## API

| Method | Path | Auth | Description |
| ------ | ---- | ---- | ----------- |
| GET    | `/api/v1/health`                                  | public  | Liveness probe. |
| POST   | `/api/v1/auth/signup`                             | public  | `{username, email, password}` → `{token, user}`. |
| POST   | `/api/v1/auth/login`                              | public  | `{username, password}` (username or email) → `{token, user}`. |
| GET    | `/api/v1/auth/me`                                 | user    | Current user. |
| GET    | `/api/v1/problems`                                | optional| List with difficulty/tags and per-user solved flag. |
| GET    | `/api/v1/problem/{id}`                            | public  | One problem; 404 if missing. |
| POST   | `/api/v1/newproblem`                              | **admin** | Create a problem (`difficulty`, `tags[]`). |
| GET    | `/api/v1/testcases/{id}`                          | public  | All test cases. |
| GET    | `/api/v1/problem/{id}/sampleTestCases`            | public  | Sample test cases only. |
| POST   | `/api/v1/createTestCase`                          | **admin** | Add a test case. |
| POST   | `/api/v1/runCode`                                 | user    | Full judge; records a submission. |
| POST   | `/api/v1/runSampleCode`                           | public  | Sample-only run (no submission recorded). |
| POST   | `/api/v1/runCustomCode`                           | public  | Playground custom-stdin run. |
| GET    | `/api/v1/submissions`                             | user    | My recent submissions (optional `?problem_id=`). |
| GET    | `/api/v1/submissions/{id}`                        | user    | One submission (owner or admin). |
| GET    | `/api/v1/problem/{id}/submissions`                | user    | My submissions for one problem. |
| GET    | `/api/v1/admin/submissions`                       | **admin** | Every submission across every user. |

Auth header: `Authorization: Bearer <jwt>`. The client stores the token in `localStorage` under `tally.token`.

---

## Local development (without Docker)

Only needed if you want hot-reload. Otherwise just use Docker.

### Prereqs

- Go 1.22+
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
cp .env.example .env       # edit POSTGRES_URL & JWT_SECRET if needed
go run main.go
```

### Client

```bash
cd client
npm install
echo "REACT_APP_API_URL=http://localhost:8080" > .env
npm start
```

---

## Project layout

```
.
├── client/                       # React app
│   ├── Dockerfile / nginx.conf
│   └── src/
│       ├── api.js                # fetch wrapper that injects JWT
│       ├── config.js             # API_BASE + token helpers
│       ├── auth/
│       │   ├── AuthContext.jsx   # login/signup/logout + /me hydration
│       │   └── ProtectedRoute.jsx
│       ├── components/
│       │   ├── NavBar.jsx        # nav with user dropdown / admin badge
│       │   ├── Badges.jsx        # Difficulty + Status chips
│       │   └── CodeEditor.js
│       └── pages/
│           ├── Home / Login / Signup
│           ├── Problems / Problem
│           ├── Submissions / SubmissionDetail
│           ├── Playground
│           └── AdminHome / CreateProblem / CreateTestCase
├── server/                       # Go API
│   ├── Dockerfile
│   ├── cmd/hashgen/              # tiny CLI: bcrypt-hash a password
│   ├── handler/
│   │   ├── auth.go               # signup / login / JWT middleware
│   │   ├── problem.go            # problems CRUD (+ solved flag)
│   │   ├── testCases.go
│   │   ├── runCode.go            # judge + records submission
│   │   ├── submissions.go        # history endpoints
│   │   └── dbConnection.go       # singleton pool with retry
│   ├── models/                   # struct types
│   ├── router/                   # routes + CORS
│   └── main.go
├── db/init.sql                   # schema + seed data + seed users
├── docker-compose.yml
└── start.bat / stop.bat / start.sh / stop.sh
```

---

## Hosting

See the previous `Hosting` section, with one addition: in production set a **strong random** `JWT_SECRET` and a real `CORS_ALLOWED_ORIGINS`.

### Production safety checklist

The code-execution endpoint runs arbitrary Python inside the server container. Container isolation buys you a lot, but for anything past a demo you should also:

- Run user code in a separate, throwaway container per request (Docker-in-Docker or a sidecar runner).
- Drop network access from the runner.
- Add rate limiting on `/api/v1/run*` and `/api/v1/auth/*`.
- Switch JWTs to HTTP-only cookies (currently they're in `localStorage` for simplicity).
- Lower `executionTimeout` in `server/handler/runCode.go` for stricter judging.

---

## License

MIT.
