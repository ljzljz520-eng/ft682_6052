package service

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestCollaborationVisibilityAndReport(t *testing.T) {
	app := seedService(t)
	reservation, err := app.SubmitReservation(model.Reservation{ID: "RSV-0003", VenueID: "A-101", SlotID: "A-101-SLOT-01", Applicant: "member-b", Purpose: "review", Scope: "team"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.AddTeamNote(reservation.ID, "member-b", "agenda ready"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.AddPublicNote(reservation.ID, "member-b", "public summary"); err != nil {
		t.Fatal(err)
	}
	visible, err := app.VisibleNotes(reservation.ID, "public")
	if err != nil || len(visible) != 1 || visible[0].Body != "public summary" {
		t.Fatalf("visible=%v err=%v", visible, err)
	}
	report, err := app.BuildVenueReport("A-101")
	if err != nil || report.Total != 1 || report.Collaboration != 2 {
		t.Fatalf("report=%v err=%v", report, err)
	}
}
