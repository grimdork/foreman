package main

import (
	"net/http"
	"time"

	ll "github.com/grimdork/loglines"
)

func (srv *Server) addLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		go func() {
			ll.Msg("client %s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (srv *Server) auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		k, err := srv.GetKey(id)
		if err != nil || k.Key == "" || k.Expiry.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		key := r.Header.Get("key")
		if key == "" || key != k.Key {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (srv *Server) authAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		k, err := srv.GetKey(id)
		if err != nil || k.Key == "" || k.Expiry.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		key := r.Header.Get("key")
		if key == "" || key != k.Key || !k.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// AddJSONHeaders for JSON responses.
func AddJSONHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// AddHTMLHeaders for web pages.
func AddHTMLHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// AddCORS to allow REST access from other domains.
func AddCORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
