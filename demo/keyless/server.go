package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/goodplayer/asa/core/tlssupport"
)

type DefaultHttpHandler struct {
	EcdsaPrivateKeyMap map[string]crypto.PrivateKey
	RsaPrivateKeyMap   map[string]crypto.PrivateKey
}

func (this *DefaultHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====> request")
	if r.Method == "POST" && r.RequestURI == "/ecdsa/sign" {
		this.ecdsaSignHandler(w, r)
		return
	} else if r.Method == "POST" && r.RequestURI == "/rsa/sign" {
		this.rsaSignHandler(w, r)
		return
	} else if r.Method == "POST" && r.RequestURI == "/rsa/decrypt" {
		this.rsaDecryptHandler(w, r)
		return
	}
	http.Error(w, "Page not found", http.StatusNotFound)
}

func (this *DefaultHttpHandler) ecdsaSignHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("parse form error:", err.Error())
		http.Error(w, "Error occurs", http.StatusInternalServerError)
		return
	}

	digestStr := r.PostForm["digest"][0]
	serverName := r.PostForm["server_name"][0]
	family := r.PostForm["family"][0]
	code := r.PostForm["code"][0]

	//TODO ecdsa family check
	if family == "3" {
		key, ok := this.EcdsaPrivateKeyMap[serverName]
		if !ok {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}

		signer, ok := key.(crypto.Signer)
		if !ok {
			http.Error(w, "signer not support", http.StatusNotFound)
			return
		}

		//TODO code check
		var _ = code

		digest, err := base64.StdEncoding.DecodeString(digestStr)
		if err != nil {
			http.Error(w, "digest is invalid", http.StatusInternalServerError)
			return
		}

		fmt.Println("ecdsa sign:", r.PostForm)

		result, err := signer.Sign(rand.Reader, digest, nil)
		if err != nil {
			http.Error(w, "sign error", http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}
	http.Error(w, "family not match", http.StatusNotFound)
}

func (this *DefaultHttpHandler) rsaSignHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("parse form error:", err.Error())
		http.Error(w, "Error occurs", http.StatusInternalServerError)
		return
	}

	digestStr := r.PostForm["digest"][0]
	serverName := r.PostForm["server_name"][0]
	family := r.PostForm["family"][0]
	code := r.PostForm["code"][0]
	hash := r.PostForm["hash"][0]

	//TODO rsa family check
	if family == "1" {
		key, ok := this.RsaPrivateKeyMap[serverName]
		if !ok {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}

		signer, ok := key.(crypto.Signer)
		if !ok {
			http.Error(w, "signer not support", http.StatusNotFound)
			return
		}

		//TODO code check
		var _ = code

		digest, err := base64.StdEncoding.DecodeString(digestStr)
		if err != nil {
			http.Error(w, "digest is invalid", http.StatusInternalServerError)
			return
		}

		fmt.Println("rsa sign:", r.PostForm)

		hashFunc := tlssupport.ToHashFunc(hash)
		result, err := signer.Sign(rand.Reader, digest, hashFunc)
		if err != nil {
			http.Error(w, "sign error", http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}
	http.Error(w, "family not match", http.StatusNotFound)
}

func (this *DefaultHttpHandler) rsaDecryptHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("parse form error:", err.Error())
		http.Error(w, "Error occurs", http.StatusInternalServerError)
		return
	}

	msg := r.PostForm["msg"][0]
	serverName := r.PostForm["server_name"][0]
	family := r.PostForm["family"][0]
	code := r.PostForm["code"][0]
	sessionKeyLengthStr := r.PostForm["session_key_length"][0]
	sessionKeyLength, err := strconv.Atoi(sessionKeyLengthStr)
	if err != nil {
		http.Error(w, "session key length param error.", http.StatusInternalServerError)
		return
	}

	//TODO rsa family check
	if family == "1" {
		key, ok := this.RsaPrivateKeyMap[serverName]
		if !ok {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}

		signer, ok := key.(crypto.Decrypter)
		if !ok {
			http.Error(w, "signer not support", http.StatusNotFound)
			return
		}

		//TODO code check
		var _ = code

		msgdata, err := base64.StdEncoding.DecodeString(msg)
		if err != nil {
			http.Error(w, "digest is invalid", http.StatusInternalServerError)
			return
		}

		fmt.Println("rsa decrypt:", r.PostForm)

		result, err := signer.Decrypt(rand.Reader, msgdata, &rsa.PKCS1v15DecryptOptions{sessionKeyLength})
		if err != nil {
			http.Error(w, "sign error", http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}
	http.Error(w, "family not match", http.StatusNotFound)
}

func main() {
	ecdsaCert, err := tlssupport.LoadCert("./demo/key/ecdsa_server.crt")
	fail(err)
	rsaCert, err := tlssupport.LoadCert("./demo/key/rsa_server.crt")
	fail(err)

	fmt.Println(ecdsaCert)
	fmt.Println(rsaCert)

	ecdsaPriKey, err := tlssupport.LoadPrivateKey("./demo/key/private_ecdsa_for_crt.key")
	fail(err)
	rsaPriKey, err := tlssupport.LoadPrivateKey("./demo/key/private_rsa.key")
	fail(err)

	fmt.Println(ecdsaPriKey)
	fmt.Println(reflect.TypeOf(ecdsaPriKey))
	fmt.Println(ecdsaPriKey.(*ecdsa.PrivateKey).PublicKey)
	fmt.Println(rsaPriKey)

	h := new(DefaultHttpHandler)
	h.EcdsaPrivateKeyMap = map[string]crypto.PrivateKey{
		"localhost": ecdsaPriKey,
	}
	h.RsaPrivateKeyMap = map[string]crypto.PrivateKey{
		"localhost": rsaPriKey,
	}

	// start keyless server
	err = http.ListenAndServe("0.0.0.0:8999", h)
	if err != nil {
		panic(err)
	}
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}
