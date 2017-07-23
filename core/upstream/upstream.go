package upstream

import (
	"errors"
	"net"
	"net/url"
	"time"
)

type Upstream struct {
	Key      string
	List     []string
	Strategy string
}

type UpstreamTcpConn struct {
	id int
	*net.TCPConn
	Addr string
}

type UpstreamUdpConn struct {
	id int
	*net.UDPConn
	Addr string
}

func (this *Upstream) SelectTcp(key int, dial bool) (UpstreamTcpConn, error) {
	addrStr := this.List[key%len(this.List)]
	urla, err := url.Parse(addrStr)
	if err != nil {
		return UpstreamTcpConn{}, err
	}

	if !dial {
		return UpstreamTcpConn{
			Addr: urla.Host,
		}, nil
	}

	//TODO dial timeout
	conn, err := net.DialTimeout(urla.Scheme, urla.Host, 2*time.Second)
	if err != nil {
		return UpstreamTcpConn{}, err
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return UpstreamTcpConn{}, errors.New("not tcp conn when dial.")
	}
	upstreamConn := UpstreamTcpConn{
		TCPConn: tcpConn,
		Addr:    urla.Host,
	}
	return upstreamConn, nil
}

func (this *Upstream) RecycleTcp(tcpConn UpstreamTcpConn) {
	if tcpConn.TCPConn != nil {
		tcpConn.Close()
	}
}

func (this *Upstream) SelectUdp(key int, dial bool) (UpstreamUdpConn, error) {
	//TODO
	return UpstreamUdpConn{}, errors.New("upstream udp conn unsupported.")
}

func (this *Upstream) RecycleUdp(udpConn UpstreamUdpConn) {
	//TODO
}
