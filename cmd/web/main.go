package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"snippetbox.com/pkg/models/pgx"
)

type application struct {

	errorLog *log.Logger
	infoLog  *log.Logger
	snippet  *pgx.SnippetModel

}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	connStr := flag.String("connStr", "postgresql://localhost/godb?user=postgres&password=1337", "PostGreSQL")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*connStr)

	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippet:  &pgx.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting  server on %v", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDB(dsn string) (*pgxpool.Pool, error) {

	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return db, nil

}