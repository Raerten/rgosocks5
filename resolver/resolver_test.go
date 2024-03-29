package resolver

import (
	"context"
	"fmt"
	"github.com/foxcpp/go-mockdns"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	"net"
	"net/netip"
	"rgosocks/config"
	"testing"
	"time"
)

var (
	ipv4       = "192.168.2.1"
	ipv6       = "fd02:b47a:7800:cf3c:1234:1234:1234:1111"
	ipv62      = "fd02:b47a:7800:cf3c:0000:0000:0000:0001"
	ipv62Short = "fd02:b47a:7800:cf3c::1"
)

type ResolverTestSuite struct {
	suite.Suite
	srv *mockdns.Server
}

func TestCalculatorTestSuite(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}

func (suite *ResolverTestSuite) SetupSuite() {
	suite.srv, _ = mockdns.NewServer(map[string]mockdns.Zone{
		"1.example.com.": {
			A:    []string{ipv4, "192.168.1.2"},
			AAAA: []string{ipv6},
		},
		"2.example.com.": {
			AAAA: []string{ipv62},
		},
		"3.example.com.": {
			A: []string{ipv4},
		},
		"4.example.com.": {
			A:    []string{ipv4},
			AAAA: []string{ipv6},
		},
	}, false)

	suite.srv.PatchNet(net.DefaultResolver)
}

// this function executes after all tests executed
func (suite *ResolverTestSuite) TearDownSuite() {
	suite.srv.Close()
	mockdns.UnpatchNet(net.DefaultResolver)
}

func (suite *ResolverTestSuite) TestNameResolveIpv4() {
	resolver := &DNSResolver{}

	config.Cfg.DnsHost = ""
	config.Cfg.PreferIpv6 = false

	_, ip, err := resolver.Resolve(context.Background(), "1.example.com")
	suite.NoError(err)

	suite.Equal(ip.String(), ipv4)
}

func (suite *ResolverTestSuite) TestNameResolveIpv4Miss() {
	resolver := &DNSResolver{}

	config.Cfg.DnsHost = ""
	config.Cfg.PreferIpv6 = false

	_, ip, err := resolver.Resolve(context.Background(), "nan.example.com")
	suite.Error(err)

	suite.Nil(ip)
}

func (suite *ResolverTestSuite) TestNameResolveIpv6() {
	resolver := &DNSResolver{}

	config.Cfg.DnsHost = ""
	config.Cfg.PreferIpv6 = true

	_, ip, err := resolver.Resolve(context.Background(), "1.example.com")
	suite.NoError(err)

	suite.Equal(ip.String(), ipv6)
}

func (suite *ResolverTestSuite) TestNameResolveIpv6Short() {
	resolver := &DNSResolver{}

	config.Cfg.DnsHost = ""
	config.Cfg.PreferIpv6 = true

	_, ip, err := resolver.Resolve(context.Background(), "2.example.com")
	suite.NoError(err)

	suite.Equal(ip.String(), ipv62Short)
}

func (suite *ResolverTestSuite) TestNameResolveIpv6Miss() {
	resolver := &DNSResolver{}

	config.Cfg.DnsHost = ""
	config.Cfg.PreferIpv6 = true

	_, ip, err := resolver.Resolve(context.Background(), "3.example.com")
	suite.NoError(err)

	suite.Equal(ip.String(), ipv4)
}

func (suite *ResolverTestSuite) TestNameResolveIpv4Custom() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	cacheDB := cache.New(1*time.Minute, 3*time.Minute)

	resolver := &DNSResolver{
		Cache:      cacheDB,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = false
	config.Cfg.DnsHost = "test"
	config.Cfg.DnsUseCache = true

	_, ip, err := resolver.Resolve(context.Background(), "4.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv4)

	suite.Equal(cacheDB.ItemCount(), 1)

	val, found := cacheDB.Get("4.example.com")
	suite.Equal(found, true)
	suite.Equal(val.([]net.IP)[0].String(), ipv4)

	_, ip, err = resolver.Resolve(context.Background(), "4.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv4)

	suite.Equal(cacheDB.ItemCount(), 1)

	val, found = cacheDB.Get("4.example.com")
	suite.Equal(found, true)
	suite.Equal(val.([]net.IP)[0].String(), ipv4)
}

func (suite *ResolverTestSuite) TestNameResolveIpv4CustomNoCache() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	cacheDB := cache.New(1*time.Minute, 3*time.Minute)

	resolver := &DNSResolver{
		Cache:      cacheDB,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = false
	config.Cfg.DnsHost = "test"
	config.Cfg.DnsUseCache = false

	_, ip, err := resolver.Resolve(context.Background(), "4.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv4)

	suite.Equal(cacheDB.ItemCount(), 0)

	_, found := cacheDB.Get("4.example.com")
	suite.Equal(found, false)
	suite.Equal(cacheDB.ItemCount(), 0)
}

func (suite *ResolverTestSuite) TestNameResolveIpv4CustomCacheNil() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	resolver := &DNSResolver{
		Cache:      nil,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = false
	config.Cfg.DnsHost = "test"
	config.Cfg.DnsUseCache = false

	_, ip, err := resolver.Resolve(context.Background(), "4.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv4)
}

func (suite *ResolverTestSuite) TestNameResolveIpv6Custom() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	cacheDB := cache.New(1*time.Minute, 3*time.Minute)

	resolver := &DNSResolver{
		Cache:      cacheDB,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = true
	config.Cfg.DnsHost = "test"
	config.Cfg.DnsUseCache = true

	_, ip, err := resolver.Resolve(context.Background(), "4.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv6)

	suite.Equal(cacheDB.ItemCount(), 1)

	val, found := cacheDB.Get("4.example.com")
	suite.Equal(found, true)
	suite.Equal(val.([]net.IP)[0].String(), ipv6)

	_, ip, err = resolver.Resolve(context.Background(), "4.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv6)

	suite.Equal(cacheDB.ItemCount(), 1)

	val, found = cacheDB.Get("4.example.com")
	suite.Equal(found, true)
	suite.Equal(val.([]net.IP)[0].String(), ipv6)
}

func (suite *ResolverTestSuite) TestNameResolveIpv6CustomТщСфсру() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	cacheDB := cache.New(1*time.Minute, 3*time.Minute)

	resolver := &DNSResolver{
		Cache:      cacheDB,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = true
	config.Cfg.DnsHost = "test"
	config.Cfg.DnsUseCache = false

	_, ip, err := resolver.Resolve(context.Background(), "4.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv6)

	suite.Equal(cacheDB.ItemCount(), 0)

	_, found := cacheDB.Get("4.example.com")
	suite.Equal(found, false)
}

func (suite *ResolverTestSuite) TestNameResolveIpv6CustomMiss() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	cacheDB := cache.New(1*time.Minute, 3*time.Minute)

	resolver := &DNSResolver{
		Cache:      cacheDB,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = true
	config.Cfg.DnsHost = "test"
	config.Cfg.DnsUseCache = true

	_, ip, err := resolver.Resolve(context.Background(), "3.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv4)

	suite.Equal(cacheDB.ItemCount(), 1)

	val, found := cacheDB.Get("3.example.com")
	suite.Equal(found, true)
	suite.Equal(val.([]net.IP)[0].String(), ipv4)
}

func (suite *ResolverTestSuite) TestNameResolveIpv6CustomMissNoCache() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	cacheDB := cache.New(1*time.Minute, 3*time.Minute)

	resolver := &DNSResolver{
		Cache:      cacheDB,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = true
	config.Cfg.DnsHost = "test"
	config.Cfg.DnsUseCache = false

	_, ip, err := resolver.Resolve(context.Background(), "3.example.com")
	suite.NoError(err)
	suite.Equal(ip.String(), ipv4)

	suite.Equal(cacheDB.ItemCount(), 0)

	_, found := cacheDB.Get("3.example.com")
	suite.Equal(found, false)
}

func (suite *ResolverTestSuite) TestNameResolveIpv6CustomNonExist() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	cacheDB := cache.New(1*time.Minute, 3*time.Minute)

	resolver := &DNSResolver{
		Cache:      cacheDB,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = true
	config.Cfg.DnsHost = "test"

	_, ips, err := resolver.Resolve(context.Background(), "nan.example.com")
	suite.NoError(err)
	suite.Nil(ips)

	suite.Equal(cacheDB.ItemCount(), 0)

	_, found := cacheDB.Get("nan.example.com")
	suite.Equal(found, false)
}

func (suite *ResolverTestSuite) TestNameResolveIpv4CustomNonExist() {
	addr, err := netip.ParseAddrPort(suite.srv.LocalAddr().String())
	suite.NoError(err)

	cacheDB := cache.New(1*time.Minute, 3*time.Minute)

	resolver := &DNSResolver{
		Cache:      cacheDB,
		DNSClient:  new(dns.Client),
		DNSAddress: net.JoinHostPort(addr.Addr().String(), fmt.Sprintf("%d", addr.Port())),
	}

	config.Cfg.PreferIpv6 = false
	config.Cfg.DnsHost = "test"

	_, ips, err := resolver.Resolve(context.Background(), "nan.example.com")
	suite.NoError(err)
	suite.Nil(ips)

	suite.Equal(cacheDB.ItemCount(), 0)

	_, found := cacheDB.Get("nan.example.com")
	suite.Equal(found, false)
}
