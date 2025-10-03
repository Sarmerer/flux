#!/usr/bin/env bash
set -euo pipefail

export DATABASE_URL=${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/flow?sslmode=disable}
export WEB_ORIGIN=${WEB_ORIGIN:-http://localhost:5173}
export PORT=${PORT:-8080}

docker compose -f deploy/docker-compose.yml up -d db

echo "Waiting for database..."
until docker exec $(docker compose -f deploy/docker-compose.yml ps -q db) pg_isready -U postgres >/dev/null 2>&1; do
  sleep 1
done

echo "Starting backend..."
pushd backend >/dev/null
if [ ! -f go.sum ]; then
  go mod tidy
fi
GOFLAGS="-tags=embed" go run ./cmd/server &
BACKEND_PID=$!
popd >/dev/null

echo "Starting frontend..."
pushd frontend >/dev/null
npm install --silent
npm run dev -- --port 5173
popd >/dev/null

trap "kill $BACKEND_PID || true" EXIT


