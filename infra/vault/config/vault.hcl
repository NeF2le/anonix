ui            = true

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = true
}

storage "raft" {
  path    = "/vault/data"
  node_id = "vault-node-1"
}

api_addr      = "${VAULT_API_ADDR}"
cluster_addr  = "${VAULT_CLUSTER_ADDR}"
disable_mlock = true
