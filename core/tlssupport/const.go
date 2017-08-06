package tlssupport

import (
	"crypto"
	"crypto/tls"
)

type CipherSuite struct {
	Code uint16
}

var ecdsaCipherClass1Map = map[uint16]*CipherSuite{}
var ecdsaCipherClass2Map = map[uint16]*CipherSuite{}
var ecdsaCipherClass3Map = map[uint16]*CipherSuite{}
var rsaCipherClass1Map = map[uint16]*CipherSuite{}
var rsaCipherClass2Map = map[uint16]*CipherSuite{}
var rsaCipherClass3Map = map[uint16]*CipherSuite{}

var ecdsaCipherClass1 = []*CipherSuite{
	{tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305},
	{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
	{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
}

var ecdsaCipherClass2 = []*CipherSuite{
	{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256},
	{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
	{tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA},
}

var ecdsaCipherClass3 = []*CipherSuite{}

var rsaCipherClass1 = []*CipherSuite{
	{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305},
	{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
	{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
}

var rsaCipherClass2 = []*CipherSuite{
	{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256},
	{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
	{tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
}

var rsaCipherClass3 = []*CipherSuite{
	{tls.TLS_RSA_WITH_AES_128_GCM_SHA256},
	{tls.TLS_RSA_WITH_AES_256_GCM_SHA384},
	{tls.TLS_RSA_WITH_AES_128_CBC_SHA256},
	{tls.TLS_RSA_WITH_AES_128_CBC_SHA},
	{tls.TLS_RSA_WITH_AES_256_CBC_SHA},
	{tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA},
	{tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA},
}

var (
	fromHashFuncMap = map[uint]string{}
	toHashFuncMap   = map[string]uint{}
)

func init() {
	genCipherClassMap(ecdsaCipherClass1, ecdsaCipherClass1Map)
	genCipherClassMap(ecdsaCipherClass2, ecdsaCipherClass2Map)
	genCipherClassMap(ecdsaCipherClass3, ecdsaCipherClass3Map)
	genCipherClassMap(rsaCipherClass1, rsaCipherClass1Map)
	genCipherClassMap(rsaCipherClass2, rsaCipherClass2Map)
	genCipherClassMap(rsaCipherClass3, rsaCipherClass3Map)

	fromHashFuncMap[uint(crypto.MD4)] = "MD4"
	fromHashFuncMap[uint(crypto.MD5)] = "MD5"
	fromHashFuncMap[uint(crypto.SHA1)] = "SHA1"
	fromHashFuncMap[uint(crypto.SHA224)] = "SHA224"
	fromHashFuncMap[uint(crypto.SHA256)] = "SHA256"
	fromHashFuncMap[uint(crypto.SHA384)] = "SHA384"
	fromHashFuncMap[uint(crypto.SHA512)] = "SHA512"
	fromHashFuncMap[uint(crypto.MD5SHA1)] = "MD5SHA1"
	fromHashFuncMap[uint(crypto.RIPEMD160)] = "RIPEMD160"
	fromHashFuncMap[uint(crypto.SHA3_224)] = "SHA3_224"
	fromHashFuncMap[uint(crypto.SHA3_256)] = "SHA3_256"
	fromHashFuncMap[uint(crypto.SHA3_384)] = "SHA3_384"
	fromHashFuncMap[uint(crypto.SHA3_512)] = "SHA3_512"
	fromHashFuncMap[uint(crypto.SHA512_224)] = "SHA512_224"
	fromHashFuncMap[uint(crypto.SHA512_256)] = "SHA512_256"
	for k, v := range fromHashFuncMap {
		toHashFuncMap[v] = k
	}
}

func ToHashFunc(hash string) crypto.Hash {
	v, ok := toHashFuncMap[hash]
	if !ok {
		return crypto.Hash(0)
	}
	return crypto.Hash(v)
}

func genCipherClassMap(ciphers []*CipherSuite, cipherMap map[uint16]*CipherSuite) {
	for _, v := range ciphers {
		cipherMap[v.Code] = v
	}
}

func matchBestEcdsaCipherSuite(c []uint16) *CipherSuite {
	// first class
	for _, v := range c {
		cipher, ok := ecdsaCipherClass1Map[v]
		if ok {
			return cipher
		}
	}
	// second class
	for _, v := range c {
		cipher, ok := ecdsaCipherClass2Map[v]
		if ok {
			return cipher
		}
	}
	// third class
	for _, v := range c {
		cipher, ok := ecdsaCipherClass3Map[v]
		if ok {
			return cipher
		}
	}
	return nil
}

func matchBestRsaCipherSuite(c []uint16) *CipherSuite {
	// first class
	for _, v := range c {
		cipher, ok := rsaCipherClass1Map[v]
		if ok {
			return cipher
		}
	}
	// second class
	for _, v := range c {
		cipher, ok := rsaCipherClass2Map[v]
		if ok {
			return cipher
		}
	}
	// third class
	for _, v := range c {
		cipher, ok := rsaCipherClass3Map[v]
		if ok {
			return cipher
		}
	}
	return nil
}
