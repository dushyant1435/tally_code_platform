@echo off
setlocal

REM One-click start for tally_code_platform on Windows.
REM Double-click this file, or run it from cmd / PowerShell.

echo ============================================================
echo  HardCode - one-click start  (by Dushyant Kuntal)
echo ============================================================
echo.

REM --- check Docker is installed -------------------------------
where docker >nul 2>nul
if errorlevel 1 (
    echo [X] Docker is NOT installed on this machine.
    echo.
    echo To run this app you need Docker Desktop:
    echo   1. Download:  https://www.docker.com/products/docker-desktop/
    echo   2. Install it ^(default options are fine^).
    echo   3. Open Docker Desktop and wait until it says "Engine running".
    echo   4. Double-click start.bat again.
    echo.
    pause
    exit /b 1
)

REM --- check Docker daemon is running --------------------------
docker info >nul 2>nul
if errorlevel 1 (
    echo [X] Docker is installed but the Docker daemon is not running.
    echo.
    echo Open "Docker Desktop" from the Start menu, wait until the icon
    echo in the system tray says "Docker Desktop is running" ^(this can
    echo take a minute on a fresh boot^), then run start.bat again.
    echo.
    pause
    exit /b 1
)

REM --- build + start -------------------------------------------
echo [OK] Docker is ready. Building images and starting HardCode...
echo.
echo (First run downloads base images and builds the React app.
echo  Expect 5-10 minutes the first time. Subsequent runs are instant.)
echo.

docker compose up --build -d
if errorlevel 1 (
    echo.
    echo [X] Build or startup failed. Recent logs:
    echo --------------------------------------------------------
    docker compose logs --tail=40
    echo --------------------------------------------------------
    echo.
    pause
    exit /b 1
)

echo.
echo ============================================================
echo  All services are up.
echo   Frontend: http://localhost:3000
echo   API:      http://localhost:8080/api/v1/health
echo.
echo  Use stop.bat to shut everything down.
echo ============================================================
echo.
pause
endlocal
