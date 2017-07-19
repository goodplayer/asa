// +build !linux

package trans

import (
	"errors"
	"net"
)

type tcpTransImpl struct {
}

func supportReuseport() bool {
	return false
}

func newTcpTrans() TcpTrans {
	return tcpTransImpl{}
}

func (this tcpTransImpl) NewTcpListener(proto, addr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr(proto, addr)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP(proto, tcpAddr)
}

func (this tcpTransImpl) NewListener(proto, addr string, config TcpTransConfig) (net.Listener, error) {
	if config.UseReuseport {
		return nil, errors.New("tcp not support reuseport - asa-trans")
	}
	return this.NewTcpListener(proto, addr)
}

func (this tcpTransImpl) PrepareAcceptedConn(conn net.Conn, config TcpTransConfig) error {
	return nil
}
