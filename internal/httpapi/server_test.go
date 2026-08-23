package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"venue-reservation/internal/service"
	"venue-reservation/internal/store"
)

func TestHTTPVenueAndHealthRoutes(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/venue.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	server := NewServer(service.New(database))
	health := httptest.NewRecorder()
	server.Handler().ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK || !strings.Contains(health.Body.String(), "ok") {
		t.Fatalf("health=%d body=%s", health.Code, health.Body.String())
	}
	request := httptest.NewRequest(http.MethodPost, "/venues", strings.NewReader(`{"id":"A-101","name":"North","address":"A","capacity":20,"enabled":true}`))
	created := httptest.NewRecorder()
	server.Handler().ServeHTTP(created, request)
	if created.Code != http.StatusCreated {
		t.Fatalf("created=%d body=%s", created.Code, created.Body.String())
	}
}
