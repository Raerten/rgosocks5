package config

import (
	"testing"
)

func TestParse(t *testing.T) {
	t.Setenv("PROXY_HOST", "1.1.1.1")
	t.Setenv("PROXY_PORT", "8080")
	t.Setenv("LOG_LEVEL_DEBUG", "true")

	Parse()

	if Cfg.ProxyHost != "1.1.1.1" {
		t.Errorf("Parse() = %q, want %q", Cfg.ProxyHost, "1.1.1.1")
	}

	if Cfg.ProxyPort != 8080 {
		t.Errorf("Parse() = %d, want %d", Cfg.ProxyPort, 8080)
	}

	if Cfg.ProxyAddress != "1.1.1.1:8080" {
		t.Errorf("Parse() = %q, want %q", Cfg.ProxyHost, "1.1.1.1:8080")
	}

	if Cfg.LogLevelDebug != true {
		t.Errorf("Parse() = %v, want %v", Cfg.LogLevelDebug, true)
	}
}
