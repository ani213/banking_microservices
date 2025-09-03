package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ani213/banking-auth/internal/auth"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://bank:password@localhost:5432/bank?sslmode=disable"
	if os.Getenv("DATABASE_URL") != "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	repo := auth.NewRepository(db)
	svc := auth.NewService(repo)
	h := auth.NewHandler(svc)

	r := mux.NewRouter()
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")

	// Example of protected route
	r.Handle("/me", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey).(string)
		w.Write([]byte("Hello user: " + userID))
	}))).Methods("GET")

	log.Println("Auth service running on :8080")
	http.ListenAndServe(":8080", r)
}
