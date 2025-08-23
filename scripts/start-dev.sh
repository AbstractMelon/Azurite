#!/usr/bin/env bash
set -e

# Always run relative to repo root
ROOT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$ROOT_DIR"

echo "Starting Azurite development environment..."

# Check dependencies
command -v go >/dev/null 2>&1 || { echo "Go is not installed"; exit 1; }
command -v pnpm >/dev/null 2>&1 || { echo "pnpm is not installed"; exit 1; }
command -v make >/dev/null 2>&1 || { echo "make is not installed"; exit 1; }

# Start backend
echo "Starting backend..."
cd backend
make dev &
BACKEND_PID=$!
cd "$ROOT_DIR"

# Start frontend
echo "Starting frontend..."
cd frontend
pnpm dev &
FRONTEND_PID=$!
cd "$ROOT_DIR"

# Cleanup on exit
trap "echo 'Stopping Azurite...'; kill $BACKEND_PID $FRONTEND_PID" EXIT

# Wait for background jobs
wait
