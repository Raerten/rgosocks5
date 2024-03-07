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
	"strings"
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
	// Prepare authenticator config
	var authenticator []socks5.Authenticator
	if config.Cfg.ProxyUser != "" && config.Cfg.ProxyPassword != "" {
		authenticator = append(authenticator, socks5.UserPassAuthenticator{
			Credentials: socks5.StaticCredentials{
				config.Cfg.ProxyUser: config.Cfg.ProxyPassword,
			},
		})
	}

	// Prepare allowed IP networks
	var allowedIPNet []*net.IPNet
	for _, cidr := range config.Cfg.AllowedIPs {
		// Process single ip as /32 network
		if !strings.Contains(cidr, "/") {
			cidr = cidr + "/32"
		}
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			slog.Error("Parse AllowedIPs", "err", err)
			os.Exit(1)
		}
		slog.Debug("Parse AllowedIPs", "ipNet", ipNet)
		allowedIPNet = append(allowedIPNet, ipNet)
	}

	// Prepare reject IP networks
	var rejectIPNet []*net.IPNet
	for _, cidr := range config.Cfg.RejectIPs {
		// Process single ip as /32 network
		if !strings.Contains(cidr, "/") {
			cidr = cidr + "/32"
		}
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			slog.Error("Parse RejectIPs", "err", err)
			os.Exit(1)
		}
		slog.Debug("Parse RejectIPs", "ipNet", ipNet)
		rejectIPNet = append(rejectIPNet, ipNet)
	}

	// Configure socks5 server
	server := socks5.NewServer(
		socks5.WithLogger(&slogger.Socks5Logger{}),
		socks5.WithAuthMethods(authenticator),
		socks5.WithRule(&rules.ProxyRulesSet{
			AllowedIPNet: allowedIPNet,
			RejectIPNet:  rejectIPNet,
		}),
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
