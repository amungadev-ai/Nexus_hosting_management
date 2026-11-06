#!/usr/bin/env bash
# Fault-injection harness for the panel test environment (owner ruling §4
# exit criteria): process kill, service stop, disk-full, network partition.
#
# usage: fault-inject.sh <injection> [args]
#   kill-panel              SIGKILL the control plane (systemd restarts it)
#   stop-agent | start-services
#   disk-full <DIR> <MB>    fill DIR with a sparse file of MB megabytes
#   disk-clear              remove the disk-full filler
#   net-partition <PORT>    drop inbound traffic to PORT (needs nft)
#   net-heal                remove the partition rules
# Destructive injections require --yes as the first argument.
set -euo pipefail

FILLER="/var/tmp/nexuspanel-diskfull.img"

confirm() {
  if [ "${1:-}" = "--yes" ]; then shift; return; fi
  echo "refusing destructive injection without --yes" >&2
  exit 2
}

case "${1:-}" in
  --yes) shift ;;
  kill-panel|stop-agent|start-services) ;; # restartable, no confirm needed
  *) confirm "$@" ;;
esac

case "${1:-}" in
  kill-panel)
    systemctl kill --signal=SIGKILL nexuspanel-control-plane.service
    echo "SIGKILL sent to control plane (Restart=on-failure should revive it)"
    ;;
  stop-agent)
    systemctl stop nexuspanel-host-agent.service
    echo "host agent stopped"
    ;;
  start-services)
    systemctl start nexuspanel-host-agent.service nexuspanel-control-plane.service
    echo "services started"
    ;;
  disk-full)
    dir="${2:?usage: fault-inject.sh disk-full <DIR> <MB>}"
    mb="${3:?usage: fault-inject.sh disk-full <DIR> <MB>}"
    fallocate -l "${mb}M" "$FILLER" 2>/dev/null || dd if=/dev/zero of="$FILLER" bs=1M count="$mb" status=none
    mv "$FILLER" "$dir/nexuspanel-diskfull.img"
    echo "injected disk-full: ${mb}M in $dir (remove with: fault-inject.sh disk-clear)"
    ;;
  disk-clear)
    rm -f /var/tmp/nexuspanel-diskfull.img "$FILLER"
    find /var/tmp /tmp -maxdepth 2 -name 'nexuspanel-diskfull.img' -delete 2>/dev/null || true
    echo "disk-full removed"
    ;;
  net-partition)
    port="${2:?usage: fault-inject.sh net-partition <PORT>}"
    command -v nft >/dev/null 2>&1 || { echo "nft not installed" >&2; exit 3; }
    nft add table inet nexuspanel-test
    nft add chain inet nexuspanel-test input '{ type filter hook input priority -10; }'
    nft add rule inet nexuspanel-test input tcp dport "$port" drop
    echo "inbound traffic to $port dropped (heal with: fault-inject.sh net-heal)"
    ;;
  net-heal)
    nft delete table inet nexuspanel-test 2>/dev/null || true
    echo "network partition healed"
    ;;
  "")
    echo "no injection given" >&2
    exit 2
    ;;
  *)
    echo "unknown injection: $1" >&2
    exit 2
    ;;
esac
