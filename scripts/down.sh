#!/usr/bin/env bash
# Tắt BE và FE đã bật bằng scripts/up.sh.
cd "$(dirname "$0")/.."
for p in be fe; do
  [ -f ".run/$p.pid" ] && kill "$(cat ".run/$p.pid")" 2>/dev/null
  rm -f ".run/$p.pid"
done
exit 0
