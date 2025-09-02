package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/ani213/banking-microservices/service"
	"github.com/gorilla/mux"
)

type Customer struct {
	Name string `json:"name" xml:"name"` // json tag for serialization
	Age  int    `json:"age" xml:"age"`   // json tag for serialization
}

type CustomersHandler struct {
	service service.CustomerService
}

func (ch *CustomersHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := ch.service.GetAllCustomers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Some error occurred"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func XMLHandler(w http.ResponseWriter, r *http.Request) {
	customer := []Customer{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	if r.Header.Get("Content-Type") == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(customer)
		return
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customer)
		return
	}

}

func CustomersByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "Customer by id:", id)
}
