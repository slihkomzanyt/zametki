package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"zametki/internal/handlers"
	"zametki/internal/storage"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := storage.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	h := handlers.New(db)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      h.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return context.Background()
		},
	}

	log.Println("HTTP server started on :8080")
	log.Fatal(srv.ListenAndServe())
}
