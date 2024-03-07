package resolver

import (
	"context"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"math/rand/v2"
	"net"
	"rgosocks/config"
	"time"
)

type DNSResolver struct {
	Cache      *cache.Cache
	DNSClient  *dns.Client
	DNSAddress string
}

func (d DNSResolver) getRandIp(ips []net.IP) net.IP {
	answerId := rand.IntN(len(ips))
	return ips[answerId]
}

func (d DNSResolver) resolve(ctx context.Context, name string, t uint16) (result []net.IP, ttl uint32, err error) {
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(name), t)
	m.RecursionDesired = true

	r, _, err := d.DNSClient.ExchangeContext(ctx, m, d.DNSAddress)
	if r == nil || r.Rcode != dns.RcodeSuccess {
		return nil, 0, err
	}

	for _, answer := range r.Answer {
		if t == dns.TypeA {
			rec := answer.(*dns.A)
			ttl = rec.Hdr.Ttl
			result = append(result, rec.A)
		} else {
			rec := answer.(*dns.AAAA)
			ttl = rec.Hdr.Ttl
			result = append(result, rec.AAAA)
		}
	}

	return result, ttl, nil
}

// Resolve implement interface NameResolver
func (d DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	if len(config.Cfg.DnsHost) == 0 {
		if config.Cfg.PreferIpv6 {
			addr6, err := net.ResolveIPAddr("ip6", name)
			if err != nil {
				slog.Debug("Resolve", "name", name, "ip", addr6.IP)

				return ctx, addr6.IP, nil
			}
		}

		addr, err := net.ResolveIPAddr("ip4", name)
		if err != nil {
			return ctx, nil, err
		}

		slog.Debug("Resolve", "name", name, "ip", addr.IP)
		return ctx, addr.IP, nil
	}

	val, expiration, found := d.Cache.GetWithExpiration(name)
	if found {
		ip := d.getRandIp(val.([]net.IP))

		slog.Debug("Resolve", "name", name, "cache", true, "expiration", expiration, "ip", ip)

		return ctx, ip, nil
	}

	var ttl uint32
	var ips []net.IP
	var err error

	if config.Cfg.PreferIpv6 {
		ips, ttl, err = d.resolve(ctx, name, dns.TypeAAAA)
		if err != nil {
			ips, ttl, err = d.resolve(ctx, name, dns.TypeA)
			if err != nil {
				return ctx, nil, err
			}
		}
	} else {
		ips, ttl, err = d.resolve(ctx, name, dns.TypeA)
		if err != nil {
			return ctx, nil, err
		}
	}

	ip := d.getRandIp(ips)
	slog.Debug("Resolve", "name", name, "cache", false, "ttl", ttl, "ip", ip)
	d.Cache.Set(name, ips, time.Duration(ttl)*time.Second)

	return ctx, ip, err
}
