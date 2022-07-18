package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "80"

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s 😏\n", port)

	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	//start server
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
