package main

import (
	"fmt"
	"krlosaren/go/cqrs/database"
	"krlosaren/go/cqrs/events"
	"krlosaren/go/cqrs/repository"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB       string `envconfig: "POSTGRES_DB"`
	PostgresUser     string `envconfig: "POSTGRES_USER"`
	PostgresPassword string `envconfig: "POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig: "NATS_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", createdFeedHandler).Methods(http.MethodPost)
	return
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatalf("%v", err)
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	repo, err := database.NewPostgresRepository(addr)

	if err != nil {
		log.Fatalf("%v", err)
	}

	repository.SetRepository(repo)

	n, err := events.NewNats(fmt.Sprintf("nats://%s", os.Getenv("NATS_ADDRESS")))

	if err != nil {
		log.Fatalf("%v", err)
	}

	events.SetEventStore(n)

	defer events.Close()

	router := newRouter()

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("%v", err)
	}
}
