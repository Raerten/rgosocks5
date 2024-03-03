package rules

import (
	"context"
	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
	"rgosocks/config"
	"slices"
)

type ProxyRulesSet struct {
}

func (r *ProxyRulesSet) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	if config.Cfg.DisableBind && req.Command == statute.CommandBind {
		return ctx, false
	}

	if config.Cfg.DisableAssociate && req.Command == statute.CommandAssociate {
		return ctx, false
	}

	result := len(config.Cfg.AllowedDestFqdn) == 0 && len(config.Cfg.AllowedIPs) == 0

	if len(config.Cfg.AllowedDestFqdn) > 0 && slices.Contains(config.Cfg.AllowedDestFqdn, req.DestAddr.FQDN) {
		result = true
	}

	if len(config.Cfg.AllowedIPs) > 0 && slices.Contains(config.Cfg.AllowedIPs, req.DestAddr.IP.String()) {
		result = true
	}

	if len(config.Cfg.RejectDestFqdn) > 0 && slices.Contains(config.Cfg.RejectDestFqdn, req.DestAddr.FQDN) {
		result = false
	}

	if len(config.Cfg.RejectIPs) > 0 && slices.Contains(config.Cfg.RejectIPs, req.DestAddr.IP.String()) {
		result = false
	}

	return ctx, result
}
