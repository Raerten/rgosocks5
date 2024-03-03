package main

import (
	"github.com/things-go/go-socks5"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"rgosocks/config"
	"rgosocks/rules"
	"syscall"
	_ "time/tzdata"
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
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithAuthMethods(authenticator),
		socks5.WithRule(&rules.ProxyRulesSet{}),
	)

	slog.Info("Starting Socks5 Proxy", "address", config.Cfg.ProxyAddress)
	if err := server.ListenAndServe("tcp", config.Cfg.ProxyAddress); err != nil {
		panic(err)
	}
}

func main() {
	config.Parse()

	go startProxy()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	slog.Info("signal", "sig", sig.String())
}
