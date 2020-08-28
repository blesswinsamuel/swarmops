package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func httpHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/sync", syncHandler)
	return r
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	err := doSync()
	if err != nil {
		log.Printf("Sync failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Sync failed: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
