package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) httpHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/sync", s.syncHandler)
	return r
}

func (s *server) syncHandler(w http.ResponseWriter, r *http.Request) {
	forceStr := r.URL.Query().Get("force")
	force := false
	if forceStr != "" {
		var err error
		force, err = strconv.ParseBool(forceStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid param force: %v", err)
			return
		}
	}
	fmt.Println(force)

	err := s.doSync(force)
	if err != nil {
		log.Infof("Sync failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Sync failed: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")
}
