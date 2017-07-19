package trans

import "net"

type TcpTransConfig struct {
	UseReuseport       bool
	MaxListenerBacklog int
}

type TcpTrans interface {
	NewTcpListener(proto, addr string) (*net.TCPListener, error)
	NewListener(proto, addr string, config TcpTransConfig) (net.Listener, error)
	PrepareAcceptedConn(conn net.Conn, config TcpTransConfig) error
}
