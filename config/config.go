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
	AllowedDestFqdn  []string `env:"PROXY_ALLOWED_DEST_FQDN" envDefault:""`
	RejectDestFqdn   []string `env:"PROXY_REJECT_DEST_FQDN" envDefault:""`
	AllowedIPs       []string `env:"PROXY_ALLOWED_IPS" envDefault:""`
	RejectIPs        []string `env:"PROXY_REJECT_IPS" envDefault:""`
	DisableBind      bool     `env:"PROXY_DISABLE_BIND" envDefault:"false"`
	DisableAssociate bool     `env:"PROXY_DISABLE_ASSOCIATE" envDefault:"false"`
}

var Cfg = Config{}

func Parse() {
	if err := env.Parse(&Cfg); err != nil {
		log.Fatal(err)
	}
}
