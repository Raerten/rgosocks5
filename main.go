package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"rgosocks/config"
	"rgosocks/resolver"
	"rgosocks/rules"
	"rgosocks/slogger"
	"rgosocks/stat"
	"rgosocks/version"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/things-go/go-socks5"
)

func startProxy(status *stat.Stat) {
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

	var dnsCache *cache.Cache = nil
	if config.Cfg.DnsUseCache {
		dnsCache = cache.New(1*time.Minute, 3*time.Minute)
	}

	// Configure socks5 server
	server := socks5.NewServer(
		socks5.WithLogger(&slogger.Socks5Logger{}),
		socks5.WithAuthMethods(authenticator),
		socks5.WithRule(&rules.ProxyRulesSet{
			AllowedIPNet: allowedIPNet,
			RejectIPNet:  rejectIPNet,
			AllowedFQDN:  config.Cfg.AllowedDestFQDN,
			RejectFQDN:   config.Cfg.RejectDestFQDN,
		}),
		socks5.WithResolver(&resolver.DNSResolver{
			Cache:      dnsCache,
			DNSClient:  new(dns.Client),
			DNSAddress: net.JoinHostPort(config.Cfg.DnsHost, fmt.Sprintf("%d", config.Cfg.DnsPort)),
		}),
		socks5.WithDial(status.Dial),
	)

	slog.Info("Starting Socks5 Proxy", "address", config.Cfg.ProxyAddress)
	if err := server.ListenAndServe("tcp", config.Cfg.ProxyAddress); err != nil {
		panic(err)
	}
}

func main() {
	ver, commit, date, goVer, arch := version.Info()
	slog.Info("Version", "version", ver, "commit", commit, "date", date, "go", goVer, "arch", arch)
	config.Parse()

	if config.Cfg.LogLevelDebug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Debug("Config", "env", config.Cfg)

	statusServer := stat.NewStat(
		config.Cfg.StatusEnabled,
		config.Cfg.StatusAddress,
		config.Cfg.StatusBearer,
	)

	go startProxy(statusServer)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	slog.Debug("Signal", "sig", sig.String())
}
