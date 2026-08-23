package service

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestReviewAndArchiveTransitions(t *testing.T) {
	app := seedService(t)
	reservation, err := app.SubmitReservation(model.Reservation{ID: "RSV-0002", VenueID: "A-101", SlotID: "A-101-SLOT-01", Applicant: "reviewer", Purpose: "forum", Scope: "public"})
	if err != nil {
		t.Fatal(err)
	}
	approved, err := app.ApproveReservation(reservation.ID, "approver")
	if err != nil || approved.Status != "approved" {
		t.Fatalf("approved=%v err=%v", approved, err)
	}
	archive, err := app.ArchiveReservation(reservation.ID, "archiver", "completed")
	if err != nil || archive.ReservationID != reservation.ID {
		t.Fatalf("archive=%v err=%v", archive, err)
	}
	detail, err := app.GetReservationDetail(reservation.ID)
	if err != nil || detail.Reservation.Status != "archived" || detail.Archive == nil {
		t.Fatalf("detail=%v err=%v", detail, err)
	}
}
