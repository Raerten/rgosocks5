package rules

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
	"net"
	"rgosocks/config"
	"testing"
)

type setupRule struct {
	reqIp       string
	reqFQDN     string
	allowedNet  []string
	rejectNet   []string
	allowedFQDN []string
	rejectFQDN  []string
	command     byte
}

func getConnectRules(setup *setupRule) (*ProxyRulesSet, *socks5.Request) {
	if setup.command == byte(0) {
		setup.command = statute.CommandConnect
	}
	req := &socks5.Request{
		Request: *&statute.Request{
			Command: setup.command,
		},
		DestAddr: &statute.AddrSpec{
			FQDN: setup.reqFQDN,
			IP:   net.ParseIP(setup.reqIp),
		},
	}
	var allowedNet []*net.IPNet
	for _, cidr := range setup.allowedNet {
		_, ipNet, _ := net.ParseCIDR(cidr)
		allowedNet = append(allowedNet, ipNet)
	}

	var rejectNet []*net.IPNet
	for _, cidr := range setup.rejectNet {
		_, ipNet, _ := net.ParseCIDR(cidr)
		rejectNet = append(rejectNet, ipNet)
	}

	rules := &ProxyRulesSet{
		AllowedIPNet: allowedNet,
		RejectIPNet:  rejectNet,
		AllowedFQDN:  setup.allowedFQDN,
		RejectFQDN:   setup.rejectFQDN,
	}

	return rules, req
}

func TestIPRejectSingle(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqIp:     "192.168.1.1",
		rejectNet: []string{"192.168.1.0/24"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}

func TestIPRejectMulti(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqIp:     "192.168.1.1",
		rejectNet: []string{"192.168.2.0/24", "192.168.1.0/24"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}

func TestIPRejectSingleMiss(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqIp:     "192.168.3.1",
		rejectNet: []string{"192.168.1.0/24"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.True(t, result)
}

func TestIPRejectSingleMissMulti(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqIp:     "192.168.3.1",
		rejectNet: []string{"192.168.2.0/24", "192.168.1.0/24"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.True(t, result)
}

func TestIPRejectWithAllow(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqIp:      "192.168.1.1",
		allowedNet: []string{"192.168.1.1/32"},
		rejectNet:  []string{"192.168.0.0/16", "172.16.1.0/24"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}

func TestIPAllowWithRejectMulti(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqIp:      "192.168.1.1",
		allowedNet: []string{"192.168.1.0/24"},
		rejectNet:  []string{"192.168.2.0/24", "172.16.3.0/24"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.True(t, result)
}

func TestIPAllowMiss(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqIp:      "192.168.1.1",
		allowedNet: []string{"192.168.2.0/24"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}

func TestIPAllow(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqIp:      "192.168.1.1",
		allowedNet: []string{"192.168.1.0/24"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.True(t, result)
}

func TestFQDNAllow(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqFQDN:     "example.com",
		allowedFQDN: []string{"example2.com"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}

func TestFQDNAllowWithRejectSameFQDN(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqFQDN:     "example.com",
		allowedFQDN: []string{"example.com"},
		rejectFQDN:  []string{"example.com"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}

func TestFQDNRejectMiss(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqFQDN:    "example.com",
		rejectFQDN: []string{"example2.com"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.True(t, result)
}

func TestFQDNReject(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		reqFQDN:     "example.com",
		allowedFQDN: []string{"example2.com"},
		rejectFQDN:  []string{"example.com"},
	})

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}

func TestDisableBind(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		command: statute.CommandBind,
	})

	config.Cfg.DisableBind = true

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}

func TestDisableAssociate(t *testing.T) {
	rules, req := getConnectRules(&setupRule{
		command: statute.CommandAssociate,
	})

	config.Cfg.DisableAssociate = true

	_, result := rules.Allow(context.Background(), req)

	assert.False(t, result)
}
