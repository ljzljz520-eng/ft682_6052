package timeline

import (
	"testing"

	"venue-reservation/internal/model"
)

func TestTimelineBuildAndSummary(t *testing.T) {
	entries := Build(model.Reservation{ID: "RSV-0001"}, []model.AuditEvent{{ID: "EVT-1", ReservationID: "RSV-0001", Action: "submitted", Actor: "member", Detail: "created", Sequence: 1}}, []model.CollaborationNote{{ID: "NOTE-1", ReservationID: "RSV-0001", Author: "member", Body: "ready", Visibility: "team", Sequence: 2}}, nil)
	filtered := Filter(entries, "team")
	if len(filtered) != 2 || !HasAction(filtered, "submitted") {
		t.Fatalf("filtered=%v", filtered)
	}
	summary := Summarize(entries)
	if summary.Entries != 2 || summary.Latest != "ready" || summary.Actors != 1 {
		t.Fatalf("summary=%v", summary)
	}
}
