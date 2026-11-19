#!/usr/bin/env bash
# Bật BE (Go, cổng BE_PORT=8080) và FE (Next.js dev, cổng FE_PORT=3000) chạy nền, chờ tới khi sẵn sàng.
# TEST_API=1 (mặc định khi chạy bằng script này): bật /api/test/* cho e2e. Tắt: bash scripts/down.sh
set -euo pipefail
cd "$(dirname "$0")/.."
BE_PORT="${BE_PORT:-8080}"
FE_PORT="${FE_PORT:-3000}"
mkdir -p .run

[ -d fe/node_modules ] || (cd fe && npm ci --no-audit --no-fund)
(cd be && go build -o ../.run/be-server ./cmd/server)
PORT="$BE_PORT" TEST_API="${TEST_API:-1}" nohup .run/be-server >.run/be.log 2>&1 &
echo $! >.run/be.pid
(cd fe && BE_URL="http://localhost:$BE_PORT" exec node_modules/.bin/next dev -p "$FE_PORT" >../.run/fe.log 2>&1) &
echo $! >.run/fe.pid

wait_for() { # url, tên
  for _ in $(seq 1 120); do curl -sf "$1" >/dev/null 2>&1 && { echo "✓ $2 $1"; return 0; }; sleep 0.5; done
  echo "✘ $2 không lên sau 60 giây — xem .run/*.log" >&2
  exit 1
}
wait_for "http://localhost:$BE_PORT/api/health" BE
wait_for "http://localhost:$FE_PORT/login" FE
