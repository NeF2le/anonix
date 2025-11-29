vault {
  address = "http://vault:8200"
  tls_skip_verify = true
}

auto_auth {
  method {
    type = "approle"

    config = {
      role_id_file_path = "/vault/data/role_id"
      secret_id_file_path = "/vault/data/secret_id"
      remove_secret_id_file_after_reading = false
    }
  }
}

api_proxy {
  use_auto_auth_token = true
}

listener "tcp" {
  address = "0.0.0.0:8100"
  tls_disable = true
}
