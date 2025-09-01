package handlers

import (
	"fmt"
	"net/http"
)

func Start() {
	http.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	http.HandleFunc("/customers", CustomersHandler)
	http.HandleFunc("/xml", XMLHandler)
	http.ListenAndServe(":8080", nil)

}
