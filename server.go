package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) httpHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/sync", s.syncHandler)
	r.HandleFunc("/api/docker/stacks", s.dockerStackListHandler)
	r.HandleFunc("/api/docker/services", s.dockerServiceListHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))
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

func (s *server) dockerStackListHandler(w http.ResponseWriter, r *http.Request) {
	stacks, err := NewDockerStackCmd().ls()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	err = json.NewEncoder(w).Encode(stacks)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
}

func (s *server) dockerServiceListHandler(w http.ResponseWriter, r *http.Request) {
	stackName := r.URL.Query().Get("stack")
	stacks, err := NewDockerStackCmd().services(stackName)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	err = json.NewEncoder(w).Encode(stacks)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
}
