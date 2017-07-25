package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
)

type DefaultHttpHandler struct {
}

func (this *DefaultHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintln("hello world!")))
}

type DefaultFcgiHttpHandler struct {
}

func (this *DefaultFcgiHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintln("hello world of fcgi!")))
}

func main() {
	ch := make(chan bool, 2)

	go func() {
		err := http.ListenAndServe("0.0.0.0:8888", new(DefaultHttpHandler))
		if err != nil {
			panic(err)
		}
		ch <- true
	}()

	go func() {
		addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8887")
		if err != nil {
			panic(err)
		}
		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			panic(err)
		}
		err = fcgi.Serve(addr, new(DefaultFcgiHttpHandler))
		if err != nil {
			panic(err)
		}
	}()

	<-ch
	<-ch
}
