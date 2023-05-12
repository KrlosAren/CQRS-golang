package main

import (
	"context"
	"encoding/json"
	"krlosaren/go/cqrs/events"
	"krlosaren/go/cqrs/models"
	"krlosaren/go/cqrs/repository"
	"krlosaren/go/cqrs/search"
	"log"
	"net/http"
)

func OnCreatedFeed(m events.CreatedFeedMessage) {
	feed := models.Feed{
		Id:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Printf("failed to create feed index for feed %s: %v", feed.Id, err)
	}

}

func ListFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	feeds, err := repository.ListFeeds(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	query := r.URL.Query().Get("q")

	if len(query) == 0 {
		http.Error(w, "Query is empty", http.StatusBadRequest)
	}

	feeds, err := search.SearchFeed(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)

}
