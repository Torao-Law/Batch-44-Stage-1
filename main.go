package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hallo Batch 44"))
	}).Methods("GET")

	port := "5000"
	fmt.Println("Server running on port", port)
	http.ListenAndServe("localhost:"+port, route)
}
