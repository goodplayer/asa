package httpproto

import (
	"crypto/tls"
	"errors"

	"github.com/goodplayer/asa/core/proto/httpproto/vhost"
	"github.com/goodplayer/asa/core/tlssupport"
)

func initCertManager(certManager *tlssupport.CertManager, vhosts map[string]*vhost.Vhost, tlsconfig *tls.Config) error {
	// init cert manager
	certManager.ParentTlsConfig = tlsconfig
	result := make(map[string]*tls.Certificate)
	for k, v := range vhosts {
		if !v.SslEnable {
			continue
		}
		if v.SslCertPath == "" {
			return errors.New("cert is empty for host: " + k)
		}
		cert, err := tlssupport.LoadCert(v.SslCertPath)
		if err != nil {
			return err
		}
		if v.SslPrivateKeyPath != "" {
			prikey, err := tlssupport.LoadPrivateKey(v.SslPrivateKeyPath)
			if err != nil {
				return err
			}
			finalcert, _ := tlssupport.MakeLocalCertAndPriKey(&cert, prikey)
			result[k] = finalcert
			continue
		}
		if v.SslKeylessServer != nil {
			finalcert, _ := tlssupport.MakeKeylessClientCertAndPriKey(&cert, v.SslKeylessServer)
			result[k] = finalcert
			continue
		}
		return errors.New("not private method set for host: " + k)
	}
	certManager.FillCerts(result)
	return nil
}
