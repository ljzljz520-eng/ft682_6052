package service

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestBusiness001Regression(t *testing.T) {
	app := seedService(t)
	for index := 1; index <= 4; index++ {
		id := "RSV-000" + string(rune('0'+index))
		if _, err := app.SubmitReservation(model.Reservation{ID: id, VenueID: "A-101", SlotID: "A-101-SLOT-01", Applicant: "member", Purpose: "session " + id, Scope: "team"}); err != nil {
			t.Fatal(err)
		}
	}
	page, err := app.FindReservations(model.ReservationFilter{VenueID: "A-101"}, model.PageRequest{Page: 1, PageSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 2 || page.Items[0].ID != "RSV-0001" || page.Items[1].ID != "RSV-0002" {
		t.Fatalf("unexpected first page: %+v", page.Items)
	}
}
