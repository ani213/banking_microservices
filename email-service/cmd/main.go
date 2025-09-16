package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ani213/email-service/internal/config"
	"github.com/ani213/email-service/internal/email"
	"github.com/ani213/email-service/internal/routes"
	"github.com/gorilla/mux"
)

func main() {
	config := config.LoadConfig()
	fmt.Printf("auth service %s \n", config.AuthService)
	svc := email.NewService(config)
	h := email.NewHandler(svc)

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	routes.Routes(api, h)

	server := &http.Server{
		Addr:    ":8083",
		Handler: r,
	}
	//

	log.Printf("Email service running on Port %s", server.Addr)

	// Create a channel to listen for OS termination signals
	done := make(chan os.Signal, 1)

	// Notify the 'done' channel when an interrupt or terminate signal is received
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the HTTP server in a separate goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error:", err)
		}
	}()

	// Block the main goroutine until a signal is received
	<-done

	// Begin graceful shutdown
	slog.Info("Shutting down the server...")

	// Create a context with a timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	} else {
		slog.Info("Server shutdown successfully")
	}
}
