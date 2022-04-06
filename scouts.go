package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// scoutsGet endpoint.
func (srv *Server) scoutsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := srv.GetScouts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	if len(list.Scouts) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}

// scoutGet endpoint.
func (srv *Server) scoutGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	hostname := r.Header.Get("hostname")
	if hostname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing hostname"}`))
		return
	}

	scout, err := srv.GetScout(hostname)
	err = json.NewEncoder(w).Encode(scout)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}

// scoutPost endpoint.
func (srv *Server) scoutPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	hostname := r.Header.Get("hostname")
	if hostname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing hostname"}`))
		return
	}

	ps := r.Header.Get("port")
	port, _ := strconv.Atoi(ps)
	if port == 0 {
		port = 443
	}

	is := r.Header.Get("interval")
	interval, _ := strconv.Atoi(is)
	if interval == 0 {
		interval = 60
	}

	err := srv.SetScout(hostname, port, interval)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}

// scoutDelete endpoint.
func (srv *Server) scoutDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	hostname := r.Header.Get("hostname")
	if hostname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing hostname"}`))
		return
	}

	err := srv.DeleteScout(hostname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}
