package main

import (
	"fmt"
	"net/http"
)

type DefaultHttpHandler struct {
}

func (this *DefaultHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintln("hello world!")))
}

func main() {
	err := http.ListenAndServe("0.0.0.0:8888", new(DefaultHttpHandler))
	if err != nil {
		panic(err)
	}
}
