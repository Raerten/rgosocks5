package stat

import (
	"net"
)

type Conn struct {
	net.Conn
	CloseCb func()
}

func (c Conn) Close() error {
	c.CloseCb()
	return c.Conn.Close()
}
