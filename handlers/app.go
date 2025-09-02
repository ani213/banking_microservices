package handlers

import (
	"net/http"

	"github.com/ani213/banking-microservices/domain"
	"github.com/ani213/banking-microservices/service"
	"github.com/gorilla/mux"
)

func Start() {
	// mux := http.NewServeMux()
	mux := mux.NewRouter()
	// wiring
	ch := CustomersHandler{service: service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	mux.HandleFunc("/customers", ch.GetAllCustomers).Methods(http.MethodGet)
	mux.HandleFunc("/xml", XMLHandler).Methods("GET")
	mux.HandleFunc("/customer/{id:[0-9]+}", CustomersByIDHandler)
	http.ListenAndServe(":8080", mux)

}
