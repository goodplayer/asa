// +build darwin

package trans

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type tcpTransImpl struct {
}

func supportReuseport() bool {
	return true
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
	if !config.UseReuseport {
		return this.NewTcpListener(proto, addr)
	}

	tcpAddr, err := net.ResolveTCPAddr(proto, addr)
	if err != nil {
		return nil, err
	}

	var family = syscall.AF_INET
	var ipv6only = false
	var sockaddr syscall.Sockaddr
	switch proto {
	case "tcp":
		family = syscall.AF_INET6
		ipv6only = false
		soaddr := &syscall.SockaddrInet6{Port: tcpAddr.Port}
		if tcpAddr.IP != nil {
			copy(soaddr.Addr[:], tcpAddr.IP)
		}
		sockaddr = soaddr
	case "tcp4":
		family = syscall.AF_INET
		ipv6only = false
		soaddr := &syscall.SockaddrInet4{Port: tcpAddr.Port}
		if tcpAddr.IP != nil {
			copy(soaddr.Addr[:], tcpAddr.IP[12:16])
		}
		sockaddr = soaddr
	case "tcp6":
		family = syscall.AF_INET6
		ipv6only = true
		soaddr := &syscall.SockaddrInet6{Port: tcpAddr.Port}
		if tcpAddr.IP != nil {
			copy(soaddr.Addr[:], tcpAddr.IP)
		}
		if tcpAddr.Zone != "" {
			iface, err := net.InterfaceByName(tcpAddr.Zone)
			if err != nil {
				return nil, err
			}
			soaddr.ZoneId = uint32(iface.Index)
		}
		sockaddr = soaddr
	default:
		return nil, errors.New("unsupported proto type - asa-trans")
	}

	syscall.ForkLock.RLock()
	s, err := syscall.Socket(family, syscall.SOCK_STREAM, 0)
	if err != nil {
		syscall.ForkLock.RUnlock()
		return nil, err
	}
	syscall.CloseOnExec(s)
	syscall.ForkLock.RUnlock()
	if err = syscall.SetNonblock(s, true); err != nil {
		syscall.Close(s)
		return nil, err
	}

	if family == syscall.AF_INET6 && ipv6only {
		err = syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 1)
		if err != nil {
			syscall.Close(s)
			return nil, err
		}
	}
	err = syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	if err != nil {
		syscall.Close(s)
		return nil, err
	}

	err = syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		syscall.Close(s)
		return nil, err
	}

	var reuseport = 0x200
	err = syscall.SetsockoptInt(s, syscall.SOL_SOCKET, reuseport, 1)
	if err != nil {
		syscall.Close(s)
		return nil, err
	}

	if err = syscall.Bind(s, sockaddr); err != nil {
		syscall.Close(s)
		return nil, err
	}

	var maxBacklogCnt = maxListenerBacklogCnt
	if config.MaxListenerBacklog > 0 && config.MaxListenerBacklog < 1<<16 {
		maxBacklogCnt = config.MaxListenerBacklog
	}

	if err = syscall.Listen(s, maxBacklogCnt); err != nil {
		syscall.Close(s)
		return nil, err
	}

	defer syscall.Close(s)

	return net.FileListener(os.NewFile(uintptr(s), fmt.Sprintf("asa-trans-reuseport.%d.%s.%s", os.Getpid(), proto, addr)))
}

func (this tcpTransImpl) PrepareAcceptedConn(conn net.Conn, config TcpTransConfig) error {
	return nil
}

var maxListenerBacklogCnt = maxListenerBacklog()

func maxListenerBacklog() int {
	var (
		n   uint32
		err error
	)
	n, err = syscall.SysctlUint32("kern.ipc.somaxconn")
	if n == 0 || err != nil {
		return syscall.SOMAXCONN
	}
	if n > 1<<16-1 {
		n = 1<<16 - 1
	}

	return int(n)
}
