package main

import (
	"fmt"
	"net/http"
)

func (srv *Server) pulse(w http.ResponseWriter, r *http.Request) {
	hostname := r.Header.Get("hostname")
	key := r.Header.Get("key")
	s := fmt.Sprintf("%s: %s", hostname, key)
	w.Write([]byte(s))
}
