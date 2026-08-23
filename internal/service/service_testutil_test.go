package service

import (
	"testing"

	"venue-reservation/internal/model"
	"venue-reservation/internal/store"
)

func newTestService(t *testing.T) (*Service, *store.Store) {
	t.Helper()
	database, err := store.Open(t.TempDir() + "/venue.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	return New(database), database
}

func seedService(t *testing.T) *Service {
	app, _ := newTestService(t)
	if err := app.RegisterVenue(model.Venue{ID: "A-101", Name: "North Workshop", Address: "A", Capacity: 40, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := app.RegisterVenue(model.Venue{ID: "B-201", Name: "South Forum", Address: "B", Capacity: 80, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := app.AddTimeSlot(model.TimeSlot{ID: "A-101-SLOT-01", VenueID: "A-101", StartsAt: "09:00", EndsAt: "10:00", State: "open", Available: true}); err != nil {
		t.Fatal(err)
	}
	if err := app.AddTimeSlot(model.TimeSlot{ID: "A-101-SLOT-02", VenueID: "A-101", StartsAt: "10:00", EndsAt: "11:00", State: "open", Available: true}); err != nil {
		t.Fatal(err)
	}
	return app
}
