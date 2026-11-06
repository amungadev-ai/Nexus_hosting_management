#!/usr/bin/env bash
# Golden snapshot of the panel test environment (owner ruling §4 exit criteria):
# tarball of state, config and unit files under /var/backups/nexuspanel-test/.
# Keep the newest 10; older snapshots are pruned.
set -euo pipefail

STATE_DIR="${NEXUS_DATA_DIR:-/var/lib/nexuspanel}"
BACKUP_DIR="${NEXUS_BACKUP_DIR:-/var/backups/nexuspanel-test}"
STAMP=$(date +%Y%m%d-%H%M%S)

mkdir -p "$BACKUP_DIR"
SNAP="$BACKUP_DIR/golden-$STAMP.tar.gz"

tar -czf "$SNAP" \
  --ignore-failed-read \
  "$STATE_DIR" \
  /etc/nexuspanel \
  /etc/systemd/system/nexuspanel-control-plane.service \
  /etc/systemd/system/nexuspanel-host-agent.service \
  2>/dev/null || true

chmod 0600 "$SNAP"
echo "snapshot: $SNAP ($(du -h "$SNAP" | cut -f1))"

# prune to the newest 10
ls -1t "$BACKUP_DIR"/golden-*.tar.gz 2>/dev/null | tail -n +11 | while read -r old; do
  rm -f "$old"
  echo "pruned: $old"
done
