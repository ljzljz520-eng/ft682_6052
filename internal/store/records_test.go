package store

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestStoreListsFilteredRecords(t *testing.T) {
	database, err := Open(t.TempDir() + "/venue.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.PutVenue(model.Venue{ID: "A-101", Name: "A", Address: "A", Capacity: 10, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := database.PutVenue(model.Venue{ID: "B-201", Name: "B", Address: "B", Capacity: 10, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	for _, item := range []model.Reservation{
		{ID: "RSV-0001", VenueID: "A-101", SlotID: "S1", Applicant: "a", Purpose: "one", Status: "pending", Scope: "team"},
		{ID: "RSV-0002", VenueID: "B-201", SlotID: "S2", Applicant: "b", Purpose: "two", Status: "approved", Scope: "public"},
	} {
		if err := database.PutReservation(item); err != nil {
			t.Fatal(err)
		}
	}
	page, err := database.ListReservations(model.ReservationFilter{VenueID: "A-101"}, model.PageRequest{Page: 1, PageSize: 1})
	if err != nil || page.Total != 1 || len(page.Items) != 0 {
		t.Fatalf("page=%v err=%v", page, err)
	}
}
