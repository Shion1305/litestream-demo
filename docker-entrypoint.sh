#!/bin/sh
set -e

# 初回 restore
if [ ! -f /tmp/db.sqlite ]; then
  echo "Restoring DB from GCS..."
  litestream restore --if-replica-exists --config /app/litestream.yml /tmp/db.sqlite
fi

echo "Starting Litestream and server..."
litestream replicate --config /app/litestream.yml &
exec /app/main
