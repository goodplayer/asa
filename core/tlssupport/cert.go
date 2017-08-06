package tlssupport

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"strings"
)

type CertManager struct {
	Certs           map[string]*tls.Certificate
	ParentTlsConfig *tls.Config
}

func (this *CertManager) GetCertificate(c *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cert, ok := this.getCert(c.ServerName)
	if !ok {
		return nil, errors.New("cannot find certificate for " + c.ServerName)
	}
	suite, err := pickOneCipherSuite(cert, c.CipherSuites)
	if err != nil {
		return nil, err
	}
	cert = repackCert(cert, suite, c.ServerName)
	if cert == nil {
		return nil, errors.New("cert is not available after repack cert")
	}
	return cert, nil
}

func (this *CertManager) GetConfigForClient(c *tls.ClientHelloInfo) (*tls.Config, error) {
	// set cipher to client and to keyless server
	cert, ok := this.getCert(c.ServerName)
	if !ok {
		return nil, errors.New("cannot find config for " + c.ServerName)
	}
	suite, err := pickOneCipherSuite(cert, c.CipherSuites)
	if err != nil {
		return nil, err
	}
	// config list:
	// 1. prefer server suites
	// 2. min tls version: tls1.0
	config := this.ParentTlsConfig.Clone()
	// 1. get cert alg  2. find first match to client  3. only set this
	config.CipherSuites = []uint16{suite.Code}
	return config, nil
}

func (this *CertManager) FillCerts(certs map[string]*tls.Certificate) {
	this.Certs = certs
}

func (this *CertManager) getCert(host string) (*tls.Certificate, bool) {
	cert, ok := this.Certs[host]
	return cert, ok
}

func repackCert(cert *tls.Certificate, cipher *CipherSuite, serverName string) *tls.Certificate {
	switch c := cert.PrivateKey.(type) {
	case *wrappedClientPrivateKey:
		// clone cert
		var certClone = *cert
		var certResult = &certClone
		var prikeyClone = *c

		family := cert.Leaf.PublicKeyAlgorithm
		// add cipher info
		prikeyClone.CipherSuiteData = &clientCipherSuiteData{
			Code:       cipher.Code,
			ServerName: serverName,
			Family:     int(family),
		}
		switch family {
		case x509.RSA:
			certResult.PrivateKey = &prikeyClone
		case x509.ECDSA:
			certResult.PrivateKey = prikeyClone.toEcdsaKey()
		default:
			return nil
		}

		return certResult
	default:
		// local cert, skip
		// clone cert
		return cert
	}
}

func pickOneCipherSuite(cert *tls.Certificate, clientCipherSuite []uint16) (*CipherSuite, error) {
	switch cert.Leaf.PublicKeyAlgorithm {
	case x509.RSA:
		// rsa alg
		suite := matchBestRsaCipherSuite(clientCipherSuite)
		if suite != nil {
			return suite, nil
		} else {
			return nil, errors.New("cannot find cipher suite matching client")
		}
	case x509.ECDSA:
		// ecdsa alg
		suite := matchBestEcdsaCipherSuite(clientCipherSuite)
		if suite != nil {
			return suite, nil
		} else {
			return nil, errors.New("cannot find cipher suite matching client")
		}
	default:
		return nil, errors.New("public key algorithm not support")
	}
}

func LoadCert(path string) (tls.Certificate, error) {
	certPEMBlock, err := ioutil.ReadFile(path)
	if err != nil {
		return tls.Certificate{}, err
	}

	fail := func(err error) (tls.Certificate, error) {
		return tls.Certificate{}, err
	}

	var cert tls.Certificate
	var skippedBlockTypes []string
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		} else {
			skippedBlockTypes = append(skippedBlockTypes, certDERBlock.Type)
		}
	}

	if len(cert.Certificate) == 0 {
		if len(skippedBlockTypes) == 0 {
			return fail(errors.New("fail to find any PEM data."))
		}
		if len(skippedBlockTypes) == 1 && strings.HasSuffix(skippedBlockTypes[0], "PRIVATE KEY") {
			return fail(errors.New("only find private key block"))
		}
		return fail(errors.New("fail to find certificate block"))
	}

	// prepare leaf cert
	x509cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fail(err)
	}
	cert.Leaf = x509cert

	return cert, nil
}
