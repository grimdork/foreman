package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/grimdork/foreman/api"
)

// keysGet endpoint.
func (srv *Server) keysGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := srv.GetKeys()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	keys := api.KeyList{Keys: list}
	if len(keys.Keys) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	err = json.NewEncoder(w).Encode(keys)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}

// keyPost endpoint.
func (srv *Server) keyPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.Header.Get("keyid")
	value := r.Header.Get("value")
	if id == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing keyid or value"}`))
		return
	}

	admin := r.Header.Get("admin")
	adminbool := false
	if admin != "" {
		adminbool, _ = strconv.ParseBool(admin)
	}

	err := srv.SetKey(id, value, adminbool)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}

// keyDelete endpoint.
func (srv *Server) keyDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.Header.Get("keyid")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing keyid"}`))
		return
	}

	err := srv.DeleteKey(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	}
}
