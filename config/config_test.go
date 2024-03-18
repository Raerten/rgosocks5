package config

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	os.Setenv("PROXY_HOST", "1.1.1.1")
	os.Setenv("PROXY_PORT", "8080")
	os.Setenv("LOG_LEVEL_DEBUG", "true")
	defer func() {
		os.Unsetenv("PROXY_HOST")
		os.Unsetenv("PROXY_PORT")
		os.Unsetenv("LOG_LEVEL_DEBUG")
	}()

	Parse()

	if Cfg.ProxyHost != "1.1.1.1" {
		t.Errorf("Parse() = %q, want %q", Cfg.ProxyHost, "1.1.1.1")
	}

	if Cfg.ProxyPort != 8080 {
		t.Errorf("Parse() = %q, want %q", Cfg.ProxyPort, 8080)
	}

	if Cfg.ProxyAddress != "1.1.1.1:8080" {
		t.Errorf("Parse() = %q, want %q", Cfg.ProxyHost, "1.1.1.1:8080")
	}

	if Cfg.LogLevelDebug != true {
		t.Errorf("Parse() = %v, want %v", Cfg.LogLevelDebug, true)
	}
}
