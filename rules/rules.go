package rules

import (
	"context"
	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
	"net"
	"rgosocks/config"
	"slices"
)

type ProxyRulesSet struct {
	AllowedIPNet []*net.IPNet
	RejectIPNet  []*net.IPNet
}

func (r *ProxyRulesSet) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	if config.Cfg.DisableBind && req.Command == statute.CommandBind {
		return ctx, false
	}

	if config.Cfg.DisableAssociate && req.Command == statute.CommandAssociate {
		return ctx, false
	}

	result := len(config.Cfg.AllowedDestFqdn) == 0 && len(r.AllowedIPNet) == 0

	if !result && len(config.Cfg.AllowedDestFqdn) > 0 && slices.Contains(config.Cfg.AllowedDestFqdn, req.DestAddr.FQDN) {
		result = true
	}

	if !result && len(r.AllowedIPNet) > 0 {
		for _, ipNet := range r.AllowedIPNet {
			if ipNet.Contains(req.DestAddr.IP) {
				result = true
				break
			}
		}
	}

	if result && len(config.Cfg.RejectDestFqdn) > 0 && slices.Contains(config.Cfg.RejectDestFqdn, req.DestAddr.FQDN) {
		result = false
	}

	if result && len(r.RejectIPNet) > 0 {
		for _, ipNet := range r.RejectIPNet {
			if ipNet.Contains(req.DestAddr.IP) {
				result = false
				break
			}
		}
	}

	return ctx, result
}
