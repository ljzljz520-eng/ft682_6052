package store

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/venue.db"
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.PutVenue(model.Venue{ID: "A-101", Name: "North Workshop", Address: "A", Capacity: 40, Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := database.PutTimeSlot(model.TimeSlot{ID: "A-101-SLOT-01", VenueID: "A-101", StartsAt: "09:00", EndsAt: "10:00", State: "open", Available: true}); err != nil {
		t.Fatal(err)
	}
	if err := database.PutReservation(model.Reservation{ID: "RSV-0001", VenueID: "A-101", SlotID: "A-101-SLOT-01", Applicant: "member", Purpose: "meeting", Status: "pending", Scope: "team", Revision: 1}); err != nil {
		t.Fatal(err)
	}
	if err := database.PutAuditEvent(model.AuditEvent{ID: "EVT-0001", ReservationID: "RSV-0001", Action: "submitted", Actor: "member", Detail: "created", Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := database.PutCollaborationNote(model.CollaborationNote{ID: "NOTE-0001", ReservationID: "RSV-0001", Author: "member", Body: "ready", Visibility: "team", Sequence: 2}); err != nil {
		t.Fatal(err)
	}
	if err := database.PutArchiveRecord(model.ArchiveRecord{ID: "ARC-0001", ReservationID: "RSV-0001", Reason: "done", ArchivedBy: "member", ArchivedAt: "sequence"}); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if _, err := reopened.GetVenue("A-101"); err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.GetTimeSlot("A-101-SLOT-01"); err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.GetReservation("RSV-0001"); err != nil {
		t.Fatal(err)
	}
	archive, err := reopened.GetArchive("RSV-0001")
	if err != nil || archive == nil {
		t.Fatalf("archive=%v err=%v", archive, err)
	}
}
