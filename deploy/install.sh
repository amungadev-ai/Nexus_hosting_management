#!/usr/bin/env bash
# NexusPanel walking-skeleton installer (DEC-011: automated and timed).
# One command installs dependencies, builds, installs systemd units, starts
# the control plane and host agent, and prints the elapsed time.
#
# usage: install.sh [--from-dir DIR] [--repo-url URL]
#   --from-dir DIR  install from an existing source tree (skips git clone)
#   --repo-url URL  repository to clone (default: the GitHub origin)
# Environment: GITHUB_TOKEN enables private clone (used transiently, never stored).
set -euo pipefail

START=$(date +%s)
FROM_DIR=""
REPO_URL="https://github.com/amungadev-ai/Nexus_hosting_management.git"

while [ $# -gt 0 ]; do
  case "$1" in
    --from-dir) FROM_DIR="$2"; shift 2 ;;
    --repo-url) REPO_URL="$2"; shift 2 ;;
    *) echo "unknown argument: $1" >&2; exit 2 ;;
  esac
done

need_cmd() { command -v "$1" >/dev/null 2>&1; }

echo "[1/7] dependencies"
if ! need_cmd git || ! need_cmd go; then
  export DEBIAN_FRONTEND=noninteractive
  apt-get update -qq
  apt-get install -y -qq git golang-go >/dev/null
fi
need_cmd curl || { apt-get update -qq; apt-get install -y -qq curl >/dev/null; }

echo "[2/7] source"
WORKDIR=$(mktemp -d /tmp/nexuspanel-install.XXXXXX)
if [ -n "$FROM_DIR" ]; then
  cp -r "$FROM_DIR"/. "$WORKDIR"/
else
  if [ -n "${GITHUB_TOKEN:-}" ]; then
    git -c credential.helper='!f(){ echo username=nexus-installer; echo password='"$GITHUB_TOKEN"'; };f' \
      clone --depth 1 "$REPO_URL" "$WORKDIR"
  else
    git clone --depth 1 "$REPO_URL" "$WORKDIR"
  fi
fi

echo "[3/7] build"
cd "$WORKDIR"
mkdir -p bin
go build -o bin/nexus-control-plane ./cmd/control-plane
go build -o bin/nexus-host-agent ./cmd/host-agent
go build -o bin/nexus-npctl ./cmd/npctl
go build -o bin/nexus-installer ./cmd/installer

echo "[4/7] preflight (fail closed: DEC-015, ADR-R0-001, ADR-R0-002)"
./bin/nexus-installer preflight

echo "[5/7] install artifacts"
install -m 0755 bin/nexus-control-plane bin/nexus-host-agent bin/nexus-npctl bin/nexus-installer /usr/local/bin/
./bin/nexus-installer write-units -bin-dir /usr/local/bin

echo "[6/7] services"
systemctl daemon-reload
systemctl enable --now nexuspanel-host-agent.service
systemctl enable --now nexuspanel-control-plane.service

echo "[7/7] health"
HEALTH=""
for _ in $(seq 1 30); do
  if curl -fsS http://127.0.0.1:8080/healthz >/dev/null 2>&1; then HEALTH=ok; break; fi
  sleep 1
done
if [ "$HEALTH" != ok ]; then
  echo "control plane did not become healthy; recent logs:" >&2
  journalctl -u nexuspanel-control-plane -n 50 --no-pager >&2 || true
  exit 1
fi

ELAPSED=$(( $(date +%s) - START ))
echo
echo "===================================================================="
echo " NexusPanel walking skeleton installed in ${ELAPSED}s (target < 600s)"
echo "   health:        curl http://127.0.0.1:8080/healthz"
echo "   capabilities:  nexus-npctl detect"
echo "   agent status:  nexus-npctl agent status"
echo "   audit verify:  nexus-npctl audit-verify"
echo "===================================================================="
