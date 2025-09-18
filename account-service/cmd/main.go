package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ani213/account-service/internal/account"
	"github.com/ani213/account-service/internal/config"
	"github.com/ani213/account-service/internal/middleware"
	"github.com/ani213/account-service/internal/routes"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {

	config := config.LoadConfig()
	// dsn := "postgres://bank:password@localhost:5432/bank?sslmode=disable"
	// if os.Getenv("DATABASE_URL") != "" {
	// 	dsn = os.Getenv("DATABASE_URL")
	// }

	db, err := sqlx.Connect("postgres", config.DBSource)

	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	// // Get absolute path to migrations/
	// wd, _ := os.Getwd()
	// migrationsPath := "file://" + filepath.Join(wd, "../migrations")
	// // migrationsPath := "file://./migrations"

	// log.Println("Migrations path:", migrationsPath)
	// m, err := migrate.New(migrationsPath, dsn)
	// if err != nil {
	// 	log.Fatalf("Could not init migrate: %v", err)
	// }

	// err = m.Up()
	// if err != nil {
	// 	if err == migrate.ErrNoChange {
	// 		log.Println("No new migrations to apply (DB already up to date).")
	// 	} else {
	// 		log.Fatalf("Migration failed: %v", err)
	// 	}
	// } else {
	// 	log.Println("Migrations applied successfully!")
	// }

	// version, dirty, err := m.Version()
	// if err != nil && err != migrate.ErrNilVersion {
	// 	log.Fatalf("Could not get migration version: %v", err)
	// }
	// log.Printf("Current DB version: %d (dirty: %v)\n", version, dirty)

	repo := account.NewRepository(db)
	svc := account.NewService(repo, config)
	h := account.NewHandler(svc)

	r := mux.NewRouter()
	r.Use(middleware.JWTMiddleware(config))
	routes.Routes(r, h)
	server := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}
	//

	log.Println("Auth service running on :8081")
	// http.ListenAndServe(":8080", r)

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
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	} else {
		slog.Info("Server shutdown successfully")
	}

}
