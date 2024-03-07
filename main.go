package main

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/things-go/go-socks5"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"rgosocks/config"
	"rgosocks/resolver"
	"rgosocks/rules"
	"rgosocks/slogger"
	"syscall"
	"time"
	_ "time/tzdata"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func startProxy() {
	var authenticator []socks5.Authenticator
	if config.Cfg.ProxyUser != "" && config.Cfg.ProxyPassword != "" {
		authenticator = append(authenticator, socks5.UserPassAuthenticator{
			Credentials: socks5.StaticCredentials{
				config.Cfg.ProxyUser: config.Cfg.ProxyPassword,
			},
		})
	}

	server := socks5.NewServer(
		socks5.WithLogger(&slogger.Socks5Logger{}),
		socks5.WithAuthMethods(authenticator),
		socks5.WithRule(&rules.ProxyRulesSet{}),
		socks5.WithResolver(&resolver.DNSResolver{
			Cache:      cache.New(1*time.Minute, 3*time.Minute),
			DNSClient:  new(dns.Client),
			DNSAddress: net.JoinHostPort(config.Cfg.DnsHost, fmt.Sprintf("%d", config.Cfg.DnsPort)),
		}),
	)

	slog.Info("Starting Socks5 Proxy", "address", config.Cfg.ProxyAddress)
	if err := server.ListenAndServe("tcp", config.Cfg.ProxyAddress); err != nil {
		panic(err)
	}
}

func main() {
	slog.Info("Version", "version", version, "commit", commit, "date", date)
	config.Parse()

	if config.Cfg.LogLevelDebug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Debug("Config", "env", config.Cfg)

	go startProxy()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	slog.Debug("Signal", "sig", sig.String())
}
