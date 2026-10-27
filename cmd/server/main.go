package main

import (
	"log"
	"net/http"
	"time"

	"github.com/one-thousand-twoth/calendar/internal/booking"
)

func main() {
	router, err := booking.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
