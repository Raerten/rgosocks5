package slogger

import (
	"fmt"
	"log/slog"
)

type Socks5Logger struct {
}

func (sf Socks5Logger) Errorf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...))
}
