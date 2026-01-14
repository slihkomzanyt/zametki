package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/pressly/goose/v3"

	"zametki/internal/handlers"
	"zametki/internal/storage"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// 1) подключаемся к БД через твой storage (внутри должен быть *sql.DB)
	pg, err := storage.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	// 2) migrations (auto)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	// локально: ./migrations, на сервере: /opt/zametki/migrations (через env)
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./migrations"
	}

	// ВАЖНО: goose принимает *sql.DB, поэтому pg.DB
	// (pg.DB должен быть экспортируемым полем типа *sql.DB внутри storage.Postgres)
	if err := goose.Up(pg.DB, migrationsDir); err != nil {
		log.Fatal(err)
	}

	// 3) http
	h := handlers.New(pg)

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

// гарантируем, что импорт database/sql не “умрёт”,
// если IDE ругнётся до того как ты используешь sql где-то ещё.
var _ *sql.DB
