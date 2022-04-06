package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// canariesGet endpoint.
func (srv *Server) canariesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := srv.GetCanaries()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	if len(list.Canaries) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}

}

// canaryGet endpoint.
func (srv *Server) canaryGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.Header.Get("hostname")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing canaryid"}`))
		return
	}

}

// canaryPost endpoint.
func (srv *Server) canaryPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.Header.Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing name"}`))
		return
	}

	is := r.Header.Get("interval")
	interval, _ := strconv.Atoi(is)
	if interval == 0 {
		interval = 60
	}
	err := srv.SetCanary(name, interval)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}

// canaryDelete endpoint.
func (srv *Server) canaryDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.Header.Get("name")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing name"}`))
		return
	}

	err := srv.DeleteCanary(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}
