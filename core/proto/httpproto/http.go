package httpproto

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/goodplayer/asa/core/proto/httpproto/vhost"
	"github.com/goodplayer/asa/core/tlssupport"
)

func StartHttpServer(l net.Listener, handler http.Handler) error {
	srv := &http.Server{Handler: handler}
	return srv.Serve(keepAliveListener{l.(*net.TCPListener)})
}

func StartHttpTlsServer(l net.Listener, handler http.Handler, vhosts map[string]*vhost.Vhost) error {
	tlsConfig := &tls.Config{}
	// prepare tls config
	tlsConfig.PreferServerCipherSuites = true
	certManager := &tlssupport.CertManager{}
	// init certManager
	err := initCertManager(certManager, vhosts, tlsConfig)
	if err != nil {
		return err
	}
	tlsConfig.GetCertificate = certManager.GetCertificate
	tlsConfig.GetConfigForClient = certManager.GetConfigForClient
	tlsConfig.MinVersion = tls.VersionTLS10

	tlsConfig = tlsConfig.Clone()

	// add h2 and http1.1
	tlsConfig.NextProtos = append(tlsConfig.NextProtos, "h2")
	tlsConfig.NextProtos = append(tlsConfig.NextProtos, "http/1.1")

	srv := &http.Server{Handler: handler}
	tlsListener := tls.NewListener(keepAliveListener{l.(*net.TCPListener)}, tlsConfig)
	return srv.Serve(tlsListener)
}

type keepAliveListener struct {
	*net.TCPListener
}

func (ln keepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
