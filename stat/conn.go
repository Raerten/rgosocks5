package stat

import (
	"net"
)

type Conn struct {
	net.Conn
	ReadCb  func(int)
	WriteCb func(int)
	CloseCb func()
}

func (c Conn) Close() error {
	c.CloseCb()
	return c.Conn.Close()
}

func (c Conn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	c.ReadCb(n)
	return n, err
}

func (c Conn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	c.WriteCb(n)
	return n, err
}
