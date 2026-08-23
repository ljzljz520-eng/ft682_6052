package service

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestReservationSubmissionUpdateAndDetail(t *testing.T) {
	app := seedService(t)
	reservation, err := app.SubmitReservation(model.Reservation{ID: "RSV-0001", VenueID: "A-101", SlotID: "A-101-SLOT-01", Applicant: "member-a", Purpose: "training", Scope: "team"})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := app.UpdateReservation(reservation.ID, "training revised", "bring projector")
	if err != nil || updated.Revision != 2 {
		t.Fatalf("updated=%v err=%v", updated, err)
	}
	detail, err := app.GetReservationDetail(reservation.ID)
	if err != nil || detail.Reservation.Purpose != "training revised" || len(detail.Events) != 2 {
		t.Fatalf("detail=%v err=%v", detail, err)
	}
	if !IsVisibleToMember(detail.Reservation, "member-a") {
		t.Fatal("member should see own reservation")
	}
}
