#!/usr/bin/env bash
# Reset the single-host test environment (owner's Leaseweb VPS) to a clean
# walking-skeleton state: stop services, wipe panel state, start fresh.
# Part of the "environment established" exit criteria (owner ruling §4).
set -euo pipefail

STATE_DIR="${NEXUS_DATA_DIR:-/var/lib/nexuspanel}"
HEALTH_URL="${NEXUS_HEALTH_URL:-http://127.0.0.1:8080/healthz}"

echo "stopping services"
systemctl stop nexuspanel-control-plane.service 2>/dev/null || true
systemctl stop nexuspanel-host-agent.service 2>/dev/null || true

echo "wiping panel state under $STATE_DIR"
rm -rf "$STATE_DIR/sites"
rm -f "$STATE_DIR/nexus.db" "$STATE_DIR/nexus.db-wal" "$STATE_DIR/nexus.db-shm"
rm -f /run/nexuspanel/agent.sock
install -d -m 0750 "$STATE_DIR"

echo "starting services"
systemctl start nexuspanel-host-agent.service
systemctl start nexuspanel-control-plane.service

for _ in $(seq 1 30); do
  curl -fsS "$HEALTH_URL" >/dev/null 2>&1 && break
  sleep 1
done
curl -fsS "$HEALTH_URL" >/dev/null

echo "test environment reset to a clean walking-skeleton state"
