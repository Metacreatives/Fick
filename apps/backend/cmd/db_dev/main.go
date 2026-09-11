package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	postgres := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Username("fick").
			Password("fick").
			Database("fick").
			Port(5433),
	)

	if err := postgres.Start(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := postgres.Stop(); err != nil {
			log.Printf("stop PostgreSQL: %v", err)
		}
	}()

	databaseURL := "postgresql://fick:fick@127.0.0.1:5433/fick?sslmode=disable"

	ctx := context.Background()

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		os.DirFS("db/migrations"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := provider.Up(ctx); err != nil {
		log.Fatal(err)
	}

	testData, err := os.ReadFile("db/dev/test_data.sql")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, string(testData)); err != nil {
		log.Fatal(err)
	}

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(
		signalChannel,
		os.Interrupt,
		syscall.SIGTERM,
	)

	log.Println("PostgreSQL is ready")

	<-signalChannel
}
