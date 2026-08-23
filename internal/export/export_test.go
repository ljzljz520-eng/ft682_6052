package export

import (
	"strings"
	"testing"

	"venue-reservation/internal/model"
)

func TestReservationExportIsDeterministic(t *testing.T) {
	document := Build([]model.Reservation{{ID: "RSV-0002", VenueID: "B", SlotID: "S2", Applicant: "b", Purpose: "two", Scope: "team", Status: "approved", Revision: 1}, {ID: "RSV-0001", VenueID: "A", SlotID: "S1", Applicant: "a", Purpose: "one", Scope: "public", Status: "pending", Revision: 1}})
	data, err := document.CSV()
	if err != nil {
		t.Fatal(err)
	}
	value := string(data)
	if !strings.Contains(value, "RSV-0001") || strings.Index(value, "RSV-0001") > strings.Index(value, "RSV-0002") {
		t.Fatalf("unexpected csv=%s", value)
	}
	if _, err := document.JSON(); err != nil {
		t.Fatal(err)
	}
}
