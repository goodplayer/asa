package fcgiclient

import (
	"net/http"

	"github.com/goodplayer/asa/util"
)

func HandltFcgi(w http.ResponseWriter, r *http.Request, listenAddr, listenPort, serverName string) {

	//TODO document root
	docRoot := "/var/www"
	//TODO request file, uri / uri + index file
	docFile := "/helloworld/index.php"
	//TODO uri: final uri (internal redirects / index file / normalized - no relatives, %xx, query string / etc.)
	uri := "/hellworld/index.php"

	env := map[string]string{}
	env["GATEWAY_INTERFACE"] = "CGI/1.1"
	env["SERVER_SOFTWARE"] = "asa-fastcgi"
	env["QUERY_STRING"] = r.URL.RawPath
	env["REQUEST_METHOD"] = r.Method
	env["CONTENT_TYPE"] = r.Header.Get("Content-Type")
	env["CONTENT_LENGTH"] = r.ContentLength
	hostAndPort := util.SplitHostAndPort(r.RemoteAddr)
	env["REMOTE_ADDR"] = hostAndPort[0]
	env["REMOTE_PORT"] = hostAndPort[1]
	env["SERVER_ADDR"] = listenAddr
	env["SERVER_PORT"] = listenPort
	env["SERVER_NAME"] = serverName
	env["SERVER_PROTOCOL"] = r.Proto
	env["REQUEST_URI"] = r.RequestURI
	env["DOCUMENT_URI"] = uri
	env["DOCUMENT_ROOT"] = docRoot
	env["SCRIPT_NAME"] = docFile
	env["SCRIPT_FILENAME"] = docRoot + docFile

}
