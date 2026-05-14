#!/usr/bin/env bash
# One-click start for tally_code_platform on macOS / Linux.
set -e

echo "============================================================"
echo " HardCode - one-click start  (by Dushyant Kuntal)"
echo "============================================================"
echo

if ! command -v docker >/dev/null 2>&1; then
  cat <<EOF
[X] Docker is NOT installed on this machine.

To run this app you need Docker:
  - macOS / Windows: install Docker Desktop
        https://www.docker.com/products/docker-desktop/
  - Linux:           install Docker Engine + the compose plugin
        https://docs.docker.com/engine/install/

Then re-run ./start.sh.
EOF
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  cat <<EOF
[X] Docker is installed but the daemon is not running.

  - macOS / Windows: open Docker Desktop and wait until it says
                     "Docker Desktop is running".
  - Linux:           sudo systemctl start docker

Then re-run ./start.sh.
EOF
  exit 1
fi

echo "[OK] Docker is ready. Building images and starting HardCode..."
echo "(First run downloads base images and builds the React app."
echo " Expect 5-10 minutes the first time. Subsequent runs are instant.)"
echo

docker compose up --build -d

cat <<EOF

============================================================
 All services are up.
  Frontend: http://localhost:3000
  API:      http://localhost:8080/api/v1/health

 Use ./stop.sh to shut everything down.
============================================================
EOF
