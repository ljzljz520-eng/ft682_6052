package integration

import (
	"testing"

	"venue-reservation/internal/model"
	"venue-reservation/internal/service"
	"venue-reservation/internal/store"
)

func workflowService(t *testing.T) *service.Service {
	t.Helper()
	database, err := store.Open(t.TempDir() + "/venue.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	app := service.New(database)
	if err := app.RegisterVenue(model.Venue{ID: "A-101", Name: "North", Address: "A", Capacity: 20, Enabled: true}); err != nil {
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

func TestWorkflowCreateReviewArchive(t *testing.T) {
	app := workflowService(t)
	reservation, err := app.SubmitReservation(model.Reservation{VenueID: "A-101", SlotID: "A-101-SLOT-01", Applicant: "planner", Purpose: "workshop", Scope: "team"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.AddTeamNote(reservation.ID, "planner", "draft agenda"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.ApproveReservation(reservation.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.ArchiveReservation(reservation.ID, "operator", "finished"); err != nil {
		t.Fatal(err)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	app := workflowService(t)
	reservation, err := app.SubmitReservation(model.Reservation{VenueID: "A-101", SlotID: "A-101-SLOT-01", Applicant: "editor", Purpose: "planning", Scope: "team"})
	if err != nil {
		t.Fatal(err)
	}
	page, err := app.SearchReservations("planning", model.PageRequest{Page: 1, PageSize: 10})
	if err != nil || page.Total != 1 {
		t.Fatalf("page=%v err=%v", page, err)
	}
	if _, err := app.SelectReservation(reservation.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := app.UpdateReservation(reservation.ID, "planning published", "owner confirmed"); err != nil {
		t.Fatal(err)
	}
	report, err := app.PublishReport("A-101")
	if err != nil || report.Total != 1 {
		t.Fatalf("report=%v err=%v", report, err)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	app := workflowService(t)
	rows := []model.ImportRow{
		{Reference: "row-1", VenueID: "A-101", SlotID: "A-101-SLOT-01", Applicant: "importer", Purpose: "orientation"},
		{Reference: "row-2", VenueID: "unknown", SlotID: "missing", Applicant: "importer", Purpose: "rejected"},
	}
	issues := app.ValidateImport(rows)
	if len(issues) != 1 {
		t.Fatalf("issues=%v", issues)
	}
	report, err := app.ImportReservations(rows)
	if err != nil || report.Accepted != 1 || report.Rejected != 1 {
		t.Fatalf("report=%v err=%v", report, err)
	}
	if _, err := app.PublishReport("A-101"); err != nil {
		t.Fatal(err)
	}
}
