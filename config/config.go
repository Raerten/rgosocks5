package config

import (
	"github.com/caarlos0/env/v10"
	"log"
)

type Config struct {
	ProxyHost        string   `env:"PROXY_HOST" envDefault:"0.0.0.0"`
	ProxyPort        int      `env:"PROXY_PORT" envDefault:"1080"`
	ProxyUser        string   `env:"PROXY_USER" envDefault:""`
	ProxyPassword    string   `env:"PROXY_PASS" envDefault:""`
	ProxyAddress     string   `env:"PROXY_ADDRESS,expand" envDefault:"$PROXY_HOST:$PROXY_PORT"`
	AllowedDestFQDN  []string `env:"PROXY_ALLOWED_DEST_FQDN" envDefault:""`
	RejectDestFQDN   []string `env:"PROXY_REJECT_DEST_FQDN" envDefault:""`
	AllowedIPs       []string `env:"PROXY_ALLOWED_IPS" envDefault:""`
	RejectIPs        []string `env:"PROXY_REJECT_IPS" envDefault:""`
	DisableBind      bool     `env:"PROXY_DISABLE_BIND" envDefault:"false"`
	DisableAssociate bool     `env:"PROXY_DISABLE_ASSOCIATE" envDefault:"false"`
	DnsHost          string   `env:"DNS_HOST" envDefault:""`
	DnsPort          int      `env:"DNS_PORT" envDefault:"53"`
	DnsUseCache      bool     `env:"DNS_USE_CACHE" envDefault:"true"`
	PreferIpv6       bool     `env:"PREFER_IPV6" envDefault:"false"`
	LogLevelDebug    bool     `env:"LOG_LEVEL_DEBUG" envDefault:"false"`

	StatusEnabled bool   `env:"STATUS_ENABLED" envDefault:"false"`
	StatusHost    string `env:"STATUS_HOST" envDefault:"0.0.0.0"`
	StatusPort    int    `env:"STATUS_PORT" envDefault:"2080"`
	StatusAddress string `env:"STATUS_ADDRESS,expand" envDefault:"$STATUS_HOST:$STATUS_PORT"`
	StatusBearer  string `env:"STATUS_TOKEN" envDefault:""`
}

var Cfg = Config{}

func Parse() {
	if err := env.Parse(&Cfg); err != nil {
		log.Fatal(err)
	}
}
