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
	fmt.Println("====> new request")
	fmt.Println(r.URL.Path)
	fmt.Println(r.URL.Query())
	fmt.Println(r.Host)
	fmt.Println(r.Header)

	if r.URL.Path == "/404" {
		http.NotFound(w, r)
		fmt.Println("result: 404")
	} else if r.URL.Path == "/302" {
		http.Redirect(w, r, "http://localhost/http", http.StatusFound)
		fmt.Println("result: 302")
	} else {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintln("hello world!")))
		fmt.Println("result: 200")
	}

	fmt.Println("----> end request")
}

type DefaultFcgiHttpHandler struct {
}

func (this *DefaultFcgiHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====> new request")
	fmt.Println(r.URL.Path)
	fmt.Println(r.URL.Query())
	fmt.Println(r.Host)
	fmt.Println(r.Header)

	if r.URL.Path == "/404" {
		http.NotFound(w, r)
		fmt.Println("result: 404")
	} else if r.URL.Path == "/302" {
		http.Redirect(w, r, "http://localhost/http", http.StatusFound)
		fmt.Println("result: 302")
	} else {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintln("hello world of fcgi!")))
		fmt.Println("result: 200")
	}

	fmt.Println("----> end request")
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
		err = fcgi.Serve(l, new(DefaultFcgiHttpHandler))
		if err != nil {
			panic(err)
		}
	}()

	<-ch
	<-ch
}
