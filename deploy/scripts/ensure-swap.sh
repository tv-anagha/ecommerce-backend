#!/bin/sh
set -eu

SWAP_SIZE="${SWAP_SIZE:-2G}"
SWAP_FILE="${SWAP_FILE:-/swapfile}"

if swapon --show 2>/dev/null | grep -q .; then
  echo "Swap already enabled:"
  swapon --show
  free -h
  exit 0
fi

echo "No swap detected. Creating ${SWAP_SIZE} swap at ${SWAP_FILE}..."

if [ ! -f "$SWAP_FILE" ]; then
  if ! fallocate -l "$SWAP_SIZE" "$SWAP_FILE" 2>/dev/null; then
    dd if=/dev/zero of="$SWAP_FILE" bs=1M count=2048 status=progress
  fi
  chmod 600 "$SWAP_FILE"
  mkswap "$SWAP_FILE"
fi

swapon "$SWAP_FILE"
echo "Swap enabled:"
swapon --show
free -h

if ! grep -q "$SWAP_FILE" /etc/fstab 2>/dev/null; then
  echo "$SWAP_FILE none swap sw 0 0" >> /etc/fstab
  echo "Added $SWAP_FILE to /etc/fstab for persistence across reboots."
fi
