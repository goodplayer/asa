package tlssupport

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/goodplayer/asa/core/upstream"
)

func sendEcdsaSign(digest []byte, ups *upstream.Upstream, data *clientCipherSuiteData) ([]byte, error) {
	//TODO support keyless server authentication
	conn, err := ups.SelectTcp(-1, false)

	digestStr := base64.StdEncoding.EncodeToString(digest)
	serverName := data.ServerName
	family := fmt.Sprint(data.Family)
	code := fmt.Sprint(data.Code)

	form := url.Values{}
	form["digest"] = []string{digestStr}
	form["server_name"] = []string{serverName}
	form["family"] = []string{family}
	form["code"] = []string{code}

	resp, err := http.PostForm("http://"+conn.Addr+"/ecdsa/sign", form)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		resultData, err := ioutil.ReadAll(resp.Body)
		if len(resultData) == 0 && err != nil {
			return resultData, err
		}
		return resultData, nil
	} else {
		return nil, errors.New("ecdsa sign no body replied.")
	}
}

func sendRsaSign(digest []byte, ups *upstream.Upstream, data *clientCipherSuiteData, hashFunc uint) ([]byte, error) {
	//TODO support keyless server authentication
	conn, err := ups.SelectTcp(-1, false)

	digestStr := base64.StdEncoding.EncodeToString(digest)
	serverName := data.ServerName
	family := fmt.Sprint(data.Family)
	code := fmt.Sprint(data.Code)
	hash := fromHashFuncMap[hashFunc]

	form := url.Values{}
	form["digest"] = []string{digestStr}
	form["server_name"] = []string{serverName}
	form["family"] = []string{family}
	form["code"] = []string{code}
	form["hash"] = []string{hash}

	resp, err := http.PostForm("http://"+conn.Addr+"/rsa/sign", form)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		resultData, err := ioutil.ReadAll(resp.Body)
		if len(resultData) == 0 && err != nil {
			return resultData, err
		}
		return resultData, nil
	} else {
		return nil, errors.New("ecdsa sign no body replied.")
	}
}

func sendRsaDecrypt(msg []byte, ups *upstream.Upstream, data *clientCipherSuiteData, sessionKeyLength int) ([]byte, error) {
	//TODO support keyless server authentication
	conn, err := ups.SelectTcp(-1, false)

	msgstr := base64.StdEncoding.EncodeToString(msg)
	serverName := data.ServerName
	family := fmt.Sprint(data.Family)
	code := fmt.Sprint(data.Code)
	sessionKeyLen := fmt.Sprint(sessionKeyLength)

	form := url.Values{}
	form["msg"] = []string{msgstr}
	form["server_name"] = []string{serverName}
	form["family"] = []string{family}
	form["code"] = []string{code}
	form["session_key_length"] = []string{sessionKeyLen}

	resp, err := http.PostForm("http://"+conn.Addr+"/rsa/decrypt", form)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		resultData, err := ioutil.ReadAll(resp.Body)
		if len(resultData) == 0 && err != nil {
			return resultData, err
		}
		return resultData, nil
	} else {
		return nil, errors.New("ecdsa sign no body replied.")
	}
}
