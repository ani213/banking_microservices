package handlers

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Customer struct {
	Name string `json:"name" xml:"name"` // json tag for serialization
	Age  int    `json:"age" xml:"age"`   // json tag for serialization
}

func CustomersHandler(w http.ResponseWriter, r *http.Request) {
	customer := []Customer{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(customer)
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
