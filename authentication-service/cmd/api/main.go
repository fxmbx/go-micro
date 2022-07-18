package main

import (
	data "authentication/Data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	app := &Config{}

	conn := connectDb()
	if conn == nil {
		log.Panic("Couldn't connect to postgres 😞")
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectDb() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDb(dsn)
		if err != nil {
			log.Println("postrgres not conntected yet 😬")
			counts++
		} else {
			log.Println("postrgres conntected 😁")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Chilling for 2 seconds before tryng again ⏲️")
		time.Sleep(2 * time.Second)
		continue
	}

}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
