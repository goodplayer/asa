package tlssupport

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/goodplayer/asa/core/upstream"
)

func LoadPrivateKey(path string) (crypto.PrivateKey, error) {
	keyPEMBlock, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fail := func(err error) (*tls.Certificate, error) {
		return nil, err
	}
	skippedBlockTypes := []string{}
	var keyDERBlock *pem.Block
	for {
		keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock)
		if keyDERBlock == nil {
			if len(skippedBlockTypes) == 0 {
				return fail(errors.New("fail to find PEM data."))
			}
			if len(skippedBlockTypes) == 1 && skippedBlockTypes[0] == "CERTIFICATE" {
				return fail(errors.New("only find certificate block."))
			}
			return fail(errors.New("fail to find private key block."))
		}
		if keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY") {
			break
		}
		skippedBlockTypes = append(skippedBlockTypes, keyDERBlock.Type)
	}

	if key, err := x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("unknown type in PCKS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(keyDERBlock.Bytes); err == nil {
		return key, nil
	}

	return nil, errors.New("fail to parse private key.")
}

type wrappedClientPrivateKey struct {
	PublicKey       crypto.PublicKey
	KeylessServer   *upstream.Upstream
	CipherSuiteData *clientCipherSuiteData
}

type clientCipherSuiteData struct {
	ServerName string
	Code       uint16
	Family     int
}

var _ crypto.Signer = new(wrappedClientPrivateKey)
var _ crypto.Decrypter = new(wrappedClientPrivateKey)

func (this *wrappedClientPrivateKey) Public() crypto.PublicKey {
	return this.PublicKey
}

func (this *wrappedClientPrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	data := this.CipherSuiteData
	switch x509.PublicKeyAlgorithm(data.Family) {
	case x509.RSA:
		// rsa support
		return sendRsaSign(digest, this.KeylessServer, data, uint(opts.(crypto.Hash)))
	case x509.ECDSA:
		//note: not care opts
		return sendEcdsaSign(digest, this.KeylessServer, data)
	}
	return nil, errors.New("client private key sign error.")
}

func (this *wrappedClientPrivateKey) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
	data := this.CipherSuiteData
	switch x509.PublicKeyAlgorithm(data.Family) {
	case x509.RSA:
		// rsa support
		return sendRsaDecrypt(msg, this.KeylessServer, data, opts.(*rsa.PKCS1v15DecryptOptions).SessionKeyLen)
	}
	return nil, errors.New("client private key decrypt error.")
}

type wrappedClientEcdsaPrivateKey struct {
	wrapped *wrappedClientPrivateKey
}

func (this wrappedClientEcdsaPrivateKey) Public() crypto.PublicKey {
	return this.wrapped.Public()
}

func (this wrappedClientEcdsaPrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return this.wrapped.Sign(rand, digest, opts)
}

func (this *wrappedClientPrivateKey) toEcdsaKey() crypto.PrivateKey {
	return wrappedClientEcdsaPrivateKey{
		wrapped: this,
	}
}

func (this *wrappedClientPrivateKey) toRsaKey() crypto.PrivateKey {
	return this
}

//====for debug====

type wrappedLocalEcdsaPrivateKey struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

func (this wrappedLocalEcdsaPrivateKey) Public() crypto.PublicKey {
	return this.PublicKey
}

func (this wrappedLocalEcdsaPrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	result, err := this.PrivateKey.(crypto.Signer).Sign(rand, digest, opts)
	return result, err
}

type wrappedLocalRsaPrivateKey struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

func (this wrappedLocalRsaPrivateKey) Public() crypto.PublicKey {
	return this.PublicKey
}

func (this wrappedLocalRsaPrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	result, err := this.PrivateKey.(crypto.Signer).Sign(rand, digest, opts)
	return result, err
}

func (this wrappedLocalRsaPrivateKey) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
	result, err := this.PrivateKey.(crypto.Decrypter).Decrypt(rand, msg, opts)
	return result, err
}
