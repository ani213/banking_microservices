package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	// mux := http.NewServeMux()
	mux := mux.NewRouter()
	mux.HandleFunc("/customers", CustomersHandler)
	mux.HandleFunc("/xml", XMLHandler)
	http.ListenAndServe(":8080", mux)

}
