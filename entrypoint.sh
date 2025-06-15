#!/usr/bin/env bash
set -euo pipefail

WG_IF=wg0
CONF=/etc/wireguard/${WG_IF}.conf
KEY_DIR=/etc/wireguard

# Generate server keys if they don't exist
test -f "$KEY_DIR/server.key" || wg genkey | tee "$KEY_DIR/server.key" | wg pubkey > "$KEY_DIR/server.pub"

# Create basic config
echo "[Interface]" > "$CONF"
echo "Address = 10.0.0.1/24" >> "$CONF"
echo "ListenPort = 51820" >> "$CONF"
echo "PrivateKey = $(cat $KEY_DIR/server.key)" >> "$CONF"

echo "# Client peer (injected via env)" >> "$CONF"
echo "[Peer]" >> "$CONF"
echo "PublicKey = $CLIENT_PUBKEY" >> "$CONF"
echo "AllowedIPs = 0.0.0.0/0" >> "$CONF"

# Bring up interface
wg-quick up "$CONF"

# Keep container alive
tail -f /dev/null