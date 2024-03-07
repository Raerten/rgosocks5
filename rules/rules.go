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
	AllowedFQDN  []string
	RejectFQDN   []string
}

func (r *ProxyRulesSet) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	if config.Cfg.DisableBind && req.Command == statute.CommandBind {
		return ctx, false
	}

	if config.Cfg.DisableAssociate && req.Command == statute.CommandAssociate {
		return ctx, false
	}

	result := len(r.AllowedFQDN) == 0 && len(r.AllowedIPNet) == 0

	if !result && len(r.AllowedFQDN) > 0 && slices.Contains(r.AllowedFQDN, req.DestAddr.FQDN) {
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

	if result && len(r.RejectFQDN) > 0 && slices.Contains(r.RejectFQDN, req.DestAddr.FQDN) {
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
