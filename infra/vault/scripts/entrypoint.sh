#!/bin/sh
set -euo pipefail

vault server -config=/vault/config/vault.hcl &
VAULT_PID=$!

echo "Running init script..."
/vault/scripts/init-vault.sh || true

wait $VAULT_PID
