package service

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestVenueRegistrationAndSlotLookup(t *testing.T) {
	app := seedService(t)
	venues, err := app.ListVenues(true)
	if err != nil || len(venues) != 2 {
		t.Fatalf("venues=%v err=%v", venues, err)
	}
	slots, err := app.ListTimeSlots("A-101", true)
	if err != nil || len(slots) != 2 {
		t.Fatalf("slots=%v err=%v", slots, err)
	}
	if err := app.SetSlotAvailability(slots[0].ID, false); err != nil {
		t.Fatal(err)
	}
	available, err := app.ListTimeSlots("A-101", true)
	if err != nil || len(available) != 1 {
		t.Fatalf("available=%v err=%v", available, err)
	}
	if err := app.UpdateVenue(model.Venue{ID: "A-101", Name: "North Workshop Updated", Address: "A", Capacity: 45, Enabled: true}); err != nil {
		t.Fatal(err)
	}
}
