package tlssupport

import (
	"crypto"
	"crypto/tls"

	"github.com/goodplayer/asa/core/upstream"
)

func MakeLocalCertAndPriKey(cert *tls.Certificate, prikey crypto.PrivateKey) (*tls.Certificate, crypto.PrivateKey) {
	cert.PrivateKey = prikey
	return cert, prikey
}

func MakeKeylessClientCertAndPriKey(cert *tls.Certificate, keylessServer *upstream.Upstream) (*tls.Certificate, crypto.PrivateKey) {
	pri := &wrappedClientPrivateKey{
		PublicKey:     cert.Leaf.PublicKey,
		KeylessServer: keylessServer,
	}
	cert.PrivateKey = pri
	return cert, pri
}

func MakeKeylessServerCertAndPriKey(cert *tls.Certificate, prikey crypto.PrivateKey) (*tls.Certificate, crypto.PrivateKey) {
	return MakeLocalCertAndPriKey(cert, prikey)
}
