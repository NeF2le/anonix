#!/bin/sh
set -euo pipefail

CONVERGENT_KEY="my-kek-convergent"
RANDOM_KEY="my-kek-random"
HMAC_KEY="my-hmac-key"
POLICY_NAME="tokenizer-policy"
ROLE_NAME="my-app"

DATA_DIR="/vault/data"
ROLE_ID_FILE="$DATA_DIR/role_id"
SECRET_ID_FILE="$DATA_DIR/secret_id"
INIT_JSON="$DATA_DIR/init.json"
POLICY_FILE_HOST="$DATA_DIR/${POLICY_NAME}.hcl"

VAULT_HEALTH_URL="${VAULT_ADDR:-http://127.0.0.1:8200}/v1/sys/health"

mkdir -p "$DATA_DIR"

if ! command -v jq >/dev/null 2>&1; then
  echo "ERROR: 'jq' not found on host. Install it (e.g. 'brew install jq' on macOS) and re-run the script."
  exit 1
fi

die() { echo "ERROR: $*" >&2; exit 1; }

# -----------------------------
# Waiting for Vault to become healthy
# -----------------------------
echo "Waiting for Vault to become healthy at ${VAULT_HEALTH_URL}..."
until curl -s --max-time 1 "$VAULT_HEALTH_URL" >/dev/null 2>&1; do
  sleep 1
done
echo "Vault is up!"

# -----------------------------
# Initialization
# -----------------------------
INITIALIZED_JSON="$(vault status -format=json 2>/dev/null || true)"
INITIALIZED="$(printf '%s' "$INITIALIZED_JSON" | jq -r '.initialized // false')"

if [ "$INITIALIZED" = "true" ]; then
  echo "Vault is already initialized"
  if [ ! -f "$INIT_JSON" ]; then
    echo "Warning: init.json not found. Automatic unseal/login will be skipped."
    INIT_JSON=""
  fi
else
  echo "Vault not initialized — initializing now..."
  vault operator init -format=json > "$INIT_JSON"
  echo "Init json saved to $INIT_JSON"
fi


# -----------------------------
# Unseal
# -----------------------------
_get_seal_status() {
  curl -s --max-time 2 "${VAULT_ADDR:-http://127.0.0.1:8200}/v1/sys/seal-status" || echo '{}'
}

is_sealed() {
  _get_seal_status | jq -r '.sealed // true'
}

get_progress() {
  _get_seal_status | jq -r '.progress // 0'
}

if [ "$(is_sealed)" = "true" ]; then
  if [ -n "$INIT_JSON" ] && [ -f "$INIT_JSON" ]; then
    echo "Vault sealed, will try unseal keys..."
    i=0
    for key in $(jq -r '.unseal_keys_b64[]' "$INIT_JSON"); do
      i=$((i+1))
      echo "Applying key #$i ..."
      vault operator unseal "$key" || echo "vault operator unseal returned non-zero"
      sleep 1
      prog=$(get_progress)
      sealed=$(is_sealed)
      echo "status sealed=$sealed progress=$prog"
      [ "$sealed" = "false" ] && break
    done

    if [ "$(is_sealed)" = "true" ]; then
      echo "Warning: after trying keys Vault is still sealed (progress=$(get_progress))."
    else
      echo "Vault is now unsealed."
    fi
  else
    echo "No init.json -> cannot auto-unseal"
  fi
else
  echo "Vault already unsealed."
fi

# -----------------------------
# Root login
# -----------------------------
if [ -n "$INIT_JSON" ]; then
  ROOT_TOKEN=$(jq -r '.root_token' "$INIT_JSON")
  vault login "$ROOT_TOKEN"
fi

# -----------------------------
# Setting transit secrets engine
# -----------------------------
echo "Enabling transit secrets engine (if not enabled)..."
vault secrets enable -path=transit transit >/dev/null 2>&1 || true
echo "Creating convergent key ${CONVERGENT_KEY}"
vault write -f transit/keys/${CONVERGENT_KEY} type=aes256-gcm96 derived=true convergent_encryption=true exportable=false >/dev/null 2>&1 || true
echo "Creating random key ${RANDOM_KEY}"
vault write -f transit/keys/${RANDOM_KEY} type=aes256-gcm96 exportable=false >/dev/null 2>&1 || true
echo "Creating hmac key ${HMAC_KEY}"
vault write -f transit/keys/${HMAC_KEY} type=hmac exportable=false >/dev/null 2>&1 || true
echo "Transit engine ready"

# -----------------------------
# Creating policy
# -----------------------------
cat > "$POLICY_FILE_HOST" <<EOF
path "transit/encrypt/${CONVERGENT_KEY}" {
  capabilities = ["create", "update"]
}
path "transit/decrypt/${CONVERGENT_KEY}" {
  capabilities = ["update"]
}
path "transit/datakey/plaintext/${CONVERGENT_KEY}" {
  capabilities = ["create", "update"]
}

path "transit/encrypt/${RANDOM_KEY}" {
  capabilities = ["create"]
}
path "transit/decrypt/${RANDOM_KEY}" {
  capabilities = ["create"]
}
path "transit/datakey/plaintext/${RANDOM_KEY}" {
  capabilities = ["create", "update"]
}

path "transit/hmac/${HMAC_KEY}" {
  capabilities = ["create"]
}
EOF

vault policy write "$POLICY_NAME" "$POLICY_FILE_HOST"
echo "Policy '${POLICY_NAME}' written"

# -----------------------------
# Creating app role
# -----------------------------
vault auth enable approle >/dev/null 2>&1 || echo "AppRole auth method already enabled."

vault write -f auth/approle/role/$ROLE_NAME \
  token_policies="${POLICY_NAME}" \
  token_ttl="24h" \
  token_max_ttl="48h" \
  secret_id_ttl="0" \
  secret_id_num_uses=0 \
  enforce_hostnames=false \
  bind_secret_id=true >/dev/null 2>&1 || echo "AppRole '${ROLE_NAME}' already exists or updated."

echo "AppRole created or verified."

if [ -f "$ROLE_ID_FILE" ]; then
  echo "Role ID file already exists at $ROLE_ID_FILE — skipping regeneration."
else
  vault read -format=json auth/approle/role/${ROLE_NAME}/role-id \
    | jq -r '.data.role_id' > "${ROLE_ID_FILE}.tmp" \
    && chmod 600 "${ROLE_ID_FILE}.tmp" && mv "${ROLE_ID_FILE}.tmp" "$ROLE_ID_FILE"
  echo "Saved role_id to $ROLE_ID_FILE"
fi

if [ -f "$SECRET_ID_FILE" ]; then
  echo "Wrap token file already exists at $SECRET_ID_FILE — skipping regeneration."
else
  WRAP_JSON=$(vault write -format=json -f auth/approle/role/${ROLE_NAME}/secret-id)
  echo "$WRAP_JSON" | jq -r '.data.secret_id' > "${SECRET_ID_FILE}.tmp" \
    && chmod 600 "${SECRET_ID_FILE}.tmp" && mv "${SECRET_ID_FILE}.tmp" "$SECRET_ID_FILE"
  echo "Saved wrap token to $SECRET_ID_FILE"
fi

echo "AppRole setup complete (files protected 600)."

# -----------------------------
# Ending
# -----------------------------
echo "Vault init successfully ended."

ROOT_TOKEN=""

exit 0