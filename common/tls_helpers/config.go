package tls_helpers

type Config struct {
	Enabled           bool     `yaml:"enabled" env:"ENABLED" env-default:"false"`
	AllowAutoGenerate bool     `yaml:"allow_auto_generate" env:"ALLOW_AUTO_GENERATE" env-default:"true"`
	RootPublicKey     string   `yaml:"root_public_key" env:"ROOT_PUBLIC_KEY" env-default:"certs/ca.pem"`
	RootPrivateKey    string   `yaml:"root_private_key" env:"ROOT_PRIVATE_KEY" env-default:"certs/ca.key"`
	ServerPublicKey   string   `yaml:"server_public_key" env:"SERVER_PUBLIC_KEY" env-default:"certs/server.pem"`
	ServerPrivateKey  string   `yaml:"server_private_key" env:"SERVER_PRIVATE_KEY" env-default:"certs/server.key"`
	ClientPublicKey   string   `yaml:"client_public_key" env:"CLIENT_PUBLIC_KEY" env-default:"certs/client.pem"`
	ClientPrivateKey  string   `yaml:"client_private_key" env:"CLIENT_PRIVATE_KEY" env-default:"certs/client.key"`
	DNSNames          []string `yaml:"dns_names" env:"DNS_NAMES" env-default:"localhost"`
}
